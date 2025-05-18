package zlb

// NewPriority returns a new instance of priority-based load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [Priority].
func NewPriority[T Target](targets ...T) *Priority[T] {
	lb := &Priority[T]{
		baseLB: &baseLB[T]{},
	}
	lb.Add(targets...)
	return lb
}

// Priority is a load balancer that uses a priority-based algorithm.
// It is very simple algorithm that select a target with maximum
// priority, or maximum weight.
// The load balancer considers both active status and weight of targets.
// For n targets, computational complexity is O(n).
//
// Algorithm:
//
//	Ti = T[i] = Target[i]
//	Wi = W[i] = Weight[i]
//
//	function Get():
//	  maxWeight <-- 0
//	  index <-- -1
//	  for i range {T0 ... Tn-1}:
//	    if CWi > maxWeight:
//	      maxWeight <-- CWi
//	      index <-- i
//	  return T[index]
//
//
//	              Select the target with max weight.
//	              First found one is selected if same weight targets exist.
//	                            ↓
//	┌─────────────┬────────┬────────┬────────┬────────┬────────┬────────┬────────┐
//	│ Target T[i] │   T0   │   T1   │   T2   │   T3   │   T4   │  ....  │  Tn-1  │
//	│ Weight W[i] │  W0=2  │  W1=3  │  W2=0  │  W3=1  │  W4=1  │  ....  │ Wn-1=2 │
//	│ Active      │  true  │  true  │  true  │ false  │  true  │  ....  │  true  │
//	└─────────────┴────────┴────────┴───┬────┴───┬────┴────────┴────────┴────────┘
//	                                    │        │
//	                                    │        └── Non-active target is ignored.
//	                                    └─────────── Target with 0 weight is ignored.
type Priority[T Target] struct {
	*baseLB[T]
}

// Get returns a target.
// Returned value found is true when an active target was found.
func (lb *Priority[T]) Get(_ uint64) (t T, found bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	maxWeight := uint16(0)
	for _, target := range lb.targets {
		w := target.Weight()
		if !target.Active() || w == 0 {
			continue
		}
		found = true
		if w > maxWeight {
			maxWeight = w
			t = target
		}
	}
	return t, found
}
