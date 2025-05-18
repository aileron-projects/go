package internal_test

import (
	"slices"
	"testing"

	"github.com/aileron-projects/go/znet/internal"
	"github.com/aileron-projects/go/ztesting"
)

func TestUniqueStore(t *testing.T) {
	t.Parallel()
	t.Run("iterate all", func(t *testing.T) {
		s := internal.UniqueStore[string]{}
		s.Set("foo")
		s.Set("bar")
		s.Set("baz")
		want := []string{"bar", "baz", "foo"}
		got := []string{}
		for v := range s.Values() {
			got = append(got, v)
		}
		slices.Sort(got) // Compare sorted slice.
		ztesting.AssertEqualSlice(t, "values not match", want, got)
		ztesting.AssertEqual(t, "length not match", 3, s.Length())
	})
	t.Run("delete while iterate", func(t *testing.T) {
		s := internal.UniqueStore[string]{}
		s.Set("foo")
		s.Set("bar")
		s.Set("baz")
		got := []string{}
		for v := range s.Values() {
			s.Delete("foo") // Delete all.
			s.Delete("bar") // Delete all.
			s.Delete("baz") // Delete all.
			got = append(got, v)
		}
		ztesting.AssertEqual(t, "number of values not match", 1, len(got))
		ztesting.AssertEqual(t, "length not match", 0, s.Length())
	})
	t.Run("cancel", func(t *testing.T) {
		s := internal.UniqueStore[string]{}
		s.Set("foo")
		s.Set("bar")
		s.Set("baz")
		got := []string{}
		for v := range s.Values() {
			if v != "" {
				break
			}
			got = append(got, v)
		}
		ztesting.AssertEqualSlice(t, "values not match", []string{}, got)
		ztesting.AssertEqual(t, "length not match", 3, s.Length())
	})
}
