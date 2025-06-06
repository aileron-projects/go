package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

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
	// You can use curl to check if it works.
	// curl --unix-socket '/var/run/example.sock' http://localhost:8080/example
	svr := &ztcp.Server{
		Addr:    "unix:///var/run/example.sock",
		Handler: ztcp.HandlerFunc(handleConn),
	}
	log.Println("starting tcp server at " + svr.Addr)
	if err := svr.ListenAndServe(); err != nil && err != ztcp.ErrServerClosed {
		panic(err)
	}
}

// END main
