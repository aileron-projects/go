package zmath_test

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/aileron-projects/go/zmath"
)

func FuzzAddElems(f *testing.F) {
	f.Fuzz(func(t *testing.T, n uint16) {
		arr1 := make([]int, 2*n)
		arr2 := make([]int, n)
		want := make([]int, 3*n)
		for i := 0; i < int(n); i++ {
			arr1[i], arr2[i] = rand.IntN(1000), rand.IntN(1000)
			want[i], want[2*int(n)+i] = arr1[i], arr2[i]
		}
		slices.Sort(want)
		set := zmath.NewSet(arr1)
		set.AddElems(arr2...)
		if !slices.Equal(set, want) {
			t.Logf("want: %v", want)
			t.Logf("set : %v", set)
			t.Error("Set not match")
		}
	})
}
