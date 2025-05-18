package znet

import (
	"net"
	"sync"
)

// NewWhiteListListener returns a new instance of [WhiteListListener].
// Returned listener allows new connections only from given networks.
// Given networks must be valid form for [net/netip.ParsePrefix].
// Both IPv4 and IPV6 are acceptable. For example, "0.0.0.0/0" allows all
// new connections and "127.0.0.1/32" allows only specified IP.
// Invalid addresses such as "/16" (CIDR only) or "127.0.0.1" (address only)
// will results in an error because they are not accepted by [net/netip.ParsePrefix].
// See also [WhiteListListener] and [WhiteList].
func NewWhiteListListener(ln net.Listener, allow ...string) (*WhiteListListener, error) {
	wl := NewWhiteList()
	if err := wl.Allow(allow...); err != nil {
		return nil, err
	}
	return &WhiteListListener{
		Listener: ln,
		WL:       wl,
	}, nil
}

// WhiteListListener is the listener that accepts new
// connections allowed by the whitelist.
// Whitelist is based on IP addresses.
type WhiteListListener struct {
	net.Listener
	// WL is the whitelist.
	// New connections allowed by the list is acceptable.
	// Otherwise, the new connection will be force closed.
	// WL must not be nil.
	WL *WhiteList
	// Allow optionally judges if the connection is allowed or not.
	// If non-nil, it is called with host and port of the remote addr.
	// Given host is allowed by the whitelist.
	Allow func(host, port string) bool
}

func (l *WhiteListListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, err
	}
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		_ = conn.Close()
		return conn, nil
	}
	if host == "" || !l.WL.Allowed(host) {
		_ = conn.Close()
		return conn, nil
	}
	if allow := l.Allow; allow != nil && !allow(host, port) {
		_ = conn.Close()
		return conn, nil
	}
	return conn, nil
}

// NewBlackListListener returns a new instance of [BlackListListener].
// Returned listener does not allow new connections from given networks.
// Given networks must be valid form for [net/netip.ParsePrefix].
// Both IPv4 and IPV6 are acceptable. For example, "0.0.0.0/0" does not allow
// any new connections and "127.0.0.1/32" does not allow only specified IP.
// Invalid addresses such as "/16" (CIDR only) or "127.0.0.1" (address only)
// will results in an error because they are not accepted by [net/netip.ParsePrefix].
// See also [BlackListListener] and [BlackList].
func NewBlackListListener(ln net.Listener, disallow ...string) (*BlackListListener, error) {
	bl := NewBlackList()
	if err := bl.Disallow(disallow...); err != nil {
		return nil, err
	}
	return &BlackListListener{
		Listener: ln,
		BL:       bl,
	}, nil
}

// BlackListListener is the listener that accepts new
// connections allowed by the blacklist.
// BlackList is based on IP addresses.
type BlackListListener struct {
	net.Listener
	// BL is the blacklist.
	// New connections allowed by the list is acceptable.
	// Otherwise, the new connection will be force closed.
	// BL must not be nil.
	BL *BlackList
	// Allow optionally judges if the connection is allowed or not.
	// If non-nil, it is called with host and port of the remote addr.
	// Given host is allowed by the blacklist.
	Allow func(host, port string) bool
}

func (l *BlackListListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, err
	}
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		_ = conn.Close()
		return conn, nil
	}
	if host == "" || !l.BL.Allowed(host) {
		_ = conn.Close()
		return conn, nil
	}
	if allow := l.Allow; allow != nil && !allow(host, port) {
		_ = conn.Close()
		return conn, nil
	}
	return conn, nil
}

// NewLimitListener returns a net listener with maximum concurrent
// connections. For limit<1, 1 is used.
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
	// A semaphore is obtained by sending struct{} to the sem
	// and is released by removing struct{} from the sem.
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
