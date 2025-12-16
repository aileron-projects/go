package zblake2b_test

import (
	"testing"

	"github.com/aileron-projects/go/zcrypto/zblake2b"
)

var (
	msg = []byte("Go is an open source programming language that makes it simple to build secure, scalable systems.")
	key = []byte("12345678901234567890123456789012")
)

func BenchmarkSum256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zblake2b.Sum256(msg)
	}
}

func BenchmarkSum384(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zblake2b.Sum384(msg)
	}
}

func BenchmarkSum512(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zblake2b.Sum512(msg)
	}
}

func BenchmarkHMACSum256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zblake2b.HMACSum256(msg, key)
	}
}

func BenchmarkHMACSum384(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zblake2b.HMACSum384(msg, key)
	}
}

func BenchmarkHMACSum512(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zblake2b.HMACSum512(msg, key)
	}
}
