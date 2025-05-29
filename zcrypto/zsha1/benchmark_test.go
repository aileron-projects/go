package zsha1_test

import (
	"testing"

	"github.com/aileron-projects/go/zcrypto/zsha1"
)

var (
	msg = []byte("Go is an open source programming language that makes it simple to build secure, scalable systems.")
	key = []byte("12345678901234567890123456789012")
)

func BenchmarkSum(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		zsha1.Sum(msg)
	}
}

func BenchmarkHMACSum(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		zsha1.HMACSum(msg, key)
	}
}
