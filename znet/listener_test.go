package znet_test

import (
	"errors"
	"io"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/aileron-projects/go/znet"
	"github.com/aileron-projects/go/ztesting"
)

type testConn struct {
	net.Conn
	closed      bool
	closedCount int
	remote      net.Addr
}

func (c *testConn) Close() error {
	c.closed = true
	c.closedCount++
	return errors.New(strconv.Itoa(c.closedCount))
}

func (c *testConn) RemoteAddr() net.Addr {
	return c.remote
}

type testListener struct {
	net.Listener
	err    error
	remote net.Addr
}

func (l *testListener) Accept() (net.Conn, error) {
	return &testConn{remote: l.remote}, l.err
}

func TestNewWhiteListListener(t *testing.T) {
	t.Parallel()
	t.Run("error", func(t *testing.T) {
		_, err := znet.NewWhiteListListener(nil, "127.0.0.1")
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("no error", func(t *testing.T) {
		_, err := znet.NewWhiteListListener(nil, "127.0.0.1/32")
		ztesting.AssertEqual(t, "error should be nil", true, err == nil)
	})
}

func TestNewBlackListListener(t *testing.T) {
	t.Parallel()
	t.Run("error", func(t *testing.T) {
		_, err := znet.NewBlackListListener(nil, "127.0.0.1")
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("no error", func(t *testing.T) {
		_, err := znet.NewBlackListListener(nil, "127.0.0.1/32")
		ztesting.AssertEqual(t, "error should be nil", true, err == nil)
	})
}

func TestWhiteListListener(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ln        net.Listener
		allow     []string
		allowFunc func(host, port string) bool
		closed    bool  // want
		err       error // want
	}{
		"allowed": {
			ln:    &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}},
			allow: []string{"127.0.0.0/24"},
		},
		"not allowed": {
			ln:     &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}},
			allow:  []string{"127.0.1.0/24"},
			closed: true,
		},
		"listener error": {
			ln:    &testListener{err: io.EOF}, // dummy error
			allow: []string{"127.0.1.0/24"},
			err:   io.EOF,
		},
		"host port error": {
			ln:     &testListener{remote: &net.UnixAddr{Name: "@example"}},
			closed: true,
		},
		"optionally allowed": {
			ln:        &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}},
			allow:     []string{"127.0.0.0/24"},
			allowFunc: func(host, port string) bool { return port == "80" },
		},
		"optionally not allowed": {
			ln:        &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}},
			allow:     []string{"127.0.0.0/24"},
			allowFunc: func(host, port string) bool { return port == "8080" },
			closed:    true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ln, _ := znet.NewWhiteListListener(tc.ln, tc.allow...)
			ln.Allow = tc.allowFunc
			conn, err := ln.Accept()
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
			ztesting.AssertEqual(t, "closed not match", tc.closed, conn.(*testConn).closed)
		})
	}
}

func TestBlackListListener(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ln        net.Listener
		disallow  []string
		allowFunc func(host, port string) bool
		closed    bool  // want
		err       error // want
	}{
		"allowed": {
			ln:       &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}},
			disallow: []string{"127.0.1.0/24"},
		},
		"not allowed": {
			ln:       &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}},
			disallow: []string{"127.0.0.0/24"},
			closed:   true,
		},
		"listener error": {
			ln:       &testListener{err: io.EOF}, // dummy error
			disallow: []string{"127.0.1.0/24"},
			err:      io.EOF,
		},
		"host port error": {
			ln:     &testListener{remote: &net.UnixAddr{Name: "@example"}},
			closed: true,
		},
		"optionally allowed": {
			ln:        &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}},
			allowFunc: func(host, port string) bool { return port == "80" },
		},
		"optionally not allowed": {
			ln:        &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}},
			allowFunc: func(host, port string) bool { return port == "8080" },
			closed:    true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ln, _ := znet.NewBlackListListener(tc.ln, tc.disallow...)
			ln.Allow = tc.allowFunc
			conn, err := ln.Accept()
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
			ztesting.AssertEqual(t, "closed not match", tc.closed, conn.(*testConn).closed)
		})
	}
}

func TestLimitListener(t *testing.T) {
	t.Parallel()
	t.Run("limit 0", func(t *testing.T) {
		inner := &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}}
		ln := znet.NewLimitListener(inner, 0)
		conn, err := ln.Accept()
		ztesting.AssertEqualErr(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "closed count not match", "1", conn.Close().Error())
	})
	t.Run("limit 1", func(t *testing.T) {
		inner := &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}}
		ln := znet.NewLimitListener(inner, 1)
		conn, err := ln.Accept()
		ztesting.AssertEqualErr(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "closed count not match", "1", conn.Close().Error())
	})
	t.Run("limit 2", func(t *testing.T) {
		inner := &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}}
		ln := znet.NewLimitListener(inner, 2)
		_, _ = ln.Accept() // Discard first
		conn, err := ln.Accept()
		ztesting.AssertEqualErr(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "closed count not match", "1", conn.Close().Error())
	})
	t.Run("wait", func(t *testing.T) {
		inner := &testListener{remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}}
		ln := znet.NewLimitListener(inner, 1)
		conn, _ := ln.Accept()
		time.AfterFunc(100*time.Millisecond, func() { conn.Close() })
		conn, err := ln.Accept()
		ztesting.AssertEqualErr(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "closed count not match", "1", conn.Close().Error())
	})
	t.Run("inner error", func(t *testing.T) {
		inner := &testListener{err: io.EOF, remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}}
		ln := znet.NewLimitListener(inner, 1)
		conn, err := ln.Accept()
		ztesting.AssertEqualErr(t, "error not match", io.EOF, err)
		ztesting.AssertEqual(t, "closed count not match", "1", conn.Close().Error())
	})
}
