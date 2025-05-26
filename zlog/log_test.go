package zlog_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/ztesting"
)

func TestContextWithAttrs(t *testing.T) {
	t.Parallel()
	t.Run("nil context", func(t *testing.T) {
		ctx := zlog.ContextWithAttrs(nil, "foo", "bar")
		attrs := zlog.AttrsFromContext(ctx)
		ztesting.AssertEqual(t, "invalid number of attributes.", 2, len(attrs))
		ztesting.AssertEqual(t, "invalid content of attributes.", []any{"foo", "bar"}, attrs)
	})
	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = zlog.ContextWithAttrs(ctx, "foo", "bar")
		attrs := zlog.AttrsFromContext(ctx)
		ztesting.AssertEqual(t, "invalid number of attributes.", 2, len(attrs))
		ztesting.AssertEqual(t, "invalid content of attributes.", []any{"foo", "bar"}, attrs)
	})
	t.Run("non empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = zlog.ContextWithAttrs(ctx, "foo")
		ctx = zlog.ContextWithAttrs(ctx, "bar")
		attrs := zlog.AttrsFromContext(ctx)
		ztesting.AssertEqual(t, "invalid number of attributes.", 2, len(attrs))
		ztesting.AssertEqual(t, "invalid content of attributes.", []any{"foo", "bar"}, attrs)
	})
}

func TestAttrsFromContext(t *testing.T) {
	t.Parallel()
	t.Run("nil context", func(t *testing.T) {
		attrs := zlog.AttrsFromContext(nil)
		ztesting.AssertEqual(t, "invalid number of attributes.", 0, len(attrs))
	})
	t.Run("empty context", func(t *testing.T) {
		attrs := zlog.AttrsFromContext(context.Background())
		ztesting.AssertEqual(t, "invalid number of attributes.", 0, len(attrs))
	})
	t.Run("non empty context", func(t *testing.T) {
		ctx := zlog.ContextWithAttrs(context.Background(), "foo", "bar")
		attrs := zlog.AttrsFromContext(ctx)
		ztesting.AssertEqual(t, "invalid number of attributes.", 2, len(attrs))
		ztesting.AssertEqual(t, "invalid content of attributes.", []any{"foo", "bar"}, attrs)
	})
}

func TestBuildLog(t *testing.T) {
	t.Cleanup(func() {
		zlog.BuildLogFunc = slog.Default().InfoContext
	})
	var w bytes.Buffer
	h := slog.NewTextHandler(&w, &slog.HandlerOptions{})
	zlog.BuildLogFunc = slog.New(h).InfoContext
	zlog.BuildLog(context.Background(), "test")
	out := w.String()
	ztesting.AssertEqual(t, "output not match", zlog.ExportedBuildlogEnabled, out != "")
}
