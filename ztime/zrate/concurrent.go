package zrate

import (
	"context"
)

// NewConcurrentLimiter returns a new instance of ConcurrentLimiter
// configured with the given parameter.
// The limit is the maximum number of concurrency.
// If limit is 0 or negative, the limiter always returns token that indicates dis-allow.
func NewConcurrentLimiter(limit int) Limiter {
	if limit <= 0 {
		return NoopLimiter(false)
	}
	return &ConcurrentLimiter{
		bucket: make(chan struct{}, limit),
	}
}

// ConcurrentLimiter limits the max concurrency.
// When using limiter as API rate limiting, implementation would be like below.
//
// Example usage:
//
//	type handler struct {
//		limiter zrate.Limiter
//	}
//
//	func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//		token := h.limiter.AllowNow() // Basically, use AllowNow for ConcurrentLimiter.
//		defer token.Release()         // Do not forget releasing the token.
//		if !token.OK() {
//			w.WriteHeader(http.StatusTooManyRequests)
//			return
//		}
//		// Some process.
//		w.WriteHeader(http.StatusOK)
//	}
type ConcurrentLimiter struct {
	// bucket holds tokens currently occupied.
	// cap(bucket) == limit.
	// len(bucket) == current number of concurrency.
	bucket chan struct{}
}

func (lim *ConcurrentLimiter) AllowNow() Token {
	select {
	case lim.bucket <- struct{}{}:
		return &token{ok: true, releaseFunc: func() { <-lim.bucket }}
	default:
		return TokenNG
	}
}

func (lim *ConcurrentLimiter) WaitNow(ctx context.Context) Token {
	select {
	case lim.bucket <- struct{}{}:
		return &token{ok: true, releaseFunc: func() { <-lim.bucket }}
	case <-ctx.Done():
		return &token{err: ctx.Err()}
	}
}
