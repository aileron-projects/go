package zcrc64_test

import (
	"testing"

	"github.com/aileron-projects/go/zhash/zcrc64"
)

var benchData = []byte("Hello Go!")

func BenchmarkSumISO(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zcrc64.SumISO(benchData)
	}
}

func BenchmarkSumECMA(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zcrc64.SumECMA(benchData)
	}
}
