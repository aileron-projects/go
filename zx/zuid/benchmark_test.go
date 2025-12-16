package zuid_test

import (
	"testing"

	"github.com/aileron-projects/go/zx/zuid"
)

func BenchmarkNewTime(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zuid.NewTime()
	}
}

func BenchmarkNewHost(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zuid.NewHost()
	}
}

func BenchmarkNewCount(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zuid.NewCount()
	}
}
