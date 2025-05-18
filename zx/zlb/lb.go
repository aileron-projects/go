package zlb

import (
	"slices"
	"sync"
)

// Target is the load balance target interface.
type Target interface {
	// ID returns the identifier of this target.
	// It depends on implementation how the ID is used.
	// In some load balancers, IDs are used for calculating
	// hash value in their load balancing algorithms.
	// Therefore, IDs should be globally unique.
	// ID must be safe for concurrent call.
	ID() uint64
	// Active returns if this target is active or not.
	// It depends on implementation if the active status is used.
	// Some load balancing algorithms ignores the status.
	// Active must be safe for concurrent call.
	Active() bool
	// Weight returns the target weight, or priority.
	// It depends on implementation how the weight is used.
	// Some load balancing algorithms ignores the weight.
	// Weight must be safe for concurrent call.
	Weight() uint16
}

// LoadBalancer is the load balancer interface.
type LoadBalancer[T Target] interface {
	// Targets returns all registered targets.
	// It is safe for concurrent call.
	Targets() (targets []T)
	// Add adds targets to the load balancer.
	// It may change the internal state.
	// It is safe for concurrent call.
	Add(targets ...T) (err error)
	// Remove removes targets with the given id.
	// It may change the internal state.
	// It is safe for concurrent call.
	Remove(id uint64)
	// Get returns a next target.
	// Get is safe to concurrent call.
	Get(hint uint64) (t T, found bool)
}

// baseLB is the base struct for some load balancers.
type baseLB[T Target] struct {
	mu sync.Mutex
	// targets is the list of load balancing target.
	targets []T
	// removeMode is the mode of target removal.
	// The slice of targets will be updated with the selected mode.
	// Consistent hash based load balancers should use 1.
	// 	0: keep the order of targets.
	// 	1: fill the removed position with the target at the last of targets.
	// 	others: mode 0 is used.
	removeMode int
}

// Targets returns registered targets.
func (lb *baseLB[T]) Targets() []T {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.targets
}

// Add adds targets.
func (lb *baseLB[T]) Add(targets ...T) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.targets = append(lb.targets, targets...)
}

// Remove removes targets.
func (lb *baseLB[T]) Remove(id uint64) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.remove(id)
}

func (lb *baseLB[T]) remove(id uint64) {
	for i := 0; i < len(lb.targets); i++ { // len() must be evaluated every loop.
		if lb.targets[i].ID() != id {
			continue
		}
		n := len(lb.targets)
		switch lb.removeMode {
		case 1:
			lb.targets[i] = lb.targets[n-1]
			lb.targets = slices.Delete(lb.targets, n-1, n)
		default:
			lb.targets = slices.Delete(lb.targets, i, i+1)
		}
		i-- // Repeat current index.
	}
}

// mix64 returns a pseudo-random number.
// It is based on the SplitMix64 algorithm.
// mix64 is called from some load balancers.
//
// References:
//   - https://arxiv.org/abs/1805.01407
//   - https://en.wikipedia.org/wiki/Xorshift
//   - https://rosettacode.org/wiki/Pseudo-random_numbers/Splitmix64
//   - https://rosettacode.org/wiki/Pseudo-random_numbers/Xorshift_star
//   - https://en.wikipedia.org/wiki/Universal_hashing
func mix64(z uint64) uint64 {
	z = z + 0x9e3779b97f4a7c15
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}
