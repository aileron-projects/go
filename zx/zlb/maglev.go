package zlb

import (
	"math"
)

// NewMaglev returns a new instance of maglev hash load balancer.
// Targets can be added or removed after instantiation.
// See the comments on [Maglev].
func NewMaglev[T Target](targets ...T) *Maglev[T] {
	lb := &Maglev[T]{
		baseLB: &baseLB[T]{},
	}
	lb.Add(targets...)
	lb.Update()
	return lb
}

// Maglev is a load balancer that uses the Maglev hashing algorithm.
// For n targets, lookup complexity is O(1).
// However, it may try up to MaxRetry+1 times to find an active,
// non-zero-weight target.
//
// Note: The load balancer may return false even when active and
// non-zero-weight targets exist.
//
// Target weights are evaluated only when the lookup table is rebuilt,
// and cannot be updated dynamically.
// Inactive targets are excluded from the table during reconstruction.
//
// Algorithm:
//   - See references for the original Maglev hashing design.
//
// References:
//   - https://research.google/pubs/maglev-a-fast-and-reliable-software-network-load-balancer/
//   - https://www.usenix.org/sites/default/files/conference/protected-files/nsdi16_slides_eisenbud.pdf
type Maglev[T Target] struct {
	*baseLB[T]
	// table is the lookup table.
	table []int
	// MaxRetry is the maximum retry count.
	// If zero, no retries are applied.
	// The load balancer returns first found target
	// with active and non-zero weight.
	MaxRetry int
	// Size is the lookup table size.
	// The actual lookup table size wil be adjusted
	// based on the target weights.
	Size int
}

// Remove removes target.
// lb.Update is internally called after removing the target.
func (lb *Maglev[T]) Remove(id uint64) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.remove(id)
	lb.update()
}

// Update update lookup table.
func (lb *Maglev[T]) Update() {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.update()
}

func (lb *Maglev[T]) update() {
	targets := make([]int, 0, len(lb.targets))
	sumWeight := 0
	for i, t := range lb.targets {
		if !t.Active() || t.Weight() == 0 {
			continue
		}
		targets = append(targets, i)
		sumWeight += int(t.Weight())
	}
	if len(targets) == 0 {
		clear(lb.table)
		lb.table = nil
		return
	}
	tableSize := genPrimeEuler(max(2*sumWeight, lb.Size))

	offset := make([]int, len(targets))
	skip := make([]int, len(targets))
	for j, i := range targets {
		key := mix64(lb.targets[i].ID())
		offset[j] = int(key % uint64(tableSize))
		skip[j] = int(key%uint64(tableSize-1)) + 1
	}

	permTable := func(i, j int) int { // Permutation tables.
		return (offset[j] + skip[j]*i) % tableSize
	}
	lb.table = make([]int, tableSize) // Lookup table.
	for i := range tableSize {
		lb.table[i] = -1 // Initialize with -1 as non-filled marker.
	}

	// Fill the lookup table with considering target weights.
	indexes := make([]int, len(targets))
	total := 0
loop:
	for range tableSize {
		for j, index := range targets {
			i := indexes[j]
			count := 0 // Assigned count.
			for i < tableSize {
				pos := permTable(i, j)
				if lb.table[pos] < 0 {
					lb.table[pos] = index
					count += 1
					total += 1
					if total >= tableSize {
						break loop
					}
					if count >= int(lb.targets[index].Weight()) {
						break
					}
				}
				i++
			}
			indexes[j] = i
		}
	}
}

// Get returns a target.
// Returned value found is true when an active target was found.
func (lb *Maglev[T]) Get(key uint64) (t T, found bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if len(lb.table) == 0 {
		return t, false
	}

	maxIter := max(1, lb.MaxRetry+1)
	for range maxIter {
		key = mix64(key)
		index := int(key % uint64(len(lb.table)))
		i := lb.table[index]
		t = lb.targets[i]
		if t.Active() && t.Weight() > 0 {
			return t, true
		}
	}
	return t, false
}

// genPrimeEuler returns a prime number grater than
// the given min value using euler's formula.
// genPrimeEuler returns a prime value grater than 3.
//
// See https://en.wikipedia.org/wiki/Lucky_numbers_of_Euler
//   - x*x + x + 41
//   - x*x - x + 41
func genPrimeEuler(min int) int {
	if min < 41 {
		for _, p := range []int{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37} {
			if min <= p {
				return p
			}
		}
	}
	var val int
	i := -1
	for {
		i++
		val = i*i + i + 41 // (i*i - i + 41)
		if val >= min {
			break
		}
	}
	for {
		if isPrime(val) {
			return val
		}
		i++
		val = i*i + i + 41 // Equal to i*i - i + 41
	}
}

// isPrime returns if the given number is
// a prime number of not.
func isPrime(n int) bool {
	switch {
	case n <= 1:
		return false
	case n == 2:
		return true
	case n%2 == 0:
		return false
	}
	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}
