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
	ztesting.AssertEqual(t, "caller", attr.Key)

	values := attr.Value.Group()
	ztesting.AssertEqual(t, 4, len(values))
	if len(values) < 4 {
		return
	}
	ztesting.AssertEqual(t, "pkg", values[0].Key)
	ztesting.AssertEqual(t, "file", values[1].Key)
	ztesting.AssertEqual(t, "func", values[2].Key)
	ztesting.AssertEqual(t, "line", values[3].Key)
}

func TestDateTimeAttr(t *testing.T) {
	t.Parallel()
	attr := DateTimeAttr()
	ztesting.AssertEqual(t, "datetime", attr.Key)

	values := attr.Value.Group()
	ztesting.AssertEqual(t, 2, len(values))
	if len(values) < 2 {
		return
	}
	ztesting.AssertEqual(t, "date", values[0].Key)
	ztesting.AssertEqual(t, "time", values[1].Key)
}

func TestFramesAttr(t *testing.T) {
	t.Parallel()
	attr := FramesAttr(0)
	ztesting.AssertEqual(t, "frames", attr.Key)

	values, ok := attr.Value.Any().([]string)
	ztesting.AssertEqual(t, true, ok)
	ztesting.AssertEqual(t, true, len(values) > 0)
}

func TestStackTraceAttrs(t *testing.T) {
	t.Parallel()
	attr := StackTraceAttrs(0)
	ztesting.AssertEqual(t, "stack", attr.Key)
	ztesting.AssertEqual(t, true, len(attr.Value.String()) > 0)
}

func TestErrorAttr(t *testing.T) {
	t.Parallel()
	attr := ErrorAttr(io.EOF)
	ztesting.AssertEqual(t, "error", attr.Key)
	ztesting.AssertEqual(t, "map[message:EOF]", fmt.Sprint(attr.Value.Any()))
}
