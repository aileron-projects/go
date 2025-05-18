package zmath_test

import (
	"testing"

	"github.com/aileron-projects/go/zmath"
	"github.com/aileron-projects/go/ztesting"
)

func TestNewSet(t *testing.T) {
	t.Parallel()
	t.Run("empty slice", func(t *testing.T) {
		s := []string{}
		ss := zmath.NewSet(s)
		ztesting.AssertEqualSlice(t, "set not matched", []string{}, ss)
	})
	t.Run("1 elem slice", func(t *testing.T) {
		s := []string{"a"}
		ss := zmath.NewSet(s)
		ztesting.AssertEqualSlice(t, "set not matched", []string{"a"}, ss)
	})
	t.Run("2 elems slice", func(t *testing.T) {
		s := []string{"a", "b"}
		ss := zmath.NewSet(s)
		ztesting.AssertEqualSlice(t, "set not matched", []string{"a", "b"}, ss)
	})
	t.Run("string slice", func(t *testing.T) {
		s := []string{"c", "a", "b"}
		ss := zmath.NewSet(s)
		ztesting.AssertEqualSlice(t, "set not matched", []string{"a", "b", "c"}, ss)
	})
	t.Run("int slice", func(t *testing.T) {
		s := []int{3, -1, 2, 0, 1}
		ss := zmath.NewSet(s)
		ztesting.AssertEqualSlice(t, "set not matched", []int{-1, 0, 1, 2, 3}, ss)
	})
}

func TestSet_Has(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		set  zmath.Set[int]
		elem int
		want bool
	}{
		"empty set":               {zmath.NewSet([]int{}), 1, false},
		"1 elem, not contain":     {zmath.NewSet([]int{2}), 1, false},
		"1 elem, contain":         {zmath.NewSet([]int{1}), 1, true},
		"2 elems, not contain 01": {zmath.NewSet([]int{1, 3}), 0, false},
		"2 elems, not contain 02": {zmath.NewSet([]int{1, 3}), 2, false},
		"2 elems, not contain 03": {zmath.NewSet([]int{1, 3}), 4, false},
		"2 elems, contain 01":     {zmath.NewSet([]int{1, 3}), 1, true},
		"2 elems, contain 02":     {zmath.NewSet([]int{1, 3}), 3, true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.set.Has(tc.elem)
			ztesting.AssertEqual(t, "has not matched", tc.want, got)
		})
	}
}

func TestSet_Add(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		set  zmath.Set[int]
		elem int
		want zmath.Set[int]
	}{
		"empty set":                 {zmath.NewSet([]int{}), 1, zmath.NewSet([]int{1})},
		"1 elem, add at first":      {zmath.NewSet([]int{3}), 1, zmath.NewSet([]int{1, 3})},
		"1 elem, add same elem":     {zmath.NewSet([]int{3}), 3, zmath.NewSet([]int{3, 3})},
		"1 elem, add at last":       {zmath.NewSet([]int{3}), 4, zmath.NewSet([]int{3, 4})},
		"2 elems, add at first":     {zmath.NewSet([]int{3, 5}), 1, zmath.NewSet([]int{1, 3, 5})},
		"2 elems, add same elem 01": {zmath.NewSet([]int{3, 5}), 3, zmath.NewSet([]int{3, 3, 5})},
		"2 elems, add same elem 02": {zmath.NewSet([]int{3, 5}), 5, zmath.NewSet([]int{3, 5, 5})},
		"2 elems, add at last":      {zmath.NewSet([]int{3, 5}), 6, zmath.NewSet([]int{3, 5, 6})},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.set.Add(tc.elem)
			ztesting.AssertEqualSlice(t, "has not matched", tc.want, tc.set)
		})
	}
}

func TestSet_AddElems(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		set   zmath.Set[int]
		elems []int
		want  zmath.Set[int]
	}{
		"empty set 01":          {zmath.NewSet([]int{}), []int{1}, zmath.NewSet([]int{1})},
		"empty set 02":          {zmath.NewSet([]int{}), []int{2, 1}, zmath.NewSet([]int{1, 2})},
		"1 elem, add before 01": {zmath.NewSet([]int{3}), []int{1}, zmath.NewSet([]int{1, 3})},
		"1 elem, add before 02": {zmath.NewSet([]int{3}), []int{1, 0}, zmath.NewSet([]int{0, 1, 3})},
		"1 elem, add same":      {zmath.NewSet([]int{3}), []int{3, 3}, zmath.NewSet([]int{3, 3, 3})},
		"1 elem, add after 01":  {zmath.NewSet([]int{3}), []int{5}, zmath.NewSet([]int{3, 5})},
		"1 elem, add after 02":  {zmath.NewSet([]int{3}), []int{5, 8}, zmath.NewSet([]int{3, 5, 8})},
		"2 elem, add before 01": {zmath.NewSet([]int{3, 6}), []int{1}, zmath.NewSet([]int{1, 3, 6})},
		"2 elem, add before 02": {zmath.NewSet([]int{3, 6}), []int{1, 0}, zmath.NewSet([]int{0, 1, 3, 6})},
		"2 elem, add same":      {zmath.NewSet([]int{3, 6}), []int{3, 3}, zmath.NewSet([]int{3, 3, 3, 6})},
		"2 elem, add after 01":  {zmath.NewSet([]int{3, 6}), []int{5}, zmath.NewSet([]int{3, 5, 6})},
		"2 elem, add after 02":  {zmath.NewSet([]int{3, 6}), []int{5, 8}, zmath.NewSet([]int{3, 5, 6, 8})},
		"2 elem, add after 03":  {zmath.NewSet([]int{3, 6}), []int{8, 9}, zmath.NewSet([]int{3, 6, 8, 9})},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.set.AddElems(tc.elems...)
			ztesting.AssertEqualSlice(t, "has not matched", tc.want, tc.set)
		})
	}
}

func TestSet_Remove(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		set  zmath.Set[int]
		elem int
		want zmath.Set[int]
	}{
		"empty set":                {zmath.NewSet([]int{}), 1, zmath.NewSet([]int{})},
		"1 elem, not found":        {zmath.NewSet([]int{3}), 1, zmath.NewSet([]int{3})},
		"1 elem, found":            {zmath.NewSet([]int{3}), 3, zmath.NewSet([]int{})},
		"2 elems, found at first":  {zmath.NewSet([]int{3, 5}), 3, zmath.NewSet([]int{5})},
		"2 elems, found at last":   {zmath.NewSet([]int{3, 5}), 5, zmath.NewSet([]int{3})},
		"3 elems, found at first":  {zmath.NewSet([]int{3, 5, 8}), 3, zmath.NewSet([]int{5, 8})},
		"3 elems, found at middle": {zmath.NewSet([]int{3, 5, 8}), 5, zmath.NewSet([]int{3, 8})},
		"3 elems, found at last":   {zmath.NewSet([]int{3, 5, 8}), 8, zmath.NewSet([]int{3, 5})},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.set.Remove(tc.elem)
			ztesting.AssertEqualSlice(t, "has not matched", tc.want, tc.set)
		})
	}
}

func TestSet_RemoveElems(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		set   zmath.Set[int]
		elems []int
		want  zmath.Set[int]
	}{
		"empty set 01":          {zmath.NewSet([]int{}), []int{1}, zmath.NewSet([]int{})},
		"empty set 02":          {zmath.NewSet([]int{}), []int{2, 1}, zmath.NewSet([]int{})},
		"1 elem, not found 01":  {zmath.NewSet([]int{3}), []int{1}, zmath.NewSet([]int{3})},
		"1 elem, not found 02":  {zmath.NewSet([]int{3}), []int{5}, zmath.NewSet([]int{3})},
		"1 elem, found":         {zmath.NewSet([]int{3}), []int{3}, zmath.NewSet([]int{})},
		"2 elems, not found 01": {zmath.NewSet([]int{3, 5}), []int{1, 0}, zmath.NewSet([]int{3, 5})},
		"2 elems, not found 02": {zmath.NewSet([]int{3, 5}), []int{4, 6}, zmath.NewSet([]int{3, 5})},
		"2 elems, not found 03": {zmath.NewSet([]int{3, 5}), []int{7, 8}, zmath.NewSet([]int{3, 5})},
		"2 elems, found 01":     {zmath.NewSet([]int{3, 5}), []int{3}, zmath.NewSet([]int{5})},
		"2 elems, found 02":     {zmath.NewSet([]int{3, 5}), []int{5}, zmath.NewSet([]int{3})},
		"2 elems, found 03":     {zmath.NewSet([]int{3, 5}), []int{1, 3}, zmath.NewSet([]int{5})},
		"2 elems, found 04":     {zmath.NewSet([]int{3, 5}), []int{1, 5}, zmath.NewSet([]int{3})},
		"2 elems, found 05":     {zmath.NewSet([]int{3, 5}), []int{3, 5}, zmath.NewSet([]int{})},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.set.RemoveElems(tc.elems...)
			ztesting.AssertEqualSlice(t, "has not matched", tc.want, tc.set)
		})
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		setA, setB []int
		equal      bool
	}{
		"case01": {[]int{}, []int{}, true},
		"case02": {[]int{1}, []int{}, false},
		"case03": {[]int{}, []int{1}, false},
		"case04": {[]int{1}, []int{1}, true},
		"case05": {[]int{1}, []int{2}, false},
		"case06": {[]int{1, 1}, []int{1}, true},
		"case07": {[]int{1}, []int{1, 1}, true},
		"case08": {[]int{1, 1, 2}, []int{1, 2}, true},
		"case09": {[]int{1, 1, 2, 3}, []int{1, 2, 3}, true},
		"case10": {[]int{1, 1, 2, 3}, []int{2, 3}, false},
		"case11": {[]int{1, 1, 2, 3}, []int{1, 2, 3, 4}, false},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			a := zmath.NewSet(tc.setA)
			b := zmath.NewSet(tc.setB)
			equal := zmath.Equal(a, b)
			ztesting.AssertEqual(t, "equal result not match", tc.equal, equal)
		})
	}
}

func Test_Subset_SuperSet(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		setA, setB []int
		subset     bool // setA is subset of setB ?
	}{
		"case01": {[]int{}, []int{}, true},
		"case02": {[]int{1}, []int{}, false},
		"case03": {[]int{}, []int{1}, true},
		"case04": {[]int{1}, []int{0}, false},
		"case05": {[]int{1}, []int{1}, true},
		"case06": {[]int{1}, []int{2}, false},
		"case07": {[]int{1}, []int{0, 1}, true},
		"case08": {[]int{1}, []int{1, 2}, true},
		"case09": {[]int{1}, []int{2, 3}, false},
		"case10": {[]int{1, 2}, []int{0}, false},
		"case11": {[]int{1, 2}, []int{1}, false},
		"case12": {[]int{1, 2}, []int{2}, false},
		"case13": {[]int{1, 2}, []int{3}, false},
		"case14": {[]int{1, 2}, []int{0, 1}, false},
		"case15": {[]int{1, 2}, []int{1, 2}, true},
		"case16": {[]int{1, 2}, []int{2, 3}, false},
		"case17": {[]int{1, 2}, []int{3, 4}, false},
		"case18": {[]int{1, 2}, []int{0, 1, 2, 3}, true},
		"case19": {[]int{1, 2}, []int{0, 1, 2, 4}, true},
		"case20": {[]int{1, 2, 3}, []int{0}, false},
		"case21": {[]int{1, 2, 3}, []int{1}, false},
		"case22": {[]int{1, 2, 3}, []int{2}, false},
		"case23": {[]int{1, 2, 3}, []int{0, 1}, false},
		"case24": {[]int{1, 2, 3}, []int{0, 3}, false},
		"case25": {[]int{1, 2, 3}, []int{1, 2}, false},
		"case26": {[]int{1, 2, 3}, []int{1, 2, 3}, true},
		"case27": {[]int{1, 2, 3}, []int{0, 1, 2, 3}, true},
		"case28": {[]int{1, 2, 3}, []int{0, 1, 2, 3, 4}, true},
		"case29": {[]int{1, 2, 2, 2, 3}, []int{1, 2, 3}, true},
		"case30": {[]int{1, 2, 3}, []int{1, 2, 2, 2, 3, 4}, true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			a := zmath.NewSet(tc.setA)
			b := zmath.NewSet(tc.setB)
			subSet := zmath.Subset(a, b)
			superSet := zmath.Superset(b, a)
			ztesting.AssertEqual(t, "subset result not match", tc.subset, subSet)
			ztesting.AssertEqual(t, "superset result not match", tc.subset, superSet)
		})
	}
}

func Test_ProperSubset_ProperSuperset(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		setA, setB []int
		subset     bool // setA is proper subset of setB ?
	}{
		"case01": {[]int{}, []int{}, false},
		"case02": {[]int{1}, []int{}, false},
		"case03": {[]int{}, []int{1}, true},
		"case04": {[]int{1}, []int{0}, false},
		"case05": {[]int{1}, []int{1}, false},
		"case06": {[]int{1}, []int{2}, false},
		"case07": {[]int{1}, []int{0, 1}, true},
		"case08": {[]int{1}, []int{1, 1}, false},
		"case09": {[]int{1}, []int{1, 2}, true},
		"case10": {[]int{1}, []int{2, 3}, false},
		"case11": {[]int{1, 2}, []int{0}, false},
		"case12": {[]int{1, 2}, []int{1}, false},
		"case13": {[]int{1, 2}, []int{2}, false},
		"case14": {[]int{1, 2}, []int{3}, false},
		"case15": {[]int{1, 2}, []int{0, 1}, false},
		"case16": {[]int{1, 2}, []int{1, 2}, false},
		"case17": {[]int{1, 2}, []int{1, 3}, false},
		"case18": {[]int{1, 2}, []int{0, 1, 2}, true},
		"case19": {[]int{1, 2}, []int{1, 2, 3}, true},
		"case20": {[]int{1, 2}, []int{0, 1, 2, 3}, true},
		"case21": {[]int{1, 2, 3}, []int{0}, false},
		"case22": {[]int{1, 2, 3}, []int{1}, false},
		"case23": {[]int{1, 2, 3}, []int{4}, false},
		"case24": {[]int{1, 2, 3}, []int{0, 1, 2}, false},
		"case25": {[]int{1, 2, 3}, []int{1, 2, 3}, false},
		"case26": {[]int{1, 2, 3}, []int{0, 1, 2, 3}, true},
		"case27": {[]int{1, 2, 3}, []int{1, 2, 3, 4}, true},
		"case28": {[]int{1, 2, 3}, []int{0, 1, 2, 3, 4}, true},
		"case29": {[]int{1, 2, 2, 2, 3}, []int{0, 1, 2, 3, 4}, true},
		"case30": {[]int{1, 2, 3}, []int{0, 1, 2, 2, 2, 3, 4}, true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			a := zmath.NewSet(tc.setA)
			b := zmath.NewSet(tc.setB)
			subset := zmath.ProperSubset(a, b)
			superset := zmath.ProperSuperset(b, a)
			ztesting.AssertEqual(t, "proper subset result not match", tc.subset, subset)
			ztesting.AssertEqual(t, "proper superset result not match", tc.subset, superset)
		})
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		setA, setB []int
		union      []int // setA is proper subset of setB ?
	}{
		"case01": {[]int{}, []int{}, []int{}},
		"case02": {[]int{1}, []int{}, []int{1}},
		"case03": {[]int{}, []int{1}, []int{1}},
		"case04": {[]int{1, 2}, []int{}, []int{1, 2}},
		"case05": {[]int{}, []int{1, 2}, []int{1, 2}},
		"case06": {[]int{1}, []int{2}, []int{1, 2}},
		"case07": {[]int{1}, []int{0, 1}, []int{0, 1}},
		"case08": {[]int{1}, []int{1, 1}, []int{1}},
		"case09": {[]int{1}, []int{1, 2}, []int{1, 2}},
		"case10": {[]int{1, 3}, []int{0}, []int{0, 1, 3}},
		"case11": {[]int{1, 3}, []int{1}, []int{1, 3}},
		"case12": {[]int{1, 3}, []int{2}, []int{1, 2, 3}},
		"case13": {[]int{1, 3}, []int{3}, []int{1, 3}},
		"case14": {[]int{1, 3}, []int{4}, []int{1, 3, 4}},
		"case15": {[]int{1, 3}, []int{0, 2, 4}, []int{0, 1, 2, 3, 4}},
		"case16": {[]int{1, 3, 5}, []int{0, 1}, []int{0, 1, 3, 5}},
		"case17": {[]int{1, 3, 5}, []int{5, 6}, []int{1, 3, 5, 6}},
		"case18": {[]int{1, 3, 5}, []int{0, 2, 4}, []int{0, 1, 2, 3, 4, 5}},
		"case19": {[]int{1, 3, 3, 3, 5}, []int{0, 2, 4}, []int{0, 1, 2, 3, 4, 5}},
		"case20": {[]int{1, 3, 5}, []int{0, 2, 2, 2, 4}, []int{0, 1, 2, 3, 4, 5}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			a := zmath.NewSet(tc.setA)
			b := zmath.NewSet(tc.setB)
			union := zmath.Union(a, b)
			ztesting.AssertEqualSlice(t, "unexpected union", tc.union, union)
		})
	}
}

func TestIntersection(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		setA, setB   []int
		intersection []int
	}{
		"case01": {[]int{}, []int{}, []int{}},
		"case02": {[]int{1}, []int{}, []int{}},
		"case03": {[]int{}, []int{1}, []int{}},
		"case04": {[]int{1, 2}, []int{}, []int{}},
		"case05": {[]int{}, []int{1, 2}, []int{}},
		"case06": {[]int{1}, []int{2}, []int{}},
		"case07": {[]int{1}, []int{0, 1}, []int{1}},
		"case08": {[]int{1}, []int{1, 1}, []int{1}},
		"case09": {[]int{1}, []int{1, 2}, []int{1}},
		"case10": {[]int{1, 3}, []int{0}, []int{}},
		"case11": {[]int{1, 3}, []int{1}, []int{1}},
		"case12": {[]int{1, 3}, []int{2}, []int{}},
		"case13": {[]int{1, 3}, []int{3}, []int{3}},
		"case14": {[]int{1, 3}, []int{4}, []int{}},
		"case15": {[]int{1, 3}, []int{1, 3}, []int{1, 3}},
		"case16": {[]int{1, 3}, []int{0, 1}, []int{1}},
		"case17": {[]int{1, 3}, []int{3, 4}, []int{3}},
		"case18": {[]int{1, 3, 5}, []int{1, 3, 5}, []int{1, 3, 5}},
		"case19": {[]int{1, 3, 5}, []int{1, 3}, []int{1, 3}},
		"case20": {[]int{1, 3, 5}, []int{3, 5}, []int{3, 5}},
		"case21": {[]int{1, 3, 3, 3, 5}, []int{3, 5}, []int{3, 5}},
		"case22": {[]int{1, 3, 5}, []int{3, 3, 3, 5}, []int{3, 5}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			a := zmath.NewSet(tc.setA)
			b := zmath.NewSet(tc.setB)
			intersection := zmath.Intersection(a, b)
			ztesting.AssertEqualSlice(t, "unexpected intersection", tc.intersection, intersection)
		})
	}
}

func TestDifference(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		setA, setB []int
		difference []int
	}{
		"case01": {[]int{}, []int{}, []int{}},
		"case02": {[]int{1}, []int{}, []int{1}},
		"case03": {[]int{}, []int{1}, []int{}},
		"case04": {[]int{1, 2}, []int{}, []int{1, 2}},
		"case05": {[]int{}, []int{1, 2}, []int{}},
		"case06": {[]int{1}, []int{2}, []int{1}},
		"case07": {[]int{1}, []int{0, 1}, []int{}},
		"case08": {[]int{1}, []int{1, 1}, []int{}},
		"case09": {[]int{1}, []int{1, 2}, []int{}},
		"case10": {[]int{1, 3}, []int{0}, []int{1, 3}},
		"case11": {[]int{1, 3}, []int{1}, []int{3}},
		"case12": {[]int{1, 3}, []int{2}, []int{1, 3}},
		"case13": {[]int{1, 3}, []int{3}, []int{1}},
		"case14": {[]int{1, 3}, []int{4}, []int{1, 3}},
		"case15": {[]int{1, 3}, []int{1, 3}, []int{}},
		"case16": {[]int{1, 3}, []int{0, 1}, []int{3}},
		"case17": {[]int{1, 3}, []int{3, 4}, []int{1}},
		"case18": {[]int{1, 3, 5}, []int{1, 3, 5}, []int{}},
		"case19": {[]int{1, 3, 5}, []int{1, 3}, []int{5}},
		"case20": {[]int{1, 3, 5}, []int{3, 5}, []int{1}},
		"case21": {[]int{1, 3, 3, 3, 5}, []int{3, 5}, []int{1}},
		"case22": {[]int{1, 3, 5}, []int{3, 3, 3, 5}, []int{1}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			a := zmath.NewSet(tc.setA)
			b := zmath.NewSet(tc.setB)
			difference := zmath.Difference(a, b)
			ztesting.AssertEqualSlice(t, "unexpected difference", tc.difference, difference)
		})
	}
}
