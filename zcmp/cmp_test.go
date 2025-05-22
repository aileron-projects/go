package zcmp_test

import (
	"testing"

	"github.com/aileron-projects/go/zcmp"
	"github.com/aileron-projects/go/ztesting"
)

func TestTrue(t *testing.T) {
	t.Parallel()

	t.Run("test int", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			v := zcmp.True(true, 1, 2)
			ztesting.AssertEqual(t, "invalid response value.", 1, v)
		})
		t.Run("false", func(t *testing.T) {
			v := zcmp.True(false, 1, 2)
			ztesting.AssertEqual(t, "invalid response value.", 2, v)
		})
	})

	t.Run("test nil", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			var np *struct{}
			v := zcmp.True(true, nil, np)
			ztesting.AssertEqual(t, "invalid response value.", nil, v)
		})
		t.Run("false", func(t *testing.T) {
			var np *struct{}
			v := zcmp.True(false, nil, np)
			ztesting.AssertEqual(t, "invalid response value.", np, v)
		})
	})
}

func TestFalse(t *testing.T) {
	t.Parallel()

	t.Run("test int", func(t *testing.T) {
		t.Run("false", func(t *testing.T) {
			v := zcmp.False(false, 1, 2)
			ztesting.AssertEqual(t, "invalid response value.", 1, v)
		})
		t.Run("true", func(t *testing.T) {
			v := zcmp.False(true, 1, 2)
			ztesting.AssertEqual(t, "invalid response value.", 2, v)
		})
	})

	t.Run("test nil", func(t *testing.T) {
		t.Run("false", func(t *testing.T) {
			var np *struct{}
			v := zcmp.False(false, nil, np)
			ztesting.AssertEqual(t, "invalid response value.", nil, v)
		})
		t.Run("true", func(t *testing.T) {
			var np *struct{}
			v := zcmp.False(true, nil, np)
			ztesting.AssertEqual(t, "invalid response value.", np, v)
		})
	})
}

func TestOrSlice(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		vals [][]int
		want []int
	}{
		"no slice": {
			vals: [][]int{},
			want: nil,
		},
		"nil slice": {
			vals: [][]int{{}},
			want: nil,
		},
		"non-nil slice": {
			vals: [][]int{{1}},
			want: []int{1},
		},
		"non-nil at 0": {
			vals: [][]int{{1}, {}},
			want: []int{1},
		},
		"non-nil at 1": {
			vals: [][]int{{}, {1}},
			want: []int{1},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v := zcmp.OrSlice(tc.vals...)
			ztesting.AssertEqual(t, "wrong element returned.", tc.want, v)
		})
	}
}

func TestOrMap(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		vals []map[int]string
		want map[int]string
	}{
		"no map": {
			vals: []map[int]string{},
			want: nil,
		},
		"nil map": {
			vals: []map[int]string{{}},
			want: nil,
		},
		"non-nil map": {
			vals: []map[int]string{{1: "1"}},
			want: map[int]string{1: "1"},
		},
		"non-nil at 0": {
			vals: []map[int]string{{1: "1"}, {}},
			want: map[int]string{1: "1"},
		},
		"non-nil at 1": {
			vals: []map[int]string{{}, {1: "1"}},
			want: map[int]string{1: "1"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v := zcmp.OrMap(tc.vals...)
			ztesting.AssertEqual(t, "wrong element returned.", tc.want, v)
		})
	}
}
