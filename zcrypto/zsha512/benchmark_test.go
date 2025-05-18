package zsha512_test

import (
	"testing"

	"github.com/aileron-projects/go/zcrypto/zsha512"
)

var (
	msg = []byte("Go is an open source programming language that makes it simple to build secure, scalable systems.")
	key = []byte("12345678901234567890123456789012")
)

func BenchmarkSum224(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha512.Sum224(msg)
	}
}

func BenchmarkSum256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha512.Sum256(msg)
	}
}

func BenchmarkSum384(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha512.Sum384(msg)
	}
}

func BenchmarkSum512(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha512.Sum512(msg)
	}
}

func BenchmarkHMACSum224(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha512.HMACSum224(msg, key)
	}
}

func BenchmarkHMACSum256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha512.HMACSum256(msg, key)
	}
}

func BenchmarkHMACSum384(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha512.HMACSum384(msg, key)
	}
}

func BenchmarkHMACSum512(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha512.HMACSum512(msg, key)
	}
}
