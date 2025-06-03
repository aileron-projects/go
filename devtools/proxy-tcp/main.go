package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aileron-projects/go/znet/ztcp"
)

// curl --unix-socket '/var/run/example.sock' http://localhost:8080/debug
// curl --abstract-unix-socket 'example' http://localhost:8080/debug
func main() {
	// proxy := ztcp.NewProxy("localhost:9090", "localhost:9091")
	proxy := ztcp.NewProxy("localhost:9090")
	svr := &ztcp.Server{
		// Addr:    "@example",
		Addr:    ":8080",
		Handler: proxy,
	}

	runner := &ztcp.ServerRunner{
		Serve:           svr.ListenAndServe,
		Shutdown:        svr.Shutdown,
		Close:           svr.Close,
		ShutdownTimeout: 30 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("starting tcp server at " + svr.Addr)
	if err := runner.Run(ctx); err != nil {
		panic(err)
	}
}
