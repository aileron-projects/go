package zfnv_test

import (
	"testing"

	"github.com/aileron-projects/go/zhash/zfnv"
)

var benchData = []byte("Hello Go!")

func BenchmarkSum32(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zfnv.Sum32(benchData)
	}
}

func BenchmarkSum32a(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zfnv.Sum32a(benchData)
	}
}

func BenchmarkSum64(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zfnv.Sum64(benchData)
	}
}

func BenchmarkSum64a(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zfnv.Sum64a(benchData)
	}
}

func BenchmarkSum128(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zfnv.Sum128(benchData)
	}
}

func BenchmarkSum128a(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zfnv.Sum128a(benchData)
	}
}
