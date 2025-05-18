package zslog_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zlog"
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
	setLv    zlog.Level
	setCtxLv zlog.Level
}{
	"trace":           {zlog.LvTrace, zlog.LvUndef},
	"debug":           {zlog.LvDebug, zlog.LvUndef},
	"info":            {zlog.LvInfo, zlog.LvUndef},
	"warn":            {zlog.LvWarn, zlog.LvUndef},
	"error":           {zlog.LvError, zlog.LvUndef},
	"fatal":           {zlog.LvFatal, zlog.LvUndef},
	"trace ctx=trace": {zlog.LvTrace, zlog.LvTrace},
	"trace ctx=debug": {zlog.LvTrace, zlog.LvDebug},
	"trace ctx=info":  {zlog.LvTrace, zlog.LvInfo},
	"trace ctx=warn":  {zlog.LvTrace, zlog.LvWarn},
	"trace ctx=error": {zlog.LvTrace, zlog.LvError},
	"trace ctx=fatal": {zlog.LvTrace, zlog.LvFatal},
	"debug ctx=trace": {zlog.LvDebug, zlog.LvTrace},
	"debug ctx=debug": {zlog.LvDebug, zlog.LvDebug},
	"debug ctx=info":  {zlog.LvDebug, zlog.LvInfo},
	"debug ctx=warn":  {zlog.LvDebug, zlog.LvWarn},
	"debug ctx=error": {zlog.LvDebug, zlog.LvError},
	"debug ctx=fatal": {zlog.LvDebug, zlog.LvFatal},
	"info ctx=trace":  {zlog.LvInfo, zlog.LvTrace},
	"info ctx=debug":  {zlog.LvInfo, zlog.LvDebug},
	"info ctx=info":   {zlog.LvInfo, zlog.LvInfo},
	"info ctx=warn":   {zlog.LvInfo, zlog.LvWarn},
	"info ctx=error":  {zlog.LvInfo, zlog.LvError},
	"info ctx=fatal":  {zlog.LvInfo, zlog.LvFatal},
	"warn ctx=trace":  {zlog.LvWarn, zlog.LvTrace},
	"warn ctx=debug":  {zlog.LvWarn, zlog.LvDebug},
	"warn ctx=info":   {zlog.LvWarn, zlog.LvInfo},
	"warn ctx=warn":   {zlog.LvWarn, zlog.LvWarn},
	"warn ctx=error":  {zlog.LvWarn, zlog.LvError},
	"warn ctx=fatal":  {zlog.LvWarn, zlog.LvFatal},
	"error ctx=trace": {zlog.LvError, zlog.LvTrace},
	"error ctx=debug": {zlog.LvError, zlog.LvDebug},
	"error ctx=info":  {zlog.LvError, zlog.LvInfo},
	"error ctx=warn":  {zlog.LvError, zlog.LvWarn},
	"error ctx=error": {zlog.LvError, zlog.LvError},
	"error ctx=fatal": {zlog.LvError, zlog.LvFatal},
	"fatal ctx=trace": {zlog.LvFatal, zlog.LvTrace},
	"fatal ctx=debug": {zlog.LvFatal, zlog.LvDebug},
	"fatal ctx=info":  {zlog.LvFatal, zlog.LvInfo},
	"fatal ctx=warn":  {zlog.LvFatal, zlog.LvWarn},
	"fatal ctx=error": {zlog.LvFatal, zlog.LvError},
	"fatal ctx=fatal": {zlog.LvFatal, zlog.LvFatal},
}

func TestZSLogger_Enabled(t *testing.T) {
	t.Parallel()

	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.setCtxLv != zlog.LvUndef {
				threshold = tc.setCtxLv
				ctx = zlog.ContextWithLevel(ctx, tc.setCtxLv)
			}

			lg := zslog.NewJSON(nil, &slog.HandlerOptions{Level: slog.Level(tc.setLv - 9)}) // Convert zlog.Level to slog.Level
			for i := -1; i < int(threshold); i++ {
				ztesting.AssertEqual(t, "log level is not enabled.", false, lg.Enabled(ctx, zlog.Level(i)))
			}
			for i := int(threshold); i < 25; i++ {
				ztesting.AssertEqual(t, "log level is enabled.", true, lg.Enabled(ctx, zlog.Level(i)))
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
			if tc.setCtxLv != zlog.LvUndef {
				threshold = tc.setCtxLv
				ctx = zlog.ContextWithLevel(ctx, tc.setCtxLv)
			}

			var buf bytes.Buffer
			lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.Level(tc.setLv - 9)})
			lg.Debug(ctx, "test message", "arg1", "arg2")
			result := buf.String()
			if threshold > zlog.LvDebug {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				ztesting.AssertEqual(t, "log line is not written.", true, strings.Contains(result, "test message"))
			}
		})
	}

	t.Run("add caller", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		lg.AddCaller = zlog.Debug
		lg.Debug(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})

	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		lg.AddFrames = zlog.Debug
		lg.Debug(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestZSLogger_Info(t *testing.T) {
	t.Parallel()

	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.setCtxLv != zlog.LvUndef {
				threshold = tc.setCtxLv
				ctx = zlog.ContextWithLevel(ctx, tc.setCtxLv)
			}

			var buf bytes.Buffer
			lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.Level(tc.setLv - 9)})
			lg.Info(ctx, "test message", "arg1", "arg2")
			result := buf.String()
			if threshold > zlog.LvInfo {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				ztesting.AssertEqual(t, "log line is not written.", true, strings.Contains(result, "test message"))
			}
		})
	}

	t.Run("add caller", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
		lg.AddCaller = zlog.Info
		lg.Info(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})

	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
		lg.AddFrames = zlog.Info
		lg.Info(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestZSLogger_Warn(t *testing.T) {
	t.Parallel()

	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.setCtxLv != zlog.LvUndef {
				threshold = tc.setCtxLv
				ctx = zlog.ContextWithLevel(ctx, tc.setCtxLv)
			}

			var buf bytes.Buffer
			lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.Level(tc.setLv - 9)})
			lg.Warn(ctx, "test message", "arg1", "arg2")
			result := buf.String()
			if threshold > zlog.LvWarn {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				ztesting.AssertEqual(t, "log line is not written.", true, strings.Contains(result, "test message"))
			}
		})
	}

	t.Run("add caller", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
		lg.AddCaller = zlog.Warn
		lg.Warn(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})

	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
		lg.AddFrames = zlog.Warn
		lg.Warn(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestZSLogger_Error(t *testing.T) {
	t.Parallel()

	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.setCtxLv != zlog.LvUndef {
				threshold = tc.setCtxLv
				ctx = zlog.ContextWithLevel(ctx, tc.setCtxLv)
			}

			var buf bytes.Buffer
			lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.Level(tc.setLv - 9)})
			lg.Error(ctx, "test message", "arg1", "arg2")
			result := buf.String()
			if threshold > zlog.LvError {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				ztesting.AssertEqual(t, "log line is not written.", true, strings.Contains(result, "test message"))
			}
		})
	}

	t.Run("add caller", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelError})
		lg.AddCaller = zlog.Error
		lg.Error(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain caller.", true, strings.Contains(buf.String(), `"caller"`))
	})

	t.Run("add frames", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.LevelError})
		lg.AddFrames = zlog.Error
		lg.Error(nil, "test message")
		ztesting.AssertEqual(t, "log line does not contain frames.", true, strings.Contains(buf.String(), `"frames"`))
	})
}

func TestZSLogger_nilContext(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.Level(zlog.LvInfo - 9)})
		ztesting.AssertEqual(t, "log level is unexpectedly enabled.", false, lg.Enabled(nil, zlog.LvDebug))
		ztesting.AssertEqual(t, "log level is unexpectedly disabled.", true, lg.Enabled(nil, zlog.LvInfo))
		ztesting.AssertEqual(t, "log level is unexpectedly disabled.", true, lg.Enabled(nil, zlog.LvWarn))
	})

	t.Run("debug", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.Level(zlog.LvTrace - 9)})
		lg.Debug(nil, "test message")
		ztesting.AssertNotEqual(t, "log line is not written.", "", buf.String())
	})

	t.Run("info", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.Level(zlog.LvTrace - 9)})
		lg.Info(nil, "test message")
		ztesting.AssertNotEqual(t, "log line is not written.", "", buf.String())
	})

	t.Run("warn", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.Level(zlog.LvTrace - 9)})
		lg.Warn(nil, "test message")
		ztesting.AssertNotEqual(t, "log line is not written.", "", buf.String())
	})

	t.Run("error", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.NewJSON(&buf, &slog.HandlerOptions{Level: slog.Level(zlog.LvTrace - 9)})
		lg.Error(nil, "test message")
		ztesting.AssertNotEqual(t, "log line is not written.", "", buf.String())
	})
}
