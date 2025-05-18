package zmaps_test

import (
	"slices"
	"testing"

	"github.com/aileron-projects/go/zmaps"
	"github.com/aileron-projects/go/ztesting"
)

func TestKeys(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input map[string]int
		want  []string
	}{
		"nil map":   {nil, nil},
		"empty map": {map[string]int{}, nil},
		"1 key":     {map[string]int{"a": 1}, []string{"a"}},
		"2 keys":    {map[string]int{"a": 1, "b": 2}, []string{"a", "b"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v := zmaps.Keys(tc.input)
			slices.Sort(v)
			ztesting.AssertEqualSlice(t, "wrong element returned.", tc.want, v)
		})
	}
}

func TestValues(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input map[string]int
		want  []int
	}{
		"nil map":   {nil, nil},
		"empty map": {map[string]int{}, nil},
		"1 key":     {map[string]int{"a": 1}, []int{1}},
		"2 keys":    {map[string]int{"a": 1, "b": 2}, []int{1, 2}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v := zmaps.Values(tc.input)
			slices.Sort(v)
			ztesting.AssertEqualSlice(t, "wrong element returned.", tc.want, v)
		})
	}
}
