package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aileron-projects/go/znet/zhttp"
	"github.com/davecgh/go-spew/spew"
)

var (
	addr     = flag.String("addr", ":8080", "listen address")
	timeout  = flag.Int("timeout", 10, "graceful shutdown timeout second")
	certFile = flag.String("certFile", "", "tls cert file path")
	keyFile  = flag.String("keyFile", "", "tls key file path")
)

func main() {
	flag.Parse()
	svr := &http.Server{
		Addr: *addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// time.Sleep(30 * time.Second)
			fmt.Fprintln(w, "Hello Go!!")
			fmt.Fprintln(w, "It's", time.Now())
			fmt.Fprintln(w, strings.Repeat("-", 50))
			fmt.Fprintln(w, spew.Sdump(r))
		}),
		ReadTimeout: 30 * time.Second,
	}
	r := &zhttp.ServerRunner{
		Serve:           svr.ListenAndServe,
		Shutdown:        svr.Shutdown,
		Close:           svr.Close,
		ShutdownTimeout: time.Duration(*timeout) * time.Second,
	}

	// Replace serve function.
	if *certFile != "" || *keyFile != "" {
		r.Serve = func() error { return svr.ListenAndServeTLS(*certFile, *keyFile) }
	}

	sigCtx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	log.Println("starting http server at " + *addr)
	if err := r.Run(sigCtx); err != nil {
		panic(err)
	}
}
