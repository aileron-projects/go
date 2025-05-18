package zrate

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztesting"
)

func TestLeakyBucketLimiter_AllowNow(t *testing.T) {
	t.Parallel()
	t.Run("queueSize=-1", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(-1, time.Second)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("queueSize=0", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(0, time.Second)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("interval=-1sec", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(1, -time.Second)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
	})
	t.Run("interval=0", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(1, 0)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
	})
	t.Run("queueSize=1", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(1, time.Second)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("queueSize=2", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(2, time.Second)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("queueSize=2", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(2, time.Second)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("interval passed", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(2, time.Second).(*LeakyBucketLimiter)
		tm := time.Now()
		now := &tm
		lim.timeNow = func() time.Time { return *now }
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		*now = now.Add(time.Second)
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
	})
	t.Run("someone waiting", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(5, time.Second)
		wg := &sync.WaitGroup{}
		for range 5 {
			wg.Add(1)
			go func() {
				wg.Done()
				_ = lim.WaitNow(context.Background())
			}()
		}
		wg.Wait() // All 5 goroutine started.
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", false, t1.OK())
	})
}

func TestLeakyBucketLimiter_WaitNow(t *testing.T) {
	t.Parallel()
	t.Run("queueSize=-1", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(-1, time.Second)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process allowed", false, t1.OK())
		t2 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("queueSize=0", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(0, time.Second)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process allowed", false, t1.OK())
		t2 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("interval=-1sec", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(1, -time.Second)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
	})
	t.Run("interval=0", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(1, 0)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
	})
	t.Run("queueSize=1, wait", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(1, 100*time.Millisecond)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
	})
	t.Run("queueSize=1, context deadline error", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(1, time.Second)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		dc, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Millisecond))
		defer cancel()
		t2 := lim.WaitNow(dc)
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
		ztesting.AssertEqualErr(t, "wrong error reason", context.DeadlineExceeded, t2.Err())
	})
	t.Run("interval passed", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(2, time.Second).(*LeakyBucketLimiter)
		tm := time.Now()
		now := &tm
		lim.timeNow = func() time.Time { return *now }
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		*now = now.Add(time.Second)
		t2 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
	})
	t.Run("fully queue", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(5, time.Second)
		wg := &sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				wg.Done()
				_ = lim.WaitNow(context.Background())
			}()
		}
		wg.Wait() // All 10 goroutine started.
		dc, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Millisecond))
		defer cancel()
		t1 := lim.WaitNow(dc)
		ztesting.AssertEqual(t, "process allowed", false, t1.OK())
	})
	t.Run("canceled context", func(t *testing.T) {
		lim := NewLeakyBucketLimiter(1, time.Second)
		dc, cancel := context.WithDeadline(context.Background(), time.Now())
		cancel()
		t1 := lim.WaitNow(dc)
		ztesting.AssertEqual(t, "process allowed", false, t1.OK())
		ztesting.AssertEqualErr(t, "wrong error reason", context.DeadlineExceeded, t1.Err())
	})
}
