package zlb

// NewJumpHash returns a new instance of jump hash load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [JumpHash].
func NewJumpHash[T Target](targets ...T) *JumpHash[T] {
	lb := &JumpHash[T]{
		baseLB: &baseLB[T]{removeMode: 1},
	}
	lb.Add(targets...)
	return lb
}

// JumpHash is the load balancer that uses jump hash algorithm.
// For n targets, computational complexity is O(1).
// However, it attempts up to MaxRetry+1 times to find an active
// target with non-zero weight.
// Note that the load balancer may return false even if active targets
// with non-zero weight exist.
//
// Algorithm:
//
//	Ti = T[i] = Target[i]
//	Wi = W[i] = Weight[i]
//
//	function Get(key):
//	  for range {0 ... MaxRetry}:
//	    key <-- SplitMix64(key)
//	    index <-- JumpHash(key, n)
//	    if T[index] is not active:
//	      continue
//	    if W[index] == 0:
//	      continue
//	    return T[index]
//	  return nil
//
//
//	                    Select JumpHash(key, n).
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
//
// References:
//   - https://arxiv.org/abs/1406.2294
//   - https://arxiv.org/ftp/arxiv/papers/1406/1406.2294.pdf
//   - https://en.wikipedia.org/wiki/Consistent_hashing
type JumpHash[T Target] struct {
	*baseLB[T]
	// MaxRetry is the maximum retry count.
	// If zero, no retries are applied.
	// The load balancer returns first found target
	// with active and non-zero weight.
	MaxRetry int
}

func (lb *JumpHash[T]) Get(key uint64) (t T, found bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	n := len(lb.targets)
	if n == 0 {
		return t, false
	}

	maxIter := max(1, lb.MaxRetry+1)
	for range maxIter {
		key = mix64(key)
		i := jumpConsistentHash(key, n)
		t = lb.targets[i]
		if t.Active() && t.Weight() > 0 {
			return t, true
		}
	}
	return t, false
}

// jumpConsistentHash implements jump consistent hash.
// It returns a value of (0, num).
func jumpConsistentHash(key uint64, num int) int {
	b := -1 // b is less than num.
	j := 0  // j is less than num.
	for j < num {
		b = j
		key = key*2862933555777941757 + 1
		j = int(float64(b+1) * (float64(1<<31) / float64(key>>33+1)))
	}
	return b
}
