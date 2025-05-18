package zlb_test

import (
	"fmt"

	"github.com/aileron-projects/go/zx/zlb"
)

type Target struct {
	name   string
	id     uint64
	weight uint16
	active bool
}

func (t *Target) ID() uint64 {
	return t.id
}

func (t *Target) Weight() uint16 {
	return t.weight
}

func (t *Target) Active() bool {
	return t.active
}

func ExamplePriority() {
	t0 := &Target{name: "t0", weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", weight: 1, active: true}
	t2 := &Target{name: "t2", weight: 2, active: true}
	t3 := &Target{name: "t3", weight: 3, active: true}
	t4 := &Target{name: "t4", weight: 4, active: false} // Non-active
	lb := zlb.NewPriority(t0, t1, t2, t3, t4)

	count := map[string]int{}
	for range 100 {
		target, found := lb.Get(0)
		if found {
			count["total"] += 1
			count[target.name] += 1
		}
	}

	fmt.Println(count)
	// Output:
	// map[t3:100 total:100]
}

func ExampleBasicRoundRobin() {
	t0 := &Target{name: "t0", weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", weight: 1, active: true}
	t2 := &Target{name: "t2", weight: 2, active: true}
	t3 := &Target{name: "t3", weight: 3, active: true}
	t4 := &Target{name: "t4", weight: 4, active: false} // Non-active
	lb := zlb.NewBasicRoundRobin(t0, t1, t2, t3, t4)

	count := map[string]int{}
	for range 1200 {
		target, found := lb.Get(0)
		if found {
			count["total"] += 1
			count[target.name] += 1
		}
	}

	fmt.Println(count)
	// Output:
	// map[t1:200 t2:400 t3:600 total:1200]
}

func ExampleRoundRobin() {
	t0 := &Target{name: "t0", weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", weight: 1, active: true}
	t2 := &Target{name: "t2", weight: 2, active: true}
	t3 := &Target{name: "t3", weight: 3, active: true}
	t4 := &Target{name: "t4", weight: 4, active: false} // Non-active
	lb := zlb.NewRoundRobin(t0, t1, t2, t3, t4)

	count := map[string]int{}
	for range 1200 {
		target, found := lb.Get(0)
		if found {
			count["total"] += 1
			count[target.name] += 1
		}
	}

	fmt.Println(count)
	// Output:
	// map[t1:200 t2:400 t3:600 total:1200]
}

func ExampleRendezvousHash() {
	t0 := &Target{name: "t0", id: 123, weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", id: 456, weight: 1, active: true}
	t2 := &Target{name: "t2", id: 789, weight: 2, active: true}
	t3 := &Target{name: "t3", id: 321, weight: 3, active: true}
	t4 := &Target{name: "t4", id: 654, weight: 4, active: false} // Non-active
	lb := zlb.NewRendezvousHash(t0, t1, t2, t3, t4)

	count := map[string]int{}
	for i := range 1200 {
		target, found := lb.Get(uint64(i+1) * 0x9e3779b97f4a7c15) // Give a very simple hash.
		if found {
			count["total"] += 1
			count[target.name] += 1
		}
	}

	fmt.Println(count)
	// Output:
	// map[t1:201 t2:387 t3:612 total:1200]
}

func ExampleRandom() {
	t0 := &Target{name: "t0", weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", weight: 1, active: true}
	t2 := &Target{name: "t2", weight: 2, active: true}
	t3 := &Target{name: "t3", weight: 3, active: true}
	t4 := &Target{name: "t4", weight: 4, active: false} // Non-active
	lb := zlb.NewRandom(t0, t1, t2, t3, t4)
	lb.MaxRetry = 1

	count := map[string]int{}
	for range 1200 {
		target, found := lb.Get(0)
		if found {
			count["total"] += 1 // Total won't be 1200.
			count[target.name] += 1
		}
	}

	fmt.Println(count)
}

func ExampleRandomW() {
	t0 := &Target{name: "t0", weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", weight: 1, active: true}
	t2 := &Target{name: "t2", weight: 2, active: true}
	t3 := &Target{name: "t3", weight: 3, active: true}
	t4 := &Target{name: "t4", weight: 4, active: false} // Non-active
	lb := zlb.NewRandomW(t0, t1, t2, t3, t4)

	count := map[string]int{}
	for range 1200 {
		target, found := lb.Get(0)
		if found {
			count["total"] += 1
			count[target.name] += 1
		}
	}

	fmt.Println(count)
}

func ExampleDirectHash() {
	t0 := &Target{name: "t0", weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", weight: 1, active: true}
	t2 := &Target{name: "t2", weight: 2, active: true}
	t3 := &Target{name: "t3", weight: 3, active: true}
	t4 := &Target{name: "t4", weight: 4, active: false} // Non-active
	lb := zlb.NewDirectHash(t0, t1, t2, t3, t4)
	lb.MaxRetry = 1

	count := map[string]int{}
	for i := range 1200 {
		target, found := lb.Get(uint64(i+1) * 0x9e3779b97f4a7c15) // Give a very simple hash.
		if found {
			count["total"] += 1 // Total won't be 1200.
			count[target.name] += 1
		}
	}

	fmt.Println(count)
	// Output:
	// map[t1:329 t2:319 t3:351 total:999]
}

func ExampleDirectHashW() {
	t0 := &Target{name: "t0", weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", weight: 1, active: true}
	t2 := &Target{name: "t2", weight: 2, active: true}
	t3 := &Target{name: "t3", weight: 3, active: true}
	t4 := &Target{name: "t4", weight: 4, active: false} // Non-active
	lb := zlb.NewDirectHashW(t0, t1, t2, t3, t4)

	count := map[string]int{}
	for i := range 1200 {
		target, found := lb.Get(uint64(i+1) * 0x9e3779b97f4a7c15) // Give a very simple hash.
		if found {
			count["total"] += 1
			count[target.name] += 1
		}
	}

	fmt.Println(count)
	// Output:
	// map[t1:199 t2:384 t3:617 total:1200]
}

func ExampleJumpHash() {
	t0 := &Target{name: "t0", id: 123, weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", id: 456, weight: 1, active: true}
	t2 := &Target{name: "t2", id: 789, weight: 2, active: true}
	t3 := &Target{name: "t3", id: 123, weight: 3, active: true}
	t4 := &Target{name: "t4", id: 926, weight: 4, active: false} // Non-active
	lb := zlb.NewJumpHash(t0, t1, t2, t3, t4)
	lb.MaxRetry = 1

	count := map[string]int{}
	for i := range 1200 {
		target, found := lb.Get(uint64(i+1) * 0x9e3779b97f4a7c15) // Give a very simple hash.
		if found {
			count["total"] += 1 // Total won't be 1200.
			count[target.name] += 1
		}
	}

	fmt.Println(count)
	// Output:
	// map[t1:324 t2:359 t3:326 total:1009]
}

func ExampleRingHash() {
	t0 := &Target{name: "t0", id: 123, weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", id: 456, weight: 10, active: true}
	t2 := &Target{name: "t2", id: 789, weight: 20, active: true}
	t3 := &Target{name: "t3", id: 123, weight: 30, active: true}
	t4 := &Target{name: "t4", id: 926, weight: 40, active: false} // Non-active
	lb := zlb.NewRingHash(t0, t1, t2, t3, t4)
	lb.Update()

	count := map[string]int{}
	for i := range 1200 {
		target, found := lb.Get(uint64(i+1) * 0x9e3779b97f4a7c15) // Give a very simple hash.
		if found {
			count["total"] += 1
			count[target.name] += 1
		}
	}

	fmt.Println(count)
	// Output:
	// map[t1:171 t2:442 t3:587 total:1200]
}

func ExampleMaglev() {
	t0 := &Target{name: "t0", id: 123, weight: 0, active: true} // 0 weight
	t1 := &Target{name: "t1", id: 456, weight: 1, active: true}
	t2 := &Target{name: "t2", id: 789, weight: 2, active: true}
	t3 := &Target{name: "t3", id: 123, weight: 3, active: true}
	t4 := &Target{name: "t4", id: 926, weight: 4, active: false} // Non-active
	lb := zlb.NewMaglev(t0, t1, t2, t3, t4)
	lb.MaxRetry = 1
	lb.Size = 111 // Prime number grater than sum(weight). It's good to choose near N*sum(weight).
	lb.Update()   // Update table.

	count := map[string]int{}
	for i := range 1200 {
		target, found := lb.Get(uint64(i+1) * 0x9e3779b97f4a7c15) // Give a very simple hash.
		if found {
			count["total"] += 1
			count[target.name] += 1
		}
	}

	fmt.Println(count)
	// Output:
	// map[t1:209 t2:410 t3:581 total:1200]
}
