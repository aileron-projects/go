package zlb

import (
	"math"
	"math/rand/v2"
)

// NewRandom returns a new instance of random load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [Random].
func NewRandom[T Target](targets ...T) *Random[T] {
	lb := &Random[T]{
		baseLB: &baseLB[T]{},
	}
	lb.Add(targets...)
	return lb
}

// NewRandomW returns a new instance of weighted random load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [RandomW].
func NewRandomW[T Target](targets ...T) *RandomW[T] {
	lb := &RandomW[T]{
		baseLB: &baseLB[T]{},
	}
	lb.Add(targets...)
	return lb
}

// Random is the load balancer that uses random algorithm.
// For n targets, computational complexity is O(1).
// But it tries MaxRetry+1 times at most until it finds an active
// and non-zero weight target.
// Note that the load balancer may return false even active and
// non-zero weight targets are exist.
// Weighted random load balancer is available with [NewRandomW].
//
// Algorithm:
//
//	Ti = T[i] = Target[i]
//	Wi = W[i] = Weight[i]
//
//	function Get(_):
//	  for i range {0 ... MaxRetry}:
//	    index <-- random(0, n)
//	    if T[index] is not active:
//	      continue
//	    if W[index] == 0:
//	      continue
//	    return T[index]
//	  return nil
//
//
//	                     Select random()%n.
//	              Retry if active=false or weight=0.
//	                            ↓
//	┌─────────────┬────────┬────────┬────────┬────────┬────────┬────────┬────────┐
//	│ Target T[i] │   T0   │   T1   │   T2   │   T3   │   T4   │  ....  │  Tn-1  │
//	│ Weight W[i] │  W0=2  │  W1=3  │  W2=0  │  W3=1  │  W4=1  │  ....  │ Wn-1=2 │
//	│ Active      │  true  │  true  │  true  │ false  │  true  │  ....  │  true  │
//	└─────────────┴────────┴────────┴───┬────┴───┬────┴────────┴────────┴────────┘
//	                                    │        │
//	                                    │        └── Non-active target is ignored.
//	                                    └─────────── Target with 0 weight is ignored.
type Random[T Target] struct {
	*baseLB[T]
	MaxRetry int
}

// Get returns a target.
// Returned value found is true when an active target was found.
func (lb *Random[T]) Get(_ uint64) (t T, found bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	n := len(lb.targets)
	if n == 0 {
		return t, false
	}

	maxIter := max(1, lb.MaxRetry+1)
	for range maxIter {
		i := rand.Int() % n
		t = lb.targets[i]
		if t.Active() && t.Weight() > 0 {
			return t, true
		}
	}
	return t, false
}

// RandomW is the load balancer that uses weighted random algorithm.
// The load balancer considers both active status and weight of targets.
// For n targets, computational complexity is O(n).
// Non-weighted random load balancer is available with [NewRandom].
//
// Algorithm:
//
//	Ti = T[i] = Target[i]
//	Wi = W[i] = Weight[i]
//
//	function Get(_):
//	  maxScore <-- 0
//	  target <-- nil
//	  for i range {T0 ... Tn-1}:
//	    score <-- Wi / -log( random() )
//	    if score > maxScore:
//	      maxScore <-- score
//	      target <-- Ti
//	  return target
//
//
//	              Select the target with max score.
//	                            ↓
//	  Calculate score  -------→ -------→ -------→ -------→ -------→ -------→
//	  Score  S[i] │  1.25  │  3.92  │   --   │   --   │  2.04  │  ....  │  1.95  │
//	┌─────────────┬────────┬────────┬────────┬────────┬────────┬────────┬────────┐
//	│ Target T[i] │   T0   │   T1   │   T2   │   T3   │   T4   │  ....  │  Tn-1  │
//	│ Weight W[i] │  W0=2  │  W1=3  │  W2=0  │  W3=1  │  W4=1  │  ....  │ Wn-1=2 │
//	│ Active      │  true  │  true  │  true  │ false  │  true  │  ....  │  true  │
//	└─────────────┴────────┴────────┴───┬────┴───┬────┴────────┴────────┴────────┘
//	                                    │        │
//	                                    │        └── Non-active target is ignored.
//	                                    └─────────── Target with 0 weight is ignored.
//
// References:
//   - https://en.wikipedia.org/wiki/Reservoir_sampling
//   - https://utopia.duth.gr/~pefraimi/research/data/2007EncOfAlg.pdf
type RandomW[T Target] struct {
	*baseLB[T]
}

// Get returns a target.
// Returned value found is true when an active target was found.
func (lb *RandomW[T]) Get(_ uint64) (t T, found bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	n := len(lb.targets)
	if n == 0 {
		return t, false
	}

	maxScore := -float64(math.MaxUint64)
	index := 0 // Target index of max score.
	for i, target := range lb.targets {
		w := float64(target.Weight())
		if !target.Active() || w == 0 {
			continue
		}
		found = true
		score := w / -math.Log(rand.Float64())
		if score > maxScore {
			maxScore = score
			index = i
		}
	}

	if !found {
		return t, false
	}
	return lb.targets[index], true
}
