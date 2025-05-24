package zudp

import (
	"cmp"
	"context"
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

const mtu = 65535 // UDP maximum transmission unit

var (
	// ErrAbortHandler is a sentinel panic value to abort a handler.
	// While any panic from ServeUDP aborts the response to the client,
	// panicking with ErrAbortHandler also suppresses logging the error and stacktraces.
	ErrAbortHandler = errors.New("znet/zudp: abort Handler")
	// ErrSkipHandler is a sentinel error to skip serving clients' connection.
	// If ErrSkipHandler is returned from the net.Listener, received packet is
	// immediately discarded and not provided to ServeUDP.
	ErrSkipHandler = errors.New("znet/zudp: skip Handler")
)

// A Handler responds to a UDP request.
//
// ServeUDP is invoked in a new goroutine for each new incoming connections.
// A connection is created for each client's address:port pair.
// ServeUDP does not need to close the connection because the [Server] ensures to.
// If ServeUDP panics, the server (the caller of ServeUDP) assumes that
// the effect of the panic was isolated to the active connections.
// It recovers the panic, logs a stack trace to the server error log,
// and closes the network connection. Panicking [ErrAbortHandler] suppresses
// logging the error and stack traces.
type Handler interface {
	ServeUDP(ctx context.Context, conn Conn)
}

// HandlerFunc implements [Handler] interface to the function.
type HandlerFunc func(ctx context.Context, conn Conn)

func (f HandlerFunc) ServeUDP(ctx context.Context, conn Conn) {
	f(ctx, conn)
}

// Server is a UDP server.
type Server struct {
	// Addr is the address to listen to.
	// Network prefix "udp", "udp4", "udp6" and "unixgram"
	// can be specified with the form of "<PREFIX>://<ADDRESS>".
	// For example "udp4://localhost:8080".
	// "udp" is assumed when no network prefix found.
	// Addr is used by [Server.ListenAndServe].
	Addr string

	// Handler to invoke.
	// Handler must not be nil.
	Handler Handler

	// BaseContext optionally specifies a function that returns
	// the base context for incoming requests on this server.
	// The provided Listener is the specific Listener that's
	// about to start accepting requests.
	// If BaseContext is nil, the default is context.Background().
	// If non-nil, it must return a non-nil context.
	BaseContext func(net.PacketConn) context.Context

	// PanicHandler optionally handles panic.
	// Recovered value, which always non-nil, and remote and local addresses
	// are provided. The sentinel error [ErrAbortHandler] is also given to
	// the PanicHandler. It bypasses default logging of stacktraces.
	PanicHandler func(recovered any, local, remote net.Addr)

	shutdown    atomic.Bool
	packetConns internal.CloserStore[*ocPacketConn]
	conns       internal.CloserStore[*ocConn]

	// serveNotify notifies the inner listener is working
	// and the Server.Serve is called.
	// serveNotify is used for testing.
	serveNotify chan struct{}
}

func (s *Server) ListenAndServe() error {
	if s.shutdown.Load() {
		return net.ErrClosed
	}
	pc, err := newPacketConn(s.Addr)
	if err != nil {
		return err
	}
	return s.Serve(pc)
}

// Serve accepts incoming packets on the PacketConn p.
// Serve creates a new service goroutine for each remote address.
// The received packets are passed to s.Handler through the connection.
//
// Serve always returns a non-nil error.
// After [Server.Shutdown] or [Server.Close], the returned error is [net.ErrClosed].
func (s *Server) Serve(p net.PacketConn) error {
	if s.shutdown.Load() {
		return net.ErrClosed
	}

	ctx := context.Background()
	if bc := s.BaseContext; bc != nil {
		ctx = bc(p)
	}

	ocp := &ocPacketConn{PacketConn: p}
	s.packetConns.Store(ocp)
	defer func() {
		_ = ocp.Close()
		s.packetConns.Delete(ocp) // Delete after close.
	}()

	if s.serveNotify != nil {
		s.serveNotify <- struct{}{}
	}

	var channels sync.Map
	wait := int64(1)
	buf := make([]byte, mtu)
	for {
		n, addr, err := ocp.ReadFrom(buf)
		if err != nil {
			if err == ErrSkipHandler {
				continue
			}
			if s.shutdown.Load() { // Error is caused by shutdown.
				return net.ErrClosed
			}
		}
		if n > 0 {
			c, isNew := getChannel(&channels, addr.String())
			if isNew {
				conn := &conn{pc: p, raddr: addr, packets: c, channels: &channels}
				go s.serve(ctx, p.LocalAddr(), addr, conn)
			}
			packet := make([]byte, n)
			copy(packet, buf[:n])
			select {
			case c <- packet:
			default:
				// Channel is full. Discard the packet.
			}
			wait = 1 // Reset
		}
		if err != nil {
			// Check if this is caused by timeout or not. [net.Error] implements the interface.
			if to, ok := err.(interface{ Timeout() bool }); ok && to.Timeout() {
				wait = min(wait*2, 1<<9) // Up to 512 msec.
				time.Sleep(time.Duration(wait) * time.Millisecond)
				continue
			}
			return err
		}
	}
}

func (s *Server) serve(ctx context.Context, laddr, raddr net.Addr, conn *conn) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		if ph := s.PanicHandler; ph != nil {
			ph(err, laddr, raddr)
			return
		}
		if err != ErrAbortHandler {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("znet/zudp: panic serving %v: %v\n%s", raddr, err, buf)
		}
	}()

	occ := &ocConn{Conn: conn}
	s.conns.Store(occ)
	defer func() {
		_ = occ.Close()
		s.conns.Delete(occ) // Delete after close.
	}()
	s.Handler.ServeUDP(ctx, occ)
}

// Close immediately closes all active [net.PacketConn] and any connections.
// For a graceful shutdown, use [Server.Shutdown].
//
// When Close is called, [Server.Serve] and [Server.ListenAndServe] immediately
// return [net.ErrClosed].
//
// Once Close has been called on a server, it may not be reused;
// future calls to methods such as [Server.Serve] will return [net.ErrClosed].
func (s *Server) Close() error {
	s.shutdown.Store(true)
	err1 := s.packetConns.CloseAll()
	err2 := s.conns.CloseAll()
	return cmp.Or(err1, err2)
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
// Shutdown works by first closing all open [net.PacketConn]s,
// and then waiting for connections to be closed and then shut down.
// If the provided context expires before the shutdown is complete, Shutdown returns
// the context's error, otherwise it returns all errors returned from closing the
// Server's underlying PacketConn(s). Non-nil errors occurred while shutting down
// are returned after joined with [errors.Join].
//
// When Shutdown is called, [Server.Serve] and [Server.ListenAndServe]
// immediately return [net.ErrClosed].
// Make sure the program doesn't exit and waits instead for Shutdown to return.
//
// Once Shutdown has been called on a server, it may not be reused;
// future calls to methods such as [Server.Serve] will return [net.ErrClosed].
func (s *Server) Shutdown(ctx context.Context) error {
	if s.shutdown.Swap(true) {
		return net.ErrClosed
	}
	err := s.packetConns.CloseAll()
	for s.conns.Length() > 0 {
		select {
		case <-time.After(100 * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}

func newPacketConn(addr string) (pc net.PacketConn, err error) {
	network, address := znet.ParseNetAddr(addr)
	switch network {
	case "":
		network = "udp" // Assume "udp".
		fallthrough
	case "udp", "udp4", "udp6":
		laddr, resolvErr := net.ResolveUDPAddr(network, address)
		if resolvErr != nil {
			return nil, resolvErr
		}
		return net.ListenUDP(network, laddr)
	case "unixgram":
		laddr := &net.UnixAddr{Name: address, Net: network}
		return net.ListenUnixgram(network, laddr)
	default:
		return net.ListenPacket("udp", addr) // Fallback. May be invalid addr.
	}
}

// ocPacketConn is once close PacketConn that wraps a [net.PacketConn],
// protecting it from multiple Close calls.
type ocPacketConn struct {
	net.PacketConn
	once     sync.Once
	closeErr error
}

func (oc *ocPacketConn) Close() error {
	oc.once.Do(func() {
		oc.closeErr = oc.PacketConn.Close()
		addr := oc.LocalAddr()
		if _, ok := addr.(*net.UnixAddr); !ok {
			return // Non unix socket.
		}
		address := addr.String()
		if len(address) > 0 && address[0] == '@' {
			return // Abstract socket.
		}
		_ = os.Remove(address) // Remove socket file.
	})
	return oc.closeErr
}

// ocConn is once close Conn that wraps a [Conn],
// protecting it from multiple Close calls.
type ocConn struct {
	Conn
	once     sync.Once
	closeErr error
}

func (oc *ocConn) Close() error {
	oc.once.Do(func() {
		oc.closeErr = oc.Conn.Close()
	})
	return oc.closeErr
}

// getChannel returns a bytes channel.
// If a channel found in the store with the key addr, getChannel returns it and false.
// If not found, getChannel creates a new byte channel and register it to the
// store with the key addr and returns the new channel and true.
func getChannel(store *sync.Map, addr string) (c chan []byte, isNew bool) {
	const size = 256 // 256*8 = 2048 B/addr.
	v, ok := store.Load(addr)
	if !ok {
		c := make(chan []byte, size)
		store.Store(addr, c)
		return c, true
	}
	return v.(chan []byte), false
}
