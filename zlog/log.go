package zlog

import (
	"context"
)

// Logger is the basic logger interface.
type Logger interface {
	DebugEnabled(ctx context.Context) bool
	InfoEnabled(ctx context.Context) bool
	WarnEnabled(ctx context.Context) bool
	ErrorEnabled(ctx context.Context) bool
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

// ctxKey is the context key type.
type ctxKey struct{ string }

// ctxAttrsKey is the key to store attributes in a context.
var ctxAttrsKey = &ctxKey{"attrs"}

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
