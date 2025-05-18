package zuid_test

import (
	"testing"

	"github.com/aileron-projects/go/zx/zuid"
)

func BenchmarkNewTimeBase(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zuid.NewTimeBase()
	}
}

func BenchmarkNewHostBase(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zuid.NewHostBase()
	}
}

func BenchmarkNewCountBase(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zuid.NewCountBase()
	}
}
