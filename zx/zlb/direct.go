package zlb

import "math"

// NewDirectHash returns a new instance of direct hash load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [DirectHash].
func NewDirectHash[T Target](targets ...T) *DirectHash[T] {
	lb := &DirectHash[T]{
		baseLB: &baseLB[T]{removeMode: 1},
	}
	lb.Add(targets...)
	return lb
}

// NewDirectHashW returns a new instance of weighted direct hash load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [DirectHashW].
func NewDirectHashW[T Target](targets ...T) *DirectHashW[T] {
	lb := &DirectHashW[T]{
		baseLB: &baseLB[T]{removeMode: 1},
	}
	lb.Add(targets...)
	return lb
}

// DirectHash is a load balancer that uses a direct-hash algorithm.
// For n targets, the computational complexity is O(1).
// However, it attempts up to MaxRetry+1 times to find an active
// target with non-zero weight.
// Note that the load balancer may return false even if active targets
// with non-zero weight exist.
// A weighted variant of the direct-hash load balancer is available
// via [NewDirectHashW].
//
// Algorithm:
//
//	Ti = T[i] = Target[i]
//	Wi = W[i] = Weight[i]
//
//	function Get(key):
//	  for range {0 ... MaxRetry}:
//	    key <-- SplitMix64(key)
//	    index <-- key % n
//	    if T[index] is not active:
//	      continue
//	    if W[index] == 0:
//	      continue
//	    return T[index]
//	  return nil
//
//
//	                   Select SplitMix64(key)%n.
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
type DirectHash[T Target] struct {
	*baseLB[T]
	// MaxRetry is the maximum retry count.
	// If zero, no retries are applied.
	// The load balancer returns first found target
	// with active and non-zero weight.
	MaxRetry int
}

// Get returns a target.
// Returned value found is true when an active target was found.
func (lb *DirectHash[T]) Get(key uint64) (t T, found bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	n := len(lb.targets)
	if n == 0 {
		return t, false
	}

	maxIter := max(1, lb.MaxRetry+1)
	for range maxIter {
		key = mix64(key)
		i := int(key % uint64(n))
		t = lb.targets[i]
		if t.Active() && t.Weight() > 0 {
			return t, true
		}
	}
	return t, false
}

// DirectHashW is a load balancer that uses a weighted direct-hash algorithm.
// For n targets, computational complexity is O(n).
// Non-weighted direct-hash load balancer is available with [NewDirectHash].
//
// Algorithm:
//
//	Ti = T[i] = Target[i]
//	Wi = W[i] = Weight[i]
//
//	function Get(key):
//	  maxScore <-- 0
//	  target <-- nil
//	  for i range {T0 ... Tn-1}:
//	    key <-- SplitMix64(key)
//	    score <-- Wi / -log( key/MaxUint64 )
//	    if score > maxScore:
//	      maxScore <-- score
//	      target <-- Ti
//	  return target
//
//
//	              Select the target with max score.
//	                            ↓
//	  Calculate scores -------→ -------→ -------→ -------→ -------→ -------→
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
type DirectHashW[T Target] struct {
	*baseLB[T]
}

// Get returns a target.
// Returned value found is true when an active target was found.
func (lb *DirectHashW[T]) Get(key uint64) (t T, found bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	n := len(lb.targets)
	if n == 0 {
		return t, false
	}

	maxScore := -float64(math.MaxUint64)
	index := 0 // Target index of max score.
	for i, target := range lb.targets {
		w := target.Weight()
		if !target.Active() || w == 0 {
			continue
		}
		found = true
		key = mix64(key)
		score := float64(w) / -math.Log(float64(key)/math.MaxUint64)
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
