package zstrings

import (
	"strings"
	"unicode/utf8"
)

// CutLeftByte slices the given string s at the first instance of c,
// returning the text before and after the separator of c.
// The found result reports whether c appears in s.
// The separator c is not contained both before and after.
// If the separator c does not appear in s, cut returns s, "", false.
// See also [strings.Cut] and [strings.IndexByte].
func CutLeftByte(s string, c byte) (before, after string, found bool) {
	i := strings.IndexByte(s, c)
	if i < 0 {
		return s, "", false
	}
	return s[:i], s[i+1:], true
}

// CutRightByte slices the given string s at the last instance of c,
// returning the text before and after the separator of c.
// The found result reports whether c appears in s.
// The separator c is not contained both before and after.
// If the separator c does not appear in s, cut returns s, "", false.
// See also [strings.Cut] and [strings.LastIndexByte].
func CutRightByte(s string, c byte) (before, after string, found bool) {
	i := strings.LastIndexByte(s, c)
	if i < 0 {
		return s, "", false
	}
	return s[:i], s[i+1:], true
}

// CutLeftRune slices the given string s at the first instance of r,
// returning the text before and after the separator of r.
// The found result reports whether r appears in s.
// The separator r is not contained both before and after.
// If the separator r does not appear in s, cut returns s, "", false.
// If the given r is not valid UTF-8 rune, it also returns s, "", false.
// CutLeftRune uses Brute Force strategy for all s.
// See also [strings.IndexRune].
func CutLeftRune(s string, r rune) (before, after string, found bool) {
	switch {
	case 0 <= r && r < utf8.RuneSelf:
		return CutLeftByte(s, byte(r))
	case !utf8.ValidRune(r):
		return s, "", false
	}
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if rs[i] == r {
			return string(rs[:i]), string(rs[i+1:]), true
		}
	}
	return s, "", false
}

// CutRightRune slices the given string s at the last instance of r,
// returning the text before and after the separator of r.
// The found result reports whether r appears in s.
// The separator r is not contained both before and after.
// If the separator r does not appear in s, cut returns s, "", false.
// If the given r is not valid UTF-8 rune, it also returns s, "", false.
// CutLeftRune uses Brute Force strategy for all s.
// See also [strings.IndexRune].
func CutRightRune(s string, r rune) (before, after string, found bool) {
	switch {
	case 0 <= r && r < utf8.RuneSelf:
		return CutRightByte(s, byte(r))
	case !utf8.ValidRune(r):
		return s, "", false
	}
	rs := []rune(s)
	for i := len(rs) - 1; i >= 0; i-- {
		if rs[i] == r {
			return string(rs[:i]), string(rs[i+1:]), true
		}
	}
	return s, "", false
}
