package main

import (
	"context"
	"log/slog"

	"github.com/aileron-projects/go/zlog/zslog"
)

func main() {
	lg := zslog.New(slog.Default().Handler()) // Info level logger.
	ctx := context.Background()

	// Try to output debug log.
	lg.DebugContext(ctx, "log should not be output")

	// Set this context to debug level.
	ctx = zslog.ContextWithLevel(ctx, slog.LevelDebug)

	// Once again, try to output debug log.
	lg.DebugContext(ctx, "log should be output")
}
