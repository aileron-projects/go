package zslog

import (
	"log/slog"
)

// RemoveTime removes time attribute from slog record.
// RemoveTime is intended to be used as a options of
// [slog.HandlerOptions.ReplaceAttr].
func RemoveTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}
