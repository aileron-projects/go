package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/aileron-projects/go/znet"
	"github.com/aileron-projects/go/znet/zudp"
)

// handleConn reads and prints UDP packets received from the conn.
func handleConn(ctx context.Context, conn zudp.Conn) {
	buf := make([]byte, 1<<10)
	for {
		n, err := conn.Read(buf)
		fmt.Println(string(buf[:n]))
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			panic(err)
		}
	}
}

type WhitelistPacketConn struct {
	net.PacketConn
	wl *znet.WhiteList
}

func (c *WhitelistPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, addr, err = c.PacketConn.ReadFrom(p)
	if err == nil {
		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			host = addr.String() // Fallback
		}
		if !c.wl.Allowed(host) {
			return n, addr, zudp.ErrSkipHandler // Return zudp.ErrSkipHandler.
		}
	}
	return n, addr, err
}

func main() {
	svr := &zudp.Server{
		Addr:    "",
		Handler: zudp.HandlerFunc(handleConn),
	}

	var lc net.ListenConfig
	p, err := lc.ListenPacket(context.Background(), "udp", ":8080")
	if err != nil {
		panic(err)
	}
	wl := znet.NewWhiteList()
	_ = wl.Allow("127.0.0.1/32", "::1/128")
	pc := &WhitelistPacketConn{PacketConn: p, wl: wl}

	log.Println("starting udp server at " + pc.LocalAddr().String())
	if err := svr.Serve(pc); err != nil && err != zudp.ErrServerClosed {
		panic(err)
	}
}
