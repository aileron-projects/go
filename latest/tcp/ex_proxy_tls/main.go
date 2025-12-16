package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"

	"github.com/aileron-projects/go/znet/ztcp"
)

func main() {
	pem, _ := os.ReadFile("cert.pem")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)

	dialer := &tls.Dialer{
		Config: &tls.Config{RootCAs: pool},
	}
	svr := &ztcp.Server{
		Addr: ":8080",
		Handler: &ztcp.Proxy{
			Dial: func(ctx context.Context, dc net.Conn) (uc net.Conn, err error) {
				return dialer.DialContext(context.Background(), "tcp", "localhost:9090")
			},
		},
	}

	log.Println("starting tcp proxy server at " + svr.Addr)
	if err := svr.ListenAndServe(); err != nil {
		panic(err)
	}
}
