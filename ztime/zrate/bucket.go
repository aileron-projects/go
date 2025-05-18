package zrate

import (
	"context"
	"math"
	"sync"
	"time"
)

// NewFixedWindowLimiter returns a new instance of
// a limiter that works with fixed-window algorithm.
// NewFixedWindowLimiter is short for NewFixedWindowLimiterWidth(limit, time.Second).
func NewFixedWindowLimiter(limit int) Limiter {
	return NewFixedWindowLimiterWidth(limit, time.Second)
}

// NewFixedWindowLimiterWidth returns a new instance of
// a limiter that works with fixed-window algorithm.
// The limit is the maximum count allowed within the width.
// For limit<=0, the limiter always returns token that indicates dis-allow.
// For width<=0, the limiter always returns token that indicates allow.
func NewFixedWindowLimiterWidth(limit int, width time.Duration) Limiter {
	if limit <= 0 {
		return NoopLimiter(false)
	}
	if width <= 0 {
		return NoopLimiter(true)
	}
	return &BucketLimiter{
		getToken: getTokenFunc(limit, limit, width, time.Now),
	}
}

// NewTokenBucketLimiter returns a new limiter instance
// that works with token bucket algorithm.
// NewTokenBucketLimiter is short for NewTokenBucketInterval(bucketSize, fillRate, time.Second).
func NewTokenBucketLimiter(bucketSize, fillRate int) Limiter {
	return NewTokenBucketInterval(bucketSize, fillRate, time.Second)
}

// NewTokenBucketLimiter returns a new limiter instance
// that works with token bucket algorithm.
// For bucketSize<=0, the limiter always returns token that indicates dis-allow.
// For fillInterval<=0, the limiter always returns token that indicates allow.
func NewTokenBucketInterval(bucketSize, fillRate int, fillInterval time.Duration) Limiter {
	if bucketSize <= 0 {
		return NoopLimiter(false)
	}
	if fillInterval <= 0 {
		return NoopLimiter(true)
	}
	return &BucketLimiter{
		getToken: getTokenFunc(bucketSize, max(0, fillRate), fillInterval, time.Now),
	}
}

// BucketLimiter limits the rate of something to be proceeded using bucket algorithm.
// BucketLimiter can be used for Fixed Window Algorithm and Token Bucket Algorithm.
//
//   - Bucket Size == Token Fill Rate  -->  Fixed Window Algorithm (Use [BucketLimiter.AllowNow])
//   - Bucket Size >= Token Fill Rate  -->  Token Bucket Algorithm (Use [BucketLimiter.WaitNow])
//
// When using limiter as API rate limiting, implementation would be like below.
//
// Example usage:
//
//	type handler struct {
//		limiter zrate.Limiter
//	}
//
//	func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//		token := h.limiter.AllowNow() // Basically, use AllowNow for BucketLimiter.
//		defer token.Release()         // Release is not required for BucketLimiter.
//		if !token.OK() {
//			w.WriteHeader(http.StatusTooManyRequests)
//			return
//		}
//		// Some process.
//		w.WriteHeader(http.StatusOK)
//	}
type BucketLimiter struct {
	getToken func() (ok bool, retryAfter time.Duration)
}

func (lim *BucketLimiter) AllowNow() Token {
	ok, _ := lim.getToken()
	if ok {
		return TokenOK
	}
	return TokenNG
}

func (lim *BucketLimiter) WaitNow(ctx context.Context) Token {
	for {
		ok, retryAfter := lim.getToken()
		if ok {
			return TokenOK
		}
		select {
		case <-time.After(retryAfter):
		case <-ctx.Done():
			return &token{err: ctx.Err()}
		}
	}
}

// getTokenFunc returns function to get tokens.
func getTokenFunc(bucketSize, fillRate int, interval time.Duration, timeNow func() time.Time) func() (ok bool, retryAfter time.Duration) {
	if bucketSize <= 0 {
		return func() (bool, time.Duration) { return false, math.MaxInt64 }
	}
	if interval <= 0 {
		return func() (bool, time.Duration) { return true, 0 }
	}
	return (&bucket{
		tokens:       int64(bucketSize),
		bucketSize:   float64(bucketSize),
		fillRate:     float64(max(fillRate, 0)),
		fillInterval: interval,
		lastFilled:   timeNow(),
		timeNow:      timeNow,
	}).getToken
}

// bucket is the token bucket for limiters.
// - Bucket Size == Fill Rate  ---> Fixed Window Algorithm
// - Bucket Size >= Fill Rate  ---> Token Bucket Algorithm
type bucket struct {
	mu sync.Mutex
	// tokens is the number of tokens
	// available for consume.
	tokens int64
	// bucketSize is the size of bucket.
	// It is held as float64 for calculation.
	bucketSize float64
	// fillRate is the number of tokens
	// that will be re-filled to the bucket
	// at every fillInterval.
	// It is held as float64 for calculation.
	fillRate float64
	// fillInterval is the token re-fill interval.
	// If bucket is empty, tokens are re-filled by following equation.
	// tokens = max(bucketSize, fillRate*(now-lastFilled)/fillInterval )
	// fillInterval must not be zero.
	fillInterval time.Duration
	// lastFilled is the last time that tokens are
	// re-filled to the bucket.
	// If bucket is empty, tokens are re-filled by following equation.
	// tokens = max(bucketSize, fillRate*(now-lastFilled)/fillInterval )
	lastFilled time.Time

	// timeNow returns the current time.
	// This can be replaced for testing.
	timeNow func() time.Time
}

func (b *bucket) getToken() (ok bool, retryAfter time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.tokens > 0 {
		b.tokens -= 1
		return true, 0
	}

	now := b.timeNow()
	passed := now.Sub(b.lastFilled)
	if passed < b.fillInterval {
		return false, b.fillInterval - passed
	}

	x := b.fillRate * float64(passed) / float64(b.fillInterval)
	x = min(x, b.bucketSize) // Limit up to the bucket size.

	b.tokens = int64(x) - 1 // Consume 1 for this request.
	b.lastFilled = now
	return true, 0
}
