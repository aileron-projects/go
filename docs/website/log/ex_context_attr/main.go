package main

import (
	"context"
	"log/slog"

	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/zlog/zslog"
)

func main() {
	lg := zslog.New(slog.Default().Handler())

	ctx := context.Background()
	ctx = zlog.ContextWithAttrs(ctx, "key", "value")

	if lg.InfoEnabled(ctx) {
		lg.InfoContext(ctx, "output context attributes")
	}
}
