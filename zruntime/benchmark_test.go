package zruntime_test

import (
	"testing"

	"github.com/aileron-projects/go/zruntime"
)

func BenchmarkCallerFrame(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = zruntime.CallerFrame(0)
	}
}

func BenchmarkCallerFrames(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = zruntime.CallerFrames(0)
	}
}

func BenchmarkConvertFrame(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = zruntime.ConvertFrame(zruntime.CallerFrame(0))
	}
}

func BenchmarkConvertFrames(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = zruntime.ConvertFrames(zruntime.CallerFrames(0))
	}
}
