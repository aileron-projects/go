package zlb

import (
	"slices"
	"sync"
)

// NewRoundRobin returns a new instance of round-robin load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [RoundRobin].
func NewRoundRobin[T Target](targets ...T) *RoundRobin[T] {
	lb := &RoundRobin[T]{}
	lb.Add(targets...)
	return lb
}

// RoundRobin is the load balancer that uses smooth round-robin algorithm.
// The load balancer considers both active status and weight.
// For n targets, computational complexity is O(n).
//
// Algorithm:
//
//	Ti = T[i] = Target[i]
//	Wi = W[i] = Weight[i]
//
//	States:
//	  CWi = CW[i] = Current weight of T[i]
//
//	function Get():
//	  sumWeight <-- 0
//	  maxWeight <-- 0
//	  index <-- -1
//	  for i range {T0 ... Tn-1}:
//	    sumWeight <-- sumWeight + Wi
//	    CWi <-- CWi + Wi
//	    if CWi > maxWeight:
//	      maxWeight <-- CWi
//	      index <-- i
//	  CW[index] = CW[index] - sumWeight
//	  return T[index]
//
//
//	              Select the target with max current weight.
//	                            ↓
//	Calculate each     -------→ -------→ -------→ -------→ -------→ -------→
//	        CW[i] │    4   │    8   │    0   │    0   │    2   │  ....  │    5   │
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
//   - https://en.wikipedia.org/wiki/Weighted_round_robin
//   - https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35
//   - https://dubbo.apache.org/en/overview/what/core-features/load-balance/
type RoundRobin[T Target] struct {
	mu      sync.Mutex
	targets []T
	// weights is the list of current weight.
	weights []float64
}

// Targets returns registered targets.
func (lb *RoundRobin[T]) Targets() []T {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.targets
}

// Add adds targets.
func (lb *RoundRobin[T]) Add(targets ...T) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	for _, t := range targets {
		lb.targets = append(lb.targets, t)
		lb.weights = append(lb.weights, 0)
	}
}

// Remove removes targets.
func (lb *RoundRobin[T]) Remove(id uint64) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	for i := 0; i < len(lb.targets); i++ { // len() must be evaluated every loop.
		if lb.targets[i].ID() != id {
			continue
		}
		lb.targets = slices.Delete(lb.targets, i, i+1)
		lb.weights = slices.Delete(lb.weights, i, i+1)
		i-- // Repeat current index.
	}
}

// Get returns a target.
// Returned value found is true when an active target was found.
func (lb *RoundRobin[T]) Get(_ uint64) (t T, found bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	n := len(lb.targets)
	if n == 0 {
		return t, false
	}

	sumWeight := float64(0)
	maxWeight := float64(0)
	index := 0
	for i, target := range lb.targets {
		w := float64(target.Weight())
		if !target.Active() || w == 0 {
			continue
		}
		found = true
		sumWeight += w
		lb.weights[i] += w
		if lb.weights[i] > maxWeight {
			maxWeight = lb.weights[i]
			index = i
		}
	}

	if !found {
		return t, false
	}
	lb.weights[index] -= sumWeight
	return lb.targets[index], true
}
