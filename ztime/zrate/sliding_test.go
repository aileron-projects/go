package zrate

import (
	"context"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztesting"
)

func TestNewSlidingWindowLimiterWidth(t *testing.T) {
	t.Parallel()
	t.Run("limit=-1", func(t *testing.T) {
		lim := NewSlidingWindowLimiter(-1)
		for range 5 { // Token always should be false.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
	t.Run("limit=0", func(t *testing.T) {
		lim := NewSlidingWindowLimiter(0)
		for range 5 { // Token always should be false.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
	t.Run("limit=1", func(t *testing.T) {
		lim := NewSlidingWindowLimiter(1)
		token1, token2 := lim.AllowNow(), lim.AllowNow()
		ztesting.AssertEqual(t, "incorrect token status", true, token1.OK())
		ztesting.AssertEqual(t, "incorrect token status", false, token2.OK())
	})
	t.Run("width=0", func(t *testing.T) {
		lim := NewSlidingWindowLimiterWidth(1, 0)
		for range 5 { // Token always should be true.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
		}
	})
}

func TestSlidingWindowLimiter(t *testing.T) {
	t.Parallel()
	t.Run("reached limit", func(t *testing.T) {
		lim := &SlidingWindowLimiter{
			limit:      100,
			sum:        100,
			lastUpdate: time.Now(),
			subWidth:   time.Minute,
			timeNow:    time.Now,
		}
		for range 5 { // Token always should be false.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
	t.Run("AllowNow returns OK token", func(t *testing.T) {
		lim := &SlidingWindowLimiter{
			limit:    5,
			subWidth: time.Minute,
			timeNow:  time.Now,
		}
		for range 5 { // Token always should be true.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
		}
	})
	t.Run("WaitNow returns OK token", func(t *testing.T) {
		lim := &SlidingWindowLimiter{
			limit:      5,
			lastUpdate: time.Now(),
			subWidth:   time.Minute,
			timeNow:    time.Now,
		}
		for range 5 { // Token always should be true.
			println(lim.limit, lim.sum)
			token := lim.WaitNow(context.Background())
			ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
		}
		for range 5 { // Token always should be false.
			println(lim.limit, lim.sum)
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
	t.Run("retry after worked", func(t *testing.T) {
		nowTime := time.Now()
		now := &nowTime
		lim := &SlidingWindowLimiter{
			limit:      1,
			sum:        0,
			lastUpdate: nowTime,
			subWidth:   100 * time.Millisecond,
			timeNow:    func() time.Time { return *now },
		}
		token := lim.AllowNow() // Remove first token.
		ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
		token = lim.AllowNow() // Now token should be false.
		ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		time.AfterFunc(time.Second, func() { *now = (*now).Add(100 * 100 * time.Millisecond) }) // Forward time after 1sec.
		token = lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
	})
	t.Run("context canceled", func(t *testing.T) {
		now := time.Now()
		lim := &SlidingWindowLimiter{
			limit:      1,
			sum:        0,
			lastUpdate: time.Now(),
			subWidth:   100 * time.Millisecond,
			timeNow:    func() time.Time { return now },
		}
		token := lim.AllowNow() // Remove first token.
		ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
		defer cancel()
		token = lim.WaitNow(ctx)
		ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		ztesting.AssertEqual(t, "incorrect error", context.DeadlineExceeded, token.Err())
	})
}
