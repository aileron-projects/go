package zruntime

import (
	"runtime"
	"strings"
)

// Frame holds stack frame information.
// See also [runtime.Frame].
type Frame struct {
	// Pkg is go package name of the caller.
	Pkg string `json:"pkg" msgpack:"pkg" xml:"pkg" yaml:"pkg"`
	// File is the file name of the caller.
	File string `json:"file" msgpack:"file" xml:"file" yaml:"file"`
	// Func is the function name of the caller.
	Func string `json:"func" msgpack:"func" xml:"func" yaml:"func"`
	// Line is the line number of the caller.
	Line int `json:"line" msgpack:"line" xml:"line" yaml:"line"`
}

// ConvertFrame converts [runtime.Frame] into [Frame].
// ConvertFrame returns zero value of [Frame] if the given frame
// is empty.
func ConvertFrame(f runtime.Frame) Frame {
	if f.PC == 0 {
		return Frame{}
	}
	pkg := ""
	fn := f.Function // fn is "<Pkg>.<Func>"
	if i := strings.LastIndexByte(fn, '.'); i > 0 {
		pkg = fn[:i]
		fn = fn[i+1:]
	}
	file := f.File
	if i := strings.LastIndexByte(file, '/'); i > 0 {
		file = file[i+1:]
	}
	return Frame{
		Pkg:  pkg,
		File: file,
		Line: f.Line,
		Func: fn,
	}
}

// ConvertFrames converts a slice of [runtime.Frame] into a slice of [Frame].
// ConvertFrames returns nil slice if the given slice is empty.
func ConvertFrames(fs []runtime.Frame) []Frame {
	n := len(fs)
	if n == 0 {
		return nil
	}
	frames := make([]Frame, n)
	for i := range fs {
		frames[i] = ConvertFrame(fs[i])
	}
	return frames
}

// CallerFrame returns single caller frame.
// The argument skip is the number of stack frames to skip
// before recording frames, skip=0 means the caller frame
// and skip=1 means the caller of the caller.
// CallerFrame returns zero-value of [runtime.Frame]
// when there is no frames to report.
// See also [runtime.Caller].
func CallerFrame(skip int) runtime.Frame {
	pc := make([]uintptr, 1)
	n := runtime.Callers(skip+2, pc)
	if n < 1 {
		return runtime.Frame{} // No frame to report.
	}
	frame, _ := runtime.CallersFrames(pc).Next()
	return frame
}

// CallerFrames returns a slice of caller frames.
// The argument skip is the number of stack frames to skip
// before recording frames, skip=0 means the caller frame
// and skip=1 means the caller of the caller.
// CallerFrames returns nil slice of [runtime.Frame]
// when there is no frames to report.
// CallerFrames returns 64 frames at maximum.
// See also [runtime.CallersFrames].
func CallerFrames(skip int) []runtime.Frame {
	pcs := make([]uintptr, 64) // Max 64 frames.
	n := runtime.Callers(skip+2, pcs)
	if n < 1 {
		return nil // No frames to report.
	}
	frames := runtime.CallersFrames(pcs[:n])
	fs := make([]runtime.Frame, n)
	for i := range n {
		fs[i], _ = frames.Next()
	}
	return fs
}
