package zrate_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztime/zrate"
)

func BenchmarkMaxConcurrentLimiter(b *testing.B) {
	lim := zrate.NewConcurrentLimiter(math.MaxInt)
	b.ResetTimer()
	for b.Loop() {
		lim.AllowNow()
	}
}

func BenchmarkFixedWindowLimiter(b *testing.B) {
	lim := zrate.NewFixedWindowLimiterWidth(math.MaxInt, time.Nanosecond)
	b.ResetTimer()
	for b.Loop() {
		lim.AllowNow()
	}
}

func BenchmarkNewSlidingWindowLimiter(b *testing.B) {
	lim := zrate.NewSlidingWindowLimiterWidth(math.MaxInt, time.Nanosecond)
	b.ResetTimer()
	for b.Loop() {
		lim.AllowNow()
	}
}

func BenchmarkTokenBucketLimiter(b *testing.B) {
	lim := zrate.NewTokenBucketInterval(math.MaxInt, math.MaxInt, time.Nanosecond)
	b.ResetTimer()
	for b.Loop() {
		lim.AllowNow()
	}
}

func BenchmarkLeakyBucketLimiter(b *testing.B) {
	lim := zrate.NewLeakyBucketLimiter(math.MaxInt, time.Nanosecond)
	b.ResetTimer()
	for b.Loop() {
		lim.WaitNow(context.Background())
	}
}
