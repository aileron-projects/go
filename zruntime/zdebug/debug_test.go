package zdebug_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zruntime"
	"github.com/aileron-projects/go/zruntime/zdebug"
	"github.com/aileron-projects/go/ztesting"
)

func TestDumpTo(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		a     []any
		wants []string
	}{
		"nil": {
			a:     nil,
			wants: []string{`| Nothing to dump.`},
		},
		"int": {
			a:     []any{int(1), int32(2), int64(3)},
			wants: []string{`| (int) 1`, `| (int32) 2`, `| (int64) 3`},
		},
		"uint": {
			a:     []any{uint(1), uint32(2), uint64(3)},
			wants: []string{`| (uint) 1`, `| (uint32) 2`, `| (uint64) 3`},
		},
		"float": {
			a:     []any{float32(1.23), float64(4.56)},
			wants: []string{`| (float32) 1.23`, `| (float64) 4.56`},
		},
		"bool": {
			a:     []any{true, false},
			wants: []string{`| (bool) true`, `| (bool) false`},
		},
		"complex": {
			a:     []any{complex64(1 + 2i), complex128(3 + 4i)},
			wants: []string{`| (complex64) (1+2i)`, `| (complex128) (3+4i)`},
		},
		"slice": {
			a:     []any{[]int{1, 2}},
			wants: []string{`| ([]int) (len=2 cap=2) {`, `|  (int) 1,`, `|  (int) 2`},
		},
		"map": {
			a:     []any{map[int]string{1: "a", 2: "b"}},
			wants: []string{`| (map[int]string) (len=2) {`, `|  (int) 1: (string) (len=1) "a"`, `|  (int) 2: (string) (len=1) "b"`},
		},
		"struct": {
			a:     []any{struct{ x int }{}},
			wants: []string{`| (struct { x int }) {`, `|  x: (int) 0`},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			zdebug.DumpTo(&buf, tc.a...)
			result := buf.String()
			ztesting.AssertEqual(t, "dump result does not contain date and time.", true, strings.Contains(result, "1970-01-01 00:00:00 [DUMP]"))
			for _, w := range tc.wants {
				ztesting.AssertEqual(t, "expected dump result not output.", true, strings.Contains(result, w))
			}
		})
	}
}

func TestHookDumpFunc(t *testing.T) {
	returnBool := false
	hookDumpFunc := zdebug.HookDumpFunc
	zdebug.HookDumpFunc = func(w io.Writer, f zruntime.Frame, a ...any) bool {
		ztesting.AssertEqual(t, "length does not match.", 2, len(a))
		ztesting.AssertEqual(t, "unexpected value.", int(0), a[0].(int))
		ztesting.AssertEqual(t, "unexpected value.", int(1), a[1].(int))
		return returnBool
	}
	defer func() {
		zdebug.HookDumpFunc = hookDumpFunc // Reset to original.
	}()

	t.Run("func returns true", func(t *testing.T) {
		var buf bytes.Buffer
		returnBool = true // Set value true.
		zdebug.DumpTo(&buf, int(0), int(1))
		ztesting.AssertEqual(t, "dump continued after HookDumpFunc returned true.", "", buf.String())
	})
	t.Run("func returns false", func(t *testing.T) {
		var buf bytes.Buffer
		returnBool = false // Set value false.
		zdebug.DumpTo(&buf, int(0), int(1))
		ztesting.AssertEqual(t, "dump should output information.", false, buf.String() == "")
	})
}
