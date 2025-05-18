package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aileron-projects/go/ztime/zrate"
	"gopkg.in/yaml.v3"
)

var (
	file = flag.String("file", "config.yaml", "config file path")
)

func main() {
	flag.Parse()
	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}
	b, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}
	config := &config{}
	if err := yaml.Unmarshal(b, config); err != nil {
		panic(err)
	}

	var limiter zrate.Limiter
	var waitMode bool
	switch config.UseLimiter {
	case "concurrentLimiter":
		c := config.CLimiter
		limiter = zrate.NewConcurrentLimiter(c.Limit)
		fmt.Println("Use concurrentLimiter:")
		fmt.Println("    |", "limit=", c.Limit)
	case "fixedWindowLimiter":
		c := config.FLimiter
		limiter = zrate.NewFixedWindowLimiterWidth(c.Limit, time.Duration(c.Width)*time.Millisecond)
		fmt.Println("Use fixedWindowLimiter:")
		fmt.Println("    |", "limit=", c.Limit)
		fmt.Println("    |", "width=", c.Width, " millisecond")
	case "slidingWindowLimiter":
		c := config.SLimiter
		limiter = zrate.NewSlidingWindowLimiterWidth(c.Limit, time.Duration(c.Width)*time.Millisecond)
		fmt.Println("Use slidingWindowLimiter:")
		fmt.Println("    |", "limit=", c.Limit)
		fmt.Println("    |", "width=", c.Width, " millisecond")
	case "tokenBucketLimiter":
		c := config.TLimiter
		limiter = zrate.NewTokenBucketInterval(c.BucketSize, c.FillRate, time.Duration(c.FillInterval)*time.Millisecond)
		fmt.Println("Use tokenBucketLimiter:")
		fmt.Println("    |", "bucketSize=", c.BucketSize)
		fmt.Println("    |", "fillRate=", c.FillRate)
		fmt.Println("    |", "fillInterval=", c.FillInterval, " millisecond")
	case "leakyBucketLimiter":
		c := config.LLimiter
		limiter = zrate.NewLeakyBucketLimiter(c.QueueSize, time.Duration(c.Interval)*time.Millisecond)
		fmt.Println("Use leakyBucketLimiter:")
		fmt.Println("    |", "queueSize=", c.QueueSize)
		fmt.Println("    |", "interval=", c.Interval, " millisecond")
		waitMode = true
	default:
		panic("limiter not defined.")
	}

	fmt.Println("\nserver listens on:", config.Listen)

	h := &testHandler{
		limiter:  limiter,
		waitMode: waitMode,
	}
	go startReporting(h, time.Duration(config.ReportIntervalSec)*time.Second)
	svr := http.Server{
		Addr:    config.Listen,
		Handler: h,
	}
	if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func startReporting(h *testHandler, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	fmt.Println("\nTimestamp", " - ", "ALL", "OK", "NG")
	for {
		h.mu.Lock()
		nReq, nOK, nNG := h.nReq, h.nOK, h.nNG
		h.nReq, h.nOK, h.nNG = 0, 0, 0
		h.mu.Unlock()
		fmt.Println(time.Now().Local().Format(time.TimeOnly), " - ", nReq, nOK, nNG)
		<-ticker.C
	}
}

type testHandler struct {
	mu       sync.RWMutex
	nReq     uint64
	nOK      uint64
	nNG      uint64
	limiter  zrate.Limiter
	waitMode bool
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.nReq += 1
	h.mu.Unlock()

	var token zrate.Token
	if h.waitMode {
		token = h.limiter.WaitNow(r.Context())
		defer token.Release()
	} else {
		token = h.limiter.AllowNow()
		defer token.Release()
	}

	h.mu.Lock()
	if token.OK() {
		h.nOK += 1
		w.WriteHeader(http.StatusOK)
	} else {
		h.nNG += 1
		w.WriteHeader(http.StatusTooManyRequests)
	}
	h.mu.Unlock()

	if !token.OK() {
		return
	}
	time.Sleep(time.Millisecond * time.Duration(rand.Int64N(100)))
}

type config struct {
	Listen            string                `yaml:"listen"`
	ReportIntervalSec int                   `yaml:"reportIntervalSec"`
	UseLimiter        string                `yaml:"useLimiter"`
	CLimiter          *ConcurrentLimiter    `yaml:"concurrentLimiter"`
	FLimiter          *FixedWindowLimiter   `yaml:"fixedWindowLimiter"`
	SLimiter          *SlidingWindowLimiter `yaml:"slidingWindowLimiter"`
	TLimiter          *TokenBucketLimiter   `yaml:"tokenBucketLimiter"`
	LLimiter          *LeakyBucketLimiter   `yaml:"leakyBucketLimiter"`
}

type ConcurrentLimiter struct {
	Limit int `yaml:"limit"`
}

type FixedWindowLimiter struct {
	Limit int `yaml:"limit"`
	Width int `yaml:"width"`
}

type SlidingWindowLimiter struct {
	Limit int `yaml:"limit"`
	Width int `yaml:"width"`
}

type TokenBucketLimiter struct {
	BucketSize   int `yaml:"bucketSize"`
	FillInterval int `yaml:"fillInterval"`
	FillRate     int `yaml:"fillRate"`
}

type LeakyBucketLimiter struct {
	QueueSize int `yaml:"queueSize"`
	Interval  int `yaml:"interval"`
}
