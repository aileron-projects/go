package zudp

import (
	"net"
	"sync"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestConn_Read(t *testing.T) {
	t.Parallel()
	// Obtain dummy PacketConn.
	pc, _ := net.ListenPacket("udp4", ":0")
	pc.Close()
	t.Run("read", func(t *testing.T) {
		packets := make(chan []byte, 1)
		c := &conn{
			pc:      pc,
			packets: packets,
		}
		ztesting.AssertEqual(t, "local addr not match", pc.LocalAddr(), c.LocalAddr())

		go func() {
			packets <- []byte("test")
		}()
		buf := make([]byte, 10)
		n, err := c.Read(buf)
		ztesting.AssertEqual(t, "num packet not match", 4, n)
		ztesting.AssertEqual(t, "packet content not match", "test", string(buf[:n]))
		ztesting.AssertEqualErr(t, "error not match", nil, err)
	})
	t.Run("read from closed conn", func(t *testing.T) {
		packets := make(chan []byte, 1)
		raddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}
		c := &conn{
			pc:       pc,
			packets:  packets,
			channels: &sync.Map{},
			raddr:    raddr,
		}
		ztesting.AssertEqual(t, "local addr not match", pc.LocalAddr(), c.LocalAddr())

		c.Close()
		go func() {
			packets <- []byte("test")
		}()
		buf := make([]byte, 10)
		n, err := c.Read(buf)
		ztesting.AssertEqual(t, "num packet not match", 0, n)
		ztesting.AssertEqual(t, "packet content not match", "", string(buf[:n]))
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
}

type recordPacketConn struct {
	net.PacketConn
	raddr   net.Addr
	written []byte
}

func (pc *recordPacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	pc.raddr = addr
	pc.written = p
	return len(p), nil
}

func TestConn_Write(t *testing.T) {
	t.Parallel()
	// Obtain dummy PacketConn.
	pc, _ := net.ListenPacket("udp4", ":0")
	pc.Close()

	t.Run("write", func(t *testing.T) {
		packets := make(chan []byte, 1)
		rpc := &recordPacketConn{PacketConn: pc}
		raddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}
		c := &conn{
			pc:      rpc,
			packets: packets,
			raddr:   raddr,
		}
		ztesting.AssertEqual(t, "remote addr not match", raddr.String(), c.RemoteAddr().String())

		n, err := c.Write([]byte("test"))
		ztesting.AssertEqual(t, "num packet not match", 4, n)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "packet content not match", "test", string(rpc.written))
		ztesting.AssertEqual(t, "written address not match", raddr.String(), rpc.raddr.String())
	})
	t.Run("write to closed conn", func(t *testing.T) {
		packets := make(chan []byte, 1)
		rpc := &recordPacketConn{PacketConn: pc}
		raddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}
		c := &conn{
			pc:       rpc,
			packets:  packets,
			channels: &sync.Map{},
			raddr:    raddr,
		}
		ztesting.AssertEqual(t, "remote addr not match", raddr.String(), c.RemoteAddr().String())

		c.Close()
		go func() {
			packets <- []byte("test")
		}()
		n, err := c.Write([]byte("test"))
		ztesting.AssertEqual(t, "num packet not match", 0, n)
		ztesting.AssertEqual(t, "packet content not match", "", string(rpc.written))
		ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
	})
}

func TestConn_Close(t *testing.T) {
	t.Parallel()
	t.Run("close once", func(t *testing.T) {
		channels := &sync.Map{}
		addr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}
		c := &conn{
			raddr:    &addr,
			channels: channels,
		}
		channels.Store(addr.String(), "test")
		_, ok := channels.Load(addr.String())
		ztesting.AssertEqual(t, "value not found", true, ok)
		c.Close()
		_, ok = channels.Load(addr.String())
		ztesting.AssertEqual(t, "value unexpectedly found", false, ok)
	})
	t.Run("close multiple times", func(t *testing.T) {
		channels := &sync.Map{}
		addr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}
		c := &conn{
			raddr:    &addr,
			channels: channels,
		}
		channels.Store(addr.String(), "test")
		_, ok := channels.Load(addr.String())
		ztesting.AssertEqual(t, "value not found", true, ok)
		c.Close()
		_, ok = channels.Load(addr.String())
		ztesting.AssertEqual(t, "key found", false, ok)

		channels.Store(addr.String(), "test") // Store value again.
		c.Close()                             // Value should not be removed.
		_, ok = channels.Load(addr.String())
		ztesting.AssertEqual(t, "value not found", true, ok)
	})
}
