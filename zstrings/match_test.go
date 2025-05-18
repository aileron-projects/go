package zstrings

import (
	"strconv"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestMatch_error(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		pattern string
		str     string
	}{
		{"\\", ""},
		{"\\\\\\", ""},
		{"x\\", ""},
		{"x\\\\\\", ""},
		{"x*\\", ""},
		{"x\\*\\", ""},
		{"\\", "\\"},
		{"\\\\\\", "\\\\\\"},
		{"x\\", "x\\"},
		{"x\\\\\\", "x\\\\\\"},
		{"x*\\", "x\\"},
		{"x\\*\\", "x*\\"},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			matched, err := Match(tc.pattern, tc.str)
			ztesting.AssertEqual(t, "unexpectedly mismatched.", false, matched)
			ztesting.AssertEqual(t, "unexpectedly error.", ErrBadPattern, err)
		})
	}
}

func TestMatch(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		pattern string
		str     string
		matched bool
		err     error
	}{
		{"", "", true, nil},
		{"", "x", false, nil},
		{"x", "", false, nil},
		{"*", "", true, nil},
		{"?", "", false, nil},
		{"*", "x", true, nil},
		{"?", "x", true, nil},
		{"*", "foo", true, nil},
		{"?", "foo", false, nil},
		{"**", "foo", true, nil},
		{"***", "foo", true, nil},
		{"****", "foo", true, nil},
		{"??", "foo", false, nil},
		{"???", "foo", true, nil},
		{"*foo", "foo", true, nil},
		{"?foo", "foo", false, nil},
		{"*foo", "_foo", true, nil},
		{"?foo", "_foo", true, nil},
		{"foo*", "foo", true, nil},
		{"foo?", "foo", false, nil},
		{"foo*", "foo_", true, nil},
		{"foo?", "foo_", true, nil},
		{"f*o", "foo", true, nil},
		{"f?o", "foo", true, nil},
		{"*.bar", "foo.bar", true, nil},
		{"f*.bar*", "foo.bar", true, nil},
		{"f*.ba?*", "foo.bar", true, nil},
		{"f**.bar", "foo.bar.bar", true, nil},
		{"f**?bar", "foo.bar.bar", true, nil},
		{"f**???", "foo.bar.bar", true, nil},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			matched, err := Match(tc.pattern, tc.str)
			ztesting.AssertEqual(t, "matched mismatch.", tc.matched, matched)
			ztesting.AssertEqual(t, "err mismatch.", tc.err, err)
		})
	}
}

func TestScanChunk(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		pattern string
		star    bool
		chunk   string
		rest    string
	}{
		{"x", false, "x", ""},
		{"foo", false, "foo", ""},
		{"*", true, "", ""},
		{"*foo", true, "foo", ""},
		{"*foo*", true, "foo", "*"},
		{"*foo?", true, "foo?", ""},
		{"*foo*bar", true, "foo", "*bar"},
		{"*foo?bar", true, "foo?bar", ""},
		{"?", false, "?", ""},
		{"?foo", false, "?foo", ""},
		{"?foo*", false, "?foo", "*"},
		{"?foo*bar", false, "?foo", "*bar"},
		{"?foo?bar", false, "?foo?bar", ""},
		{"\\x", false, "\\x", ""},
		{"\\*", false, "\\*", ""},
		{"\\*foo", false, "\\*foo", ""},
		{"\\*foo*bar", false, "\\*foo", "*bar"},
		{"\\?", false, "\\?", ""},
		{"\\?foo", false, "\\?foo", ""},
		{"\\?foo?bar", false, "\\?foo?bar", ""},
		{"\\f\\o\\o", false, "\\f\\o\\o", ""},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			star, chunk, rest := scanChunk(tc.pattern)
			ztesting.AssertEqual(t, "star mismatch.", tc.star, star)
			ztesting.AssertEqual(t, "chunk mismatch.", tc.chunk, chunk)
			ztesting.AssertEqual(t, "rest mismatch.", tc.rest, rest)
		})
	}
}

func TestMatchChunk(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		chunk string
		str   string
		ok    bool
		rest  string
	}{
		{"", "", true, ""},
		{" ", " ", true, ""},
		{"", "s", true, "s"},
		{"", "str", true, "str"},
		{"?", "", false, ""},
		{"?", "s", true, ""},
		{"?", "str", true, "tr"},
		{"s", "str", true, "tr"},
		{"s", "", false, ""},
		{"\\s", "str", true, "tr"},
		{"str", "s", false, ""},
		{"str", "str", true, ""},
		{"str", "str ", true, " "},
		{"str ", "str ", true, ""},
		{"ss", "str", false, "tr"},
		{"\\s\\s", "str", false, "tr"},
		{"\\\\", "\\", true, ""},
		{"\\\\s", "\\str", true, "tr"},
		{"\\*", "*", true, ""},
		{"\\*s", "*str", true, "tr"},
		{"\\?s", "?str", true, "tr"},
		{"s\\*r", "s*rx", true, "x"},
		{"s\\?r", "s?rx", true, "x"},
		{"foo\\? bar\\?", "foo? bar?", true, ""},
		{"foo\\? bar\\?", "foo? bar? baz", true, " baz"},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ok, rest := matchChunk(tc.chunk, tc.str)
			ztesting.AssertEqual(t, "unexpected match result.", tc.ok, ok)
			ztesting.AssertEqual(t, "rest string is wrong.", tc.rest, rest)
		})
	}
}
