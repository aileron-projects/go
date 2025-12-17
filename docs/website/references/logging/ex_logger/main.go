package main

import (
	"context"
	"log/slog"

	"github.com/aileron-projects/go/zlog/zslog"
)

func main() {
	// Create logger from default slogger handler.
	lg := zslog.New(slog.Default().Handler())

	ctx := context.Background()
	if lg.DebugEnabled(ctx) {
		lg.DebugContext(ctx, "debug enabled")
	}
	if lg.InfoEnabled(ctx) {
		lg.InfoContext(ctx, "info enabled")
	}
	if lg.WarnEnabled(ctx) {
		lg.WarnContext(ctx, "warn enabled")
	}
	if lg.ErrorEnabled(ctx) {
		lg.ErrorContext(ctx, "error enabled")
	}
}
