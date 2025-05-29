package zuid_test

import (
	"testing"

	"github.com/aileron-projects/go/zx/zuid"
)

func BenchmarkNewTime(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		zuid.NewTime()
	}
}

func BenchmarkNewHost(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		zuid.NewHost()
	}
}

func BenchmarkNewCount(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		zuid.NewCount()
	}
}
