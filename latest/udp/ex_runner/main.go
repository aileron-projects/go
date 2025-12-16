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
		Addr:    ":8080",
		Handler: zudp.HandlerFunc(handleConn),
	}

	runner := &zudp.ServerRunner{
		Serve:           svr.ListenAndServe,
		Shutdown:        svr.Shutdown,
		Close:           svr.Close,
		ShutdownTimeout: 10 * time.Second,
	}

	// Receive SIGINT and SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("starting udp server at " + svr.Addr)
	if err := runner.Run(ctx); err != nil && err != zudp.ErrServerClosed {
		panic(err)
	}
}
