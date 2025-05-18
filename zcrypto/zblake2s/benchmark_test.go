package zblake2s_test

import (
	"testing"

	"github.com/aileron-projects/go/zcrypto/zblake2s"
)

var (
	msg = []byte("Go is an open source programming language that makes it simple to build secure, scalable systems.")
	key = []byte("12345678901234567890123456789012")
)

func BenchmarkSum256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zblake2s.Sum256(msg)
	}
}

func BenchmarkHMACSum256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zblake2s.HMACSum256(msg, key)
	}
}
