package zbackoff

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztesting"
)

func TestBackoffError(t *testing.T) {
	t.Parallel()
	t.Run("equality", func(t *testing.T) {
		e1 := &BackoffError{Type: "foo", Info: "bar"}
		e2 := &BackoffError{Type: "foo", Info: "bar"}
		e3 := &BackoffError{Type: "bar", Info: "foo"}
		ztesting.AssertEqual(t, "error is not the same", true, errors.Is(e1, e2))
		ztesting.AssertEqual(t, "error is the same", false, errors.Is(e1, e3))
		ztesting.AssertEqual(t, "error is the same", false, errors.Is(e1, nil))
	})
	t.Run("wrapped error equality", func(t *testing.T) {
		e1 := &BackoffError{Type: "foo", Info: "bar"}
		e2 := fmt.Errorf("this is e2 [%w]", &BackoffError{Type: "foo", Info: "bar"})
		e3 := fmt.Errorf("this is e3 [%w]", &BackoffError{Type: "bar", Info: "foo"})
		e4 := fmt.Errorf("this is e4 [%w]", nil)
		ztesting.AssertEqual(t, "error is not the same", true, errors.Is(e1, e2))
		ztesting.AssertEqual(t, "error is the same", false, errors.Is(e1, e3))
		ztesting.AssertEqual(t, "error is the same", false, errors.Is(e1, e4))
	})
}

func TestFixedBackoff(t *testing.T) {
	t.Parallel()
	t.Run("panic when value<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "value must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewFixedBackoff(-1)
	})

	testCases := map[string]struct {
		value   time.Duration
		attempt int
		want    time.Duration
	}{
		"value=0_case01":      {0, math.MinInt, 0},
		"value=0_case02":      {0, -1, 0},
		"value=0_case03":      {0, 0, 0},
		"value=0_case04":      {0, 1, 0},
		"value=0_case05":      {0, math.MaxInt, 0},
		"value=1_case01":      {1, math.MinInt, 1},
		"value=1_case02":      {1, -1, 1},
		"value=1_case03":      {1, 0, 1},
		"value=1_case04":      {1, 1, 1},
		"value=1_case05":      {1, math.MaxInt, 1},
		"value=MaxInt_case01": {math.MaxInt, math.MinInt, math.MaxInt},
		"value=MaxInt_case02": {math.MaxInt, -1, math.MaxInt},
		"value=MaxInt_case03": {math.MaxInt, 0, math.MaxInt},
		"value=MaxInt_case04": {math.MaxInt, 1, math.MaxInt},
		"value=MaxInt_case05": {math.MaxInt, math.MaxInt, math.MaxInt},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			backoff := NewFixedBackoff(tc.value)
			got := backoff.Attempt(tc.attempt)
			ztesting.AssertEqual(t, "wrong backoff duration for attempt="+strconv.Itoa(tc.attempt), tc.want, got)
		})
	}
}

func TestRandomBackoff(t *testing.T) {
	t.Parallel()
	t.Run("panic when offset<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewRandomBackoff(-1, 1)
	})
	t.Run("panic when fluctuation==0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "fluctuation must be positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewRandomBackoff(1, 0)
	})
	t.Run("panic when fluctuation<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "fluctuation must be positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewRandomBackoff(1, 0)
	})
	t.Run("panic when offset+fluctuation>MaxInt64", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset+fluctuation must be under MaxInt64 (9,223,372,036,854,775,807)"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewRandomBackoff(math.MaxInt64, 1)
	})

	testCases := map[string]struct {
		offset, fluctuation time.Duration
		attempt             int
	}{
		"offset=0,flc=10_case01":     {0, 10, math.MinInt},
		"offset=0,flc=10_case02":     {0, 10, -1},
		"offset=0,flc=10_case03":     {0, 10, 0},
		"offset=0,flc=10_case04":     {0, 10, 1},
		"offset=0,flc=10_case05":     {0, 10, math.MaxInt},
		"offset=1,flc=10_case01":     {1, 10, math.MinInt},
		"offset=1,flc=10_case02":     {1, 10, -1},
		"offset=1,flc=10_case03":     {1, 10, 0},
		"offset=1,flc=10_case04":     {1, 10, 1},
		"offset=1,flc=10_case05":     {1, 10, math.MaxInt},
		"offset=10000,flc=10_case01": {10000, 10, math.MinInt},
		"offset=10000,flc=10_case02": {10000, 10, -1},
		"offset=10000,flc=10_case03": {10000, 10, 0},
		"offset=10000,flc=10_case04": {10000, 10, 1},
		"offset=10000,flc=10_case05": {10000, 10, math.MaxInt},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			backoff := NewRandomBackoff(tc.offset, tc.fluctuation)
			got := backoff.Attempt(tc.attempt)
			inRange := (got - tc.offset) <= tc.fluctuation
			ztesting.AssertEqual(t, "wrong backoff duration range for attempt="+strconv.Itoa(tc.attempt), true, inRange)
		})
	}
}

func TestLinearBackoff(t *testing.T) {
	t.Parallel()
	t.Run("panic when offset<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewLinearBackoff(-1, 1, 1)
	})
	t.Run("panic when fluctuation<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "fluctuation must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewLinearBackoff(1, -1, 1)
	})
	t.Run("panic when coeff<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "coeff must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewLinearBackoff(1, 1, -1)
	})
	t.Run("panic when offset+fluctuation>MaxInt64", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset+fluctuation must be under MaxInt64 (9,223,372,036,854,775,807)"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewLinearBackoff(math.MaxInt64, 1, 1)
	})

	testCases := map[string]struct {
		offset, fluctuation, coeff time.Duration
		attempt                    int
		want                       time.Duration
	}{
		"offset=0,flc=10,coeff=1_01":  {0, 10, 1, -1, -1},
		"offset=0,flc=10,coeff=1_02":  {0, 10, 1, 1, 1},
		"offset=0,flc=10,coeff=1_03":  {0, 10, 1, 0, 0},
		"offset=0,flc=10,coeff=1_04":  {0, 10, 1, 1, 1},
		"offset=0,flc=10,coeff=1_05":  {0, 10, 1, 20, 10},
		"offset=0,flc=10,coeff=2_01":  {0, 10, 2, -1, -2},
		"offset=0,flc=10,coeff=2_02":  {0, 10, 2, 1, 2},
		"offset=0,flc=10,coeff=2_03":  {0, 10, 2, 0, 0},
		"offset=0,flc=10,coeff=2_04":  {0, 10, 2, 1, 2},
		"offset=0,flc=10,coeff=2_05":  {0, 10, 2, 20, 10},
		"offset=10,flc=10,coeff=1_01": {10, 10, 1, -1, 9},
		"offset=10,flc=10,coeff=1_02": {10, 10, 1, 1, 11},
		"offset=10,flc=10,coeff=1_03": {10, 10, 1, 0, 10},
		"offset=10,flc=10,coeff=1_04": {10, 10, 1, 1, 11},
		"offset=10,flc=10,coeff=1_05": {10, 10, 1, 20, 20},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			backoff := NewLinearBackoff(tc.offset, tc.fluctuation, tc.coeff)
			got := backoff.Attempt(tc.attempt)
			ztesting.AssertEqual(t, "wrong backoff duration for attempt="+strconv.Itoa(tc.attempt), tc.want, got)
		})
	}
}

func TestPolynomialBackoff(t *testing.T) {
	t.Parallel()
	t.Run("panic when offset<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewPolynomialBackoff(-1, 1, 1, 1)
	})
	t.Run("panic when fluctuation<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "fluctuation must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewPolynomialBackoff(1, -1, 1, 1)
	})
	t.Run("panic when coeff<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "coeff must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewPolynomialBackoff(1, 1, -1, 1)
	})
	t.Run("panic when offset+fluctuation>MaxInt64", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset+fluctuation must be under MaxInt64 (9,223,372,036,854,775,807)"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewPolynomialBackoff(math.MaxInt64, 1, 1, 1)
	})

	testCases := map[string]struct {
		offset, fluctuation, coeff time.Duration
		exponent                   float64
		attempt                    int
		want                       time.Duration
	}{
		"offset=0,flc=10,coeff=1,exp=1_case01":  {0, 10, 1, 1, -10, 1}, // attempt<=0 modified to 1
		"offset=0,flc=10,coeff=1,exp=1_case02":  {0, 10, 1, 1, -1, 1},  // attempt<=0 modified to 1
		"offset=0,flc=10,coeff=1,exp=1_case03":  {0, 10, 1, 1, 0, 1},   // attempt<=0 modified to 1
		"offset=0,flc=10,coeff=1,exp=1_case04":  {0, 10, 1, 1, 1, 1},
		"offset=0,flc=10,coeff=1,exp=1_case05":  {0, 10, 1, 1, 2, 2},
		"offset=0,flc=10,coeff=1,exp=1_case06":  {0, 10, 1, 1, 20, 10},
		"offset=0,flc=10,coeff=2,exp=1_case01":  {0, 10, 2, 1, -10, 2},
		"offset=0,flc=10,coeff=2,exp=1_case02":  {0, 10, 2, 1, -1, 2},
		"offset=0,flc=10,coeff=2,exp=1_case03":  {0, 10, 2, 1, 0, 2},
		"offset=0,flc=10,coeff=2,exp=1_case04":  {0, 10, 2, 1, 1, 2},
		"offset=0,flc=10,coeff=2,exp=1_case05":  {0, 10, 2, 1, 2, 4},
		"offset=0,flc=10,coeff=2,exp=1_case06":  {0, 10, 2, 1, 20, 10},
		"offset=10,flc=10,coeff=1,exp=1_case01": {10, 10, 1, 1, -10, 11},
		"offset=10,flc=10,coeff=1,exp=1_case02": {10, 10, 1, 1, -1, 11},
		"offset=10,flc=10,coeff=1,exp=1_case03": {10, 10, 1, 1, 0, 11},
		"offset=10,flc=10,coeff=1,exp=1_case04": {10, 10, 1, 1, 1, 11},
		"offset=10,flc=10,coeff=1,exp=1_case05": {10, 10, 1, 1, 2, 12},
		"offset=10,flc=10,coeff=1,exp=1_case06": {10, 10, 1, 1, 20, 20},
		"offset=10,flc=10,coeff=1,exp=2_case01": {10, 10, 1, 2, -10, 11},
		"offset=10,flc=10,coeff=1,exp=2_case02": {10, 10, 1, 2, -1, 11},
		"offset=10,flc=10,coeff=1,exp=2_case03": {10, 10, 1, 2, 0, 11},
		"offset=10,flc=10,coeff=1,exp=2_case04": {10, 10, 1, 2, 1, 11},
		"offset=10,flc=10,coeff=1,exp=2_case05": {10, 10, 1, 2, 2, 14},
		"offset=10,flc=10,coeff=1,exp=2_case06": {10, 10, 1, 2, 20, 20},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			backoff := NewPolynomialBackoff(tc.offset, tc.fluctuation, tc.coeff, tc.exponent)
			got := backoff.Attempt(tc.attempt)
			ztesting.AssertEqual(t, "wrong backoff duration for attempt="+strconv.Itoa(tc.attempt), tc.want, got)
		})
	}
}

func TestExponentialBackoff(t *testing.T) {
	t.Parallel()
	t.Run("panic when offset<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoff(-1, 1, 1)
	})
	t.Run("panic when fluctuation<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "fluctuation must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoff(1, -1, 1)
	})
	t.Run("panic when coeff<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "coeff must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoff(1, 1, -1)
	})
	t.Run("panic when offset+fluctuation>MaxInt64", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset+fluctuation must be under MaxInt64 (9,223,372,036,854,775,807)"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoff(math.MaxInt64, 1, 1)
	})

	testCases := map[string]struct {
		offset, fluctuation, coeff time.Duration
		attempt                    int
		want                       time.Duration
	}{
		"offset=0,flc=100,coeff=100_case01":  {0, 100, 100, -2, 25},
		"offset=0,flc=100,coeff=100_case02":  {0, 100, 100, -1, 50},
		"offset=0,flc=100,coeff=100_case03":  {0, 100, 100, 0, 100},
		"offset=0,flc=100,coeff=100_case04":  {0, 100, 100, 1, 100},
		"offset=0,flc=100,coeff=100_case05":  {0, 100, 100, 2, 100},  // fluctuation exceeded.
		"offset=0,flc=100,coeff=100_case06":  {0, 100, 100, 20, 100}, // fluctuation exceeded.
		"offset=0,flc=100,coeff=200_case01":  {0, 100, 200, -2, 50},
		"offset=0,flc=100,coeff=200_case02":  {0, 100, 200, -1, 100},
		"offset=0,flc=100,coeff=200_case03":  {0, 100, 200, 0, 100},  // fluctuation exceeded
		"offset=0,flc=100,coeff=200_case04":  {0, 100, 200, 1, 100},  // fluctuation exceeded
		"offset=0,flc=100,coeff=200_case05":  {0, 100, 200, 2, 100},  // fluctuation exceeded.
		"offset=0,flc=100,coeff=200_case06":  {0, 100, 200, 20, 100}, // fluctuation exceeded.
		"offset=10,flc=100,coeff=100_case01": {10, 100, 100, -2, 35},
		"offset=10,flc=100,coeff=100_case02": {10, 100, 100, -1, 60},
		"offset=10,flc=100,coeff=100_case03": {10, 100, 100, 0, 110},
		"offset=10,flc=100,coeff=100_case04": {10, 100, 100, 1, 110},  // fluctuation exceeded.
		"offset=10,flc=100,coeff=100_case05": {10, 100, 100, 2, 110},  // fluctuation exceeded.
		"offset=10,flc=100,coeff=100_case06": {10, 100, 100, 20, 110}, // fluctuation exceeded.
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			backoff := NewExponentialBackoff(tc.offset, tc.fluctuation, tc.coeff)
			got := backoff.Attempt(tc.attempt)
			ztesting.AssertEqual(t, "wrong backoff duration for attempt="+strconv.Itoa(tc.attempt), tc.want, got)
		})
	}
}

func TestExponentialBackoffFullJitter(t *testing.T) {
	t.Parallel()
	t.Run("panic when offset<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoffFullJitter(-1, 1, 1)
	})
	t.Run("panic when fluctuation<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "fluctuation must be positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoffFullJitter(1, -1, 1)
	})
	t.Run("panic when coeff<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "coeff must be positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoffFullJitter(1, 1, -1)
	})
	t.Run("panic when offset+fluctuation>MaxInt64", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset+fluctuation must be under MaxInt64 (9,223,372,036,854,775,807)"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoffFullJitter(math.MaxInt64, 1, 1)
	})

	testCases := map[string]struct {
		offset, fluctuation, coeff time.Duration
		attempt                    int
	}{
		"offset=0,flc=10,coeff=1_case01":  {0, 10, 1, -2},
		"offset=0,flc=10,coeff=1_case02":  {0, 10, 1, -1},
		"offset=0,flc=10,coeff=1_case03":  {0, 10, 1, 0},
		"offset=0,flc=10,coeff=1_case04":  {0, 10, 1, 1},
		"offset=0,flc=10,coeff=1_case05":  {0, 10, 1, 2},
		"offset=0,flc=10,coeff=1_case06":  {0, 10, 1, 10},
		"offset=0,flc=10,coeff=2_case01":  {0, 10, 2, -2},
		"offset=0,flc=10,coeff=2_case02":  {0, 10, 2, -1},
		"offset=0,flc=10,coeff=2_case03":  {0, 10, 2, 0},
		"offset=0,flc=10,coeff=2_case04":  {0, 10, 2, 1},
		"offset=0,flc=10,coeff=2_case05":  {0, 10, 2, 2},
		"offset=0,flc=10,coeff=2_case06":  {0, 10, 2, 10},
		"offset=10,flc=10,coeff=1_case01": {10, 10, 1, -2},
		"offset=10,flc=10,coeff=1_case02": {10, 10, 1, -1},
		"offset=10,flc=10,coeff=1_case03": {10, 10, 1, 0},
		"offset=10,flc=10,coeff=1_case04": {10, 10, 1, 1},
		"offset=10,flc=10,coeff=1_case05": {10, 10, 1, 2},
		"offset=10,flc=10,coeff=1_case06": {10, 10, 1, 10},
		"offset=100,flc=1,coeff=1_case01": {100, 1, 1, -2},
		"offset=100,flc=1,coeff=1_case02": {100, 1, 1, -1},
		"offset=100,flc=1,coeff=1_case03": {100, 1, 1, 0},
		"offset=100,flc=1,coeff=1_case04": {100, 1, 1, 1},
		"offset=100,flc=1,coeff=1_case05": {100, 1, 1, 2},
		"offset=100,flc=1,coeff=1_case06": {100, 1, 1, 10},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			backoff := NewExponentialBackoffFullJitter(tc.offset, tc.fluctuation, tc.coeff)
			got := backoff.Attempt(tc.attempt)
			inRange := (got - tc.offset) <= tc.fluctuation
			ztesting.AssertEqual(t, "wrong backoff duration range for attempt="+strconv.Itoa(tc.attempt), true, inRange)
		})
	}
}

func TestExponentialBackoffEqualJitter(t *testing.T) {
	t.Parallel()
	t.Run("panic when offset<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoffEqualJitter(-1, 1, 1)
	})
	t.Run("panic when fluctuation<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "fluctuation must be positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoffEqualJitter(1, -1, 1)
	})
	t.Run("panic when coeff<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "coeff must be positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoffEqualJitter(1, 1, -1)
	})
	t.Run("panic when offset+fluctuation>MaxInt64", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset+fluctuation must be under MaxInt64 (9,223,372,036,854,775,807)"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewExponentialBackoffEqualJitter(math.MaxInt64, 1, 1)
	})

	testCases := map[string]struct {
		offset, fluctuation, coeff time.Duration
		attempt                    int
	}{
		"offset=0,flc=10,coeff=1_case01":  {0, 10, 1, -2},
		"offset=0,flc=10,coeff=1_case02":  {0, 10, 1, -1},
		"offset=0,flc=10,coeff=1_case03":  {0, 10, 1, 0},
		"offset=0,flc=10,coeff=1_case04":  {0, 10, 1, 1},
		"offset=0,flc=10,coeff=1_case05":  {0, 10, 1, 2},
		"offset=0,flc=10,coeff=1_case06":  {0, 10, 1, 10},
		"offset=0,flc=10,coeff=2_case01":  {0, 10, 2, -2},
		"offset=0,flc=10,coeff=2_case02":  {0, 10, 2, -1},
		"offset=0,flc=10,coeff=2_case03":  {0, 10, 2, 0},
		"offset=0,flc=10,coeff=2_case04":  {0, 10, 2, 1},
		"offset=0,flc=10,coeff=2_case05":  {0, 10, 2, 2},
		"offset=0,flc=10,coeff=2_case06":  {0, 10, 2, 10},
		"offset=10,flc=10,coeff=1_case01": {10, 10, 1, -2},
		"offset=10,flc=10,coeff=1_case02": {10, 10, 1, -1},
		"offset=10,flc=10,coeff=1_case03": {10, 10, 1, 0},
		"offset=10,flc=10,coeff=1_case04": {10, 10, 1, 1},
		"offset=10,flc=10,coeff=1_case05": {10, 10, 1, 2},
		"offset=10,flc=10,coeff=1_case06": {10, 10, 1, 10},
		"offset=100,flc=1,coeff=1_case01": {100, 1, 1, -2},
		"offset=100,flc=1,coeff=1_case02": {100, 1, 1, -1},
		"offset=100,flc=1,coeff=1_case03": {100, 1, 1, 0},
		"offset=100,flc=1,coeff=1_case04": {100, 1, 1, 1},
		"offset=100,flc=1,coeff=1_case05": {100, 1, 1, 2},
		"offset=100,flc=1,coeff=1_case06": {100, 1, 1, 10},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			backoff := NewExponentialBackoffEqualJitter(tc.offset, tc.fluctuation, tc.coeff)
			got := backoff.Attempt(tc.attempt)
			inRange := (got - tc.offset - tc.fluctuation/2) <= tc.fluctuation
			ztesting.AssertEqual(t, "wrong backoff duration range for attempt="+strconv.Itoa(tc.attempt), true, inRange)
		})
	}
}

func TestFibonacciBackoff(t *testing.T) {
	t.Parallel()
	t.Run("panic when offset<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewFibonacciBackoff(-1, 1, 1)
	})
	t.Run("panic when fluctuation<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "fluctuation must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewFibonacciBackoff(1, -1, 1)
	})
	t.Run("panic when coeff<0", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "coeff must be zero or positive"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewFibonacciBackoff(1, 1, -1)
	})
	t.Run("panic when offset+fluctuation>MaxInt64", func(t *testing.T) {
		defer func() {
			got, _ := recover().(error)
			want := &BackoffError{Type: errPram, Info: "offset+fluctuation must be under MaxInt64 (9,223,372,036,854,775,807)"}
			ztesting.AssertEqualErr(t, "error not matched", want, got)
		}()
		NewFibonacciBackoff(math.MaxInt64, 1, 1)
	})

	testCases := map[string]struct {
		offset, fluctuation, coeff time.Duration
		attempt                    int
		want                       time.Duration
	}{
		"offset=0,flc=100,coeff=10_case01":  {0, 100, 10, -2, 0},
		"offset=0,flc=100,coeff=10_case02":  {0, 100, 10, -1, 0},
		"offset=0,flc=100,coeff=10_case03":  {0, 100, 10, 0, 0},
		"offset=0,flc=100,coeff=10_case04":  {0, 100, 10, 1, 10},
		"offset=0,flc=100,coeff=10_case05":  {0, 100, 10, 2, 10},
		"offset=0,flc=100,coeff=10_case06":  {0, 100, 10, 3, 20},
		"offset=0,flc=100,coeff=10_case07":  {0, 100, 10, 4, 30},
		"offset=0,flc=100,coeff=10_case08":  {0, 100, 10, 5, 50},
		"offset=0,flc=100,coeff=10_case09":  {0, 100, 10, 6, 80},
		"offset=0,flc=100,coeff=10_case10":  {0, 100, 10, 7, 100},           // fluctuation exceeded.
		"offset=0,flc=100,coeff=10_case11":  {0, 100, 10, 46, 100},          // fluctuation exceeded.
		"offset=0,flc=100,coeff=10_case12":  {0, 100, 10, math.MaxInt, 100}, // fluctuation exceeded.
		"offset=10,flc=100,coeff=10_case01": {10, 100, 10, -2, 10},
		"offset=10,flc=100,coeff=10_case02": {10, 100, 10, -1, 10},
		"offset=10,flc=100,coeff=10_case03": {10, 100, 10, 0, 10},
		"offset=10,flc=100,coeff=10_case04": {10, 100, 10, 1, 20},
		"offset=10,flc=100,coeff=10_case05": {10, 100, 10, 2, 20},
		"offset=10,flc=100,coeff=10_case06": {10, 100, 10, 3, 30},
		"offset=10,flc=100,coeff=10_case07": {10, 100, 10, 4, 40},
		"offset=10,flc=100,coeff=10_case08": {10, 100, 10, 5, 60},
		"offset=10,flc=100,coeff=10_case09": {10, 100, 10, 6, 90},
		"offset=10,flc=100,coeff=10_case10": {10, 100, 10, 7, 110},           // fluctuation exceeded.
		"offset=10,flc=100,coeff=10_case11": {10, 100, 10, 46, 110},          // fluctuation exceeded.
		"offset=10,flc=100,coeff=10_case12": {10, 100, 10, math.MaxInt, 110}, // fluctuation exceeded.
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			backoff := NewFibonacciBackoff(tc.offset, tc.fluctuation, tc.coeff)
			got := backoff.Attempt(tc.attempt)
			ztesting.AssertEqual(t, "wrong backoff duration for attempt="+strconv.Itoa(tc.attempt), tc.want, got)
		})
	}
}

func TestFibonacci(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input int
		want  int
	}{
		"-10": {-10, 0},
		"-1":  {-1, 0},
		"0":   {0, 0},
		"1":   {1, 1},
		"2":   {2, 1},
		"3":   {3, 2},
		"4":   {4, 3},
		"5":   {5, 5},
		"10":  {10, 55},
		"20":  {20, 6765},
		"30":  {30, 832040},
		"40":  {40, 102334155},
		"46":  {46, 1836311903}, // Int32 limit.
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := fibonacci(tc.input)
			ztesting.AssertEqual(t, "wrong fibonacci for "+strconv.Itoa(tc.input), tc.want, got)
		})
	}
}
