package zcrc32_test

import (
	"testing"

	"github.com/aileron-projects/go/zhash/zcrc32"
)

var benchData = []byte("Hello Go!")

func BenchmarkSumIEEE(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		zcrc32.SumIEEE(benchData)
	}
}

func BenchmarkSumCastagnoli(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		zcrc32.SumCastagnoli(benchData)
	}
}

func BenchmarkSumKoopman(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		zcrc32.SumKoopman(benchData)
	}
}
