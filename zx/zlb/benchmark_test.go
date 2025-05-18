package zlb_test

import (
	"testing"

	"github.com/aileron-projects/go/zx/zlb"
)

type target struct {
	id     uint64
	weight uint16
	active bool
}

func (t *target) ID() uint64 {
	return t.id
}

func (t *target) Weight() uint16 {
	return t.weight
}

func (t *target) Active() bool {
	return t.active
}

var testTargets = func() []*target {
	ts := []*target{}
	for i := range 1000 {
		for j := range 10 {
			ts = append(ts, &target{
				id:     uint64(i)<<32 + uint64(j),
				weight: uint16(j / 2),
				active: []bool{true, true, false}[j%3],
			})
		}
	}
	return ts
}()

func BenchmarkPriority(b *testing.B) {
	lb := zlb.NewPriority(testTargets...)
	b.ResetTimer()
	for b.Loop() {
		lb.Get(0)
	}
}

func BenchmarkRandom(b *testing.B) {
	lb := zlb.NewRandom(testTargets...)
	lb.MaxRetry = 2
	b.ResetTimer()
	for b.Loop() {
		lb.Get(0)
	}
}

func BenchmarkRandomW(b *testing.B) {
	lb := zlb.NewRandomW(testTargets...)
	b.ResetTimer()
	for b.Loop() {
		lb.Get(0)
	}
}

func BenchmarkBasicRoundRobin(b *testing.B) {
	lb := zlb.NewBasicRoundRobin(testTargets...)
	b.ResetTimer()
	for b.Loop() {
		lb.Get(0)
	}
}

func BenchmarkRoundRobin(b *testing.B) {
	lb := zlb.NewRoundRobin(testTargets...)
	b.ResetTimer()
	for b.Loop() {
		lb.Get(0)
	}
}

func BenchmarkRendezvousHash(b *testing.B) {
	lb := zlb.NewRendezvousHash(testTargets...)
	b.ResetTimer()
	for i := range b.N {
		lb.Get(mulShift(i))
	}
}

func BenchmarkJumpHash(b *testing.B) {
	lb := zlb.NewJumpHash(testTargets...)
	lb.MaxRetry = 0
	b.ResetTimer()
	for i := range b.N {
		lb.Get(mulShift(i))
	}
}

func BenchmarkDirectHash(b *testing.B) {
	lb := zlb.NewDirectHash(testTargets...)
	lb.MaxRetry = 0
	b.ResetTimer()
	for i := range b.N {
		lb.Get(mulShift(i))
	}
}

func BenchmarkDirectHashW(b *testing.B) {
	lb := zlb.NewDirectHashW(testTargets...)
	b.ResetTimer()
	for i := range b.N {
		lb.Get(mulShift(i))
	}
}

func BenchmarkRingHash(b *testing.B) {
	lb := zlb.NewRingHash(testTargets...)
	b.ResetTimer()
	for i := range b.N {
		lb.Get(mulShift(i))
	}
}

func BenchmarkRingHash_Update(b *testing.B) {
	lb := zlb.NewRingHash(testTargets...)
	b.ResetTimer()
	for b.Loop() {
		lb.Update()
	}
}

func BenchmarkMaglev(b *testing.B) {
	lb := zlb.NewMaglev(testTargets...)
	lb.MaxRetry = 0
	b.ResetTimer()
	for i := range b.N {
		lb.Get(mulShift(i))
	}
}

func BenchmarkMaglev_Update(b *testing.B) {
	lb := zlb.NewMaglev(testTargets...)
	lb.MaxRetry = 0
	b.ResetTimer()
	for b.Loop() {
		lb.Update()
	}
}

// mulShift returns a hash value based on the
// multiply-shift algorithm.
//
// References:
//   - https://hjemmesider.diku.dk/~jyrki/Paper/CP-11.4.1997.pdf
//   - https://arxiv.org/pdf/1504.06804
//   - https://en.wikipedia.org/wiki/Universal_hashing
func mulShift(x int) uint64 {
	z := uint64(x)
	return (z * 0x9e3779b97f4a7c15)
}
