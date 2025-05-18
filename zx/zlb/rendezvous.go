package zlb

import (
	"math"
)

// NewRendezvousHash returns a new instance of rendezvous hash load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [RendezvousHash].
func NewRendezvousHash[T Target](targets ...T) *RendezvousHash[T] {
	lb := &RendezvousHash[T]{
		baseLB: &baseLB[T]{},
	}
	lb.Add(targets...)
	return lb
}

// RendezvousHash is a load balancer that use rendezvous-hash algorithm.
// Rendezvous Hashing also known as Highest Random Weight (HRW).
// The load balancer considers both active status and weight of targets.
// For n targets, computational complexity is O(n).
//
// Algorithm:
//
//	Ti = T[i] = Target[i]
//	Wi = W[i] = Weight[i]
//	IDi = T[i].ID() = identifier of T[i]
//
//	function Get(key):
//	  maxScore <-- 0
//	  target <-- nil
//	  for i range {T0 ... Tn-1}:
//	    score <-- Wi / -log( Hash(key*IDi) / MaxUint64 )
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
//   - https://en.wikipedia.org/wiki/Rendezvous_hashing
//   - https://en.wikipedia.org/wiki/Consistent_hashing
//   - https://en.wikipedia.org/wiki/Reservoir_sampling
//   - https://www.ietf.org/archive/id/draft-ietf-bess-weighted-hrw-00.html
//   - https://utopia.duth.gr/~pefraimi/research/data/2007EncOfAlg.pdf
//   - Weighted Random Sampling (2005; Efraimidis, Spirakis)
//   - Using Name-Based Mappings to Increase Hit Rates
//   - -> https://www.microsoft.com/en-us/research/wp-content/uploads/2017/02/HRW98.pdf
type RendezvousHash[T Target] struct {
	*baseLB[T]
}

// Get returns a target.
// Returned value found is true when an active target was found.
func (lb *RendezvousHash[T]) Get(key uint64) (t T, found bool) {
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
		uniform := float64(mix64(key*target.ID())) / math.MaxUint64
		score := w / -math.Log(uniform)
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
