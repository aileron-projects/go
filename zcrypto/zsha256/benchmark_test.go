package zsha256_test

import (
	"testing"

	"github.com/aileron-projects/go/zcrypto/zsha256"
)

var (
	msg = []byte("Go is an open source programming language that makes it simple to build secure, scalable systems.")
	key = []byte("12345678901234567890123456789012")
)

func BenchmarkSum224(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha256.Sum224(msg)
	}
}

func BenchmarkSum256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha256.Sum256(msg)
	}
}

func BenchmarkHMACSum224(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha256.HMACSum224(msg, key)
	}
}

func BenchmarkHMACSum256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha256.HMACSum256(msg, key)
	}
}
