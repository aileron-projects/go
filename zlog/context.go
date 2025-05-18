package zlog

import (
	"context"
)

// ctxAttrs is the type to store attributes in a context.
type ctxAttrs struct{}

// ctxLevel is the type to store log level in a context.
type ctxLevel struct{}

var (
	// ctxAttrsKey is the key to store attributes in a context.
	ctxAttrsKey = &ctxAttrs{}
	// ctxLevelKey is the key to store log level in a context.
	ctxLevelKey = &ctxLevel{}
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
// AttrsFromContext returns nil slice if the given ctx is nil or
// no attributes were found in the context.
func AttrsFromContext(ctx context.Context) []any {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(ctxAttrsKey); v != nil {
		return v.([]any)
	}
	return nil
}

// ContextWithLevel returns a new context with given log level.
// Use [LevelFromContext] to extract log level from the context.
// ContextWithLevel uses a new context created with [context.Background]
// if the given ctx is nil.
func ContextWithLevel(parent context.Context, lv Level) context.Context {
	if parent == nil {
		return context.WithValue(context.Background(), ctxLevelKey, lv)
	}
	return context.WithValue(parent, ctxLevelKey, lv)
}

// LevelFromContext returns a log level stored in the context.
// Use [ContextWithLevel] to store a log level in context.
// LevelFromContext returns [LvUndef] if the given ctx is nil
// of no log levels are found in the context.
func LevelFromContext(ctx context.Context) Level {
	if ctx == nil {
		return LvUndef
	}
	if v := ctx.Value(ctxLevelKey); v != nil {
		return v.(Level)
	}
	return LvUndef
}
