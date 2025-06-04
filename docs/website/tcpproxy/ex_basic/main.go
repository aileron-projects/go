package main

import (
	"log"

	"github.com/aileron-projects/go/znet/ztcp"
)

func main() {
	svr := &ztcp.Server{
		Addr:    ":8080",
		Handler: ztcp.NewProxy("localhost:9090"),
	}

	log.Println("starting tcp proxy server at " + svr.Addr)
	if err := svr.ListenAndServe(); err != nil {
		panic(err)
	}
}
