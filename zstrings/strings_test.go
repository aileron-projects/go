package zstrings_test

import (
	"testing"
	"unicode/utf8"

	"github.com/aileron-projects/go/zstrings"
	"github.com/aileron-projects/go/ztesting"
)

func TestCutLeftByte(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		s             string
		c             byte
		before, after string
		found         bool
	}{
		"case00": {"", 0x00, "", "", false},
		"case01": {"", ' ', "", "", false},
		"case02": {"a", 0x00, "a", "", false},
		"case03": {"a", ' ', "a", "", false},
		"case04": {"a", 'z', "a", "", false},
		"case05": {"a", 'a', "", "", true},
		"case06": {"abc", 0x00, "abc", "", false},
		"case07": {"abc", ' ', "abc", "", false},
		"case08": {"abc", 'z', "abc", "", false},
		"case09": {"abc", 'a', "", "bc", true},
		"case10": {"abc", 'b', "a", "c", true},
		"case11": {"abc", 'c', "ab", "", true},
		"case12": {"a.b.c", '.', "a", "b.c", true},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			before, after, found := zstrings.CutLeftByte(tc.s, tc.c)
			ztesting.AssertEqual(t, "before not matched.", tc.before, before)
			ztesting.AssertEqual(t, "after not matched.", tc.after, after)
			ztesting.AssertEqual(t, "found not matched.", tc.found, found)
		})
	}
}

func TestCutRightByte(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		s             string
		c             byte
		before, after string
		found         bool
	}{
		"case00": {"", 0x00, "", "", false},
		"case01": {"", ' ', "", "", false},
		"case02": {"a", 0x00, "a", "", false},
		"case03": {"a", ' ', "a", "", false},
		"case04": {"a", 'z', "a", "", false},
		"case05": {"a", 'a', "", "", true},
		"case06": {"abc", 0x00, "abc", "", false},
		"case07": {"abc", ' ', "abc", "", false},
		"case08": {"abc", 'z', "abc", "", false},
		"case09": {"abc", 'a', "", "bc", true},
		"case10": {"abc", 'b', "a", "c", true},
		"case11": {"abc", 'c', "ab", "", true},
		"case12": {"a.b.c", '.', "a.b", "c", true},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			before, after, found := zstrings.CutRightByte(tc.s, tc.c)
			ztesting.AssertEqual(t, "before not matched.", tc.before, before)
			ztesting.AssertEqual(t, "after not matched.", tc.after, after)
			ztesting.AssertEqual(t, "found not matched.", tc.found, found)
		})
	}
}

func TestCutLeftRune(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		s             string
		r             rune
		before, after string
		found         bool
	}{
		"case00": {"", 0x00, "", "", false},
		"case01": {"", ' ', "", "", false},
		"case02": {"a", 0x00, "a", "", false},
		"case03": {"a", ' ', "a", "", false},
		"case04": {"a", 'z', "a", "", false},
		"case05": {"a", 'a', "", "", true},
		"case06": {"abc", 0x00, "abc", "", false},
		"case07": {"abc", ' ', "abc", "", false},
		"case08": {"abc", 'z', "abc", "", false},
		"case09": {"abc", 'a', "", "bc", true},
		"case10": {"abc", 'b', "a", "c", true},
		"case11": {"abc", 'c', "ab", "", true},
		"case12": {"a.b.c", '.', "a", "b.c", true},
		"case13": {"あいう", '　', "あいう", "", false},
		"case14": {"あいう", 'あ', "", "いう", true},
		"case15": {"あいう", 'い', "あ", "う", true},
		"case16": {"あいう", 'う', "あい", "", true},
		"case17": {"あえいえう", 'え', "あ", "いえう", true},
		"case18": {"あいう", utf8.MaxRune + 1, "あいう", "", false},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			before, after, found := zstrings.CutLeftRune(tc.s, tc.r)
			ztesting.AssertEqual(t, "before not matched.", tc.before, before)
			ztesting.AssertEqual(t, "after not matched.", tc.after, after)
			ztesting.AssertEqual(t, "found not matched.", tc.found, found)
		})
	}
}

func TestCutRightRune(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		s             string
		r             rune
		before, after string
		found         bool
	}{
		"case00": {"", 0x00, "", "", false},
		"case01": {"", ' ', "", "", false},
		"case02": {"a", 0x00, "a", "", false},
		"case03": {"a", ' ', "a", "", false},
		"case04": {"a", 'z', "a", "", false},
		"case05": {"a", 'a', "", "", true},
		"case06": {"abc", 0x00, "abc", "", false},
		"case07": {"abc", ' ', "abc", "", false},
		"case08": {"abc", 'z', "abc", "", false},
		"case09": {"abc", 'a', "", "bc", true},
		"case10": {"abc", 'b', "a", "c", true},
		"case11": {"abc", 'c', "ab", "", true},
		"case12": {"a.b.c", '.', "a.b", "c", true},
		"case13": {"あいう", '　', "あいう", "", false},
		"case14": {"あいう", 'あ', "", "いう", true},
		"case15": {"あいう", 'い', "あ", "う", true},
		"case16": {"あいう", 'う', "あい", "", true},
		"case17": {"あえいえう", 'え', "あえい", "う", true},
		"case18": {"あいう", utf8.MaxRune + 1, "あいう", "", false},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			before, after, found := zstrings.CutRightRune(tc.s, tc.r)
			ztesting.AssertEqual(t, "before not matched.", tc.before, before)
			ztesting.AssertEqual(t, "after not matched.", tc.after, after)
			ztesting.AssertEqual(t, "found not matched.", tc.found, found)
		})
	}
}
