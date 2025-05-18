package zruntime_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zruntime"
	"github.com/aileron-projects/go/ztesting"
)

func TestConvertFrame(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		frame runtime.Frame
		want  zruntime.Frame
	}{
		"empty frame": {
			frame: runtime.Frame{},
			want:  zruntime.Frame{},
		},
		"non empty frame": {
			frame: runtime.Frame{PC: 123, Function: "foo/bar.testFunc", File: "test.go", Line: 100},
			want:  zruntime.Frame{Pkg: "foo/bar", File: "test.go", Func: "testFunc", Line: 100},
		},
		"short func name": {
			frame: runtime.Frame{PC: 123, Function: "bar.testFunc", File: "test.go", Line: 100},
			want:  zruntime.Frame{Pkg: "bar", File: "test.go", Func: "testFunc", Line: 100},
		},
		"short file name": {
			frame: runtime.Frame{PC: 123, Function: "foo/bar.testFunc", File: "baz/test.go", Line: 100},
			want:  zruntime.Frame{Pkg: "foo/bar", File: "test.go", Func: "testFunc", Line: 100},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			f := zruntime.ConvertFrame(tc.frame)
			ztesting.AssertEqual(t, "pkg value does not match.", tc.want.Pkg, f.Pkg)
			ztesting.AssertEqual(t, "file value does not match.", tc.want.File, f.File)
			ztesting.AssertEqual(t, "func value does not match.", tc.want.Func, f.Func)
			ztesting.AssertEqual(t, "line value does not match.", tc.want.Line, f.Line)
		})
	}
}

func TestConvertFrames(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		frames []runtime.Frame
		want   []zruntime.Frame
	}{
		"empty frame": {
			frames: nil,
			want:   nil,
		},
		"non empty frame": {
			frames: []runtime.Frame{
				{PC: 1, Function: "foo/bar.testFunc1", File: "test1.go", Line: 101},
				{PC: 1, Function: "bar/foo.testFunc2", File: "test2.go", Line: 102},
			},
			want: []zruntime.Frame{
				{Pkg: "foo/bar", File: "test1.go", Func: "testFunc1", Line: 101},
				{Pkg: "bar/foo", File: "test2.go", Func: "testFunc2", Line: 102},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			fs := zruntime.ConvertFrames(tc.frames)
			if !(len(fs) == len(tc.want)) {
				t.Errorf("length not match. want:%d got:%d", len(tc.want), len(fs))
			}
			for i, f := range fs {
				w := tc.want[i]
				ztesting.AssertEqual(t, "pkg value does not match.", w.Pkg, f.Pkg)
				ztesting.AssertEqual(t, "file value does not match.", w.File, f.File)
				ztesting.AssertEqual(t, "func value does not match.", w.Func, f.Func)
				ztesting.AssertEqual(t, "line value does not match.", w.Line, f.Line)
			}
		})
	}
}

func TestCallerFrame(t *testing.T) {
	t.Parallel()

	t.Run("skip=0", func(t *testing.T) {
		f := zruntime.CallerFrame(0)
		ztesting.AssertEqual(t, "file does not have appropriate suffix.", true, strings.HasSuffix(f.File, "frame_test.go"))
		ztesting.AssertEqual(t, "func does not have appropriate suffix.", true, strings.HasSuffix(f.Function, "TestCallerFrame.func1"))
		ztesting.AssertEqual(t, "line number is not positive.", true, f.Line > 0)
	})

	t.Run("skip=9999", func(t *testing.T) {
		f := zruntime.CallerFrame(9999)
		ztesting.AssertEqual(t, "file has value.", "", f.File)
		ztesting.AssertEqual(t, "func has value.", "", f.Function)
		ztesting.AssertEqual(t, "line number is not 0.", 0, f.Line)
	})

	t.Run("skip=-9999", func(t *testing.T) {
		f := zruntime.CallerFrame(-9999)
		ztesting.AssertEqual(t, "file does not have appropriate suffix.", true, strings.HasSuffix(f.File, "runtime/extern.go"))
		ztesting.AssertEqual(t, "func does not have appropriate suffix.", true, strings.HasSuffix(f.Function, "runtime.Callers"))
		ztesting.AssertEqual(t, "line number is not positive.", true, f.Line > 0)
	})
}

func TestCallerFrames(t *testing.T) {
	t.Parallel()

	t.Run("skip=0", func(t *testing.T) {
		fs := zruntime.CallerFrames(0)
		ztesting.AssertEqual(t, "no frames found.", true, len(fs) > 1)
		f := fs[0]
		ztesting.AssertEqual(t, "file does not have appropriate suffix.", true, strings.HasSuffix(f.File, "frame_test.go"))
		ztesting.AssertEqual(t, "func does not have appropriate suffix.", true, strings.HasSuffix(f.Function, "TestCallerFrames.func1"))
		ztesting.AssertEqual(t, "line number is not positive.", true, f.Line > 0)
	})

	t.Run("skip=9999", func(t *testing.T) {
		fs := zruntime.CallerFrames(9999)
		ztesting.AssertEqual(t, "frames found.", 0, len(fs))
	})

	t.Run("skip=-9999", func(t *testing.T) {
		fs := zruntime.CallerFrames(-9999)
		ztesting.AssertEqual(t, "no frames found.", true, len(fs) > 1)
		f := fs[0]
		ztesting.AssertEqual(t, "file does not have appropriate suffix.", true, strings.HasSuffix(f.File, "runtime/extern.go"))
		ztesting.AssertEqual(t, "func does not have appropriate suffix.", true, strings.HasSuffix(f.Function, "runtime.Callers"))
		ztesting.AssertEqual(t, "line number is not positive.", true, f.Line > 0)
	})
}
