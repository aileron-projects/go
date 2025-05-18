package zbackoff_test

import (
	"testing"
	"time"

	"github.com/aileron-projects/go/ztime/zbackoff"
)

func BenchmarkFixedBackoff(b *testing.B) {
	backoff := zbackoff.NewFixedBackoff(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backoff.Attempt(10)
	}
}

func BenchmarkRandomBackoff(b *testing.B) {
	backoff := zbackoff.NewRandomBackoff(time.Second, 10*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backoff.Attempt(10)
	}
}

func BenchmarkLinearBackoff(b *testing.B) {
	backoff := zbackoff.NewLinearBackoff(time.Second, 10*time.Second, 100*time.Millisecond)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backoff.Attempt(10)
	}
}

func BenchmarkPolynomialBackoff(b *testing.B) {
	backoff := zbackoff.NewPolynomialBackoff(time.Second, 10*time.Second, 10*time.Millisecond, 3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backoff.Attempt(10)
	}
}

func BenchmarkExponentialBackoff(b *testing.B) {
	backoff := zbackoff.NewExponentialBackoff(time.Second, 10*time.Second, 10*time.Millisecond)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backoff.Attempt(10)
	}
}

func BenchmarkExponentialBackoffFullJitter(b *testing.B) {
	backoff := zbackoff.NewExponentialBackoffFullJitter(time.Second, 10*time.Second, 10*time.Millisecond)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backoff.Attempt(10)
	}
}

func BenchmarkExponentialBackoffEqualJitter(b *testing.B) {
	backoff := zbackoff.NewExponentialBackoffEqualJitter(time.Second, 10*time.Second, 10*time.Millisecond)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backoff.Attempt(10)
	}
}

func BenchmarkFibonacciBackoff(b *testing.B) {
	backoff := zbackoff.NewFibonacciBackoff(time.Second, 10*time.Second, 10*time.Millisecond)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backoff.Attempt(10)
	}
}
