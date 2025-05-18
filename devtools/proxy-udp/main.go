package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aileron-projects/go/znet/zudp"
)

var (
	addr = flag.String("addr", ":8080", "listen address")
)

func main() {
	// proxy := zudp.NewProxy("localhost:5001", "localhost:5002")
	proxy := zudp.NewProxy("localhost:5001")
	svr := &zudp.Server{
		Addr:    *addr,
		Handler: proxy,
	}

	runner := &zudp.ServerRunner{
		Serve:           svr.ListenAndServe,
		Shutdown:        svr.Shutdown,
		Close:           svr.Close,
		ShutdownTimeout: 30 * time.Second,
	}
	log.Println("starting udp server at " + *addr)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()
	if err := runner.Run(ctx); err != nil {
		panic(err)
	}
}
