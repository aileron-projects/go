package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aileron-projects/go/znet/zhttp"
)

var (
	addr = flag.String("addr", ":8080", "listen address")
)

func main() {
	flag.Parse()
	// p, _ := zhttp.NewProxy("http://httpbin.org/", "http://echo.free.beeceptor.com/foo?alice=bob")
	// p, _ := zhttp.NewProxy("http://sse.dev")
	p, _ := zhttp.NewProxy("http://localhost:9090")

	svr := &http.Server{
		Addr:        *addr,
		Handler:     p,
		ReadTimeout: 30 * time.Second,
	}
	runner := &zhttp.ServerRunner{
		Serve:           svr.ListenAndServe,
		Shutdown:        svr.Shutdown,
		Close:           svr.Close,
		ShutdownTimeout: 5 * time.Second,
	}
	log.Println("starting http server at " + *addr)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()
	if err := runner.Run(ctx); err != nil {
		panic(err)
	}
}
