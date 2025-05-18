package zrate

import (
	"context"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztesting"
)

func TestConcurrentLimiter_AllowNow(t *testing.T) {
	t.Parallel()
	t.Run("limit=-1", func(t *testing.T) {
		lim := NewConcurrentLimiter(-1)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("limit=0", func(t *testing.T) {
		lim := NewConcurrentLimiter(0)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("limit=1", func(t *testing.T) {
		lim := NewConcurrentLimiter(1)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("limit=2", func(t *testing.T) {
		lim := NewConcurrentLimiter(2)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
		t3 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t3.OK())
	})
	t.Run("limit=2 and release", func(t *testing.T) {
		lim := NewConcurrentLimiter(2)
		t1 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
		t3 := lim.AllowNow()
		ztesting.AssertEqual(t, "process allowed", false, t3.OK())
		t1.Release()
		t4 := lim.AllowNow()
		ztesting.AssertEqual(t, "process not allowed", true, t4.OK())
	})
}

func TestConcurrentLimiter_WaitNow(t *testing.T) {
	t.Parallel()
	t.Run("limit=-1", func(t *testing.T) {
		lim := NewConcurrentLimiter(-1)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process allowed", false, t1.OK())
		t2 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("limit=0", func(t *testing.T) {
		lim := NewConcurrentLimiter(0)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process allowed", false, t1.OK())
		t2 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
	})
	t.Run("limit=1", func(t *testing.T) {
		lim := NewConcurrentLimiter(1)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		dc, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Millisecond))
		defer cancel()
		t2 := lim.WaitNow(dc)
		ztesting.AssertEqual(t, "process allowed", false, t2.OK())
		ztesting.AssertEqualErr(t, "wrong error reason", context.DeadlineExceeded, t2.Err())
	})
	t.Run("limit=2", func(t *testing.T) {
		lim := NewConcurrentLimiter(2)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
		dc, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Millisecond))
		defer cancel()
		t3 := lim.WaitNow(dc)
		ztesting.AssertEqual(t, "process allowed", false, t3.OK())
		ztesting.AssertEqualErr(t, "wrong error reason", context.DeadlineExceeded, t3.Err())
	})
	t.Run("limit=2 and release", func(t *testing.T) {
		lim := NewConcurrentLimiter(2)
		t1 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t1.OK())
		t2 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t2.OK())
		t1.Release()
		t3 := lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "process not allowed", true, t3.OK())
	})
}
