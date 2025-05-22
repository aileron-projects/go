package ztcp

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/aileron-projects/go/znet"
)

var (
	// ErrNoTarget indicates there is no proxy target.
	ErrNoTarget = errors.New("znet/ztcp: at least 1 targe is required")
)

// NewProxy returns a new instance of [Proxy].
// Returned proxy proxies requests to the targets.
// If no targets specified, NewProxy panics [ErrNoTarge].
// Targets must be a valid TCP or Unix address.
// Proxy target is selected by round-robin algorithm without any health checks.
// Some examples of valid target addresses are listed below.
//
// Examples:
//   - "localhost:8080"
//   - "127.0.0.1:8080"
//   - "[::1]:8080"
//   - "tcp://127.0.0.1:8080"
//   - "tcp4://127.0.0.1:8080"
//   - "tcp6://[::1]:8080"
//   - "unix:///var/run/example.sock"
//   - "unix://@example"
//   - "unixpacket:///var/run/example.sock"
//   - "unixpacket://@example"
func NewProxy(targets ...string) *Proxy {
	if len(targets) == 0 {
		panic(ErrNoTarget)
	}
	return &Proxy{
		Dial: (&roundRobinDialer{addrs: targets, index: -1}).dial,
	}
}

// Proxy is a TCP proxy.
type Proxy struct {
	// Dial returns a new upstream connection
	// for the downstream connection dc.
	// Dial must not be nil, otherwise ServeTCP will panic.
	// If the returned err is nil, uc must not be nil and
	// if the returned err is non-nil, uc must be nil.
	Dial func(ctx context.Context, dc net.Conn) (uc net.Conn, err error)
	// ErrorHandler optionally handles non-nil error.
	// Provided err is always non-nil.
	// Downstream connection and upstream connection are
	// provided as dc and uc each.
	// uc can be nil when Dial returned an error.
	ErrorHandler func(dc, uc net.Conn, err error)
}

func (p *Proxy) handleError(dc, uc net.Conn, err error) {
	if err == nil {
		return
	}
	if eh := p.ErrorHandler; eh != nil {
		eh(dc, uc, err)
	}
}

func (p *Proxy) ServeTCP(ctx context.Context, conn net.Conn) {
	upConn, err := p.Dial(ctx, conn)
	if err != nil {
		p.handleError(conn, upConn, err)
		return
	}
	defer upConn.Close() // Ensure close upstream connection.

	errChan := make(chan error)
	go copyBuf(conn, upConn, errChan) // downstream --> proxy --> upstream
	go copyBuf(upConn, conn, errChan) // downstream <-- proxy <-- upstream

	if err := <-errChan; err != nil {
		p.handleError(conn, upConn, err)
		return
	}
	if err := <-errChan; err != nil {
		p.handleError(conn, upConn, err)
		return
	}
}

// pool is the buffer pool.
var pool = sync.Pool{
	New: func() any {
		buf := make([]byte, 1<<14) // 16kiB
		return &buf
	},
}

func copyBuf(dst io.Writer, src io.Reader, errChan chan<- error) {
	buf := *pool.Get().(*[]byte)
	defer pool.Put(&buf)
	_, err := io.CopyBuffer(dst, src, buf)
	errChan <- err
}

// roundRobinDialer dials to the address
// in addrs with round-robin algorithm.
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

func (d *roundRobinDialer) dial(_ context.Context, _ net.Conn) (net.Conn, error) {
	addr := d.next()
	network, address := znet.ParseNetAddr(addr)
	switch network {
	case "tcp", "tcp4", "tcp6":
		raddr, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return net.DialTCP(network, nil, raddr)
	case "unix", "unixpacket":
		raddr := &net.UnixAddr{Name: address, Net: network}
		return net.DialUnix(network, nil, raddr)
	default:
		return net.Dial("tcp", addr) // Fallback.
	}
}
