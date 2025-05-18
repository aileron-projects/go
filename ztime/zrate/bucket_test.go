package zrate

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztesting"
)

func TestNewFixedWindowLimiterWidth(t *testing.T) {
	t.Parallel()
	t.Run("bucketSize=-1", func(t *testing.T) {
		lim := NewFixedWindowLimiter(-1)
		for range 5 { // Token always should be false.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
	t.Run("bucketSize=0", func(t *testing.T) {
		lim := NewFixedWindowLimiter(0)
		for range 5 { // Token always should be false.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
	t.Run("bucketSize=1", func(t *testing.T) {
		lim := NewFixedWindowLimiter(1)
		token1, token2 := lim.AllowNow(), lim.AllowNow()
		ztesting.AssertEqual(t, "incorrect token status", true, token1.OK())
		ztesting.AssertEqual(t, "incorrect token status", false, token2.OK())
	})
	t.Run("fillInterval=0", func(t *testing.T) {
		lim := NewFixedWindowLimiterWidth(1, 0)
		for range 5 { // Token always should be true.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
		}
	})
}

func TestNewTokenBucketInterval(t *testing.T) {
	t.Parallel()
	t.Run("bucketSize=-1", func(t *testing.T) {
		lim := NewTokenBucketLimiter(-1, 1)
		for range 5 { // Token always should be false.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
	t.Run("bucketSize=0", func(t *testing.T) {
		lim := NewTokenBucketLimiter(0, 1)
		for range 5 { // Token always should be false.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
	t.Run("bucketSize=1", func(t *testing.T) {
		lim := NewTokenBucketLimiter(1, 1)
		token1, token2 := lim.AllowNow(), lim.AllowNow()
		ztesting.AssertEqual(t, "incorrect token status", true, token1.OK())
		ztesting.AssertEqual(t, "incorrect token status", false, token2.OK())
	})
	t.Run("fillInterval=0", func(t *testing.T) {
		lim := NewTokenBucketInterval(1, 1, 0)
		for range 5 { // Token always should be true.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
		}
	})
}

func TestNewLimiter(t *testing.T) {
	t.Parallel()
	t.Run("NewFixedWindowLimiter", func(t *testing.T) {
		lim := &BucketLimiter{
			getToken: getTokenFunc(0, 10, time.Second, time.Now),
		}
		for range 5 { // Token always should be false.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
}

func TestBucketLimiter(t *testing.T) {
	t.Parallel()
	t.Run("AllowNow returns NG token", func(t *testing.T) {
		lim := &BucketLimiter{
			getToken: getTokenFunc(0, 10, time.Second, time.Now),
		}
		for range 5 { // Token always should be false.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
	t.Run("AllowNow returns OK token", func(t *testing.T) {
		lim := &BucketLimiter{
			getToken: getTokenFunc(10, 10, 0, time.Now),
		}
		for range 5 { // Token always should be true.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
		}
	})
	t.Run("WaitNow returns OK token", func(t *testing.T) {
		lim := &BucketLimiter{
			getToken: getTokenFunc(5, 5, time.Minute, time.Now),
		}
		for range 5 { // Token always should be true.
			token := lim.WaitNow(context.Background())
			ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
		}
		for range 5 { // Token always should be false.
			token := lim.AllowNow()
			ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		}
	})
	t.Run("retry after worked", func(t *testing.T) {
		nowTime := time.Now()
		now := &nowTime
		lim := &BucketLimiter{
			getToken: getTokenFunc(1, 1, time.Second, func() time.Time { return *now }),
		}
		token := lim.AllowNow() // Remove first token.
		ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
		token = lim.AllowNow() // Now token should be false.
		ztesting.AssertEqual(t, "incorrect token status", false, token.OK())
		time.AfterFunc(time.Second, func() { *now = (*now).Add(time.Second) }) // Forward time after 1sec.
		token = lim.WaitNow(context.Background())
		ztesting.AssertEqual(t, "incorrect token status", true, token.OK())
	})
	t.Run("context canceled", func(t *testing.T) {
		now := time.Now()
		lim := &BucketLimiter{
			getToken: getTokenFunc(1, 1, time.Second, func() time.Time { return now }),
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

func TestBucket(t *testing.T) {
	t.Parallel()
	t.Run("bucketSize=-1", func(t *testing.T) {
		getToken := getTokenFunc(-1, 10, time.Second, time.Now)
		ok, after := getToken()
		ztesting.AssertEqual(t, "token unexpectedly true", false, ok)
		ztesting.AssertEqual(t, "retry after time incorrect", math.MaxInt64, after)
	})
	t.Run("bucketSize=0", func(t *testing.T) {
		getToken := getTokenFunc(0, 10, time.Second, time.Now)
		ok, after := getToken()
		ztesting.AssertEqual(t, "token unexpectedly true", false, ok)
		ztesting.AssertEqual(t, "retry after time incorrect", math.MaxInt64, after)
	})
	t.Run("interval=-1", func(t *testing.T) {
		getToken := getTokenFunc(1, 10, -1, time.Now)
		ok, after := getToken()
		ztesting.AssertEqual(t, "token unexpectedly false", true, ok)
		ztesting.AssertEqual(t, "retry after time incorrect", 0, after)
	})
	t.Run("interval=0", func(t *testing.T) {
		getToken := getTokenFunc(1, 10, 0, time.Now)
		ok, after := getToken()
		ztesting.AssertEqual(t, "token unexpectedly false", true, ok)
		ztesting.AssertEqual(t, "retry after time incorrect", 0, after)
	})
	t.Run("fillRate=-1", func(t *testing.T) {
		now := time.Now()
		getToken := getTokenFunc(1, -1, time.Hour, func() time.Time { return now })
		ok1, after1 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly false", true, ok1)
		ztesting.AssertEqual(t, "retry after time incorrect", 0, after1)
		ok2, after2 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly true", false, ok2)
		ztesting.AssertEqual(t, "retry after time incorrect", time.Hour, after2)
	})
	t.Run("fillRate=0", func(t *testing.T) {
		now := time.Now()
		getToken := getTokenFunc(1, 0, time.Hour, func() time.Time { return now })
		ok1, after1 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly false", true, ok1)
		ztesting.AssertEqual(t, "retry after time incorrect", 0, after1)
		ok2, after2 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly true", false, ok2)
		ztesting.AssertEqual(t, "retry after time incorrect", time.Hour, after2)
	})
	t.Run("fillRate=1", func(t *testing.T) {
		nowTime := time.Now()
		now := &nowTime
		getToken := getTokenFunc(1, 1, time.Second, func() time.Time { return *now })
		ok1, after1 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly false", true, ok1)
		ztesting.AssertEqual(t, "retry after time incorrect", 0, after1)
		*now = (*now).Add(time.Second) // Forward 1 sec.
		ok2, after2 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly false", true, ok2)
		ztesting.AssertEqual(t, "retry after time incorrect", 0, after2)
	})
	t.Run("fillRate=1,try3times", func(t *testing.T) {
		nowTime := time.Now()
		now := &nowTime
		getToken := getTokenFunc(1, 1, time.Second, func() time.Time { return *now })
		ok1, after1 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly false", true, ok1)
		ztesting.AssertEqual(t, "retry after time incorrect", 0, after1)
		*now = (*now).Add(time.Second) // Forward 1 sec.
		ok2, after2 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly false", true, ok2)
		ztesting.AssertEqual(t, "retry after time incorrect", 0, after2)
		*now = (*now).Add(100 * time.Millisecond) // Forward 0.1 sec.
		ok3, after3 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly true", false, ok3)
		ztesting.AssertEqual(t, "retry after time incorrect", 900*time.Millisecond, after3)
	})
	t.Run("fillRate=2", func(t *testing.T) {
		nowTime := time.Now()
		now := &nowTime
		getToken := getTokenFunc(2, 2, time.Second, func() time.Time { return *now })
		_, _ = getToken() // Remove first token.
		_, _ = getToken() // Remove second token. Now the bucket is empty.
		ok1, after1 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly true", false, ok1)
		ztesting.AssertEqual(t, "retry after time incorrect", time.Second, after1)
		*now = (*now).Add(time.Second) // Forward 1 sec.
		ok2, after2 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly false", true, ok2)
		ztesting.AssertEqual(t, "retry after time incorrect", 0, after2)
		ok3, after3 := getToken()
		ztesting.AssertEqual(t, "token unexpectedly false", true, ok3)
		ztesting.AssertEqual(t, "retry after time incorrect", 0, after3)
	})
}
