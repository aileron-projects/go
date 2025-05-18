package zrate

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// NewLeakyBucketLimiter returns a new instance of
// a limiter that works with leaky bucket algorithm.
// The queueSize is the size of waiting queue.
// The second argument interval specifies the dequeue interval.
// For queueSize<=0, the limiter always returns token that indicates dis-allow.
// For interval<=0, the limiter always returns token that indicates allow.
func NewLeakyBucketLimiter(queueSize int, interval time.Duration) Limiter {
	if queueSize <= 0 {
		return NoopLimiter(false)
	}
	if interval <= 0 {
		return NoopLimiter(true)
	}
	return &LeakyBucketLimiter{
		interval: interval,
		timeNow:  time.Now,
		queue:    make(chan struct{}, queueSize),
		notifier: make(chan struct{}, 1),
	}
}

// LeakyBucketLimiter limits the rate of something
// to process using leaky bucket algorithm.
// When using limiter as API rate limiting, implementation would be like below.
//
// Example usage:
//
//	type handler struct {
//		limiter zrate.Limiter
//	}
//
//	func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//		token := h.limiter.WaitNow(r.Context()) // Basically, use WaitNow for LeakyBucketLimiter.
//		defer token.Release()                   // Release is not required for LeakyBucket algorithm.
//		if !token.OK() {
//			w.WriteHeader(http.StatusTooManyRequests)
//			return
//		}
//		// Some process.
//		w.WriteHeader(http.StatusOK)
//	}
type LeakyBucketLimiter struct {
	// mu protects lastLeak.
	mu sync.Mutex
	// lastLeak is the last dequeue time.
	lastLeak time.Time
	// interval is the leak, or dequeue interval.
	// interval can be zero or positive.
	interval time.Duration

	// queue is the wait queue.
	// The input struct{}{} is notified through notifier.
	// Waiters wait until they accept a struct{}{}
	// through notifier.
	// cap(queue) == actual queue size.
	// len(queue) == number of waiters.
	queue chan struct{}
	// queue is the channel for notifying the waiters.
	// Its capacity is 1.
	notifier chan struct{}
	// dequeueWorking is the flag if the dequeue worker
	// is working in a goroutine or not.
	dequeueWorking atomic.Bool

	// timeNow returns the current time.
	// This can be replaced for testing.
	timeNow func() time.Time
}

func (lim *LeakyBucketLimiter) AllowNow() Token {
	if len(lim.queue) > 0 {
		return TokenNG // Someone already waiting.
	}
	lim.mu.Lock()
	defer lim.mu.Unlock()
	now := lim.timeNow()
	if now.Sub(lim.lastLeak) >= lim.interval {
		lim.lastLeak = now
		return TokenOK
	}
	return TokenNG
}

func (lim *LeakyBucketLimiter) WaitNow(ctx context.Context) Token {
	select {
	case lim.queue <- struct{}{}:
	default:
		return TokenNG
	}
	lim.notifyWorker()
	select {
	case <-lim.notifier:
		return TokenOK
	case <-ctx.Done():
		select {
		case <-lim.queue: // Discard, or cancel this request.
		default: // It might already been dequeued by the worker.
		}
		return &token{err: ctx.Err()}
	}
}

func (lim *LeakyBucketLimiter) notifyWorker() {
	if lim.dequeueWorking.Swap(true) {
		return
	}
	go func() {
		defer lim.dequeueWorking.Store(false)
		timer := time.NewTimer(time.Second)
		defer timer.Stop()
		for {
			var s struct{}
			select {
			case s = <-lim.queue:
			case <-time.After(10 * time.Second): // Keep at least 10 seconds. No reason to the value.
				return
			}
			wait := lim.interval - lim.timeNow().Sub(lim.lastLeak)
			if wait > 0 {
				timer.Reset(wait)
				<-timer.C
			}
			select {
			case lim.notifier <- s: // Notify to the next waiter.
				lim.mu.Lock()
				lim.lastLeak = lim.timeNow()
				lim.mu.Unlock()
			default: // Waiter might canceled.
				// This case, do not update lastLeak time.
			}
		}
	}()
}
