package zslog

import (
	"context"
	"log/slog"
	"math"
)

// Range is the log range type.
// Defined ranges are:
//   - [LvUndef]
//   - [LvDebug]
//   - [LvInfo]
//   - [LvWarn]
//   - [LvError]
type Level uint

const (
	LvUndef Level = 1 << iota
	LvDebug
	LvInfo
	LvWarn
	LvError
)

// ctxKey is the context key type.
type ctxKey struct{ string }

// ctxLevelKey is the key to store log level in a context.
var ctxLevelKey = &ctxKey{"level"}

// ContextWithLevel returns a new context with given log level.
// Use [LevelFromContext] to extract log level from the context.
// ContextWithLevel uses a new context created with [context.Background]
// if the given ctx is nil.
func ContextWithLevel(parent context.Context, lv slog.Level) context.Context {
	if parent == nil {
		return context.WithValue(context.Background(), ctxLevelKey, lv)
	}
	return context.WithValue(parent, ctxLevelKey, lv)
}

// LevelFromContext returns a log level stored in the context.
// Use [ContextWithLevel] to store a log level in context.
// LevelFromContext returns [math.MinInt] if the given ctx is nil
// or no log levels are found in the context.
func LevelFromContext(ctx context.Context) slog.Level {
	if ctx == nil {
		return math.MinInt
	}
	if v := ctx.Value(ctxLevelKey); v != nil {
		return v.(slog.Level)
	}
	return math.MinInt
}
