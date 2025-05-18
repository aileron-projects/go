package zmd5_test

import (
	"testing"

	"github.com/aileron-projects/go/zcrypto/zmd5"
)

var (
	msg = []byte("Go is an open source programming language that makes it simple to build secure, scalable systems.")
	key = []byte("12345678901234567890123456789012")
)

func BenchmarkSum(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zmd5.Sum(msg)
	}
}

func BenchmarkHMACSum(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zmd5.HMACSum(msg, key)
	}
}
