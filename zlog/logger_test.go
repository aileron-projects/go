package zlog_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"testing"

	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/ztesting"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("nil writer", func(t *testing.T) {
		t.Run("nil option", func(t *testing.T) {
			lg := zlog.New(nil, nil)
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stdout), lg.Writer())
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stdout), lg.Logger().Writer())
			ztesting.AssertEqual(t, "logger uses wrong flag.", log.LstdFlags, lg.Logger().Flags())
		})
		t.Run("non-nil option", func(t *testing.T) {
			lg := zlog.New(nil, &zlog.LoggerOption{Lv: zlog.LvDebug, Flag: log.LUTC})
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stdout), lg.Writer())
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stdout), lg.Logger().Writer())
			ztesting.AssertEqual(t, "logger uses wrong flag.", log.LUTC, lg.Logger().Flags())
		})
	})

	t.Run("non-nil writer", func(t *testing.T) {
		t.Run("nil option", func(t *testing.T) {
			lg := zlog.New(os.Stderr, nil)
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stderr), lg.Writer())
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stderr), lg.Logger().Writer())
			ztesting.AssertEqual(t, "logger uses wrong flag.", log.LstdFlags, lg.Logger().Flags())
		})
		t.Run("non-nil option", func(t *testing.T) {
			lg := zlog.New(os.Stderr, &zlog.LoggerOption{Lv: zlog.LvDebug, Flag: log.LUTC})
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stderr), lg.Writer())
			ztesting.AssertEqual(t, "logger uses wrong io writer.", io.Writer(os.Stderr), lg.Logger().Writer())
			ztesting.AssertEqual(t, "logger uses wrong flag.", log.LUTC, lg.Logger().Flags())
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

func TestZLogger_Enabled(t *testing.T) {
	t.Parallel()
	for name, tc := range testLevels {
		t.Run(name, func(t *testing.T) {
			threshold := tc.setLv
			ctx := context.Background()
			if tc.setCtxLv != zlog.LvUndef {
				threshold = tc.setCtxLv
				ctx = zlog.ContextWithLevel(ctx, tc.setCtxLv)
			}

			lg := zlog.New(nil, &zlog.LoggerOption{Lv: tc.setLv})
			for i := -1; i < int(threshold); i++ {
				ztesting.AssertEqual(t, "log level is unexpectedly enabled.", false, lg.Enabled(ctx, zlog.Level(i)))
			}
			for i := int(threshold); i < 25; i++ {
				ztesting.AssertEqual(t, "log level is unexpectedly disabled.", true, lg.Enabled(ctx, zlog.Level(i)))
			}
		})
	}
}

func TestZLogger_Debug(t *testing.T) {
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
			lg := zlog.New(&buf, &zlog.LoggerOption{Lv: tc.setLv})
			lg.Debug(ctx, "test message", "arg1", "arg2")
			if threshold > zlog.LvDebug {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				want := "level=DEBUG msg=test message arg1 arg2\n"
				ztesting.AssertEqual(t, "unexpected log message.", want, buf.String())
			}
		})
	}
}

func TestZLogger_Info(t *testing.T) {
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
			lg := zlog.New(&buf, &zlog.LoggerOption{Lv: tc.setLv})
			lg.Info(ctx, "test message", "arg1", "arg2")
			if threshold > zlog.LvInfo {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				want := "level=INFO msg=test message arg1 arg2\n"
				ztesting.AssertEqual(t, "unexpected log message.", want, buf.String())
			}
		})
	}
}

func TestZLogger_Warn(t *testing.T) {
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
			lg := zlog.New(&buf, &zlog.LoggerOption{Lv: tc.setLv})
			lg.Warn(ctx, "test message", "arg1", "arg2")
			if threshold > zlog.LvWarn {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				want := "level=WARN msg=test message arg1 arg2\n"
				ztesting.AssertEqual(t, "unexpected log message.", want, buf.String())
			}
		})
	}
}

func TestZLogger_Error(t *testing.T) {
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
			lg := zlog.New(&buf, &zlog.LoggerOption{Lv: tc.setLv})
			lg.Error(ctx, "test message", "arg1", "arg2")
			if threshold > zlog.LvError {
				ztesting.AssertEqual(t, "log line is written.", "", buf.String())
			} else {
				want := "level=ERROR msg=test message arg1 arg2\n"
				ztesting.AssertEqual(t, "unexpected log message.", want, buf.String())
			}
		})
	}
}

func TestZLogger_nilContext(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zlog.New(&buf, &zlog.LoggerOption{Lv: zlog.LvInfo})
		ztesting.AssertEqual(t, "log levels unexpectedly enabled.", false, lg.Enabled(nil, zlog.LvDebug))
		ztesting.AssertEqual(t, "log levels unexpectedly enabled.", true, lg.Enabled(nil, zlog.LvInfo))
		ztesting.AssertEqual(t, "log levels unexpectedly enabled.", true, lg.Enabled(nil, zlog.LvWarn))
	})

	t.Run("debug", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zlog.New(&buf, &zlog.LoggerOption{Lv: zlog.LvTrace})
		lg.Debug(nil, "test message")
		ztesting.AssertEqual(t, "wrong log line.", "level=DEBUG msg=test message\n", buf.String())
	})

	t.Run("info", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zlog.New(&buf, &zlog.LoggerOption{Lv: zlog.LvTrace})
		lg.Info(nil, "test message")
		ztesting.AssertEqual(t, "wrong log line.", "level=INFO msg=test message\n", buf.String())
	})

	t.Run("warn", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zlog.New(&buf, &zlog.LoggerOption{Lv: zlog.LvTrace})
		lg.Warn(nil, "test message")
		ztesting.AssertEqual(t, "wrong log line.", "level=WARN msg=test message\n", buf.String())
	})

	t.Run("error", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zlog.New(&buf, &zlog.LoggerOption{Lv: zlog.LvTrace})
		lg.Error(nil, "test message")
		ztesting.AssertEqual(t, "wrong log line.", "level=ERROR msg=test message\n", buf.String())
	})
}
