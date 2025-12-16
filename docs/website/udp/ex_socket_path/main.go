package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

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

func main() {
	svr := &zudp.Server{
		Addr:    "unixgram:///var/run/example.sock",
		Handler: zudp.HandlerFunc(handleConn),
	}
	log.Println("starting udp server at " + svr.Addr)
	if err := svr.ListenAndServe(); err != nil && err != zudp.ErrServerClosed {
		panic(err)
	}
}
