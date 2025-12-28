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

func main() {
	svr := &ztcp.Server{
		Addr:      ":8080",
		Handler:   ztcp.HandlerFunc(handleConn),
		TLSConfig: nil,
	}
	log.Println("starting tcp server at " + svr.Addr)
	if err := svr.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != ztcp.ErrServerClosed {
		panic(err)
	}
}
