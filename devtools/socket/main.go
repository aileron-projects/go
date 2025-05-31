package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/aileron-projects/go/zsyscall"
)

// On linux, socket options can be checked with the command:
// 	strace -C -f -e trace=setsockopt go run ./

var opts = &zsyscall.SockOption{
	SO: &zsyscall.SockSOOption{
		BindToIFindex: 3,
		ReuseAddr:     true,
		ReusePort:     true,
		KeepAlive:     true,
	},
	IP: &zsyscall.SockIPOption{
		// LocalPortRangeLower: 10000,
		// LocalPortRangeUpper: 10010,
	},
	IPV6: &zsyscall.SockIPV6Option{},
	TCP: &zsyscall.SockTCPOption{
		NoDelay: true,
	},
	UDP: &zsyscall.SockUDPOption{},
}

func main() {
	targets := zsyscall.SockOptSO | zsyscall.SockOptIP | zsyscall.SockOptIPV6 | zsyscall.SockOptTCP
	lc := &net.ListenConfig{
		Control: opts.ControlFunc(targets),
	}

	ln, err := lc.Listen(context.Background(), "tcp", ":8080")
	if err != nil {
		panic(err)
	}
	svr := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello Go!!")
			fmt.Fprintln(w, "It's", time.Now())
		}),
		ReadTimeout: 30 * time.Second,
	}

	log.Println("starting http server at", ln.Addr())
	if err := svr.Serve(ln); err != nil {
		panic(err)
	}
}
