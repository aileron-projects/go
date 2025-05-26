package zslog

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestBuildLogFunc(t *testing.T) {
	t.Parallel()
	var w bytes.Buffer
	h := slog.NewJSONHandler(&w, &slog.HandlerOptions{})
	f := buildLogFunc(h)
	f(context.Background(), "test", "foo", "bar")
	out := w.String()
	ztesting.AssertEqual(t, "message not output", true, strings.Contains(out, `"msg":"test"`))
	ztesting.AssertEqual(t, "args not output", true, strings.Contains(out, `"foo":"bar"`))
}
