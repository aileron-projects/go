package zudp

import (
	"net"
	"sync"
	"sync/atomic"
)

// Conn is a virtual connection for UDP.
// A connection is bounded to a remote address.
// Read and write operation is always from/to the
// bounded remote address.
type Conn interface {
	// Read reads datagram packet from the connection.
	// Received packets are sent by the remote address.
	// If the length of b is shorter than received packet size,
	// the rest of packets are discarded without returning an error.
	Read(b []byte) (n int, err error)
	// Write writes a datagram packet to the connection.
	// Written packets are sent to remote address.
	Write(b []byte) (n int, err error)
	// LocalAddr returns the local network address, if known.
	LocalAddr() net.Addr
	// RemoteAddr returns the remote network address, if known.
	RemoteAddr() net.Addr
	// Close closes the connection.
	// Calling [Conn.Read] or [Conn.Write] after Close is called
	// returns an [net.ErrClosed].
	Close() error
}

// conn is a virtual UDP connection.
// It read and write to a same remote address.
type conn struct {
	// pc is the PacketConn that this connection is
	// associated to. pc.WriteTo is used to write to
	// the bounded remote address.
	pc net.PacketConn
	// raddr is the remote address that this connection write to.
	raddr net.Addr
	// packets is the buffered channel to receive UDP packets from.
	// Received packets are the packets sent from the remote address.
	packets <-chan []byte
	// closed keeps if this connection is closed or not.
	// Once closed, later read or write returns [net.ErrClosed].
	closed atomic.Bool
	// channels stores packet channels with the key of raddr.String().
	// The entry must be deleted from the map when this connection was
	// closed.
	channels *sync.Map
}

func (c *conn) Read(b []byte) (n int, err error) {
	if c.closed.Load() {
		return 0, net.ErrClosed
	}
	packet := <-c.packets
	return copy(b, packet), nil
}

func (c *conn) Write(b []byte) (n int, err error) {
	if c.closed.Load() {
		return 0, net.ErrClosed
	}
	return c.pc.WriteTo(b, c.raddr)
}

func (c *conn) LocalAddr() net.Addr {
	return c.pc.LocalAddr()
}

func (c *conn) RemoteAddr() net.Addr {
	return c.raddr
}

func (c *conn) Close() error {
	if c.closed.Swap(true) {
		return nil
	}
	c.channels.Delete(c.raddr.String())
	return nil
}
