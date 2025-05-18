package zslog

import (
	"fmt"
	"io"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestCallerAttr(t *testing.T) {
	t.Parallel()
	attr := CallerAttr(0)
	ztesting.AssertEqual(t, "invalid attribute key.", "caller", attr.Key)

	values := attr.Value.Group()
	ztesting.AssertEqual(t, "invalid number of values.", 4, len(values))
	if len(values) < 4 {
		return
	}
	ztesting.AssertEqual(t, "invalid key name.", "pkg", values[0].Key)
	ztesting.AssertEqual(t, "invalid key name.", "file", values[1].Key)
	ztesting.AssertEqual(t, "invalid key name.", "func", values[2].Key)
	ztesting.AssertEqual(t, "invalid key name.", "line", values[3].Key)
}

func TestDateTimeAttr(t *testing.T) {
	t.Parallel()
	attr := DateTimeAttr()
	ztesting.AssertEqual(t, "invalid attribute key.", "datetime", attr.Key)

	values := attr.Value.Group()
	ztesting.AssertEqual(t, "invalid number of values.", 2, len(values))
	if len(values) < 2 {
		return
	}
	ztesting.AssertEqual(t, "invalid key name.", "date", values[0].Key)
	ztesting.AssertEqual(t, "invalid key name.", "time", values[1].Key)
}

func TestFramesAttr(t *testing.T) {
	t.Parallel()
	attr := FramesAttr(0)
	ztesting.AssertEqual(t, "invalid attribute key.", "frames", attr.Key)

	values, ok := attr.Value.Any().([]string)
	ztesting.AssertEqual(t, "no frames found.", true, ok)
	ztesting.AssertEqual(t, "no frames found.", true, len(values) > 0)
}

func TestStackTraceAttrs(t *testing.T) {
	t.Parallel()
	attr := StackTraceAttrs(0)
	ztesting.AssertEqual(t, "invalid attribute key.", "stack", attr.Key)
	ztesting.AssertEqual(t, "stack trace is empty.", true, len(attr.Value.String()) > 0)
}

func TestErrorAttr(t *testing.T) {
	t.Parallel()
	attr := ErrorAttr(io.EOF)
	ztesting.AssertEqual(t, "invalid attribute key.", "error", attr.Key)
	ztesting.AssertEqual(t, "unexpected error attribute.", "map[msg:EOF]", fmt.Sprint(attr.Value.Any()))
}
