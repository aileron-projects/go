package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aileron-projects/go/znet/zudp"
)

func main() {
	// proxy := zudp.NewProxy("localhost:5001", "localhost:5002")
	proxy := zudp.NewProxy("localhost:5001")
	svr := &zudp.Server{
		Addr:    ":8080",
		Handler: proxy,
	}

	runner := &zudp.ServerRunner{
		Serve:           svr.ListenAndServe,
		Shutdown:        svr.Shutdown,
		Close:           svr.Close,
		ShutdownTimeout: 30 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("starting udp server at " + svr.Addr)
	if err := runner.Run(ctx); err != nil {
		panic(err)
	}
}
