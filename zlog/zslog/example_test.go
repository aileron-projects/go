package zslog_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/aileron-projects/go/zlog/zslog"
)

func ExampleNew_jsonHandler() {
	opts := &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		ReplaceAttr: RemoveTime, // Remove time to avoid test failure.
	}
	lg := zslog.New(slog.NewJSONHandler(os.Stdout, opts))

	ctx := context.Background()
	lg.InfoContext(ctx, "this is info")                // Will be output.
	lg.DebugContext(ctx, "this is debug")              // Won't be output.
	ctx = zslog.ContextWithLevel(ctx, slog.LevelDebug) // Update log level through the context.
	lg.DebugContext(ctx, "this is debug again")        // Will be output.

	// Output:
	// {"level":"INFO","msg":"this is info"}
	// {"level":"DEBUG","msg":"this is debug again"}
}

func ExampleNew_textHandler() {
	opts := &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		ReplaceAttr: RemoveTime, // Remove time to avoid test failure.
	}
	lg := zslog.New(slog.NewTextHandler(os.Stdout, opts))

	ctx := context.Background()
	lg.InfoContext(ctx, "this is info")                // Will be output.
	lg.DebugContext(ctx, "this is debug")              // Won't be output.
	ctx = zslog.ContextWithLevel(ctx, slog.LevelDebug) // Update log level through the context.
	lg.DebugContext(ctx, "this is debug again")        // Will be output.

	// Output:
	// level=INFO msg="this is info"
	// level=DEBUG msg="this is debug again"
}

func RemoveTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}
