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
