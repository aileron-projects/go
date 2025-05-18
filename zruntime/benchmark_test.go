package zruntime_test

import (
	"testing"

	"github.com/aileron-projects/go/zruntime"
)

func BenchmarkCallerFrame(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = zruntime.CallerFrame(0)
	}
}

func BenchmarkCallerFrames(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = zruntime.CallerFrames(0)
	}
}

func BenchmarkConvertFrame(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = zruntime.ConvertFrame(zruntime.CallerFrame(0))
	}
}

func BenchmarkConvertFrames(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = zruntime.ConvertFrames(zruntime.CallerFrames(0))
	}
}
