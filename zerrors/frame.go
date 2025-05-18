package zerrors

import (
	"runtime"
	"strings"
)

// Frame holds stack frame information.
// See also [runtime.Frame].
type Frame struct {
	// Pkg is go package name of the caller.
	Pkg string `json:"pkg" msgpack:"pkg" toml:"pkg" xml:"pkg" yaml:"pkg"`
	// File is the file name of the caller.
	File string `json:"file" msgpack:"file" toml:"file" xml:"file" yaml:"file"`
	// Func is the function name of the caller.
	Func string `json:"func" msgpack:"func" toml:"func" xml:"func" yaml:"func"`
	// Line is the line number of the caller.
	Line int `json:"line" msgpack:"line" toml:"line" xml:"line" yaml:"line"`
}

// callerFrames returns a slice of caller frames.
// The argument skip is the number of stack frames to skip
// before recording frames, skip=0 means the caller frame
// and skip=1 means the caller of the caller.
// Second argument size is the maximum number of frames to report.
// callerFrames returns nil slice of [Frame] if there is no frames to report.
// See also [runtime.CallersFrames].
func callerFrames(skip int, size uint) []Frame {
	pcs := make([]uintptr, size)
	n := runtime.Callers(skip+2, pcs)
	if n < 1 {
		return nil // No frames to report.
	}
	frames := runtime.CallersFrames(pcs[:n])
	fs := make([]Frame, n)
	for i := range n {
		f, _ := frames.Next()
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
		fs[i] = Frame{
			Pkg:  pkg,
			File: file,
			Line: f.Line,
			Func: fn,
		}
	}
	return fs
}
