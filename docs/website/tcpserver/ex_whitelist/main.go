package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/aileron-projects/go/znet"
	"github.com/aileron-projects/go/znet/ztcp"
)

// handleConn reads and prints TCP data received from the conn.
func handleConn(ctx context.Context, conn net.Conn) {
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

// START main
func main() {
	svr := &ztcp.Server{
		Addr:    "", // This is not used when we call [Server.Serve].
		Handler: ztcp.HandlerFunc(handleConn),
	}

	ln, _ := net.Listen("tcp", ":8080")                   // Create a new TCP listener.
	ln, _ = znet.NewWhiteListListener(ln, "127.0.0.1/32") // Apply whitelist.

	log.Println("starting tcp server at " + ln.Addr().String())
	if err := svr.Serve(ln); err != nil && err != ztcp.ErrServerClosed {
		panic(err)
	}
}

// END main
