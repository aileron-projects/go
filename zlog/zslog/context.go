package zslog

import (
	"context"
	"log/slog"
	"math"
)

// Range is the log range type.
// Defined ranges are:
//   - [RangeUndef]
//   - [RangeDebug]
//   - [RangeInfo]
//   - [RangeWarn]
//   - [RangeError]
type Range uint

const (
	RangeUndef Range = 1 << iota
	RangeDebug
	RangeInfo
	RangeWarn
	RangeError
)

// ctxKey is the context key type.
type ctxKey struct{ string }

var (
	ctxAttrsKey   = &ctxKey{"attrs"}   // ctxAttrsKey is the key to store log attributes.
	ctxLevelKey   = &ctxKey{"level"}   // ctxLevelKey is the key to store log level.
	ctxHandlerKey = &ctxKey{"handler"} // ctxHandlerKey is the key to store slog handler.
)

// ContextWithAttrs returns a new context with given attributes.
// Use [AttrsFromContext] to extract attributes from the context.
// ContextWithAttrs uses a new context created with [context.Background]
// if the given ctx is nil.
func ContextWithAttrs(parent context.Context, attrs ...any) context.Context {
	if parent == nil {
		return context.WithValue(context.Background(), ctxAttrsKey, attrs)
	}
	if v := parent.Value(ctxAttrsKey); v != nil {
		return context.WithValue(parent, ctxAttrsKey, append(v.([]any), attrs...))
	}
	return context.WithValue(parent, ctxAttrsKey, attrs)
}

// AttrsFromContext returns a new context with given log levels.
// Use [ContextWithAttrs] to store log attributes in context.
// AttrsFromContext returns nil slice when no attributes found in the context.
// It panics when nil context was given.
func AttrsFromContext(ctx context.Context) []any {
	if v := ctx.Value(ctxAttrsKey); v != nil {
		return v.([]any)
	}
	return nil
}

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
// LevelFromContext returns [math.MinInt] when no levels found in the context.
// It panics when nil context was given.
func LevelFromContext(ctx context.Context) slog.Level {
	if v := ctx.Value(ctxLevelKey); v != nil {
		return v.(slog.Level)
	}
	return math.MinInt
}

// ContextWithHandler returns a new context with given slog handler.
// Use [HandlerFromContext] to extract handler from context.
// ContextWithHandler uses a new context created with [context.Background]
// if the given ctx is nil.
func ContextWithHandler(parent context.Context, h slog.Handler) context.Context {
	if parent == nil {
		return context.WithValue(context.Background(), ctxHandlerKey, h)
	}
	return context.WithValue(parent, ctxHandlerKey, h)
}

// HandlerFromContext returns a slog handler stored in the context.
// Use [ContextWithHandler] to store a slog handler in context.
// HandlerFromContext returns nil when no handler found in the context.
// It panics when nil context was given.
func HandlerFromContext(ctx context.Context) slog.Handler {
	if v := ctx.Value(ctxHandlerKey); v != nil {
		return v.(slog.Handler)
	}
	return nil
}
