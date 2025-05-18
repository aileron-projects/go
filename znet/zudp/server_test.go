package zudp

import (
	"net"
	"os"
	"testing"

	"github.com/aileron-projects/go/znet/internal"
	"github.com/aileron-projects/go/ztesting"
)

func TestNewPacketConn(t *testing.T) {
	t.Parallel()
	// Obtain available address.
	pc, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		panic(err)
	}
	pc.Close()
	addr := pc.LocalAddr().String()

	t.Run("udp without prefix", func(t *testing.T) {
		ln, err := newPacketConn("" + addr)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		defer ln.Close()
		cn, err := net.Dial("udp", addr)
		ztesting.AssertEqual(t, "dial failed", nil, err)
		cn.Close()
	})
	t.Run("listen udp4 success", func(t *testing.T) {
		ln, err := newPacketConn("udp4://" + addr)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		defer ln.Close()
		cn, err := net.Dial("udp4", addr)
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
