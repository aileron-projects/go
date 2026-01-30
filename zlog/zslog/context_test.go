package zslog_test

import (
	"context"
	"log/slog"
	"math"
	"testing"

	"github.com/aileron-projects/go/zlog/zslog"
	"github.com/aileron-projects/go/ztesting"
)

func TestContextWithAttrs(t *testing.T) {
	t.Parallel()
	t.Run("nil context", func(t *testing.T) {
		ctx := zslog.ContextWithAttrs(nil, "foo", "bar")
		attrs := zslog.AttrsFromContext(ctx)
		ztesting.AssertEqual(t, "invalid number of attributes.", 2, len(attrs))
		ztesting.AssertEqual(t, "invalid content of attributes.", []any{"foo", "bar"}, attrs)
	})
	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = zslog.ContextWithAttrs(ctx, "foo", "bar")
		attrs := zslog.AttrsFromContext(ctx)
		ztesting.AssertEqual(t, "invalid number of attributes.", 2, len(attrs))
		ztesting.AssertEqual(t, "invalid content of attributes.", []any{"foo", "bar"}, attrs)
	})
	t.Run("non empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = zslog.ContextWithAttrs(ctx, "foo")
		ctx = zslog.ContextWithAttrs(ctx, "bar")
		attrs := zslog.AttrsFromContext(ctx)
		ztesting.AssertEqual(t, "invalid number of attributes.", 2, len(attrs))
		ztesting.AssertEqual(t, "invalid content of attributes.", []any{"foo", "bar"}, attrs)
	})
}

func TestAttrsFromContext(t *testing.T) {
	t.Parallel()
	t.Run("empty context", func(t *testing.T) {
		attrs := zslog.AttrsFromContext(context.Background())
		ztesting.AssertEqual(t, "invalid number of attributes.", 0, len(attrs))
	})
	t.Run("non empty context", func(t *testing.T) {
		ctx := zslog.ContextWithAttrs(context.Background(), "foo", "bar")
		attrs := zslog.AttrsFromContext(ctx)
		ztesting.AssertEqual(t, "invalid number of attributes.", 2, len(attrs))
		ztesting.AssertEqual(t, "invalid content of attributes.", []any{"foo", "bar"}, attrs)
	})
}

func TestContextWithLevel(t *testing.T) {
	t.Parallel()
	t.Run("nil context", func(t *testing.T) {
		ctx := zslog.ContextWithLevel(nil, slog.LevelError)
		lv := zslog.LevelFromContext(ctx)
		ztesting.AssertEqual(t, "level mismatch.", slog.LevelError, lv)
	})
	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = zslog.ContextWithLevel(ctx, slog.LevelError)
		lv := zslog.LevelFromContext(ctx)
		ztesting.AssertEqual(t, "level mismatch.", slog.LevelError, lv)
	})
	t.Run("non empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = zslog.ContextWithLevel(ctx, slog.LevelDebug)
		ctx = zslog.ContextWithLevel(ctx, slog.LevelError)
		lv := zslog.LevelFromContext(ctx)
		ztesting.AssertEqual(t, "level mismatch.", slog.LevelError, lv)
	})
}

func TestLevelFromContext(t *testing.T) {
	t.Parallel()
	t.Run("empty context", func(t *testing.T) {
		lv := zslog.LevelFromContext(context.Background())
		ztesting.AssertEqual(t, "level mismatch.", slog.Level(math.MinInt), lv)
	})
	t.Run("non empty context", func(t *testing.T) {
		ctx := zslog.ContextWithLevel(context.Background(), slog.LevelError)
		lv := zslog.LevelFromContext(ctx)
		ztesting.AssertEqual(t, "level mismatch.", slog.LevelError, lv)
	})
}

func TestContextWithHandler(t *testing.T) {
	t.Parallel()
	testHandler := struct {
		slog.Handler
		s string
	}{nil, "test"}
	t.Run("nil context", func(t *testing.T) {
		ctx := zslog.ContextWithHandler(nil, testHandler)
		got := zslog.HandlerFromContext(ctx)
		ztesting.AssertEqual(t, "handler level mismatch.", slog.Handler(testHandler), got)
	})
	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = zslog.ContextWithHandler(ctx, testHandler)
		got := zslog.HandlerFromContext(ctx)
		ztesting.AssertEqual(t, "handler mismatch.", slog.Handler(testHandler), got)
	})
	t.Run("non empty context", func(t *testing.T) {
		testHandler2 := struct {
			slog.Handler
			s string
		}{nil, "test2"}
		ctx := context.Background()
		ctx = zslog.ContextWithHandler(ctx, testHandler)
		ctx = zslog.ContextWithHandler(ctx, testHandler2)
		got := zslog.HandlerFromContext(ctx)
		ztesting.AssertEqual(t, "handler mismatch.", slog.Handler(testHandler2), got)
	})
}

func TestHandlerFromContext(t *testing.T) {
	t.Parallel()
	testHandler := struct {
		slog.Handler
		s string
	}{nil, "test"}
	t.Run("empty context", func(t *testing.T) {
		got := zslog.HandlerFromContext(context.Background())
		ztesting.AssertEqual(t, "handler mismatch.", nil, got)
	})
	t.Run("non empty context", func(t *testing.T) {
		ctx := zslog.ContextWithHandler(context.Background(), testHandler)
		got := zslog.HandlerFromContext(ctx)
		ztesting.AssertEqual(t, "handler mismatch.", slog.Handler(testHandler), got)
	})
}
