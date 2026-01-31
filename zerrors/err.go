package zerrors

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

// NewErr returns new [Err].
// Detail is the format string for [smt.Sprintf].
func NewErr(cause error, message, detail string, a ...any) *Err {
	return &Err{
		Cause:   cause,
		Message: message,
		Detail:  fmt.Sprintf(detail, a...),
	}
}

// Err is the simple general error object.
// Use [NewErr] to create instances.
type Err struct {
	// Cause is the cause of this error.
	Cause error
	// Message is the fixed error message.
	// Message is compared in the [Err.Is].
	Message string
	// Detail is the error detail.
	// Detail is NOT compared in the [Err.Is].
	Detail string
}

// Error implements [error] interface.
func (e *Err) Error() string {
	var b strings.Builder
	b.Grow(len(e.Message) + len(e.Detail) + 1)
	_, _ = b.WriteString(e.Message)
	if e.Detail != "" {
		_, _ = b.WriteString(" ")
		_, _ = b.WriteString(e.Detail)
	}
	if e.Cause != nil {
		_, _ = b.WriteString(" [")
		_, _ = b.WriteString(e.Cause.Error())
		_, _ = b.WriteString("]")
	}
	return b.String()
}

// Unwrap returns the inner error if any.
func (e *Err) Unwrap() error {
	return e.Cause
}

// Is returns if this error is identical to the given error.
// This can be used with [errors.Is].
func (e *Err) Is(err error) bool {
	if err == nil || e == nil {
		return e == err
	}
	for err != nil {
		ee, ok := err.(*Err)
		if ok {
			return e.Message == ee.Message
		}
		err = errors.Unwrap(err)
	}
	return false
}

// Map returns error information in map.
func (e *Err) Map() map[string]any {
	m := map[string]any{
		"message": e.Message,
		"detail":  e.Detail,
	}
	if cause := ToMap(e.Cause); cause != nil {
		m["cause"] = cause
	}
	return m
}

// Map returns error information in slice of [slog.Attr].
func (e *Err) SlogAttrs() []slog.Attr {
	a := []slog.Attr{
		slog.String("message", e.Message),
		slog.String("detail", e.Detail),
	}
	if causes := ToSlogAttrs(e.Cause); len(causes) > 0 {
		a = append(a, slog.GroupAttrs("cause", causes...))
	}
	return a
}
