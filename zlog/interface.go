package zlog

import (
	"context"
)

// Range is the log range type.
// Defined ranges are:
//   - Undefined
//   - Trace
//   - Debug
//   - Info
//   - Warn
//   - Error
//   - Fatal
type Range uint

const (
	Undefined Range = 1 << iota
	Trace
	Debug
	Info
	Warn
	Error
	Fatal
)

// String returns string representation of this log level range.
func (r Range) String() string {
	switch r {
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "UNDEFINED"
	}
}

// Level is the log level type.
// Level implements [fmt.Stringer].
// Defined levels are:
//   - LvTrace
//   - LvDebug
//   - LvInfo
//   - LvWarn
//   - LvError
//   - LvFatal
type Level int

const (
	lvMin    Level = iota
	lvTrace1       // TRACE1
	lvTrace2       // TRACE2
	lvTrace3       // TRACE3
	lvTrace4       // TRACE4
	lvDebug1       // DEBUG1
	lvDebug2       // DEBUG2
	lvDebug3       // DEBUG3
	lvDebug4       // DEBUG4
	lvInfo1        // INFO1
	lvInfo2        // INFO2
	lvInfo3        // INFO3
	lvInfo4        // INFO4
	lvWarn1        // WARN1
	lvWarn2        // WARN2
	lvWarn3        // WARN3
	lvWarn4        // WARN4
	lvError1       // ERROR1
	lvError2       // ERROR2
	lvError3       // ERROR3
	lvError4       // ERROR4
	lvFatal1       // FATAL1
	lvFatal2       // FATAL2
	lvFatal3       // FATAL3
	lvFatal4       // FATAL4
	lvMax

	LvUndef = lvMin    // UNDEFINED
	LvTrace = lvTrace1 // TRACE
	LvDebug = lvDebug1 // DEBUG
	LvInfo  = lvInfo1  // INFO
	LvWarn  = lvWarn1  // WARN
	LvError = lvError1 // ERROR
	LvFatal = lvFatal1 // FATAL
)

// HigherEqual returns the result of lv>=target.
func (lv Level) HigherEqual(target Level) bool {
	return lv >= target
}

// HigherThan returns the result of lv>target.
func (lv Level) HigherThan(target Level) bool {
	return lv > target
}

// LessEqual returns the result of lv<=target.
func (lv Level) LessEqual(target Level) bool {
	return lv <= target
}

// LessThan returns the result of lv<target.
func (lv Level) LessThan(target Level) bool {
	return lv < target
}

// Range returns this log leven [Range].
func (lv Level) Range() Range {
	if lv <= 0 {
		return Undefined
	}
	switch {
	case lv-lvTrace1 < 4:
		return Trace
	case lv-lvDebug1 < 4:
		return Debug
	case lv-lvInfo1 < 4:
		return Info
	case lv-lvWarn1 < 4:
		return Warn
	case lv-lvError1 < 4:
		return Error
	case lv-lvFatal1 < 4:
		return Fatal
	default:
		return Undefined
	}
}

// String returns string representation of this log level.
func (lv Level) String() string {
	switch lv.Range() {
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "UNDEFINED"
	}
}

// Logger logs given message and values.
type Logger interface {
	Enabled(ctx context.Context, lv Level) bool
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}
