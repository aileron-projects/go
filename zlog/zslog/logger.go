package zslog

import (
	"cmp"
	"context"
	"io"
	"log/slog"
	"math"
	"os"
	"time"

	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/zruntime"
)

func init() {
	h := &ctxHandler{slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})}
	slog.SetDefault(slog.New(h))
	hh := &ctxHandler{slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})}
	zlog.BuildLogFunc = buildLogFunc(hh)
}

// buildLogFunc returns a log function that should be used for [zlog.BuildLogFunc].
func buildLogFunc(h slog.Handler) func(ctx context.Context, msg string, args ...any) {
	return func(ctx context.Context, msg string, args ...any) {
		r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, 0)
		r.Add(append(args, CallerAttr(1), FramesAttr(1))...)
		e := h.Handle(ctx, r)
		zruntime.ReportErr(e, "") // Report runtime error if any.
	}
}

// NewContextHandler wraps the given handler with context aware handler.
// Once the handler wrapped, [ContextWithLevel] and [zlog.ContextWithAttrs]
// will works.
func NewContextHandler(h slog.Handler) slog.Handler {
	return &ctxHandler{Handler: h}
}

// NewJSON returns a new ZSLogger instance with given handler options.
// NewJSON uses [os.Stdout] if the given w is nil.
// NewJSON uses [slog.Logger] created with [slog.NewJSONHandler] internally.
func NewJSON(w io.Writer, opts *slog.HandlerOptions) *ZSLogger {
	w = cmp.Or(w, io.Writer(os.Stdout))
	h := slog.NewJSONHandler(w, opts)
	return &ZSLogger{
		h: NewContextHandler(h),
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
		h: NewContextHandler(h),
		w: w,
	}
}

type ZSLogger struct {
	h         slog.Handler
	w         io.Writer
	AddCaller Level
	AddFrames Level
}

func (l *ZSLogger) Handler() slog.Handler {
	return l.h
}

func (l *ZSLogger) Writer() io.Writer {
	return l.w
}

func (l *ZSLogger) Enabled(ctx context.Context, level slog.Level) bool {
	return l.h.Enabled(ctx, level)
}

func (l *ZSLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	if !l.h.Enabled(ctx, slog.LevelDebug) {
		return
	}
	if l.AddCaller&LvDebug > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&LvDebug > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelDebug, msg, 0)
	r.Add(args...)
	e := l.h.Handle(ctx, r)
	zruntime.ReportErr(e, "") // Report runtime error if any.
}

func (l *ZSLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	if !l.h.Enabled(ctx, slog.LevelInfo) {
		return
	}
	if l.AddCaller&LvInfo > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&LvInfo > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, 0)
	r.Add(args...)
	e := l.h.Handle(ctx, r)
	zruntime.ReportErr(e, "") // Report runtime error if any.
}

func (l *ZSLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	if !l.h.Enabled(ctx, slog.LevelWarn) {
		return
	}
	if l.AddCaller&LvWarn > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&LvWarn > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelWarn, msg, 0)
	r.Add(args...)
	e := l.h.Handle(ctx, r)
	zruntime.ReportErr(e, "") // Report runtime error if any.
}

func (l *ZSLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	if !l.h.Enabled(ctx, slog.LevelError) {
		return
	}
	if l.AddCaller&LvError > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&LvError > 0 {
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
	if v := LevelFromContext(ctx); v != slog.Level(math.MinInt) {
		return lv >= v
	}
	return h.Handler.Enabled(ctx, lv)
}

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := zlog.AttrsFromContext(ctx)
	r.Add(attrs...)
	return h.Handler.Handle(ctx, r)
}
