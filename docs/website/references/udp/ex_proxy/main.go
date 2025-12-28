package main

import (
	"log"

	"github.com/aileron-projects/go/znet/zudp"
)

func main() {
	svr := &zudp.Server{
		Addr:    ":8080",
		Handler: zudp.NewProxy("localhost:9090"),
	}

	log.Println("starting udp proxy server at " + svr.Addr)
	if err := svr.ListenAndServe(); err != nil {
		panic(err)
	}
}
