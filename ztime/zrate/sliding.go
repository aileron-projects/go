package zrate

import (
	"context"
	"sync"
	"time"
)

// subWindowSize is the number of sub windows
// for sliding windows limiter.
const subWindowSize = 100

// NewSlidingWindowLimiter returns a new instance of
// a limiter that works with sliding-window algorithm.
// NewSlidingWindowLimiter is short for NewSlidingWindowLimiterWidth(limit, time.Second).
func NewSlidingWindowLimiter(limit int) Limiter {
	return NewSlidingWindowLimiterWidth(limit, time.Second)
}

// NewSlidingWindowLimiterWidth returns a new instance of
// a limiter that works with sliding-window algorithm.
// The limit is the maximum count allowed within the width.
// The window width is split into sub windows.
// Currently the number of sub window is 100.
// For limit<=0, the limiter always returns token that indicates dis-allow.
// For width<100, the limiter always returns token that indicates allow.
func NewSlidingWindowLimiterWidth(limit int, width time.Duration) Limiter {
	if limit <= 0 {
		return NoopLimiter(false)
	}
	if width/subWindowSize <= 0 {
		return NoopLimiter(true)
	}
	return &SlidingWindowLimiter{
		limit:      int64(limit),
		lastUpdate: time.Now(),
		subWidth:   width / subWindowSize,
		timeNow:    time.Now,
	}
}

// SlidingWindowLimiter limits the rate of something
// to process using sliding-window algorithm.
// When using limiter as API rate limiting, implementation would be like below.
//
// Example usage:
//
//	type handler struct {
//		limiter zrate.Limiter
//	}
//
//	func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//		token := h.limiter.WaitNow(r.Context()) // Basically, use WaitNow for SlidingWindowLimiter.
//		defer token.Release()                   // Release is not required for SlidingWindow algorithm.
//		if !token.OK() {
//			w.WriteHeader(http.StatusTooManyRequests)
//			return
//		}
//		// Some process.
//		w.WriteHeader(http.StatusOK)
//	}
type SlidingWindowLimiter struct {
	mu sync.Mutex
	// limit is the limit within the window.
	limit int64
	// sum is the total number of allowed tokens
	// within the windows.
	// sum is always sum<=limit.
	sum int64

	// lastUpdate is the last update time
	// of sub windows.
	lastUpdate time.Time
	// subWidth is the sub windows width.
	// subWidth is equal to width/subWindowSize.
	subWidth time.Duration
	// subWindow is the sub windows.
	// Each window holds the count of issued tokens.
	subWindow [subWindowSize]int64
	// index is the current position in subWindow array.
	index int

	// timeNow returns the current time.
	// This can be replaced for testing.
	timeNow func() time.Time
}

func (lim *SlidingWindowLimiter) updateSubWindow() {
	skip := int64(lim.timeNow().Sub(lim.lastUpdate) / lim.subWidth)
	if skip == 0 {
		return
	}
	lim.lastUpdate = lim.lastUpdate.Add(lim.subWidth * time.Duration(skip))
	for range min(skip, subWindowSize) {
		lim.index++
		if lim.index >= subWindowSize {
			lim.index = 0
		}
		lim.sum -= lim.subWindow[lim.index]
		lim.subWindow[lim.index] = 0
	}
}

func (lim *SlidingWindowLimiter) AllowNow() Token {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	lim.updateSubWindow()
	if lim.sum >= lim.limit {
		return TokenNG
	}
	lim.sum += 1
	lim.subWindow[lim.index] += 1
	return TokenOK
}

func (lim *SlidingWindowLimiter) WaitNow(ctx context.Context) Token {
	for {
		t := lim.AllowNow()
		if t.OK() {
			return TokenOK
		}
		select {
		case <-time.After(lim.subWidth):
		case <-ctx.Done():
			return &token{err: ctx.Err()}
		}
	}
}
