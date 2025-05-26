package zslog_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zlog/zslog"
	"github.com/aileron-projects/go/ztesting"
)

func TestNewJSON(t *testing.T) {
	t.Parallel()
	t.Run("nil writer", func(t *testing.T) {
		t.Run("nil option", func(t *testing.T) {
			lg := zslog.NewJSON(nil, nil)
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stdout), lg.Writer())
			ztesting.AssertEqual(t, "log level is not enabled.", true, lg.Handler().Enabled(nil, slog.LevelInfo))
		})
		t.Run("non-nil option", func(t *testing.T) {
			lg := zslog.NewJSON(nil, &slog.HandlerOptions{Level: slog.LevelDebug})
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stdout), lg.Writer())
			ztesting.AssertEqual(t, "log level is not enabled.", true, lg.Handler().Enabled(nil, slog.LevelDebug))
		})
	})
	t.Run("non-nil writer", func(t *testing.T) {
		t.Run("nil option", func(t *testing.T) {
			lg := zslog.NewJSON(os.Stderr, nil)
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stderr), lg.Writer())
			ztesting.AssertEqual(t, "log level is not enabled.", false, lg.Handler().Enabled(nil, slog.LevelDebug))
		})
		t.Run("non-nil option", func(t *testing.T) {
			lg := zslog.NewJSON(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stderr), lg.Writer())
			ztesting.AssertEqual(t, "log level is not enabled.", true, lg.Handler().Enabled(nil, slog.LevelDebug))
		})
	})
}

func TestNewText(t *testing.T) {
	t.Parallel()
	t.Run("nil writer", func(t *testing.T) {
		t.Run("nil option", func(t *testing.T) {
			lg := zslog.NewText(nil, nil)
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stdout), lg.Writer())
			ztesting.AssertEqual(t, "log level is not enabled.", true, lg.Handler().Enabled(nil, slog.LevelInfo))
		})
		t.Run("non-nil option", func(t *testing.T) {
			lg := zslog.NewText(nil, &slog.HandlerOptions{Level: slog.LevelDebug})
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stdout), lg.Writer())
			ztesting.AssertEqual(t, "log level is not enabled.", true, lg.Handler().Enabled(nil, slog.LevelDebug))
		})
	})
	t.Run("non-nil writer", func(t *testing.T) {
		t.Run("nil option", func(t *testing.T) {
			lg := zslog.NewText(os.Stderr, nil)
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stderr), lg.Writer())
			ztesting.AssertEqual(t, "log level is not enabled.", true, lg.Handler().Enabled(nil, slog.LevelInfo))
		})
		t.Run("non-nil option", func(t *testing.T) {
			lg := zslog.NewText(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stderr), lg.Writer())
			ztesting.AssertEqual(t, "log level is not enabled.", true, lg.Handler().Enabled(nil, slog.LevelDebug))
		})
	})
}

var testLevels = map[string]struct {
	setLv slog.Level
	ctxLv slog.Level
}{
	"debug-1":             {slog.LevelDebug - 1, slog.Level(math.MinInt)},
	"debug":               {slog.LevelDebug, slog.Level(math.MinInt)},
	"info":                {slog.LevelInfo, slog.Level(math.MinInt)},
	"warn":                {slog.LevelWarn, slog.Level(math.MinInt)},
	"error":               {slog.LevelError, slog.Level(math.MinInt)},
	"error+1":             {slog.LevelError + 1, slog.Level(math.MinInt)},
	"debug-1 ctx=debug-1": {slog.LevelDebug - 1, slog.LevelDebug - 1},
	"debug-1 ctx=debug":   {slog.LevelDebug - 1, slog.LevelDebug},
	"debug-1 ctx=info":    {slog.LevelDebug - 1, slog.LevelInfo},
	"debug-1 ctx=warn":    {slog.LevelDebug - 1, slog.LevelWarn},
	"debug-1 ctx=error":   {slog.LevelDebug - 1, slog.LevelError},
	"debug-1 ctx=error+1": {slog.LevelDebug - 1, slog.LevelError + 1},
	"debug ctx=debug-1":   {slog.LevelDebug, slog.LevelDebug - 1},
	"debug ctx=debug":     {slog.LevelDebug, slog.LevelDebug},
	"debug ctx=info":      {slog.LevelDebug, slog.LevelInfo},
	"debug ctx=warn":      {slog.LevelDebug, slog.LevelWarn},
	"debug ctx=error":     {slog.LevelDebug, slog.LevelError},
	"debug ctx=error+1":   {slog.LevelDebug, slog.LevelError + 1},
	"info ctx=debug-1":    {slog.LevelInfo, slog.LevelDebug - 1},
	"info ctx=debug":      {slog.LevelInfo, slog.LevelDebug},
	"info ctx=info":       {slog.LevelInfo, slog.LevelInfo},
	"info ctx=warn":       {slog.LevelInfo, slog.LevelWarn},
	"info ctx=error":      {slog.LevelInfo, slog.LevelError},
	"info ctx=error+1":    {slog.LevelInfo, slog.LevelError + 1},
	"warn ctx=debug-1":    {slog.LevelWarn, slog.LevelDebug - 1},
	"warn ctx=debug":      {slog.LevelWarn, slog.LevelDebug},
	"warn ctx=info":       {slog.LevelWarn, slog.LevelInfo},
	"warn ctx=warn":       {slog.LevelWarn, slog.LevelWarn},
	"warn ctx=error":      {slog.LevelWarn, slog.LevelError},
	"warn ctx=error+1":    {slog.LevelWarn, slog.LevelError + 1},
	"error ctx=debug-1":   {slog.LevelError, slog.LevelDebug - 1},
	"error ctx=debug":     {slog.LevelError, slog.LevelDebug},
	"error ctx=info":      {slog.LevelError, slog.LevelInfo},
	"error ctx=warn":      {slog.LevelError, slog.LevelWarn},
	"error ctx=error":     {slog.LevelError, slog.LevelError},
	"error ctx=error+1":   {slog.LevelError, slog.LevelError + 1},
	"error+1 ctx=debug-1": {slog.LevelError + 1, slog.LevelDebug - 1},
	"error+1 ctx=debug":   {slog.LevelError + 1, slog.LevelDebug},
	"error+1 ctx=info":    {slog.LevelError + 1, slog.LevelInfo},
	"error+1 ctx=warn":    {slog.LevelError + 1, slog.LevelWarn},
	"error+1 ctx=error":   {slog.LevelError + 1, slog.LevelError},
	"error+1 ctx=error+1": {slog.LevelError + 1, slog.LevelError + 1},
}

func TestZSLogger_Enabled(t *testing.T) {
	t.Parallel()
	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.ctxLv != slog.Level(math.MinInt) {
				threshold = tc.ctxLv
				ctx = zslog.ContextWithLevel(ctx, tc.ctxLv)
			}
			lg := zslog.NewJSON(nil, &slog.HandlerOptions{Level: tc.setLv})
			for i := -10; i < 10; i++ {
				want := i >= int(threshold)
				got := lg.Enabled(ctx, slog.Level(i))
				ztesting.AssertEqual(t, "enabled not match.", want, got)
			}
		})
	}
}

func TestZSLogger_Debug(t *testing.T) {
	t.Parallel()
	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.ctxLv != slog.Level(math.MinInt) {
				threshold = tc.ctxLv
				ctx = zslog.ContextWithLevel(ctx, tc.ctxLv)
			}
			var buf bytes.Buffer
			lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: tc.setLv})
			lg.DebugContext(ctx, "test message", "arg1", "arg2")
			result := buf.String()
			if threshold > slog.LevelDebug {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				ztesting.AssertEqual(t, "log line is not written.", true, strings.Contains(result, "test message"))
			}
		})
	}
	t.Run("add caller", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		lg.AddCaller = zslog.LvDebug
		lg.DebugContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})
	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		lg.AddFrames = zslog.LvDebug
		lg.DebugContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestZSLogger_Info(t *testing.T) {
	t.Parallel()
	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.ctxLv != slog.Level(math.MinInt) {
				threshold = tc.ctxLv
				ctx = zslog.ContextWithLevel(ctx, tc.ctxLv)
			}
			var buf bytes.Buffer
			lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: tc.setLv})
			lg.InfoContext(ctx, "test message", "arg1", "arg2")
			result := buf.String()
			if threshold > slog.LevelInfo {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				ztesting.AssertEqual(t, "log line is not written.", true, strings.Contains(result, "test message"))
			}
		})
	}
	t.Run("add caller", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
		lg.AddCaller = zslog.LvInfo
		lg.InfoContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})
	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
		lg.AddFrames = zslog.LvInfo
		lg.InfoContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestZSLogger_Warn(t *testing.T) {
	t.Parallel()
	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.ctxLv != slog.Level(math.MinInt) {
				threshold = tc.ctxLv
				ctx = zslog.ContextWithLevel(ctx, tc.ctxLv)
			}
			var buf bytes.Buffer
			lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: tc.setLv})
			lg.WarnContext(ctx, "test message", "arg1", "arg2")
			result := buf.String()
			if threshold > slog.LevelWarn {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				ztesting.AssertEqual(t, "log line is not written.", true, strings.Contains(result, "test message"))
			}
		})
	}
	t.Run("add caller", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
		lg.AddCaller = zslog.LvWarn
		lg.WarnContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})
	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
		lg.AddFrames = zslog.LvWarn
		lg.WarnContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestZSLogger_Error(t *testing.T) {
	t.Parallel()
	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.ctxLv != slog.Level(math.MinInt) {
				threshold = tc.ctxLv
				ctx = zslog.ContextWithLevel(ctx, tc.ctxLv)
			}
			var buf bytes.Buffer
			lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: tc.setLv})
			lg.ErrorContext(ctx, "test message", "arg1", "arg2")
			result := buf.String()
			if threshold > slog.LevelError {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				ztesting.AssertEqual(t, "log line is not written.", true, strings.Contains(result, "test message"))
			}
		})
	}
	t.Run("add caller", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelError})
		lg.AddCaller = zslog.LvError
		lg.ErrorContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})
	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelError})
		lg.AddFrames = zslog.LvError
		lg.ErrorContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestZSLogger_nilContext(t *testing.T) {
	t.Parallel()
	t.Run("enabled", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
		ztesting.AssertEqual(t, "log level is unexpectedly enabled.", false, lg.Enabled(nil, slog.LevelDebug))
		ztesting.AssertEqual(t, "log level is unexpectedly disabled.", true, lg.Enabled(nil, slog.LevelInfo))
		ztesting.AssertEqual(t, "log level is unexpectedly disabled.", true, lg.Enabled(nil, slog.LevelWarn))
	})
	t.Run("debug", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		lg.DebugContext(nil, "test message")
		ztesting.AssertEqual(t, "log line should be written", false, buf.String() == "")
	})
	t.Run("info", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		lg.InfoContext(nil, "test message")
		ztesting.AssertEqual(t, "log line should be written", false, buf.String() == "")
	})
	t.Run("warn", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		lg.WarnContext(nil, "test message")
		ztesting.AssertEqual(t, "log line should be written", false, buf.String() == "")
	})
	t.Run("error", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		lg.ErrorContext(nil, "test message")
		ztesting.AssertEqual(t, "log line should be written", false, buf.String() == "")
	})
}
