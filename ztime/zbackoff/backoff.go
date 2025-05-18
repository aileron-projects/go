package zbackoff

import (
	"errors"
	"math"
	"math/rand/v2"
	"time"
)

var (
	_ Backoff = &FixedBackoff{}
	_ Backoff = &RandomBackoff{}
	_ Backoff = &LinearBackoff{}
	_ Backoff = &PolynomialBackoff{}
	_ Backoff = &ExponentialBackoff{}
	_ Backoff = &ExponentialBackoffFullJitter{}
	_ Backoff = &ExponentialBackoffEqualJitter{}
	_ Backoff = &FibonacciBackoff{}
)

var (
	errPram = "invalid parameter."
)

func mustPositive(v float64, name string) {
	if v > 0 {
		return
	}
	panic(&BackoffError{Type: errPram, Info: name + " must be positive"})
}

func mustZeroOrPositive(v float64, name string) {
	if v >= 0 {
		return
	}
	panic(&BackoffError{Type: errPram, Info: name + " must be zero or positive"})
}

func mustUnderMaxInt64(v1, v2 time.Duration, name string) {
	if v1 < math.MaxInt64-v2 {
		return
	}
	panic(&BackoffError{Type: errPram, Info: name + " must be under MaxInt64 (9,223,372,036,854,775,807)"})
}

// BackoffError reports errors in backoff.
type BackoffError struct {
	Type string
	Info string
}

func (e *BackoffError) Error() string {
	return "zbackoff: " + e.Type + " " + e.Info
}

func (e *BackoffError) Is(target error) bool {
	for target != nil {
		ee, ok := target.(*BackoffError)
		if ok {
			return e.Type == ee.Type && e.Info == ee.Info
		}
		target = errors.Unwrap(target)
	}
	return false
}

// Backoff provides backoff algorithm.
type Backoff interface {
	// Attempt returns the n-th backoff duration.
	// Used algorithm depends on implementers.
	// Returned duration can be zero or positive.
	Attempt(attempt int) time.Duration
}

// NewFixedBackoff returns a new instance of FixedBackoff.
// See comments on the [FixedBackoff] for details.
// Allowed value range is value>=0, otherwise it panics.
func NewFixedBackoff(value time.Duration) *FixedBackoff {
	mustZeroOrPositive(float64(value), "value")
	return &FixedBackoff{
		value: value,
	}
}

// FixedBackoff provides fixed interval backoff strategy.
//
// Algorithm:
//
//	Parameter:
//	  value: fixed duration. (>=0)
//	Input:
//	  attempt: the count of attempts.
//	Output:
//	  Always return the value.
//	Value range:
//	  Min: value
//	  Max: value
//
// Graph:
//
//	     y:backoff
//	     |
//	     |
//	     |
//	     |      y=value
//	value|-------------------------------
//	     |
//	     |
//	     |
//	     |
//	     └─────────────────────────────── x:attempts
//	     0
type FixedBackoff struct {
	value time.Duration
}

func (b *FixedBackoff) Attempt(_ int) time.Duration {
	return b.value
}

// NewRandomBackoff returns a new instance of RandomBackoff.
// See comments on the [RandomBackoff] for details.
// Allowed value range is offset>=0, fluctuation>0, otherwise it panics.
func NewRandomBackoff(offset, fluctuation time.Duration) *RandomBackoff {
	mustZeroOrPositive(float64(offset), "offset")
	mustPositive(float64(fluctuation), "fluctuation")
	mustUnderMaxInt64(offset, fluctuation, "offset+fluctuation")
	return &RandomBackoff{
		offset:      offset,
		fluctuation: int64(fluctuation),
	}
}

// RandomBackoff provides random backoff strategy.
// It will panic if the parameter range is not satisfied.
//
// Algorithm:
//
//	Parameter:
//	  offset: offset of the backoff duration. (>=0)
//	  fluctuation: backoff fluctuation range. (>0)
//	Input:
//	  attempt: the count of attempts.
//	Output:
//	  Calculate backoff duration bod := offset + RandomRange(0, fluctuation)
//	  Return bod.
//	Value range:
//	  Min: offset
//	  Max: offset + fluctuation
//
// Graph:
//
//	      y:backoff
//	      |
//	offset|-------------------------------
//	  +   |      /\        /──\
//	fluc. |     /  | /|    |  |   /\  y=offset+random(fluc)
//	      |    /   | | \   |  \  /  \
//	      |   /    |/   \  /   \/    \
//	      |  /           \/           \
//	offset|-------------------------------
//	      |
//	      |
//	      └─────────────────────────────── x:attempts
//	      0
type RandomBackoff struct {
	offset      time.Duration
	fluctuation int64
}

func (b *RandomBackoff) Attempt(_ int) time.Duration {
	return b.offset + time.Duration(rand.Int64N(b.fluctuation))
}

// NewLinearBackoff returns a new instance of LinearBackoff.
// See comments on the [LinearBackoff] for details.
// Allowed value range is offset>=0, fluctuation>=0, coeff>=0, otherwise it panics.
func NewLinearBackoff(offset, fluctuation, coeff time.Duration) *LinearBackoff {
	mustZeroOrPositive(float64(offset), "offset")
	mustZeroOrPositive(float64(fluctuation), "fluctuation")
	mustZeroOrPositive(float64(coeff), "coeff")
	mustUnderMaxInt64(offset, fluctuation, "offset+fluctuation")
	return &LinearBackoff{
		offset:      offset,
		fluctuation: float64(fluctuation),
		coeff:       float64(coeff),
	}
}

// LinearBackoff provides fixed backoff strategy.
// It will panic if the parameter range is not satisfied.
//
// Algorithm:
//
//	Parameter:
//	  offset: offset duration. (>=0)
//	  coeff: coefficient for linear function. (>=0)
//	  fluctuation: value range of fluctuation part. (>=0)
//	Input:
//	  attempt : the count of attempts.
//	Output:
//	  Calculate fluctuation value flc := (coeff * attempt)
//	  If flc > fluctuation, then let flc fluctuation.
//	  Calculate backoff duration bod := offset + flc
//	  Return bod.
//	Value range:
//	  Min: offset
//	  Max: offset + fluctuation
//
// Graph:
//
//	      y:backoff
//	      |
//	offset|------------───────────────────
//	  +   |         ／
//	fluc. |       ／
//	      |     ／ y=offset+coeff*x
//	      |   ／
//	      | ／
//	offset|-------------------------------
//	      |
//	      |
//	      └─────────────────────────────── x:attempts
//	      0
type LinearBackoff struct {
	offset      time.Duration
	coeff       float64
	fluctuation float64
}

func (b *LinearBackoff) Attempt(attempt int) time.Duration {
	flc := b.coeff * float64(attempt)
	if flc > b.fluctuation || flc == math.Inf(1) {
		flc = b.fluctuation
	}
	return b.offset + time.Duration(flc)
}

// NewPolynomialBackoff returns a new instance of PolynomialBackoff.
// See comments on the [PolynomialBackoff] for details.
// Allowed value range is offset>=0, fluctuation>=0, coeff>=0, otherwise it panics.
func NewPolynomialBackoff(offset, fluctuation, coeff time.Duration, exponent float64) *PolynomialBackoff {
	mustZeroOrPositive(float64(offset), "offset")
	mustZeroOrPositive(float64(fluctuation), "fluctuation")
	mustZeroOrPositive(float64(coeff), "coeff")
	mustUnderMaxInt64(offset, fluctuation, "offset+fluctuation")
	return &PolynomialBackoff{
		offset:      offset,
		coeff:       float64(coeff),
		fluctuation: float64(fluctuation),
		exponent:    exponent,
	}
}

// PolynomialBackoff provides polynomial backoff strategy.
// It will panic if the parameter range is not satisfied.
//
// Algorithm:
//
//	Parameter:
//	  offset: offset duration. (>=0)
//	  coeff: coefficient for polynomial function. (>=0)
//	  fluctuation: value range of fluctuation part. (>=0)
//	  exponent: exponent value for exponential function.
//	Input:
//	  attempt : the count of attempts. (>0)
//	Output:
//	  Calculate fluctuation value flc := (coeff * attempt^exponent)
//	  If flc > fluctuation, then let flc fluctuation.
//	  Calculate backoff duration bod := offset + flc
//	  Return bod.
//	Value range:
//	  Min: offset
//	  Max: offset + fluctuation
//	Limitations:
//	  The "attempt" must be attempt>0. If not, 1 is used.
//
// Graph:
//
//	      y:backoff
//	      |
//	offset|------------────────────────────
//	  +   |          /
//	fluc. |         /  y=offset+coeff*x^exponent
//	      |        /
//	      |      ／
//	      |    ／
//	offset|──／----------------------------
//	      |
//	      |
//	      └─────────────────────────────── x:attempts
//	      0
type PolynomialBackoff struct {
	offset      time.Duration
	coeff       float64
	fluctuation float64
	exponent    float64
}

func (b PolynomialBackoff) Attempt(attempt int) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}
	flc := b.coeff * math.Pow(float64(attempt), b.exponent)
	if flc > b.fluctuation || flc == math.Inf(1) {
		flc = b.fluctuation
	}
	return b.offset + time.Duration(flc)
}

// NewExponentialBackoff returns a new instance of ExponentialBackoff.
// See comments on the [ExponentialBackoff] for details.
// Allowed value range is offset>=0, fluctuation>=0, coeff>=0, otherwise it panics.
func NewExponentialBackoff(offset, fluctuation, coeff time.Duration) *ExponentialBackoff {
	mustZeroOrPositive(float64(offset), "offset")
	mustZeroOrPositive(float64(fluctuation), "fluctuation")
	mustZeroOrPositive(float64(coeff), "coeff")
	mustUnderMaxInt64(offset, fluctuation, "offset+fluctuation")
	return &ExponentialBackoff{
		offset:      offset,
		coeff:       float64(coeff),
		fluctuation: float64(fluctuation),
	}
}

// ExponentialBackoff provides exponential backoff strategy.
// It will panic if the parameter range is not satisfied.
//
// Algorithm:
//
//	Parameter:
//	  offset: offset duration. (>=0)
//	  coeff: coefficient for exponential function. (>=0)
//	  fluctuation: value range of fluctuation part. (>=0)
//	Input:
//	  attempt : the count of attempts.
//	Output:
//	  Calculate fluctuation value flc := (coeff * 2^attempt)
//	  If flc > fluctuation, then let flc fluctuation.
//	  Calculate backoff duration bod := offset + flc
//	  Return bod.
//	Value range:
//	  Min: offset
//	  Max: offset + fluctuation
//
// Graph:
//
//	      y:backoff
//	      |
//	offset|------------────────────────────
//	  +   |          /
//	fluc. |         /  y=offset+coeff*2^x
//	      |        /
//	      |      ／
//	      |    ／
//	offset|──／----------------------------
//	      |
//	      |
//	      └─────────────────────────────── x:attempts
//	      0
type ExponentialBackoff struct {
	offset      time.Duration
	coeff       float64
	fluctuation float64
}

func (b *ExponentialBackoff) Attempt(attempt int) time.Duration {
	flc := b.coeff * math.Pow(2, float64(attempt))
	if flc > b.fluctuation || flc == math.Inf(1) {
		flc = b.fluctuation
	}
	return b.offset + time.Duration(flc)
}

// NewExponentialBackoffFullJitter returns a new instance of ExponentialBackoffFullJitter.
// See comments on the [ExponentialBackoffFullJitter] for details.
// Allowed value range is offset>=0, fluctuation>0, coeff>0, otherwise it panics.
func NewExponentialBackoffFullJitter(offset, fluctuation, coeff time.Duration) *ExponentialBackoffFullJitter {
	mustZeroOrPositive(float64(offset), "offset")
	mustPositive(float64(fluctuation), "fluctuation")
	mustPositive(float64(coeff), "coeff")
	mustUnderMaxInt64(offset, fluctuation, "offset+fluctuation")
	return &ExponentialBackoffFullJitter{
		offset:      offset,
		coeff:       float64(coeff),
		fluctuation: float64(fluctuation),
	}
}

// ExponentialBackoffFullJitter provides exponential backoff strategy
// with full jitter.
// It will panic if the parameter range is not satisfied.
//
// Algorithm:
//
//	Parameter:
//	  offset: offset duration. (>=0)
//	  coeff: coefficient for exponential function. (>0)
//	  fluctuation: value range of fluctuation part. (>0)
//	Input:
//	  attempt : the count of attempts.
//	Output:
//	  Calculate fluctuation value flc := (coeff * 2^attempt)
//	  If flc > fluctuation, then let flc fluctuation.
//	  Calculate backoff duration bod := offset + RandomRange(0,flc)
//	  Return bod.
//	Value range:
//	  Min: offset
//	  Max: offset + fluctuation
//
// Graph:
//
//	      y:backoff
//	      |            y=offset+coeff*2^x
//	offset|------------────────────────────
//	  +   |          ////////// ↑ /////////
//	fluc. |         /////////// | /////////
//	      |        //////////// | random //
//	      |      ／//////////// | /////////
//	      |    ／////////////// ↓ /////////
//	offset|──／----------------------------
//	      |
//	      |
//	      └─────────────────────────────── x:attempts
//	      0
type ExponentialBackoffFullJitter struct {
	offset      time.Duration
	coeff       float64
	fluctuation float64
}

func (b *ExponentialBackoffFullJitter) Attempt(attempt int) time.Duration {
	flc := b.coeff * math.Pow(2, float64(attempt))
	if flc > b.fluctuation || flc == math.Inf(1) {
		flc = b.fluctuation
	}
	v := max(1, int64(flc))
	return b.offset + time.Duration(rand.Int64N(v))
}

// NewExponentialBackoffEqualJitter returns a new instance of ExponentialBackoffEqualJitter.
// See comments on the [ExponentialBackoffEqualJitter] for details.
// Allowed value range is offset>=0, fluctuation>0, coeff>0, otherwise it panics.
func NewExponentialBackoffEqualJitter(offset, fluctuation, coeff time.Duration) *ExponentialBackoffEqualJitter {
	mustZeroOrPositive(float64(offset), "offset")
	mustPositive(float64(fluctuation), "fluctuation")
	mustPositive(float64(coeff), "coeff")
	mustUnderMaxInt64(offset, fluctuation, "offset+fluctuation")
	return &ExponentialBackoffEqualJitter{
		offset:      offset,
		coeff:       float64(coeff),
		fluctuation: float64(fluctuation),
	}
}

// ExponentialBackoffEqualJitter provides exponential backoff strategy
// with equal jitter.
// It will panic if the parameter range is not satisfied.
//
// Algorithm:
//
//	Parameter:
//	  offset: offset duration. (>=0)
//	  coeff: coefficient for exponential function. (>0)
//	  fluctuation: value range of fluctuation part. (>0)
//	Input:
//	  attempt : the count of attempts.
//	Output:
//	  Calculate fluctuation value flc := (coeff * 2^attempt)
//	  If flc > fluctuation, then let flc fluctuation.
//	  Calculate backoff duration bod := offset + (flc/2 + RandomRange(0,flc/2))
//	  Return bod.
//	Value range:
//	  Min: offset
//	  Max: offset + fluctuation
//
// Graph:
//
//	      y:backoff
//	      |            y=offset+coeff*2^x
//	offset|------------────────────────────
//	  +   |          ////////// ↑ /////////
//	fluc. |         /////////// ↓ random //
//	      |        ////────────────────────
//	      |      ／//
//	      |    ／//
//	offset|──／----------------------------
//	      |
//	      |
//	      └─────────────────────────────── x:attempts
//	      0
type ExponentialBackoffEqualJitter struct {
	offset      time.Duration
	coeff       float64
	fluctuation float64
}

func (b *ExponentialBackoffEqualJitter) Attempt(attempt int) time.Duration {
	flc := b.coeff * math.Pow(2, float64(attempt))
	if flc > b.fluctuation || flc == math.Inf(1) {
		flc = b.fluctuation
	}
	var v int64
	if flc > 2 {
		v = int64(0.5 * flc)
	} else {
		v = 1 // rand.Int64N panics when 0.
	}
	return b.offset + time.Duration(v+rand.Int64N(v))
}

// NewFibonacciBackoff returns a new instance of FibonacciBackoff.
// See comments on the [FibonacciBackoff] for details.
// Allowed value range is offset>=0, fluctuation>=0, coeff>=0, otherwise it panics.
func NewFibonacciBackoff(offset, fluctuation, coeff time.Duration) *FibonacciBackoff {
	mustZeroOrPositive(float64(offset), "offset")
	mustZeroOrPositive(float64(fluctuation), "fluctuation")
	mustZeroOrPositive(float64(coeff), "coeff")
	mustUnderMaxInt64(offset, fluctuation, "offset+fluctuation")
	return &FibonacciBackoff{
		offset:      offset,
		coeff:       float64(coeff),
		fluctuation: float64(fluctuation),
	}
}

// FibonacciBackoff provides backoff strategy
// using fibonacci sequence.
// It will panic if the parameter range is not satisfied.
//
// Algorithm:
//
//	Parameter:
//	  offset: offset duration. (>=0)
//	  coeff: coefficient for fibonacci number. (>=0)
//	  fluctuation: value range of fluctuation part. (>=0)
//	Input:
//	  attempt : the count of attempts.
//	Output:
//	  Calculate fluctuation value flc := coeff * fibonacci(attempt).
//	  If flc > fluctuation, then let flc fluctuation.
//	  Calculate backoff duration bod := offset + flc
//	  Return bod.
//	Value range:
//	  Min: offset
//	  Max: offset + fluctuation
//	Limitations:
//	  The "attempt" must be attempt<=46. If exceeded, 46 is used.
type FibonacciBackoff struct {
	offset      time.Duration
	coeff       float64
	fluctuation float64
}

func (b *FibonacciBackoff) Attempt(attempt int) time.Duration {
	if attempt > 46 { // fibonacci should not overflow for int32.
		attempt = 46
	}
	flc := b.coeff * float64(fibonacci(attempt))
	if flc > b.fluctuation || flc == math.Inf(1) {
		flc = b.fluctuation
	}
	return b.offset + time.Duration(flc)
}

// fibonacci returns the n-th fibonacci number.
// It does not check the overflow.
//
//   - fibonacci(-1) : 0
//   - fibonacci(0)  : 0
//   - fibonacci(1)  : 1
//   - fibonacci(2)  : 1
//   - fibonacci(3)  : 2
//   - fibonacci(4)  : 3
//   - fibonacci(46) : 1,836,311,903 (Limit of int32)
//   - fibonacci(47) : 2,971,215,073 (int32 overflow)
//   - fibonacci(92) : 7,540,113,804,746,346,429 (Limit of int64)
//   - fibonacci(93) : 12,200,160,415,121,876,738 (int64 overflow)
func fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n <= 2 {
		return 1
	}
	n1 := 1
	n2 := 1
	for range n - 3 {
		tmp := n1 + n2
		n1 = n2
		n2 = tmp
	}
	return n1 + n2
}
