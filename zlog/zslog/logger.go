package zslog

import (
	"context"
	"log/slog"
	"math"
	"os"
	"sync/atomic"
	"time"

	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/zruntime"
)

func init() {
	h := &ctxHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})}
	slog.SetDefault(slog.New(h))
	SetDefault(New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))
}

var (
	// [Logger] implements [zlog.Logger] interface.
	_ zlog.Logger = &Logger{}
	// defaultLogger keeps default [Logger] instance.
	defaultLogger atomic.Pointer[Logger]
)

// SetDefault sets the default [Logger] instance.
// It replaces existing logger.
// Use [Default] to obtain the default logger.
func SetDefault(lg *Logger) {
	defaultLogger.Store(lg)
}

// Default returns the default [Logger] instance.
// Use [SetDefault] to replace the default logger.
func Default() *Logger {
	return defaultLogger.Load()
}

// New returns a new [Logger] instance that uses given handler.
// Create with a handler generated
//
//	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
//	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
func New(h slog.Handler) *Logger {
	return &Logger{
		Handler: &ctxHandler{Handler: h},
	}
}

// Logger is the logger type that uses [slog.Handler].
// A Logger created with [New] is compatible with the
// [ContextWithLevel] and the [zlog.ContextWithAttrs].
type Logger struct {
	slog.Handler
	// AddCaller specifies the log level ranges that
	// caller attributes should be appended.
	// For example, set zslog.RangeWarn | zslog.RangeError
	// to add caller only for warn and error levels.
	AddCaller Range
	// AddFrames specifies the log level ranges that
	// stack frame attributes should be appended.
	// For example, set zslog.RangeWarn | zslog.RangeError
	// to add frames only for warn and error levels.
	AddFrames Range
}

func (l *Logger) DebugEnabled(ctx context.Context) bool {
	return l.Enabled(ctx, slog.LevelDebug)
}

func (l *Logger) InfoEnabled(ctx context.Context) bool {
	return l.Enabled(ctx, slog.LevelInfo)
}

func (l *Logger) WarnEnabled(ctx context.Context) bool {
	return l.Enabled(ctx, slog.LevelWarn)
}

func (l *Logger) ErrorEnabled(ctx context.Context) bool {
	return l.Enabled(ctx, slog.LevelError)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	if !l.Enabled(ctx, slog.LevelDebug) {
		return
	}
	if l.AddCaller&RangeDebug > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&RangeDebug > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelDebug, msg, 0)
	r.Add(args...)
	e := l.Handle(ctx, r)
	zruntime.ReportErr(e, "") // Report runtime error if any.
}

func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	if !l.Enabled(ctx, slog.LevelInfo) {
		return
	}
	if l.AddCaller&RangeInfo > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&RangeInfo > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, 0)
	r.Add(args...)
	e := l.Handle(ctx, r)
	zruntime.ReportErr(e, "") // Report runtime error if any.
}

func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	if !l.Enabled(ctx, slog.LevelWarn) {
		return
	}
	if l.AddCaller&RangeWarn > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&RangeWarn > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelWarn, msg, 0)
	r.Add(args...)
	e := l.Handle(ctx, r)
	zruntime.ReportErr(e, "") // Report runtime error if any.
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	if !l.Enabled(ctx, slog.LevelError) {
		return
	}
	if l.AddCaller&RangeError > 0 {
		args = append(args, CallerAttr(1))
	}
	if l.AddFrames&RangeError > 0 {
		args = append(args, FramesAttr(1))
	}
	r := slog.NewRecord(time.Now(), slog.LevelError, msg, 0)
	r.Add(args...)
	e := l.Handle(ctx, r)
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
