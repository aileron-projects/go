package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		Addr:    ":8080",
		Handler: ztcp.HandlerFunc(handleConn),
	}

	runner := &ztcp.ServerRunner{
		Serve:           svr.ListenAndServe,
		Shutdown:        svr.Shutdown,
		Close:           svr.Close,
		ShutdownTimeout: 10 * time.Second,
	}

	// Receive SIGINT and SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("starting tcp server at " + svr.Addr)
	if err := runner.Run(ctx); err != nil && err != ztcp.ErrServerClosed {
		panic(err)
	}
}
