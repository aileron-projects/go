package zslog

import (
	"log/slog"

	"github.com/aileron-projects/go/zlog"
)

func RemoveTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}

// toZLevel converts [slog.Level] to [zlog.Level].
func toZLevel(lv slog.Level) zlog.Level {
	return zlog.Level(lv + 9)
}

// toSLevel converts [zlog.Level] to [slog.Level].
func toSLevel(lv zlog.Level) slog.Level {
	return slog.Level(lv - 9)
}
