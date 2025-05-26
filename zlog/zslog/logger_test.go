package zslog_test

import (
	"bytes"
	"context"
	"log/slog"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zlog/zslog"
	"github.com/aileron-projects/go/ztesting"
)

func TestDefaultLogger(t *testing.T) {
	tmp := zslog.Default()
	t.Cleanup(func() { zslog.SetDefault(tmp) })
	var w bytes.Buffer
	lg := zslog.New(slog.NewJSONHandler(&w, &slog.HandlerOptions{}))
	zslog.SetDefault(lg)
	dlg := zslog.Default()
	dlg.InfoContext(context.Background(), "test", "foo", "bar")
	out := w.String()
	ztesting.AssertEqual(t, "message not found", true, strings.Contains(out, `"msg":"test"`))
	ztesting.AssertEqual(t, "message not found", true, strings.Contains(out, `"foo":"bar"`))
}

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("json handler", func(t *testing.T) {
		var w bytes.Buffer
		lg := zslog.New(slog.NewJSONHandler(&w, &slog.HandlerOptions{Level: slog.LevelDebug}))
		ztesting.AssertEqual(t, "log level is not enabled.", true, lg.DebugEnabled(context.Background()))
		lg.InfoContext(context.Background(), "test", "foo", "bar")
		out := w.String()
		ztesting.AssertEqual(t, "message not found", true, strings.Contains(out, `"msg":"test"`))
		ztesting.AssertEqual(t, "message not found", true, strings.Contains(out, `"foo":"bar"`))
	})
	t.Run("text handler", func(t *testing.T) {
		var w bytes.Buffer
		lg := zslog.New(slog.NewTextHandler(&w, &slog.HandlerOptions{Level: slog.LevelDebug}))
		ztesting.AssertEqual(t, "log level is not enabled.", true, lg.DebugEnabled(context.Background()))
		lg.InfoContext(context.Background(), "test", "foo", "bar")
		out := w.String()
		ztesting.AssertEqual(t, "message not found", true, strings.Contains(out, `msg=test`))
		ztesting.AssertEqual(t, "message not found", true, strings.Contains(out, `foo=bar`))
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

func TestLogger_Enabled(t *testing.T) {
	t.Parallel()
	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.ctxLv != slog.Level(math.MinInt) {
				threshold = tc.ctxLv
				ctx = zslog.ContextWithLevel(ctx, tc.ctxLv)
			}
			lg := zslog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: tc.setLv}))
			ztesting.AssertEqual(t, "debug enabled not match.", slog.LevelDebug >= threshold, lg.DebugEnabled(ctx))
			ztesting.AssertEqual(t, "info enabled not match.", slog.LevelInfo >= threshold, lg.InfoEnabled(ctx))
			ztesting.AssertEqual(t, "warn enabled not match.", slog.LevelWarn >= threshold, lg.WarnEnabled(ctx))
			ztesting.AssertEqual(t, "error enabled not match.", slog.LevelError >= threshold, lg.ErrorEnabled(ctx))
			for i := -10; i < 10; i++ {
				ztesting.AssertEqual(t, "enabled not match.", i >= int(threshold), lg.Enabled(ctx, slog.Level(i)))
			}
		})
	}
}

func TestLogger_Debug(t *testing.T) {
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
			lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: tc.setLv}))
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
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		lg.AddCaller = zslog.RangeDebug
		lg.DebugContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})
	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		lg.AddFrames = zslog.RangeDebug
		lg.DebugContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestLogger_Info(t *testing.T) {
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
			lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: tc.setLv}))
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
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
		lg.AddCaller = zslog.RangeInfo
		lg.InfoContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})
	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
		lg.AddFrames = zslog.RangeInfo
		lg.InfoContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestLogger_Warn(t *testing.T) {
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
			lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: tc.setLv}))
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
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
		lg.AddCaller = zslog.RangeWarn
		lg.WarnContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})
	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
		lg.AddFrames = zslog.RangeWarn
		lg.WarnContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestLogger_Error(t *testing.T) {
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
			lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: tc.setLv}))
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
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelError}))
		lg.AddCaller = zslog.RangeError
		lg.ErrorContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})
	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelError}))
		lg.AddFrames = zslog.RangeError
		lg.ErrorContext(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestLogger_nilContext(t *testing.T) {
	t.Parallel()
	t.Run("enabled", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
		ztesting.AssertEqual(t, "log level is unexpectedly enabled.", false, lg.Enabled(nil, slog.LevelDebug))
		ztesting.AssertEqual(t, "log level is unexpectedly disabled.", true, lg.Enabled(nil, slog.LevelInfo))
		ztesting.AssertEqual(t, "log level is unexpectedly disabled.", true, lg.Enabled(nil, slog.LevelWarn))
	})
	t.Run("debug", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		lg.DebugContext(nil, "test message")
		ztesting.AssertEqual(t, "log line should be written", false, buf.String() == "")
	})
	t.Run("info", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		lg.InfoContext(nil, "test message")
		ztesting.AssertEqual(t, "log line should be written", false, buf.String() == "")
	})
	t.Run("warn", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		lg.WarnContext(nil, "test message")
		ztesting.AssertEqual(t, "log line should be written", false, buf.String() == "")
	})
	t.Run("error", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		lg.ErrorContext(nil, "test message")
		ztesting.AssertEqual(t, "log line should be written", false, buf.String() == "")
	})
}
