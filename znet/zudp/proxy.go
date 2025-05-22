package zudp

import (
	"cmp"
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aileron-projects/go/znet"
)

var (
	// ErrNoTarget indicates there is no proxy target.
	ErrNoTarget = errors.New("znet/zudp: at least 1 targe is required")
)

// NewProxy returns a new instance of [Proxy].
// Returned proxy proxies request to the targets.
// If no targets specified, NewProxy panics [ErrNoTarge].
// Targets must be a valid UDP or Unix address.
// Proxy target is selected by round-robin algorithm without any health checks.
// Some examples of valid target addresses are listed below.
//
// Examples:
//   - "localhost:8080"
//   - "127.0.0.1:8080"
//   - "[::1]:8080"
//   - "udp://127.0.0.1:8080"
//   - "udp4://127.0.0.1:8080"
//   - "udp6://[::1]:8080"
//   - "unix:///var/run/example.sock"
//   - "unix://@example"
//   - "unixgram:///var/run/example.sock"
//   - "unixgram://@example"
func NewProxy(targets ...string) *Proxy {
	if len(targets) == 0 {
		panic(ErrNoTarget)
	}
	return &Proxy{
		Dial: (&roundRobinDialer{addrs: targets, index: -1}).dial,
	}
}

// Proxy is the UDP proxy.
type Proxy struct {
	// Dial returns a new upstream connection dc for the
	// new downstream connection dc.
	// Upstream connection dc must not be modified in Dial.
	// Dial must not be nil, otherwise ServeUDP will panic.
	// If the returned err is nil, uc must not be nil and
	// if the returned err is non-nil, uc must be nil.
	Dial func(ctx context.Context, dc Conn) (uc net.Conn, err error)
	// ErrorHandler optionally handles non-nil error.
	// The first argument err is always non-nil.
	// Downstream connection and upstream connection are
	// passed as dc and uc each.
	ErrorHandler func(dc Conn, uc net.Conn, err error)
	// IdleTimeout sets the idle timeout of the connection.
	// Duration after last read or write of upstream connection
	// and downstream connection is used for checking the timeout.
	// If zero or negative, default is 10 seconds.
	IdleTimeout time.Duration
}

func (p *Proxy) handleError(dc Conn, uc net.Conn, err error) {
	if err == nil {
		return
	}
	if eh := p.ErrorHandler; eh != nil {
		eh(dc, uc, err)
	}
}

func (p *Proxy) ServeUDP(ctx context.Context, conn Conn) {
	timeout := p.IdleTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	upConn, err := p.Dial(ctx, conn)
	if err != nil {
		p.handleError(conn, upConn, err)
		return
	}
	defer upConn.Close() // Ensure close upstream connection.

	var lastActive atomic.Pointer[time.Time]
	now := time.Now()
	lastActive.Store(&now)

	errChan := make(chan error)
	go copyBuf(conn, upConn, errChan, &lastActive) // downstream --> proxy --> upstream
	go copyBuf(upConn, conn, errChan, &lastActive) // downstream <-- proxy <-- upstream

	for {
		select {
		case <-time.After(timeout):
			if time.Since(*lastActive.Load()) > timeout {
				return
			}
		case err = <-errChan:
			p.handleError(conn, upConn, err)
			return
		}
	}
}

func copyBuf(dst io.Writer, src io.Reader, errChan chan<- error, active *atomic.Pointer[time.Time]) {
	buf := make([]byte, mtu)
	for {
		nr, err := src.Read(buf)
		if nr > 0 {
			now := time.Now()
			active.Store(&now)
			nw, ew := dst.Write(buf[:nr])
			if ew != nil || nr != nw {
				errChan <- cmp.Or(ew, io.ErrShortWrite)
				return
			}
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			errChan <- err
			return
		}
	}
}

// roundRobinDialer dials to the address in addrs
// with round-robin algorithm.
// It does not do any health checks.
type roundRobinDialer struct {
	mu    sync.Mutex
	index int
	addrs []string
}

func (d *roundRobinDialer) next() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.index++
	if d.index >= len(d.addrs) {
		d.index = 0
	}
	return d.addrs[d.index]
}

func (d *roundRobinDialer) dial(_ context.Context, _ Conn) (net.Conn, error) {
	addr := d.next()
	network, address := znet.ParseNetAddr(addr)
	switch network {
	case "udp", "udp4", "udp6":
		raddr, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return net.DialUDP(network, nil, raddr)
	case "unix", "unixgram":
		network = "unixgram"
		raddr := &net.UnixAddr{Name: address, Net: network}
		return net.DialUnix(network, nil, raddr)
	default:
		return net.Dial("udp", addr) // Fallback.
	}
}
