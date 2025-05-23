package znet

import (
	"cmp"
	"crypto/tls"
	"net"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/acme/autocert"
)

// NewWhiteListListener returns a new instance of [WhiteListListener].
// Returned listener allows new connections only from given networks.
// Given networks must be valid form for [net/netip.ParsePrefix].
// Both IPv4 and IPV6 are acceptable. For example, "0.0.0.0/0" allows all
// new connections and "127.0.0.1/32" allows only specified IP.
// Invalid addresses such as "/16" (CIDR only) or "127.0.0.1" (address only)
// will results in an error because they are not accepted by [net/netip.ParsePrefix].
// See also [WhiteListListener] and [WhiteList].
func NewWhiteListListener(ln net.Listener, allow ...string) (*AllowListener, error) {
	wl := NewWhiteList()
	if err := wl.Allow(allow...); err != nil {
		return nil, err
	}
	return &AllowListener{
		Listener: ln,
		Allowed: func(host, port string) bool {
			return wl.Allowed(host)
		},
	}, nil
}

// NewBlackListListener returns a new instance of [BlackListListener].
// Returned listener does not allow new connections from given networks.
// Given networks must be valid form for [net/netip.ParsePrefix].
// Both IPv4 and IPV6 are acceptable. For example, "0.0.0.0/0" does not allow
// any new connections and "127.0.0.1/32" does not allow only specified IP.
// Invalid addresses such as "/16" (CIDR only) or "127.0.0.1" (address only)
// will results in an error because they are not accepted by [net/netip.ParsePrefix].
// See also [BlackListListener] and [BlackList].
func NewBlackListListener(ln net.Listener, disallow ...string) (*AllowListener, error) {
	bl := NewBlackList()
	if err := bl.Disallow(disallow...); err != nil {
		return nil, err
	}
	return &AllowListener{
		Listener: ln,
		Allowed: func(host, port string) bool {
			return bl.Allowed(host)
		},
	}, nil
}

// AllowListener is the listener that accepts connections
// allowed by the client address.
// Basically, use [NewWhiteListListener] or [NewBlackListListener]
// to create a new AllowListener.
type AllowListener struct {
	net.Listener
	// Allowed returns if the connection should be allowed or not.
	// A connection will be immediately closed when Allowed returned false.
	// Allowed must not be nil.
	Allowed func(host, port string) bool
}

func (l *AllowListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, err
	}
	addr := conn.RemoteAddr().String()
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		host, port = addr, "" // Fallback
	}
	if allowed := l.Allowed; !allowed(host, port) {
		_ = conn.Close()
		return conn, nil
	}
	return conn, nil
}

// NewLimitListener returns a net listener with maximum concurrent
// connections. For limit<1, 1 is used.
// See also [LimitListener].
func NewLimitListener(ln net.Listener, limit int) *LimitListener {
	return &LimitListener{
		Listener: ln,
		sem:      make(chan struct{}, max(1, limit)),
	}
}

// LimitListener limits the number of simultaneous connection.
// Use [NewLimitListener] to create a new instance of LimitListener.
//
// Note that the linux command "netstat" or "ss" like below
// does not show the correct number of connections currently accepted.
//   - netstat -uant | grep ESTABLISHED | grep 8080 | wc
//   - ss -o state established "( dport = :8080 )" -np | wc
//
// This is described in https://github.com/golang/go/issues/36212#issuecomment-567838193
// Use "lsof" command instead. For example,
//   - lsof -i:8080 | grep foobar
type LimitListener struct {
	net.Listener
	// sem is the semaphore variable.
	// sem is a buffered channel.
	sem chan struct{}
}

func (l *LimitListener) Accept() (net.Conn, error) {
	l.sem <- struct{}{} // Obtain semaphore.
	conn, err := l.Listener.Accept()
	if err != nil {
		<-l.sem // Release semaphore.
		return conn, err
	}
	return &limitListenerConn{
		Conn:    conn,
		release: func() { <-l.sem },
	}, nil
}

// limitListenerConn is the connection generated by limitListener with release function.
// release function is called when the connection is disconnected.
// Any occupied resources such as semaphore should be released in the release function.
type limitListenerConn struct {
	net.Conn
	once    sync.Once
	release func()
}

// Close closes the connection and release resources.
func (l *limitListenerConn) Close() error {
	err := l.Conn.Close()
	l.once.Do(l.release)
	return err
}

// NewTLSListener creates a new net listener that returns
// a *tls.Conn connection.
// IP addresses provided by the nonTLS arguments are considered
// as non-TLS connection.
// [WhiteList] is used internally for nonTLS to determine the
// connection should be TLS or non-TLS.
func NewTLSListener(ln net.Listener, c *tls.Config, nonTLS ...string) (*TLSListener, error) {
	wl := NewWhiteList()
	if err := wl.Allow(nonTLS...); err != nil {
		return nil, err
	}
	return &TLSListener{
		Listener: ln,
		NonTLS: func(host, port string) bool {
			return wl.Allowed(host)
		},
	}, nil
}

// TLSListener is the listener that accepts
// new TLS connections.
type TLSListener struct {
	net.Listener
	// TLSConfig is the configuration applied for
	// new TLS connections.
	TLSConfig *tls.Config
	// NonTLS optionally judges if the connection is non-TLS.
	// Users who does not use NonTLS, use [crypto/tls.NewListener]
	// instead.
	// [BlackList] and [WhiteList] can be used.
	NonTLS func(host, port string) bool
}

func (l *TLSListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, err
	}
	addr := conn.RemoteAddr().String()
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		host, port = addr, "" // Fallback
	}
	if nonTLS := l.NonTLS; nonTLS(host, port) {
		return conn, nil
	}
	return tls.Server(conn, l.TLSConfig), nil
}

// NewACMEListener returns a net listener that returns
// *tls.Conn connections with LetsEncrypt certificates
// for the provided domain or domains.
// See also [golang.org/x/crypto/acme/autocert.NewListener].
func NewACMEListener(ln net.Listener, domains ...string) *ACMEListener {
	m := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
	}
	if len(domains) > 0 {
		m.HostPolicy = autocert.HostWhitelist(domains...)
	}
	dir, _ := os.UserCacheDir()
	dir = cmp.Or(dir, "/.cache") // Fall back to the root directory.
	dir = filepath.Join(dir, "autocert")
	println(dir)
	if err := os.MkdirAll(dir, os.ModePerm); err == nil {
		m.Cache = autocert.DirCache(dir)
	}
	return &ACMEListener{
		Listener: ln,
		Manager:  m,
	}
}

// ACMEListener is a listener that applies
// auto certification, ACME.
// See also [golang.org/x/crypto/acme/autocert.Manager]
type ACMEListener struct {
	net.Listener
	// Manager is a autocert manager that
	// provides a TLSConfig.
	// Manager must not be nil.
	Manager *autocert.Manager
	// Modifier optionally specifies a function to modify
	// the TLSConfig generated by the [autocert.Manager.TLSConfig].
	//
	// Example:
	//
	// 	func(c *tls.Config) {
	// 		c.NextProtos = []string{
	// 			"h2", "http/1.1", // enable HTTP/2
	// 			acme.ALPNProto, // enable tls-alpn ACME challenges
	// 		}
	// 	}
	Modifier func(*tls.Config)
}

func (l *ACMEListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, err
	}
	tlsConfig := l.Manager.TLSConfig()
	if modify := l.Modifier; modify != nil {
		modify(tlsConfig)
	}
	return tls.Server(conn, tlsConfig), nil
}
