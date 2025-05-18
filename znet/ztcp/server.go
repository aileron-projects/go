package ztcp

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aileron-projects/go/znet"
	"github.com/aileron-projects/go/znet/internal"
)

var (
	// ErrAbortHandler is a sentinel panic value to abort a handler.
	// While any panic from ServeTCP aborts the response to the client,
	// panicking with ErrAbortHandler also suppresses logging the error and stacktraces.
	ErrAbortHandler = errors.New("znet/ztcp: abort Handler")
	// ErrSkipHandler is a sentinel error to skip serving clients' connection.
	// If ErrSkipHandler is returned from the net.Listener, client connection is
	// immediately closed and ServeTCP is not called for the connection.
	ErrSkipHandler = errors.New("znet/ztcp: skip Handler")
)

// A Handler responds to a TCP request.
//
// ServeTCP is invoked in a new goroutine for each new incoming connections.
// ServeTCP does not need to close the connection because the [Server] ensures to.
// If ServeTCP panics, the server (the caller of ServeTCP) assumes that
// the effect of the panic was isolated to the active connections.
// It recovers the panic, logs a stack trace to the server error log,
// and closes the network connection. Panicking [ErrAbortHandler] suppresses
// logging the error and stack traces.
type Handler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}

// HandlerFunc implements [Handler] interface to the function.
type HandlerFunc func(ctx context.Context, conn net.Conn)

func (f HandlerFunc) ServeTCP(ctx context.Context, conn net.Conn) {
	f(ctx, conn)
}

// Server is a TCP server.
type Server struct {
	// Addr is the address to listen to.
	// Network prefix "tcp", "tcp4", "tcp6", "unix" and "unixpacket"
	// can be specified with the form of "<PREFIX>://<ADDRESS>".
	// For example "tcp4://localhost:8080".
	// "tcp" is assumed when no network prefix found.
	// Addr is used by [Server.ListenAndServe] and [Server.ListenAndServeTLS].
	Addr string

	// Handler to invoke.
	// Handler must not be nil.
	Handler Handler

	// TLSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS. Note that this value is
	// cloned by ServeTLS and ListenAndServeTLS, so it's not
	// possible to modify the configuration with methods like
	// tls.Config.SetSessionTicketKeys. To use
	// SetSessionTicketKeys, use Server.Serve with a TLS Listener instead.
	TLSConfig *tls.Config

	// BaseContext optionally specifies a function that returns
	// the base context for incoming requests on this server.
	// The provided Listener is the specific Listener that's
	// about to start accepting requests.
	// If BaseContext is nil, the default is context.Background().
	// If non-nil, it must return a non-nil context.
	BaseContext func(net.Listener) context.Context

	// PanicHandler optionally handles panic.
	// Recovered value, which always non-nil, and remote and local addresses
	// are provided. The sentinel error [ErrAbortHandler] is also given to
	// the PanicHandler. It bypasses default logging of stacktraces.
	PanicHandler func(recovered any, remote, local net.Addr)

	shutdown  atomic.Bool
	listeners internal.UniqueStore[*ocListener]
	conns     internal.UniqueStore[*ocConn]

	// serveNotify notifies the inner listener is working
	// and the Server.Serve is called.
	// serveNotify is used for testing.
	serveNotify chan struct{}
}

// ListenAndServe listens on the TCP network address s.Addr and then calls
// [Server.Serve] with handler to handle incoming connections.
//
// ListenAndServeTLS always returns a non-nil error.
// After [Server.Shutdown] or [Server.Close], the returned error is [net.ErrClosed].
func (s *Server) ListenAndServe() error {
	if s.shutdown.Load() {
		return net.ErrClosed
	}
	ln, err := newListener(s.Addr)
	if err != nil {
		return err
	}
	return s.Serve(ln)
}

// ListenAndServeTLS listens on the TCP network address s.Addr and then calls
// [Server.ServeTLS] with handler to handle incoming TLS connections.
//
// Filenames containing a certificate and matching private key for the server
// must be provided if neither the Server's TLSConfig.Certificates nor TLSConfig.GetCertificate
// are populated. If the certificate is signed by a certificate authority, the certFile
// should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.
//
// ListenAndServeTLS always returns a non-nil error.
// After [Server.Shutdown] or [Server.Close], the returned error is [net.ErrClosed].
func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if s.shutdown.Load() {
		return net.ErrClosed
	}
	ln, err := newListener(s.Addr)
	if err != nil {
		return err
	}
	return s.ServeTLS(ln, certFile, keyFile)
}

// ServeTLS accepts incoming connections on the Listener l, creating a new
// service goroutine for each. The service goroutines perform TLS setup and then
// read requests, calling s.Handler to reply to them.
//
// Filenames containing a certificate and matching private key for the server
// must be provided if neither the Server's TLSConfig.Certificates nor TLSConfig.GetCertificate
// are populated. If the certificate is signed by a certificate authority, the certFile
// should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.
//
// ServeTLS always returns a non-nil error.
// After [Server.Shutdown] or [Server.Close], the returned error is [net.ErrClosed].
func (s *Server) ServeTLS(l net.Listener, certFile, keyFile string) error {
	if s.shutdown.Load() {
		return net.ErrClosed
	}
	var config *tls.Config
	if s.TLSConfig != nil {
		config = s.TLSConfig.Clone()
	} else {
		config = &tls.Config{}
	}
	if certFile != "" || keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
		config.Certificates = append(config.Certificates, cert)
	}
	l = tls.NewListener(l, config)
	return s.Serve(l)
}

// Serve accepts incoming connections on the Listener l, creating a new
// service goroutine for each. The service goroutines read requests,
// calling s.Handler to reply to them.
//
// Serve always returns a non-nil error.
// After [Server.Shutdown] or [Server.Close], the returned error is [net.ErrClosed].
func (s *Server) Serve(l net.Listener) error {
	if s.shutdown.Load() {
		return net.ErrClosed
	}
	if s.serveNotify != nil {
		s.serveNotify <- struct{}{}
	}

	ctx := context.Background()
	if bc := s.BaseContext; bc != nil {
		ctx = bc(l)
	}

	ln := &ocListener{
		Listener: l,
		store:    &s.listeners,
	}
	s.listeners.Set(ln)
	defer ln.Close()

	wait := int64(1)
	for {
		conn, err := l.Accept()
		if err != nil {
			if err == ErrSkipHandler {
				_ = conn.Close()
				continue
			}
			if s.shutdown.Load() { // Error is caused by shutdown.
				return net.ErrClosed
			}
			// Check if this is caused by timeout or not. [net.Error] implements the interface.
			if to, ok := err.(interface{ Timeout() bool }); ok && to.Timeout() {
				wait = min(wait*2, 1<<9) // Up to 512 msec.
				time.Sleep(time.Duration(wait) * time.Millisecond)
				continue
			}
			return err
		}
		wait = 1 // Reset
		go s.serve(ctx, conn)
	}
}

func (s *Server) serve(ctx context.Context, conn net.Conn) {
	c := &ocConn{
		Conn:  conn,
		store: &s.conns,
	}
	s.conns.Set(c)
	defer func() {
		_ = c.Close()
		err := recover()
		if err == nil {
			return
		}
		if ph := s.PanicHandler; ph != nil {
			ph(err, conn.RemoteAddr(), conn.LocalAddr())
			return
		}
		if err != ErrAbortHandler {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("znet/ztcp: panic serving %v: %v\n%s", conn.RemoteAddr(), err, buf)
		}
	}()
	s.Handler.ServeTCP(ctx, conn)
}

// Close immediately closes all active net.Listeners and any connections.
// For a graceful shutdown, use [Server.Shutdown].
//
// When Close is called, [Server.Serve], [Server.ServeTLS], [Server.ListenAndServe]
// and [Server.ListenAndServeTLS] immediately return [net.ErrClosed].
//
// Once Close has been called on a server, it may not be reused;
// future calls to methods such as [Server.Serve] will return [net.ErrClosed].
func (s *Server) Close() error {
	s.shutdown.Store(true)
	var errs []error
	for ocl := range s.listeners.Values() {
		errs = appendNonNil(errs, ocl.Close())
	}
	for occ := range s.conns.Values() {
		errs = appendNonNil(errs, occ.Close())
	}
	return errors.Join(errs...)
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
// Shutdown works by first closing all open [net.Listener]s,
// and then waiting for connections to be closed and then shut down.
// If the provided context expires before the shutdown is complete, Shutdown returns
// the context's error, otherwise it returns all errors returned from closing the
// Server's underlying Listener(s). Non-nil errors occurred while shutting down
// are returned after joined with [errors.Join].
//
// When Shutdown is called, [Server.Serve], [Server.ServeTLS], [Server.ListenAndServe]
// and [Server.ListenAndServeTLS] immediately return [net.ErrClosed].
// Make sure the program doesn't exit and waits instead for Shutdown to return.
//
// Once Shutdown has been called on a server, it may not be reused;
// future calls to methods such as [Server.Serve] will return [net.ErrClosed].
func (s *Server) Shutdown(ctx context.Context) error {
	if s.shutdown.Swap(true) {
		return net.ErrClosed
	}
	var errs []error
	for ocl := range s.listeners.Values() {
		errs = appendNonNil(errs, ocl.Close())
	}
	for s.conns.Length() > 0 {
		select {
		case <-time.After(100 * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return errors.Join(errs...)
}

func appendNonNil(errs []error, err error) []error {
	if err == nil {
		return errs
	}
	return append(errs, err)
}

func newListener(addr string) (ln net.Listener, err error) {
	network, address := znet.ParseNetAddr(addr)
	switch network {
	case "":
		network = "tcp" // Assume "tcp".
		fallthrough
	case "tcp", "tcp4", "tcp6":
		laddr, resolvErr := net.ResolveTCPAddr(network, address)
		if resolvErr != nil {
			return nil, resolvErr
		}
		return net.ListenTCP(network, laddr)
	case "unix", "unixpacket":
		laddr := &net.UnixAddr{Name: address, Net: network}
		return net.ListenUnix(network, laddr)
	default:
		return net.Listen("tcp", addr) // Fallback. May be invalid addr.
	}
}

// ocListener is once close Listener that wraps [net.Listener],
// protecting it from multiple Close calls.
type ocListener struct {
	net.Listener
	once     sync.Once
	store    *internal.UniqueStore[*ocListener]
	closeErr error
}

func (oc *ocListener) Close() error {
	oc.once.Do(func() {
		defer oc.store.Delete(oc)
		oc.closeErr = oc.Listener.Close()
		addr := oc.Addr()
		if _, ok := addr.(*net.UnixAddr); !ok {
			return // Non unix socket.
		}
		address := addr.String()
		if len(address) > 0 && address[0] == '@' {
			return // Abstract socket.
		}
		info, err := os.Stat(address)
		if err == nil && !info.IsDir() {
			_ = os.Remove(address) // Remove socket file.
		}
	})
	return oc.closeErr
}

// ocConn is once close Conn that wraps a [net.Conn],
// protecting it from multiple Close calls.
type ocConn struct {
	net.Conn
	once     sync.Once
	store    *internal.UniqueStore[*ocConn]
	closeErr error
}

func (oc *ocConn) Close() error {
	oc.once.Do(func() {
		defer oc.store.Delete(oc)
		oc.closeErr = oc.Conn.Close()
	})
	return oc.closeErr
}
