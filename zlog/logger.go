package zlog

import (
	"cmp"
	"context"
	"io"
	"log"
	"os"
)

var (
	_ Logger = &ZLogger{}
)

// LoggerOption is the options for [ZLogger].
type LoggerOption struct {
	// Lv is the log level of the logger.
	Lv Level
	// Flag is the flag for standard logger.
	Flag int
}

// New returns a new instance of ZLogger.
// If w is nil, [os.Stdout] is used.
// If opts is nil, a default options are used.
func New(w io.Writer, opts *LoggerOption) *ZLogger {
	w = cmp.Or(w, io.Writer(os.Stdout))
	if opts == nil {
		opts = &LoggerOption{
			Lv:   LvInfo,
			Flag: log.LstdFlags,
		}
	}
	return &ZLogger{
		lg: log.New(w, "", opts.Flag),
		w:  w,
		lv: opts.Lv,
	}
}

// ZLogger is a logger type.
// ZLogger leverages [log.Logger] internally.
// ZLogger implements [Logger] interface.
type ZLogger struct {
	lg *log.Logger
	w  io.Writer
	lv Level
}

func (l *ZLogger) Logger() *log.Logger {
	return l.lg
}

func (l *ZLogger) Writer() io.Writer {
	return l.w
}

func (l *ZLogger) Enabled(ctx context.Context, level Level) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	if v := LevelFromContext(ctx); v != LvUndef {
		return level.HigherEqual(v)
	}
	return level.HigherEqual(l.lv)
}

func (l *ZLogger) Debug(ctx context.Context, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.Enabled(ctx, LvDebug) {
		return
	}
	args = append([]any{"level=DEBUG", "msg=" + msg}, args...)
	l.lg.Println(append(args, AttrsFromContext(ctx)...)...)
}

func (l *ZLogger) Info(ctx context.Context, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.Enabled(ctx, LvInfo) {
		return
	}
	args = append([]any{"level=INFO", "msg=" + msg}, args...)
	l.lg.Println(append(args, AttrsFromContext(ctx)...)...)
}

func (l *ZLogger) Warn(ctx context.Context, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.Enabled(ctx, LvWarn) {
		return
	}
	args = append([]any{"level=WARN", "msg=" + msg}, args...)
	l.lg.Println(append(args, AttrsFromContext(ctx)...)...)
}

func (l *ZLogger) Error(ctx context.Context, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.Enabled(ctx, LvError) {
		return
	}
	args = append([]any{"level=ERROR", "msg=" + msg}, args...)
	l.lg.Println(append(args, AttrsFromContext(ctx)...)...)
}
