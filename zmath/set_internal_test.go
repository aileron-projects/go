package zmath

import (
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestUniqueIterator(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		set   Set[int]
		elems []int
	}{
		"case01": {NewSet([]int{}), []int{}},
		"case02": {NewSet([]int{1}), []int{1}},
		"case03": {NewSet([]int{1, 1}), []int{1}},
		"case04": {NewSet([]int{1, 1, 2}), []int{1, 2}},
		"case05": {NewSet([]int{1, 1, 2, 2}), []int{1, 2}},
		"case06": {NewSet([]int{1, 1, 2, 2, 2}), []int{1, 2}},
		"case07": {NewSet([]int{1, 1, 2, 2, 2, 3}), []int{1, 2, 3}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			iter := uniqueIterator(tc.set)
			for _, v := range tc.elems {
				next, found := iter()
				ztesting.AssertEqual(t, "value not match", v, next)
				ztesting.AssertEqual(t, "value not found", true, found)
			}
			_, found := iter()
			ztesting.AssertEqual(t, "value found", false, found)
		})
	}
}
