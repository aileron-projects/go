package zcrc64_test

import (
	"testing"

	"github.com/aileron-projects/go/zhash/zcrc64"
)

var benchData = []byte("Hello Go!")

func BenchmarkSumISO(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		zcrc64.SumISO(benchData)
	}
}

func BenchmarkSumECMA(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		zcrc64.SumECMA(benchData)
	}
}
