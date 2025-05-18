package zstrings_test

import (
	"math/rand"
	"path"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zstrings"
)

func FuzzReverse(f *testing.F) {
	f.Fuzz(func(t *testing.T, n int8) {
		target := randString(n / 5)  // Adjust length.
		pattern := randString(n / 2) // Adjust length.
		pattern = strings.TrimRight(pattern, "\\")
		m1, _ := zstrings.Match(pattern, target)
		m2, _ := path.Match(pattern, target)
		if m1 != m2 {
			t.Errorf("Pattern:`%s` Target:`%s`", pattern, target)
			return
		}
	})
}

func randString(n int8) string {
	if n < 0 {
		n = -n
	}
	pattern := []byte("aaaaabbbbbccccc*?\\")
	m := len(pattern)
	b := make([]byte, n)
	for i := range b {
		b[i] = pattern[rand.Intn(m)]
	}
	return string(b)
}
