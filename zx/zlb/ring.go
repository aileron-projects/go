package zlb

import (
	"slices"
	"sort"
)

// NewRingHash returns a new instance of ring hash load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [RingHash].
func NewRingHash[T Target](targets ...T) *RingHash[T] {
	lb := &RingHash[T]{
		baseLB: &baseLB[T]{},
	}
	lb.Add(targets...)
	lb.Update()
	return lb
}

// RingHash is a load balancer that uses a ring-hash algorithm.
// For n targets, computational complexity is O(log(n)).
// This complexity is derived from the binary search used in the algorithm.
// Note that the load balancer may return false even active and
// non-zero weight targets are exist.
// The internal virtual ring is updated when [RingHash.Update] is called.
// Target weights are evaluated at the time of updating the virtual ring
// and cannot be updated dynamically.
//
// Algorithm:
//
//	Ti = T[i] = Target[i]
//	Wi = W[i] = Weight[i]
//	Ri = R[i] = Ring[i] // Ring[i] maps to Target[j].
//
//	function Get(key):
//	  key <-- SplitMix64(key)
//	  index <-- BinarySearch(key, Ring[i])
//	  for range {0 ... MaxTry}:
//	    if R[index] is not active:
//	      index <-- index+1
//	      continue
//	    if R[index] == 0:
//	      index <-- index+1
//	      continue
//	    return R[index]
//	  return nil
//
//	┌─────────────┬────────┬────────┬────────┬────────┬────────┬────────┬────────┐
//	│ Target T[i] │   T0   │   T1   │   T2   │   T3   │   T4   │  ....  │  Tn-1  │
//	│ Weight W[i] │  W0=20 │  W1=30 │  W2=0  │  W3=10 │  W4=10 │  ....  │ Wn-1=2 │
//	│ Active      │  true  │  true  │  true  │  false │  true  │  ....  │  true  │
//	└─────────────┴────────┴────────┴───┬────┴───┬────┴────────┴────────┴────────┘
//	                                    │        │
//	                                    │        └── Non-active target is ignored.
//	                                    └─────────── Target with 0 weight is ignored.
//
//	A concept of virtual ring.
//	Targets are mapped on the ring with size MaxUint64 based on their hash value.
//
//	                Start index determined by binary search.
//	                  ↓---→ ---→ ---→ ---→ ---→
//	┌────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┐
//	│ T1 │ T9 │ T3 │ T0 │ T1 │ T8 │ T7 │ T5 │ T6 │ T4 │....│ T5 │
//	├────┼────┴────┴────┴────┴────┴────┴────┴────┴────┴────┼────┤
//	│ T5 │                                                 │ T0 │
//	├────┤                                                 ├────┤
//	│ T6 │                                                 │ T8 │
//	├────┤                                                 ├────┤
//	│ T9 │                                                 │ T4 │
//	├────┤                                                 ├────┤
//	│ T3 │                                                 │ T1 │
//	├────┼────┬────┬────┬────┬────┬────┬────┬────┬────┬────┼────┤
//	│ T1 │....│ T0 │ T7 │ T4 │ T8 │ T1 │ T5 │ T8 │ T1 │ T7 │ T9 │
//	└────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┘
type RingHash[T Target] struct {
	*baseLB[T]
	ring [][2]uint64
}

// Remove removes target.
// lb.Update is internally called after removing the target.
func (lb *RingHash[T]) Remove(id uint64) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.remove(id)
	lb.update()
}

// Update update lookup table.
func (lb *RingHash[T]) Update() {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.update()
}

func (lb *RingHash[T]) update() {
	sumWeight := uint64(0)
	for i := range lb.targets {
		sumWeight += uint64(lb.targets[i].Weight())
	}
	clear(lb.ring)
	if uint64(cap(lb.ring)) < sumWeight {
		lb.ring = make([][2]uint64, 0, sumWeight)
	}
	lb.ring = lb.ring[:0]
	for i := range lb.targets {
		key := lb.targets[i].ID()
		for range lb.targets[i].Weight() {
			key = mix64(key)
			lb.ring = append(lb.ring, [2]uint64{key, uint64(i)})
		}
	}
	slices.SortFunc(lb.ring, func(a, b [2]uint64) int {
		switch {
		case a[0] == b[0]:
			return 0
		case a[0] > b[0]:
			return 1
		default:
			return -1
		}
	})
}

// Get returns a target.
// Returned value found is true when an active target was found.
func (lb *RingHash[T]) Get(key uint64) (t T, found bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if len(lb.ring) == 0 {
		return t, false
	}

	n := len(lb.ring)
	key = mix64(key)
	index := sort.Search(n, func(i int) bool {
		return lb.ring[i][0] >= key
	})

	for range min(10*len(lb.targets), n) { // No reason to the min value.
		if index >= n {
			index = 0
		}
		i := lb.ring[index][1]
		t = lb.targets[i]
		if !t.Active() || t.Weight() == 0 {
			index++
			continue
		}
		return t, true
	}
	return t, false
}
