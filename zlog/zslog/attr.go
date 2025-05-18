package zslog

import (
	"log/slog"
	"runtime"
	"strconv"
	"time"

	"github.com/aileron-projects/go/zerrors"
	"github.com/aileron-projects/go/zruntime"
)

// CallerAttr returns the following attributes
// with "caller" key.
//
//   - pkg: caller's package name.
//   - file: caller's file name.
//   - func: caller's file name.
//   - line: caller's line number.
func CallerAttr(skip int) slog.Attr {
	info := zruntime.ConvertFrame(zruntime.CallerFrame(skip + 1))
	return slog.Group("caller",
		slog.String("pkg", info.Pkg),
		slog.String("file", info.File),
		slog.String("func", info.Func),
		slog.Int("line", info.Line),
	)
}

// DateTimeAttr returns the following attributes
// with "datetime" key.
//
//   - date: current date in "2006-01-02" format.
//   - time: current time in "15:04:05.999" format.
func DateTimeAttr() slog.Attr {
	t := time.Now().Local()
	return slog.Group("datetime",
		slog.String("date", t.Format(time.DateOnly)),
		slog.String("time", t.Format(time.TimeOnly+".999")),
	)
}

// FramesAttr returns the following attributes
// with "frames" key.
//
//   - pkg: caller's package name.
//   - file: caller's file name.
//   - func: caller's file name.
//   - line: caller's line number.
func FramesAttr(skip int) slog.Attr {
	frames := zruntime.ConvertFrames(zruntime.CallerFrames(skip + 1))
	attrs := make([]string, len(frames))
	for i, f := range frames {
		attrs[i] = f.Pkg + ":" + f.File + ":L" + strconv.Itoa(f.Line) + "(" + f.Func + ")"
	}
	return slog.Any("frames", attrs)
}

// StackTraceAttrs returns stack trace
// as slog attributes with "error" key.
// The stack trace has 4kiB at maximum.
func StackTraceAttrs(skip int) slog.Attr {
	stack := make([]byte, 1<<12) // Read max 4kiB.
	n := runtime.Stack(stack, false)
	return slog.String("stack", string(stack[:n]))
}

// ErrorAttr returns the error
// as slog attributes with "error" key.
func ErrorAttr(err error) slog.Attr {
	return slog.Any("error", zerrors.Attrs(err))
}
