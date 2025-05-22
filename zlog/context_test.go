package zlog_test

import (
	"context"
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

func TestContextWithLevel(t *testing.T) {
	t.Parallel()

	t.Run("nil context", func(t *testing.T) {
		ctx := zlog.ContextWithLevel(nil, zlog.LvError)
		lv := zlog.LevelFromContext(ctx)
		ztesting.AssertEqual(t, "log level mismatch.", zlog.LvError, lv)
	})

	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = zlog.ContextWithLevel(ctx, zlog.LvError)
		lv := zlog.LevelFromContext(ctx)
		ztesting.AssertEqual(t, "log level mismatch.", zlog.LvError, lv)
	})

	t.Run("non empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = zlog.ContextWithLevel(ctx, zlog.LvDebug)
		ctx = zlog.ContextWithLevel(ctx, zlog.LvError)
		lv := zlog.LevelFromContext(ctx)
		ztesting.AssertEqual(t, "log level mismatch.", zlog.LvError, lv)
	})
}

func TestLevelFromContext(t *testing.T) {
	t.Parallel()

	t.Run("nil context", func(t *testing.T) {
		lv := zlog.LevelFromContext(nil)
		ztesting.AssertEqual(t, "log level mismatch.", zlog.LvUndef, lv)
	})

	t.Run("empty context", func(t *testing.T) {
		lv := zlog.LevelFromContext(context.Background())
		ztesting.AssertEqual(t, "log level mismatch.", zlog.LvUndef, lv)
	})

	t.Run("non empty context", func(t *testing.T) {
		ctx := zlog.ContextWithLevel(context.Background(), zlog.LvError)
		lv := zlog.LevelFromContext(ctx)
		ztesting.AssertEqual(t, "log level mismatch.", zlog.LvError, lv)
	})
}
