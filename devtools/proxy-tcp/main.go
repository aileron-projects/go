package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aileron-projects/go/znet/ztcp"
)

var (
	addr = flag.String("addr", ":8080", "listen address")
)

// curl --unix-socket '/var/run/example.sock' http://localhost:8080/debug
// curl --abstract-unix-socket 'example' http://localhost:8080/debug
func main() {
	// proxy := ztcp.NewProxy("localhost:9090", "localhost:9091")
	proxy := ztcp.NewProxy("localhost:9090")
	svr := &ztcp.Server{
		// Addr:    "@example",
		Addr:    *addr,
		Handler: proxy,
	}

	runner := &ztcp.ServerRunner{
		Serve:           svr.ListenAndServe,
		Shutdown:        svr.Shutdown,
		Close:           svr.Close,
		ShutdownTimeout: 60 * time.Second,
	}
	log.Println("starting tcp server at " + *addr)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()
	if err := runner.Run(ctx); err != nil {
		panic(err)
	}
}
