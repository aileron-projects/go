package zstrings_test

import (
	"testing"

	"github.com/aileron-projects/go/zstrings"
	"github.com/aileron-projects/go/ztesting"
)

func TestCountByteRight(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		s string
		b byte
		n int
	}{
		"empty s, empty b":        {"", 0x00, 0},
		"empty s, space b":        {"", ' ', 0},
		"empty s, char b":         {"", 'x', 0},
		"not match, empty b":      {"x", 0x00, 0},
		"not match, space b":      {"x", ' ', 0},
		"not match, char b":       {"x", 'y', 0},
		"match single, char b":    {"x", 'x', 1},
		"match multiple1, char b": {"xxx", 'x', 3},
		"match multiple2, char b": {"yxxx", 'x', 3},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			matched := zstrings.CountByteRight(tc.s, tc.b)
			ztesting.AssertEqual(t, "wrong match count.", tc.n, matched)
		})
	}
}

func TestCountByteLeft(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		s string
		b byte
		n int
	}{
		"empty s, empty b":        {"", 0x00, 0},
		"empty s, space b":        {"", ' ', 0},
		"empty s, char b":         {"", 'x', 0},
		"not match, empty b":      {"x", 0x00, 0},
		"not match, space b":      {"x", ' ', 0},
		"not match, char b":       {"x", 'y', 0},
		"match single, char b":    {"x", 'x', 1},
		"match multiple1, char b": {"xxx", 'x', 3},
		"match multiple2, char b": {"xxxy", 'x', 3},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			matched := zstrings.CountByteLeft(tc.s, tc.b)
			ztesting.AssertEqual(t, "wrong match count.", tc.n, matched)
		})
	}
}
