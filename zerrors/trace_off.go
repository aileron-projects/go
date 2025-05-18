//go:build !zerrorstrace

package zerrors

const (
	// traceEnabled represents the error tracing is enabled.
	// Use "-tags zerrorstrace" build tag to enable tracing.
	traceEnabled = false
)
