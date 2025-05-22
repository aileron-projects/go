package zudp

import (
	"context"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aileron-projects/go/znet/internal"
	"github.com/aileron-projects/go/ztesting"
)

func TestServer_ListenAndServe(t *testing.T) {
	t.Parallel()
	t.Run("already shutdown", func(t *testing.T) {
		s := &Server{}
		s.Shutdown(context.Background())
		err := s.ListenAndServe()
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
	t.Run("create listener error", func(t *testing.T) {
		s := &Server{Addr: "udp4://1234567890"}
		err := s.ListenAndServe()
		_, ok := err.(*net.AddrError)
		ztesting.AssertEqual(t, "addr error should be returned", true, ok)
	})
	t.Run("listen success", func(t *testing.T) {
		served := make(chan struct{})
		s := &Server{
			Addr:        "udp://:0",
			Handler:     HandlerFunc(func(ctx context.Context, conn Conn) {}),
			serveNotify: served,
		}
		shutdwon := make(chan error)
		go func() {
			<-served
			shutdwon <- s.Shutdown(context.Background())
		}()
		err := s.ListenAndServe()
		ztesting.AssertEqualErr(t, "serve error not match", net.ErrClosed, err)
		err = <-shutdwon
		ztesting.AssertEqualErr(t, "shutdown error not match", nil, err)
	})
}

type timeoutError bool

func (e timeoutError) Error() string {
	return "timeout"
}

func (e timeoutError) Timeout() bool {
	return bool(e)
}

type testPacketConn struct {
	net.PacketConn
	// Return values
	raddr    net.Addr
	closeErr error
	readErr  error
	// Recorded values
	closed int
	accept int
}

func (c *testPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	c.accept++
	if c.closed > 0 {
		return 0, nil, net.ErrClosed
	}
	if c.accept > 1 {
		time.Sleep(10 * time.Millisecond)
	}
	return copy(p, []byte("test")), c.raddr, c.readErr
}

func (c *testPacketConn) Close() error {
	c.closed++
	return c.closeErr
}

func TestServer_Serve(t *testing.T) {
	t.Parallel()
	dpc, _ := net.ListenPacket("udp", ":0")
	dpc.Close()
	t.Run("already shutdown", func(t *testing.T) {
		s := &Server{}
		s.Shutdown(context.Background())
		err := s.Serve(nil)
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
	t.Run("create listener error", func(t *testing.T) {
		s := &Server{Addr: "udp4://1234567890"}
		err := s.ListenAndServe()
		_, ok := err.(*net.AddrError)
		ztesting.AssertEqual(t, "addr error should be returned", true, ok)
	})
	t.Run("serve success", func(t *testing.T) {
		baseCtx := context.Background()
		ln := &testPacketConn{PacketConn: dpc, raddr: &net.UDPAddr{IP: net.ParseIP("127.0.0.1")}}
		invoked := make(chan struct{})
		s := &Server{
			BaseContext: func(_ net.PacketConn) context.Context { return baseCtx },
			Handler: HandlerFunc(func(ctx context.Context, conn Conn) {
				ztesting.AssertEqual(t, "context not match", baseCtx, ctx)
				invoked <- struct{}{}
			}),
		}
		go func() {
			<-invoked
			s.Close()
		}()
		err := s.Serve(ln)
		ztesting.AssertEqual(t, "error not match", net.ErrClosed, err)
	})
	t.Run("skip serving", func(t *testing.T) {
		pc := &testPacketConn{PacketConn: dpc, raddr: &net.UDPAddr{}, readErr: ErrSkipHandler}
		count := 0
		s := &Server{
			Handler: HandlerFunc(func(_ context.Context, _ Conn) { count++ }),
		}
		go func() {
			for pc.accept <= 2 {
				time.Sleep(10 * time.Millisecond)
			}
			ztesting.AssertEqual(t, "handler should not be called", 0, count)
			s.Close()
		}()
		err := s.Serve(pc)
		ztesting.AssertEqual(t, "error not match", net.ErrClosed, err)
	})
	t.Run("timeout error", func(t *testing.T) {
		pc := &testPacketConn{PacketConn: dpc, raddr: &net.UDPAddr{}, readErr: timeoutError(true)}
		count := 0
		s := &Server{
			Handler: HandlerFunc(func(_ context.Context, _ Conn) { count++ }),
		}
		go func() {
			for pc.accept <= 2 {
				time.Sleep(10 * time.Millisecond)
			}
			ztesting.AssertEqual(t, "read data should be proceeded", true, count > 0)
			s.Close()
		}()
		err := s.Serve(pc)
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
	t.Run("non-timeout error", func(t *testing.T) {
		pc := &testPacketConn{PacketConn: dpc, raddr: &net.UDPAddr{}, readErr: timeoutError(false)}
		count := 0
		s := &Server{
			Handler: HandlerFunc(func(_ context.Context, _ Conn) { count++ }),
		}
		go func() {
			for pc.accept <= 2 {
				time.Sleep(10 * time.Millisecond)
			}
			ztesting.AssertEqual(t, "handler should not be called", 0, count)
			s.Close()
		}()
		err := s.Serve(pc)
		ztesting.AssertEqualErr(t, "error not match", timeoutError(false), err)
	})
	t.Run("panic error", func(t *testing.T) {
		pc := &testPacketConn{PacketConn: dpc, raddr: &net.UDPAddr{}}
		s := &Server{
			Handler: HandlerFunc(func(_ context.Context, _ Conn) {
				panic(net.ErrWriteToConnected) // Panic dummy error.
			}),
		}
		go func() {
			for pc.accept <= 2 {
				time.Sleep(10 * time.Millisecond)
			}
			s.Close()
		}()
		err := s.Serve(pc)
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
		ztesting.AssertEqual(t, "packet conn not closed", 1, pc.closed)
	})
	t.Run("panic with handler", func(t *testing.T) {
		pc := &testPacketConn{PacketConn: dpc, raddr: &net.UDPAddr{}}
		panicked := make(chan error)
		s := &Server{
			Handler: HandlerFunc(func(_ context.Context, _ Conn) {
				panic(net.ErrWriteToConnected) // Panic dummy error.
			}),
			PanicHandler: func(recovered any, remote, local net.Addr) {
				panicked <- recovered.(error)
			},
		}
		go func() {
			err := <-panicked
			ztesting.AssertEqualErr(t, "error not match", net.ErrWriteToConnected, err)
			s.Close()
		}()
		err := s.Serve(pc)
		ztesting.AssertEqualErr(t, "serve error not match", net.ErrClosed, err)
		ztesting.AssertEqual(t, "packet conn not closed", 1, pc.closed)
	})
}

func TestServer_Close(t *testing.T) {
	t.Parallel()
	pc, _ := net.ListenPacket("udp", ":0")
	pc.Close()
	s := &Server{
		Handler:     HandlerFunc(func(_ context.Context, _ Conn) { time.Sleep(time.Second) }),
		packetConns: internal.UniqueStore[*ocPacketConn]{},
		conns:       internal.UniqueStore[*ocConn]{},
	}
	s.packetConns.Set(&ocPacketConn{PacketConn: pc, store: &s.packetConns})
	s.conns.Set(&ocConn{Conn: &net.TCPConn{}, store: &s.conns})

	ztesting.AssertEqual(t, "packetConns length not match", 1, s.packetConns.Length())
	ztesting.AssertEqual(t, "conns length not match", 1, s.conns.Length())
	s.Close()
	ztesting.AssertEqual(t, "packetConns length not match", 0, s.packetConns.Length())
	ztesting.AssertEqual(t, "conns length not match", 0, s.conns.Length())
}

func TestServer_Shutdown(t *testing.T) {
	t.Parallel()
	t.Run("already shutdown", func(t *testing.T) {
		s := &Server{}
		s.Shutdown(context.Background())
		err := s.Shutdown(context.Background())
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
	t.Run("shutdown success", func(t *testing.T) {
		pc, _ := net.ListenPacket("udp", ":0")
		invoked := make(chan struct{})
		s := &Server{
			Handler: HandlerFunc(func(ctx context.Context, conn Conn) {
				invoked <- struct{}{}
			}),
		}
		shutdown := make(chan error)
		go func() {
			conn, _ := net.Dial("udp", pc.LocalAddr().String())
			conn.Write([]byte("test"))
			defer conn.Close()
			<-invoked
			shutdown <- s.Shutdown(context.Background())
		}()
		err := s.Serve(pc)
		ztesting.AssertEqual(t, "serve error not match", net.ErrClosed, err)
		err = <-shutdown
		ztesting.AssertEqual(t, "shutdown error not match", nil, err)
		ztesting.AssertEqual(t, "listeners length not match", 0, s.packetConns.Length())
		ztesting.AssertEqual(t, "conns length not match", 0, s.conns.Length())
	})
	t.Run("shutdown context done", func(t *testing.T) {
		handlerInvoked := make(chan struct{})
		pc, _ := net.ListenPacket("udp", ":0")
		s := &Server{
			Handler: HandlerFunc(func(ctx context.Context, conn Conn) {
				handlerInvoked <- struct{}{}
				<-t.Context().Done()
			}),
		}
		shutdown := make(chan struct{})
		go func() {
			conn, err := net.DialUDP("udp", nil, pc.LocalAddr().(*net.UDPAddr))
			ztesting.AssertEqual(t, "error should be nil", nil, err)
			conn.Write([]byte("test"))
			defer conn.Close()
			<-handlerInvoked
			ctx, cancel := context.WithTimeout(context.Background(), 0)
			defer cancel()
			err = s.Shutdown(ctx)
			ztesting.AssertEqual(t, "error not match", context.DeadlineExceeded, err)
			shutdown <- struct{}{}
		}()
		err := s.Serve(pc)
		<-shutdown
		ztesting.AssertEqual(t, "listeners length not match", 0, s.packetConns.Length())
		ztesting.AssertEqual(t, "conns length not match", 1, s.conns.Length()) // Conn is yet alive.
		ztesting.AssertEqual(t, "error not match", net.ErrClosed, err)
	})
}

func TestNewPacketConn(t *testing.T) {
	t.Parallel()
	t.Run("udp without prefix", func(t *testing.T) {
		ln, err := newPacketConn(":0")
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		defer ln.Close()
		cn, err := net.Dial("udp", ln.LocalAddr().String())
		ztesting.AssertEqual(t, "dial failed", nil, err)
		cn.Close()
	})
	t.Run("listen udp4 success", func(t *testing.T) {
		ln, err := newPacketConn("udp4://:0")
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		defer ln.Close()
		cn, err := net.Dial("udp4", ln.LocalAddr().String())
		ztesting.AssertEqual(t, "dial failed", nil, err)
		cn.Close()
	})
	t.Run("listen udp4 failed", func(t *testing.T) {
		_, err := newPacketConn("udp4://1234567890")
		_, ok := err.(*net.AddrError)
		ztesting.AssertEqual(t, "addr error should be returned", true, ok)
	})
	t.Run("listen unixgram", func(t *testing.T) {
		s := t.TempDir() + "/not-exist/test.sock"
		_, err := newPacketConn("unixgram://" + s) // Make error because windows not support it.
		_, ok := err.(*net.OpError)
		t.Logf("%#v\n", err)
		ztesting.AssertEqual(t, "net op error should be returned", true, ok)
	})
	t.Run("fallback to udp", func(t *testing.T) {
		_, err := newPacketConn("tcp://1234567890")
		_, ok := err.(*net.OpError)
		t.Logf("%#v\n", err)
		ztesting.AssertEqual(t, "net op error should be returned", true, ok)
	})
}

type nopClosePacketConn struct {
	net.PacketConn
	addr  net.Addr // LocalAddr
	count int
}

func (l *nopClosePacketConn) LocalAddr() net.Addr {
	return l.addr
}

func (l *nopClosePacketConn) Close() error {
	l.count++
	return nil
}

type nopCloseConn struct {
	net.Conn
	count int
}

func (c *nopCloseConn) Close() error {
	c.count++
	return nil
}

func TestOCPacketConn(t *testing.T) {
	t.Parallel()
	t.Run("close once", func(t *testing.T) {
		store := internal.UniqueStore[*ocPacketConn]{}
		pc := &nopClosePacketConn{addr: &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}}
		ln := &ocPacketConn{PacketConn: pc, store: &store}
		store.Set(ln)
		ztesting.AssertEqual(t, "listener has not been stored", 1, store.Length())
		err := ln.Close()
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "length not match", 0, store.Length())
	})
	t.Run("close multiple", func(t *testing.T) {
		store := internal.UniqueStore[*ocPacketConn]{}
		pc := &nopClosePacketConn{addr: &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}}
		ln := &ocPacketConn{PacketConn: pc, store: &store}
		ln.Close()
		ln.Close()
		ztesting.AssertEqual(t, "close called more than once", 1, pc.count)
	})
	t.Run("close abstract socket", func(t *testing.T) {
		store := internal.UniqueStore[*ocPacketConn]{}
		pc := &nopClosePacketConn{addr: &net.UnixAddr{Net: "unixgram", Name: "@test"}}
		ln := &ocPacketConn{PacketConn: pc, store: &store}
		store.Set(ln)
		ztesting.AssertEqual(t, "listener has not been stored", 1, store.Length())
		err := ln.Close()
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "length not match", 0, store.Length())
	})
	t.Run("close path name socket", func(t *testing.T) {
		sock := t.TempDir() + "/test.sock"
		f, _ := os.Create(sock)
		f.Close()
		store := internal.UniqueStore[*ocPacketConn]{}
		pc := &nopClosePacketConn{addr: &net.UnixAddr{Net: "unixgram", Name: sock}}
		ln := &ocPacketConn{PacketConn: pc, store: &store}
		store.Set(ln)
		ztesting.AssertEqual(t, "listener has not been stored", 1, store.Length())
		err := ln.Close()
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "length not match", 0, store.Length())
		_, err = os.Stat(sock) // Socket must be removed.
		ztesting.AssertEqual(t, "socket not removed", true, os.IsNotExist(err))
	})
}

func TestOCConn(t *testing.T) {
	t.Parallel()
	t.Run("close once", func(t *testing.T) {
		store := internal.UniqueStore[*ocConn]{}
		conn := &ocConn{Conn: &nopCloseConn{}, store: &store}
		store.Set(conn)
		ztesting.AssertEqual(t, "conn has not been stored", 1, store.Length())
		err := conn.Close()
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "length not match", 0, store.Length())
	})
	t.Run("close multiple", func(t *testing.T) {
		nc := &nopCloseConn{}
		conn := &ocConn{Conn: nc, store: &internal.UniqueStore[*ocConn]{}}
		conn.Close()
		conn.Close()
		ztesting.AssertEqual(t, "close called more than once", 1, nc.count)
	})
}

func TestGetChannel(t *testing.T) {
	t.Parallel()
	t.Run("new channel", func(t *testing.T) {
		m := &sync.Map{}
		c, isNew := getChannel(m, "addr")
		ztesting.AssertEqual(t, "channel is not new", true, isNew)
		ztesting.AssertEqual(t, "channel length not match", 0, len(c))
		ztesting.AssertEqual(t, "channel capacity not match", 256, cap(c))
	})
	t.Run("existing channel", func(t *testing.T) {
		m := &sync.Map{}
		c, isNew := getChannel(m, "addr")
		ztesting.AssertEqual(t, "channel is not new", true, isNew)
		c <- []byte("1")
		c <- []byte("2")
		c <- []byte("3")
		c, isNew = getChannel(m, "addr")
		ztesting.AssertEqual(t, "channel is new", false, isNew)
		ztesting.AssertEqual(t, "channel length not match", 3, len(c))
		ztesting.AssertEqual(t, "channel capacity not match", 256, cap(c))
	})
}
