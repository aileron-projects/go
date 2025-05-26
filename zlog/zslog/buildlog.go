package zslog

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/aileron-projects/go/zruntime"
)

// buildLogHandler is the [slog.handler] used by [BuildLog].
// [BuildLog] is only enabled when the build tag -tag="zslogbuildlog" is specified.
// By default, a new instance [slog.NewJSONHandler] is used.
// Currently, this is not replaceable by users.
var buildLogHandler slog.Handler = &ctxHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})}

// BuildLog outputs log records when build-time logging is enabled.
// Build-time logging can be enabled by using build tag -tag="zslogbuildlog".
func BuildLog(ctx context.Context, msg string, args ...any) {
	if buildlogEnabled {
		r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, 0)
		r.Add(append(args, CallerAttr(1), FramesAttr(1))...)
		e := buildLogHandler.Handle(ctx, r)
		zruntime.ReportErr(e, "") // Report runtime error if any.
	}
}
