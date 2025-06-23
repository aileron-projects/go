package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/aileron-projects/go/ztime/zrate"
)

func main() {
	bucketSize := 10
	fillRate := 10
	limiter := zrate.NewTokenBucketLimiter(bucketSize, fillRate)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := limiter.AllowNow()
		defer token.Release()
		if token.OK() {
			time.Sleep(time.Duration(rand.Int64N(100)) * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("too many requests"))
		}
	})

	log.Println("server listens on localhost:8080")
	svr := &http.Server{
		Addr:        ":8080",
		Handler:     handler,
		ReadTimeout: 10 * time.Second,
	}
	if err := svr.ListenAndServe(); err != nil {
		panic(err)
	}
}
