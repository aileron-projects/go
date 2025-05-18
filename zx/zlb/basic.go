package zlb

import (
	"slices"
	"sync"
)

// BasicRoundRobin returns a new instance of basic round robin load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [BasicRoundRobin].
func NewBasicRoundRobin[T Target](targets ...T) *BasicRoundRobin[T] {
	lb := &BasicRoundRobin[T]{}
	lb.Add(targets...)
	return lb
}

// BasicRoundRobin is the load balancer that uses basic round-robin algorithm.
// The load balancer considers both active status and weight.
// Target weights and their orders are respected.
// For n targets, computational complexity is O(1).
//
// Algorithm:
//
//	Ti = T[i] = Target[i]
//	Wi = W[i] = Weight[i]
//
//	States:
//	  index = Current target index
//	  weight = Remaining weight of the current target
//
//	function Get():
//	  for range {T0 ... Tn-1}:
//	    if weight == 0:
//	      index <-- index+1 (Reset to 0 when n)
//	      weight <-- W[index]
//	    weight <-- weight-1
//	    return T[index]
//	  return nil
//
//	                       Visit each targets consuming their weights.
//	                  ┌-----------------------------------------------------┐
//	                  ↓                                                     |
//	                   -------→ -------→ -------→ -------→ -------→ -------→↑
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
type BasicRoundRobin[T Target] struct {
	mu        sync.Mutex
	targets   []T    // targets is the list of all targets.
	curIndex  int    // curIndex is the current target index.
	curWeight uint16 // curWeight is the remained weight.
}

// Targets returns registered targets.
func (lb *BasicRoundRobin[T]) Targets() []T {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.targets
}

// Add adds targets.
func (lb *BasicRoundRobin[T]) Add(targets ...T) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if len(lb.targets) == 0 && len(targets) > 0 {
		lb.curWeight = targets[0].Weight()
	}
	lb.targets = append(lb.targets, targets...)
}

// Remove removes targets.
func (lb *BasicRoundRobin[T]) Remove(id uint64) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	for i := 0; i < len(lb.targets); i++ { // len() must be evaluated every loop.
		if lb.targets[i].ID() != id {
			continue
		}
		lb.targets = slices.Delete(lb.targets, i, i+1)
		switch {
		case i < lb.curIndex:
			lb.curIndex--
		case i == lb.curIndex:
			if lb.curIndex >= len(lb.targets) {
				lb.curIndex = 0
			}
			if len(lb.targets) > 0 {
				lb.curWeight = lb.targets[lb.curIndex].Weight()
			}
		}
		i-- // Repeat current index.
	}
}

// Get returns a target.
// Returned value found is true when an active target was found.
func (lb *BasicRoundRobin[T]) Get(_ uint64) (t T, found bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	n := len(lb.targets)
	if n == 0 {
		return t, false
	}

	i := lb.curIndex
	for range n + 1 {
		t = lb.targets[i]
		if !t.Active() || lb.curWeight == 0 {
			if i++; i >= n {
				i = 0
			}
			lb.curWeight = lb.targets[i].Weight()
			lb.curIndex = i
			continue
		}
		lb.curWeight--
		lb.curIndex = i
		return t, true
	}
	lb.curIndex = i
	return t, false
}
