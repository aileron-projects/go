package zslog_test

import (
	"context"
	"log/slog"

	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/zlog/zslog"
)

func ExampleNewJSON() {
	opts := &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		ReplaceAttr: zslog.RemoveTime, // Remove time to avoid test failure.
	}
	lg := zslog.NewJSON(nil, opts)

	ctx := context.Background()
	lg.Info(ctx, "this is info")                   // Will be output.
	lg.Debug(ctx, "this is debug")                 // Won't be output.
	ctx = zlog.ContextWithLevel(ctx, zlog.LvDebug) // Update log level through the context.
	lg.Debug(ctx, "this is debug again")           // Will be output.

	// Output:
	// {"level":"INFO","msg":"this is info"}
	// {"level":"DEBUG","msg":"this is debug again"}
}

func ExampleNewText() {
	opts := &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		ReplaceAttr: zslog.RemoveTime, // Remove time to avoid test failure.
	}
	lg := zslog.NewText(nil, opts)

	ctx := context.Background()
	lg.Info(ctx, "this is info")                   // Will be output.
	lg.Debug(ctx, "this is debug")                 // Won't be output.
	ctx = zlog.ContextWithLevel(ctx, zlog.LvDebug) // Update log level through the context.
	lg.Debug(ctx, "this is debug again")           // Will be output.

	// Output:
	// level=INFO msg="this is info"
	// level=DEBUG msg="this is debug again"
}
