package zslog

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestBuildLog(t *testing.T) {
	tmp := buildLogHandler
	t.Cleanup(func() { buildLogHandler = tmp })
	var w bytes.Buffer
	buildLogHandler = slog.NewJSONHandler(&w, &slog.HandlerOptions{})
	BuildLog(context.Background(), "test", "foo", "bar")
	out := w.String()
	if !buildlogEnabled {
		ztesting.AssertEqual(t, "output not match", "", out)
	} else {
		ztesting.AssertEqual(t, "message not output", true, strings.Contains(out, `"msg":"test"`))
		ztesting.AssertEqual(t, "args not output", true, strings.Contains(out, `"foo":"bar"`))
	}
}
