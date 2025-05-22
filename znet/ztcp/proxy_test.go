package ztcp

import (
	"bytes"
	"cmp"
	"context"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestNewProxy(t *testing.T) {
	t.Parallel()
	t.Run("no targets", func(t *testing.T) {
		defer func() {
			r := recover()
			ztesting.AssertEqual(t, "recovered value not match", r.(error), ErrNoTarget)
		}()
		NewProxy()
	})
	t.Run("with targets", func(t *testing.T) {
		p := NewProxy("foo", "bar")
		ztesting.AssertEqual(t, "dialer is nil", true, p.Dial != nil)
	})
}

func TestRoundRobinDialer(t *testing.T) {
	t.Parallel()
	t.Run("test next", func(t *testing.T) {
		rrd := &roundRobinDialer{
			index: -1,
			addrs: []string{"addr1", "addr2", "addr3"},
		}
		got := []string{}
		for range 6 {
			got = append(got, rrd.next())
		}
		want := []string{"addr1", "addr2", "addr3", "addr1", "addr2", "addr3"}
		ztesting.AssertEqual(t, "address not match", want, got)
	})
	t.Run("invalid tcp address", func(t *testing.T) {
		rrd := &roundRobinDialer{addrs: []string{"tcp://12345"}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "conn should be nil", nil, conn)
		_, ok := err.(*net.AddrError)
		ztesting.AssertEqual(t, "error should be addr error", true, ok)
	})
	t.Run("dial tcp", func(t *testing.T) {
		ln, _ := net.Listen("tcp4", ":0")
		defer ln.Close()
		rrd := &roundRobinDialer{addrs: []string{"tcp4://" + ln.Addr().String()}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		_ = conn.Close()
	})
	t.Run("dial unix", func(t *testing.T) {
		s := filepath.Join(os.TempDir(), "test.sock")
		ln, _ := net.Listen("unix", s)
		defer ln.Close()
		rrd := &roundRobinDialer{addrs: []string{"unix://" + s}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		_ = conn.Close()
	})
	t.Run("dial fallback", func(t *testing.T) {
		ln, _ := net.Listen("tcp", ":0")
		defer ln.Close()
		rrd := &roundRobinDialer{addrs: []string{ln.Addr().String()}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		_ = conn.Close()
	})
}

func TestProxy_handleError(t *testing.T) {
	t.Parallel()
	t.Run("handle non-nil", func(t *testing.T) {
		var got error
		var called bool
		p := &Proxy{
			ErrorHandler: func(dc, uc net.Conn, err error) {
				called = true
				got = err
			},
		}
		p.handleError(nil, nil, io.EOF)
		ztesting.AssertEqual(t, "handler not called", true, called)
		ztesting.AssertEqualErr(t, "error not match", io.EOF, got)
	})
	t.Run("handle nil", func(t *testing.T) {
		var got error
		var called bool
		p := &Proxy{
			ErrorHandler: func(dc, uc net.Conn, err error) {
				called = true
				got = err
			},
		}
		p.handleError(nil, nil, nil)
		ztesting.AssertEqual(t, "handler called", false, called)
		ztesting.AssertEqualErr(t, "error not match", nil, got)
	})
}

type testProxyConn struct {
	net.Conn
	reader io.Reader
	writer io.Writer
	closed bool
}

func (c *testProxyConn) Read(p []byte) (n int, err error) {
	return c.reader.Read(p)
}

func (c *testProxyConn) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}

func (c *testProxyConn) Close() error {
	c.closed = true
	return nil
}

func TestProxy(t *testing.T) {
	t.Parallel()
	t.Run("proxy successfully finish", func(t *testing.T) {
		dWriter := bytes.NewBuffer(nil)
		dConn := &testProxyConn{
			reader: strings.NewReader("downstream data"),
			writer: dWriter,
		}
		uWriter := bytes.NewBuffer(nil)
		uConn := &testProxyConn{
			reader: strings.NewReader("upstream data"),
			writer: uWriter,
		}
		var handledErr error
		p := &Proxy{
			Dial: func(ctx context.Context, dc net.Conn) (uc net.Conn, err error) {
				return uConn, nil
			},
			ErrorHandler: func(dc, uc net.Conn, err error) { handledErr = err },
		}
		p.ServeTCP(context.Background(), dConn)
		ztesting.AssertEqual(t, "upstream data was not written", "upstream data", dWriter.String())
		ztesting.AssertEqual(t, "downstream data was not written", "downstream data", uWriter.String())
		ztesting.AssertEqual(t, "upstream conn was not closed", true, uConn.closed)
		ztesting.AssertEqual(t, "downstream conn was closed", false, dConn.closed)
		ztesting.AssertEqualErr(t, "error not match", nil, handledErr)
	})
	t.Run("dial error", func(t *testing.T) {
		var handledErr error
		p := &Proxy{
			Dial: func(ctx context.Context, dc net.Conn) (uc net.Conn, err error) {
				return nil, net.ErrClosed
			},
			ErrorHandler: func(dc, uc net.Conn, err error) { handledErr = err },
		}
		p.ServeTCP(context.Background(), nil)
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, handledErr)
	})
	t.Run("downstream read error", func(t *testing.T) {
		dWriter := bytes.NewBuffer(nil)
		dConn := &testProxyConn{
			reader: ziotest.ErrReader(strings.NewReader("downstream data"), 10),
			writer: dWriter,
		}
		uWriter := bytes.NewBuffer(nil)
		uConn := &testProxyConn{
			reader: strings.NewReader("upstream data"),
			writer: uWriter,
		}
		var handledErr error
		p := &Proxy{
			Dial: func(ctx context.Context, dc net.Conn) (uc net.Conn, err error) {
				return uConn, nil
			},
			ErrorHandler: func(dc, uc net.Conn, err error) { handledErr = cmp.Or(handledErr, err) },
		}
		p.ServeTCP(context.Background(), dConn)
		for dWriter.Len() == 0 || uWriter.Len() == 0 {
			time.Sleep(100 * time.Millisecond) // Wait both written.
		}
		ztesting.AssertEqual(t, "upstream data was not written", "upstream data", dWriter.String())
		ztesting.AssertEqual(t, "downstream data was not written", "downstream", uWriter.String())
		ztesting.AssertEqual(t, "upstream conn was not closed", true, uConn.closed)
		ztesting.AssertEqual(t, "downstream conn was closed", false, dConn.closed)
		ztesting.AssertEqualErr(t, "error not match", io.ErrClosedPipe, handledErr)
	})
	t.Run("upstream read error", func(t *testing.T) {
		dWriter := bytes.NewBuffer(nil)
		dConn := &testProxyConn{
			reader: strings.NewReader("downstream data"),
			writer: dWriter,
		}
		uWriter := bytes.NewBuffer(nil)
		uConn := &testProxyConn{
			reader: ziotest.ErrReader(strings.NewReader("upstream data"), 8),
			writer: uWriter,
		}
		var handledErr error
		p := &Proxy{
			Dial: func(ctx context.Context, dc net.Conn) (uc net.Conn, err error) {
				return uConn, nil
			},
			ErrorHandler: func(dc, uc net.Conn, err error) { handledErr = cmp.Or(handledErr, err) },
		}
		p.ServeTCP(context.Background(), dConn)
		for dWriter.Len() == 0 || uWriter.Len() == 0 {
			time.Sleep(100 * time.Millisecond) // Wait both written.
		}
		ztesting.AssertEqual(t, "upstream data was not written", "upstream", dWriter.String())
		ztesting.AssertEqual(t, "downstream data was not written", "downstream data", uWriter.String())
		ztesting.AssertEqual(t, "upstream conn was not closed", true, uConn.closed)
		ztesting.AssertEqual(t, "downstream conn was closed", false, dConn.closed)
		ztesting.AssertEqualErr(t, "error not match", io.ErrClosedPipe, handledErr)
	})
}
