package ztcp

import (
	"context"
	"crypto/tls"
	"io/fs"
	"net"
	"os"
	"path/filepath"
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
		s := &Server{Addr: "tcp4://1234567890"}
		err := s.ListenAndServe()
		_, ok := err.(*net.AddrError)
		ztesting.AssertEqual(t, "addr error should be returned", true, ok)
	})
	t.Run("listen success", func(t *testing.T) {
		served := make(chan struct{})
		s := &Server{
			Addr:        "tcp://:0",
			Handler:     HandlerFunc(func(ctx context.Context, conn net.Conn) {}),
			serveNotify: served,
		}
		shutdown := make(chan error)
		go func() {
			<-served
			shutdown <- s.Shutdown(context.Background())
		}()
		err := s.ListenAndServe()
		ztesting.AssertEqualErr(t, "serve error not match", net.ErrClosed, err)
		err = <-shutdown
		ztesting.AssertEqualErr(t, "shutdown error not match", nil, err)
	})
}

func TestServer_ListenAndServeTLS(t *testing.T) {
	t.Parallel()
	// Obtain available address.
	ln, err := net.Listen("tcp4", ":0")
	if err != nil {
		panic(err)
	}
	ln.Close()
	addr := ln.Addr().String()

	t.Run("already shutdown", func(t *testing.T) {
		s := &Server{}
		s.Shutdown(context.Background())
		err := s.ListenAndServeTLS("", "")
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
	t.Run("create listener error", func(t *testing.T) {
		s := &Server{Addr: "tcp4://1234567890"}
		err := s.ListenAndServeTLS("", "")
		_, ok := err.(*net.AddrError)
		ztesting.AssertEqual(t, "addr error should be returned", true, ok)
	})
	t.Run("listen success", func(t *testing.T) {
		served := make(chan struct{})
		s := &Server{
			Addr:        addr,
			Handler:     HandlerFunc(func(ctx context.Context, conn net.Conn) {}),
			serveNotify: served,
		}
		go func() {
			<-served
			cn, err := net.Dial("tcp4", addr)
			ztesting.AssertEqual(t, "dial failed", nil, err)
			cn.Close()
			s.Close()
		}()
		err := s.ListenAndServeTLS("./testdata/cert.pem", "./testdata/key.pem")
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
}

func TestServer_ServeTLS(t *testing.T) {
	t.Parallel()
	// Obtain available address.
	ln, err := net.Listen("tcp4", ":0")
	if err != nil {
		panic(err)
	}
	ln.Close()
	addr := ln.Addr().String()

	t.Run("already shutdown", func(t *testing.T) {
		s := &Server{}
		s.Shutdown(context.Background())
		err := s.ServeTLS(nil, "", "")
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
	t.Run("read cert error", func(t *testing.T) {
		s := &Server{}
		err := s.ServeTLS(nil, "./testdata/not-found.pem", "./testdata/key.pem")
		_, ok := err.(*fs.PathError)
		ztesting.AssertEqual(t, "path error should be returned", true, ok)
	})
	t.Run("non-nil config", func(t *testing.T) {
		served := make(chan struct{})
		s := &Server{
			Addr:        addr,
			TLSConfig:   &tls.Config{},
			Handler:     HandlerFunc(func(ctx context.Context, conn net.Conn) {}),
			serveNotify: served,
		}
		go func() {
			<-served
			conn, err := net.Dial("tcp4", addr)
			ztesting.AssertEqual(t, "dial failed", nil, err)
			conn.Close()
			s.Close()
		}()
		ln, _ := net.Listen("tcp4", addr)
		err := s.ServeTLS(ln, "./testdata/cert.pem", "./testdata/key.pem")
		ztesting.AssertEqual(t, "error not match", net.ErrClosed, err)
	})
}

type timeoutError bool

func (e timeoutError) Error() string {
	return "timeout"
}

func (e timeoutError) Timeout() bool {
	return bool(e)
}

type testConn struct {
	net.Conn
	// Recorded values
	closed int
}

func (c *testConn) Close() error {
	c.closed++
	return nil
}

type testListener struct {
	net.Listener
	// Return values
	conn      net.Conn
	closeErr  error
	acceptErr error
	// Recorded values
	closed int
	accept int
}

func (l *testListener) Accept() (net.Conn, error) {
	l.accept++
	if l.closed > 0 {
		return nil, net.ErrClosed
	}
	if l.accept > 1 {
		time.Sleep(10 * time.Millisecond)
	}
	return l.conn, l.acceptErr
}

func (l *testListener) Close() error {
	l.closed++
	return l.closeErr
}

func TestServer_Serve(t *testing.T) {
	t.Parallel()
	dln, _ := net.Listen("tcp", ":0")
	dln.Close()
	t.Run("already shutdown", func(t *testing.T) {
		s := &Server{}
		s.Shutdown(context.Background())
		err := s.Serve(nil)
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
	t.Run("create listener error", func(t *testing.T) {
		s := &Server{Addr: "tcp4://1234567890"}
		err := s.ListenAndServe()
		_, ok := err.(*net.AddrError)
		ztesting.AssertEqual(t, "addr error should be returned", true, ok)
	})
	t.Run("serve success", func(t *testing.T) {
		baseCtx := context.Background()
		cn := &testConn{Conn: &net.TCPConn{}}
		ln := &testListener{Listener: dln, conn: cn}
		served := make(chan struct{})
		checked := make(chan struct{})
		s := &Server{
			BaseContext: func(l net.Listener) context.Context { return baseCtx },
			Handler: HandlerFunc(func(ctx context.Context, conn net.Conn) {
				ztesting.AssertEqual(t, "context not match", baseCtx, ctx)
				ztesting.AssertEqual(t, "connection not match", net.Conn(cn), conn)
				checked <- struct{}{}
			}),
			serveNotify: served,
		}
		go func() {
			<-served
			<-checked
			s.Close()
		}()
		err := s.Serve(ln)
		ztesting.AssertEqual(t, "error not match", net.ErrClosed, err)
	})
	t.Run("skip serving", func(t *testing.T) {
		cn := &testConn{Conn: &net.TCPConn{}}
		ln := &testListener{Listener: dln, conn: cn, acceptErr: ErrSkipHandler}
		served := make(chan struct{})
		count := 0
		s := &Server{
			Handler:     HandlerFunc(func(_ context.Context, _ net.Conn) { count++ }),
			serveNotify: served,
		}
		go func() {
			<-served
			for cn.closed == 0 {
				time.Sleep(10 * time.Millisecond)
			}
			ztesting.AssertEqual(t, "handler should not be called", 0, count)
			s.Close()
		}()
		err := s.Serve(ln)
		ztesting.AssertEqual(t, "error not match", net.ErrClosed, err)
	})
	t.Run("timeout error", func(t *testing.T) {
		cn := &testConn{Conn: &net.TCPConn{}}
		ln := &testListener{Listener: dln, conn: cn, acceptErr: timeoutError(true)}
		served := make(chan struct{})
		count := 0
		s := &Server{
			Handler:     HandlerFunc(func(_ context.Context, _ net.Conn) { count++ }),
			serveNotify: served,
		}
		go func() {
			<-served
			for ln.accept <= 2 {
				time.Sleep(10 * time.Millisecond)
			}
			ztesting.AssertEqual(t, "handler should not be called", 0, count)
			s.Close()
		}()
		err := s.Serve(ln)
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
	t.Run("non-timeout error", func(t *testing.T) {
		cn := &testConn{Conn: &net.TCPConn{}}
		ln := &testListener{Listener: dln, conn: cn, acceptErr: timeoutError(false)}
		served := make(chan struct{})
		count := 0
		s := &Server{
			Handler:     HandlerFunc(func(_ context.Context, _ net.Conn) { count++ }),
			serveNotify: served,
		}
		go func() {
			<-served
			for ln.accept <= 2 {
				time.Sleep(10 * time.Millisecond)
			}
			ztesting.AssertEqual(t, "handler should not be called", 0, count)
			s.Close()
		}()
		err := s.Serve(ln)
		ztesting.AssertEqualErr(t, "error not match", timeoutError(false), err)
	})
	t.Run("panic error", func(t *testing.T) {
		cn := &testConn{Conn: &net.TCPConn{}}
		ln := &testListener{Listener: dln, conn: cn}
		served := make(chan struct{})
		s := &Server{
			Handler: HandlerFunc(func(_ context.Context, _ net.Conn) {
				panic(net.ErrWriteToConnected) // Panic dummy error.
			}),
			serveNotify: served,
		}
		go func() {
			<-served
			for cn.closed == 0 {
				time.Sleep(10 * time.Millisecond)
			}
			s.Close()
		}()
		err := s.Serve(ln)
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
		ztesting.AssertEqual(t, "connection not closed", true, cn.closed > 0)
	})
	t.Run("panic with handler", func(t *testing.T) {
		cn := &testConn{Conn: &net.TCPConn{}}
		ln := &testListener{Listener: dln, conn: cn}
		panicked := make(chan struct{})
		s := &Server{
			PanicHandler: func(recovered any, remote, local net.Addr) {
				ztesting.AssertEqualErr(t, "error not match", net.ErrWriteToConnected, recovered.(error))
				panicked <- struct{}{}
			},
		}
		s.Handler = HandlerFunc(func(_ context.Context, _ net.Conn) {
			defer s.Close()
			panic(net.ErrWriteToConnected) // Panic dummy error.
		})
		err := s.Serve(ln)
		<-panicked
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
		ztesting.AssertEqual(t, "connection not closed", 1, cn.closed)
	})
}

func TestServer_Close(t *testing.T) {
	t.Parallel()
	ln, _ := net.Listen("tcp", ":0")
	ln.Close()
	s := &Server{
		Handler:   HandlerFunc(func(_ context.Context, _ net.Conn) { time.Sleep(time.Second) }),
		listeners: internal.UniqueStore[*ocListener]{},
		conns:     internal.UniqueStore[*ocConn]{},
	}
	s.listeners.Set(&ocListener{Listener: ln, store: &s.listeners})
	s.conns.Set(&ocConn{Conn: &net.TCPConn{}, store: &s.conns})

	ztesting.AssertEqual(t, "listeners length not match", 1, s.listeners.Length())
	ztesting.AssertEqual(t, "conns length not match", 1, s.conns.Length())
	s.Close()
	ztesting.AssertEqual(t, "listeners length not match", 0, s.listeners.Length())
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
		served := make(chan struct{})
		s := &Server{
			Addr:        "tcp://:0",
			serveNotify: served,
		}
		go func() {
			<-served
			err := s.Shutdown(context.Background())
			ztesting.AssertEqual(t, "error not match", nil, err)
		}()
		err := s.ListenAndServe()
		ztesting.AssertEqual(t, "error not match", net.ErrClosed, err)
		ztesting.AssertEqual(t, "listeners length not match", 0, s.listeners.Length())
		ztesting.AssertEqual(t, "conns length not match", 0, s.conns.Length())
		ztesting.AssertEqual(t, "error not match", net.ErrClosed, err)
	})
	t.Run("shutdown context done", func(t *testing.T) {
		served := make(chan struct{})
		handlerInvoked := make(chan struct{})
		ln, _ := net.Listen("tcp", ":0")
		s := &Server{
			Handler: HandlerFunc(func(ctx context.Context, conn net.Conn) {
				handlerInvoked <- struct{}{}
				<-ctx.Done()
			}),
			serveNotify: served,
		}
		shutdown := make(chan struct{})
		go func() {
			<-served
			conn, _ := net.DialTCP("tcp", nil, ln.Addr().(*net.TCPAddr))
			defer conn.Close()
			<-handlerInvoked
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err := s.Shutdown(ctx)
			ztesting.AssertEqual(t, "error not match", context.DeadlineExceeded, err)
			shutdown <- struct{}{}
		}()
		err := s.Serve(ln)
		<-shutdown
		ztesting.AssertEqual(t, "listeners length not match", 0, s.listeners.Length())
		ztesting.AssertEqual(t, "conns length not match", 1, s.conns.Length()) // Conn is yet alive.
		ztesting.AssertEqual(t, "error not match", net.ErrClosed, err)
	})
}

func TestNewListener(t *testing.T) {
	t.Parallel()
	// Obtain available address.
	ln, err := net.Listen("tcp4", ":0")
	if err != nil {
		panic(err)
	}
	ln.Close()
	addr := ln.Addr().String()

	t.Run("listen tcp without prefix", func(t *testing.T) {
		ln, err := newListener("" + addr)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		defer ln.Close()
		cn, err := net.Dial("tcp", addr)
		ztesting.AssertEqual(t, "dial failed", nil, err)
		cn.Close()
	})
	t.Run("listen tcp4 success", func(t *testing.T) {
		ln, err := newListener("tcp4://" + addr)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		defer ln.Close()
		cn, err := net.Dial("tcp4", addr)
		ztesting.AssertEqual(t, "dial failed", nil, err)
		cn.Close()
	})
	t.Run("listen tcp4 failed", func(t *testing.T) {
		_, err := newListener("tcp4://1234567890")
		_, ok := err.(*net.AddrError)
		ztesting.AssertEqual(t, "addr error should be returned", true, ok)
	})
	t.Run("listen unix success", func(t *testing.T) {
		s := filepath.Join(os.TempDir(), "TestNewListener_test.sock") // Socket path must not be too long.
		ln, err := newListener("unix://" + s)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		cn, err := net.Dial("unix", s)
		ztesting.AssertEqual(t, "dial failed", nil, err)
		cn.Close()
		ln.Close() // Socket file should be removed.
		_, err = os.Stat(s)
		ztesting.AssertEqual(t, "socket not removed", true, os.IsNotExist(err))
	})
	t.Run("fallback to tcp", func(t *testing.T) {
		_, err := newListener("udp://1234567890")
		_, ok := err.(*net.OpError)
		t.Logf("%#v\n", err)
		ztesting.AssertEqual(t, "net op error should be returned", true, ok)
	})
}

type nopCloseListener struct {
	net.Listener
	addr  net.Addr // LocalAddr
	count int
}

func (l *nopCloseListener) Addr() net.Addr {
	return l.addr
}

func (l *nopCloseListener) Close() error {
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

func TestOCListener(t *testing.T) {
	t.Parallel()
	t.Run("close once", func(t *testing.T) {
		store := internal.UniqueStore[*ocListener]{}
		l := &nopCloseListener{addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}}
		ln := &ocListener{Listener: l, store: &store}
		store.Set(ln)
		ztesting.AssertEqual(t, "listener has not been stored", 1, store.Length())
		err := ln.Close()
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "length not match", 0, store.Length())
	})
	t.Run("close multiple", func(t *testing.T) {
		store := internal.UniqueStore[*ocListener]{}
		l := &nopCloseListener{addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}}
		ln := &ocListener{Listener: l, store: &store}
		ln.Close()
		ln.Close()
		ztesting.AssertEqual(t, "close called more than once", 1, l.count)
	})
	t.Run("close abstract socket", func(t *testing.T) {
		store := internal.UniqueStore[*ocListener]{}
		l := &nopCloseListener{addr: &net.UnixAddr{Net: "unix", Name: "@test"}}
		ln := &ocListener{Listener: l, store: &store}
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
		store := internal.UniqueStore[*ocListener]{}
		l := &nopCloseListener{addr: &net.UnixAddr{Net: "unix", Name: sock}}
		ln := &ocListener{Listener: l, store: &store}
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
