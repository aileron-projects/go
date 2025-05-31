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

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	ln, _ = znet.NewBlackListListener(ln, "127.0.0.2/32")
	svr := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello Go!!")
		}),
		ReadTimeout: 30 * time.Second,
	}

	log.Println("starting http server at ", ln.Addr().String())
	if err := svr.Serve(ln); err != nil {
		panic(err)
	}
}
