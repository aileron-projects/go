package zsha3_test

import (
	"testing"

	"github.com/aileron-projects/go/zcrypto/zsha3"
)

var (
	msg = []byte("Go is an open source programming language that makes it simple to build secure, scalable systems.")
	key = []byte("12345678901234567890123456789012")
)

func BenchmarkSum224(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha3.Sum224(msg)
	}
}

func BenchmarkSum256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha3.Sum256(msg)
	}
}

func BenchmarkSum384(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha3.Sum384(msg)
	}
}

func BenchmarkSum512(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha3.Sum512(msg)
	}
}

func BenchmarkSumShake128(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha3.SumShake128(msg)
	}
}

func BenchmarkSumShake256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha3.SumShake256(msg)
	}
}

func BenchmarkHMACSum224(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha3.HMACSum224(msg, key)
	}
}

func BenchmarkHMACSum256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha3.HMACSum256(msg, key)
	}
}

func BenchmarkHMACSum384(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha3.HMACSum384(msg, key)
	}
}

func BenchmarkHMACSum512(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zsha3.HMACSum512(msg, key)
	}
}
