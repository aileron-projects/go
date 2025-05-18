package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/aileron-projects/go/znet"
)

var (
	addr = flag.String("addr", ":8080", "listen address")
)

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	ln, _ = znet.NewBlackListListener(ln, "127.0.0.2/32")
	svr := &http.Server{
		Addr: *addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello Go!!")
		}),
		ReadTimeout: 30 * time.Second,
	}

	log.Println("starting http server at " + *addr)
	if err := svr.Serve(ln); err != nil {
		panic(err)
	}
}
