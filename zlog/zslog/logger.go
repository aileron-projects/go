package zslog

import (
	"cmp"
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/zruntime"
)

var (
	_ zlog.Logger = &ZSLogger{}
)

// NewJSON returns a new ZSLogger instance with given handler options.
// NewJSON uses [os.Stdout] if the given w is nil.
// NewJSON uses [slog.Logger] created with [slog.NewJSONHandler] internally.
func NewJSON(w io.Writer, opts *slog.HandlerOptions) *ZSLogger {
	w = cmp.Or(w, io.Writer(os.Stdout))
	h := slog.NewJSONHandler(w, opts)
	return &ZSLogger{
		h: &ctxHandler{Handler: h},
		w: w,
	}
}

// NewText returns a new ZSLogger instance with given handler options.
// NewText uses [os.Stdout] if the given w is nil.
// NewText uses [slog.Logger] created with [slog.NewTextHandler] internally.
func NewText(w io.Writer, opts *slog.HandlerOptions) *ZSLogger {
	w = cmp.Or(w, io.Writer(os.Stdout))
	h := slog.NewTextHandler(w, opts)
	return &ZSLogger{
		h: &ctxHandler{Handler: h},
		w: w,
	}
}

type ZSLogger struct {
	h         slog.Handler
	w         io.Writer
	AddCaller zlog.Range
	AddFrames zlog.Range
}

func (l *ZSLogger) Handler() slog.Handler {
	return l.h
}

func (l *ZSLogger) Writer() io.Writer {
	return l.w
}

func (l *ZSLogger) Enabled(ctx context.Context, level zlog.Level) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	return l.h.Enabled(ctx, toSLevel(level))
}

func (l *ZSLogger) Debug(ctx context.Context, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.h.Enabled(ctx, slog.LevelDebug) {
		return
	}
	if l.AddCaller&zlog.Debug > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&zlog.Debug > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelDebug, msg, 0)
	r.Add(args...)
	e := l.h.Handle(ctx, r)
	zruntime.ReportErr(e, "") // Report runtime error if any.
}

func (l *ZSLogger) Info(ctx context.Context, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.h.Enabled(ctx, slog.LevelInfo) {
		return
	}
	if l.AddCaller&zlog.Info > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&zlog.Info > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, 0)
	r.Add(args...)
	e := l.h.Handle(ctx, r)
	zruntime.ReportErr(e, "") // Report runtime error if any.
}

func (l *ZSLogger) Warn(ctx context.Context, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.h.Enabled(ctx, slog.LevelWarn) {
		return
	}
	if l.AddCaller&zlog.Warn > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&zlog.Warn > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelWarn, msg, 0)
	r.Add(args...)
	e := l.h.Handle(ctx, r)
	zruntime.ReportErr(e, "") // Report runtime error if any.
}

func (l *ZSLogger) Error(ctx context.Context, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.h.Enabled(ctx, slog.LevelError) {
		return
	}
	if l.AddCaller&zlog.Error > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&zlog.Error > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelError, msg, 0)
	r.Add(args...)
	e := l.h.Handle(ctx, r)
	zruntime.ReportErr(e, "") // Report runtime error if any.
}

// ctxHandler wraps [slog.Handler].
// ctxHandler checks if the handler is enabled or not
// based on the log level contained in context.
// ctxHandler extracts log attributes from context
// and add it to log records.
type ctxHandler struct {
	slog.Handler
}

func (h *ctxHandler) Enabled(ctx context.Context, lv slog.Level) bool {
	if v := zlog.LevelFromContext(ctx); v != zlog.LvUndef {
		return toZLevel(lv).HigherEqual(v)
	}
	return h.Handler.Enabled(ctx, lv)
}

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := zlog.AttrsFromContext(ctx)
	r.Add(attrs...)
	return h.Handler.Handle(ctx, r)
}
