package zerrors_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/aileron-projects/go/zerrors"
	"github.com/aileron-projects/go/ztesting"
)

func TestAttrs(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		m := zerrors.Attrs(nil)
		ztesting.AssertEqual(t, "returned map mismatch", nil, m)
	})
	t.Run("primitive error", func(t *testing.T) {
		m := zerrors.Attrs(io.EOF)
		want := map[string]any{"msg": "EOF"}
		ztesting.AssertEqual(t, "returned map mismatch", want, m)
	})
	t.Run("wrapped error", func(t *testing.T) {
		err := fmt.Errorf("outer error [%w]", io.EOF)
		m := zerrors.Attrs(err)
		want := map[string]any{
			"msg": "outer error [EOF]",
			"wraps": map[string]any{
				"msg": "EOF",
			},
		}
		ztesting.AssertEqual(t, "msg mismatch", want, m)
	})
}

func TestError_Unwrap(t *testing.T) {
	t.Parallel()

	e := &zerrors.Error{Inner: io.EOF}
	u := e.Unwrap()
	ztesting.AssertEqual(t, "unwrapped error is incorrect.", io.EOF, u)
}

func TestError_Is(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		use    *zerrors.Error
		target error
		same   bool
	}{
		"nil": {
			use:    &zerrors.Error{Inner: io.EOF, Code: "c", Pkg: "p", Msg: "m"},
			target: nil,
			same:   false,
		},
		"not equal": {
			use:    &zerrors.Error{Inner: io.EOF, Pkg: "p", Msg: "m"},
			target: io.EOF,
			same:   false,
		},
		"code mismatch": {
			use:    &zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "e"},
			target: &zerrors.Error{Code: "C", Pkg: "P", Msg: "m", Detail: "d", Ext: "e"},
			same:   false,
		},
		"pkg mismatch": {
			use:    &zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "e"},
			target: &zerrors.Error{Code: "c", Pkg: "P", Msg: "m", Detail: "d", Ext: "e"},
			same:   false,
		},
		"msg mismatch": {
			use:    &zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "e"},
			target: &zerrors.Error{Code: "c", Pkg: "p", Msg: "M", Detail: "d", Ext: "e"},
			same:   false,
		},
		"detail mismatch": {
			use:    &zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "e"},
			target: &zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "D", Ext: "e"},
			same:   true,
		},
		"ext mismatch": {
			use:    &zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "e"},
			target: &zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "D"},
			same:   true,
		},
		"same after unwrap": {
			use:    &zerrors.Error{Code: "c", Pkg: "p", Msg: "m"},
			target: fmt.Errorf("outer error [%w]", &zerrors.Error{Code: "c", Pkg: "p", Msg: "m"}),
			same:   true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			is := tc.use.Is(tc.target)
			ztesting.AssertEqual(t, "incorrect error identification.", tc.same, is)
		})
	}
}

func TestDefinition_Is(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		def    zerrors.Definition
		target error
		same   bool
	}{
		"nil": {
			def:    zerrors.Definition{"c", "p", "m"},
			target: nil,
			same:   false,
		},
		"not equal": {
			def:    zerrors.Definition{"c", "p", "m"},
			target: io.EOF,
			same:   false,
		},
		"code mismatch": {
			def:    zerrors.Definition{"c", "p", "m", "d", "e"},
			target: &zerrors.Error{Code: "C", Pkg: "P", Msg: "m", Detail: "d", Ext: "e"},
			same:   false,
		},
		"pkg mismatch": {
			def:    zerrors.Definition{"c", "p", "m", "d", "e"},
			target: &zerrors.Error{Code: "c", Pkg: "P", Msg: "m", Detail: "d", Ext: "e"},
			same:   false,
		},
		"msg mismatch": {
			def:    zerrors.Definition{"c", "p", "m", "d", "e"},
			target: &zerrors.Error{Code: "c", Pkg: "p", Msg: "M", Detail: "d", Ext: "e"},
			same:   false,
		},
		"detail mismatch": {
			def:    zerrors.Definition{"c", "p", "m", "d", "e"},
			target: &zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "D", Ext: "e"},
			same:   true,
		},
		"ext mismatch": {
			def:    zerrors.Definition{"c", "p", "m", "d", "e"},
			target: &zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "D"},
			same:   true,
		},
		"same after unwrap": {
			def:    zerrors.Definition{"c", "p", "m", "d", "e"},
			target: fmt.Errorf("outer error [%w]", &zerrors.Error{Code: "c", Pkg: "p", Msg: "m"}),
			same:   true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			is := tc.def.Is(tc.target)
			ztesting.AssertEqual(t, "incorrect error identification.", tc.same, is)
		})
	}
}

func TestDefinition_New(t *testing.T) {
	t.Parallel()
	t.Run("zero value", func(t *testing.T) {
		var ed zerrors.Definition
		e := ed.New(nil)
		w := zerrors.Error{}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "pkg mismatch.", w.Pkg, e.Pkg)
		ztesting.AssertEqual(t, "msg mismatch.", w.Msg, e.Msg)
		ztesting.AssertEqual(t, "ext mismatch.", w.Ext, e.Ext)
		ztesting.AssertEqual(t, "detail mismatch.", w.Detail, e.Detail)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "inner error mismatch.", nil, e.Inner)
	})
	t.Run("non inner error", func(t *testing.T) {
		e := zerrors.Definition{"c", "p", "m", "d", "e"}.New(nil)
		w := zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "e"}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "pkg mismatch.", w.Pkg, e.Pkg)
		ztesting.AssertEqual(t, "msg mismatch.", w.Msg, e.Msg)
		ztesting.AssertEqual(t, "ext mismatch.", w.Ext, e.Ext)
		ztesting.AssertEqual(t, "detail mismatch.", w.Detail, e.Detail)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "inner error mismatch.", nil, e.Inner)
	})
	t.Run("inner error", func(t *testing.T) {
		e := zerrors.Definition{"c", "p", "m", "d", "e"}.New(io.EOF)
		w := zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "e"}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "pkg mismatch.", w.Pkg, e.Pkg)
		ztesting.AssertEqual(t, "msg mismatch.", w.Msg, e.Msg)
		ztesting.AssertEqual(t, "ext mismatch.", w.Ext, e.Ext)
		ztesting.AssertEqual(t, "detail mismatch.", w.Detail, e.Detail)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "inner error mismatch.", io.EOF, e.Inner)
	})
	t.Run("format detail", func(t *testing.T) {
		e := zerrors.Definition{"c", "p", "m", "foo=%s", "e"}.New(nil, "xxx")
		w := zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "foo=xxx", Ext: "e"}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "pkg mismatch.", w.Pkg, e.Pkg)
		ztesting.AssertEqual(t, "msg mismatch.", w.Msg, e.Msg)
		ztesting.AssertEqual(t, "ext mismatch.", w.Ext, e.Ext)
		ztesting.AssertEqual(t, "detail mismatch.", w.Detail, e.Detail)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "inner error mismatch.", nil, e.Inner)
	})
}

func TestDefinition_NewStack(t *testing.T) {
	t.Parallel()
	t.Run("zero value", func(t *testing.T) {
		var ed zerrors.Definition
		e := ed.NewStack(nil)
		w := zerrors.Error{}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "pkg mismatch.", w.Pkg, e.Pkg)
		ztesting.AssertEqual(t, "msg mismatch.", w.Msg, e.Msg)
		ztesting.AssertEqual(t, "ext mismatch.", w.Ext, e.Ext)
		ztesting.AssertEqual(t, "detail mismatch.", w.Detail, e.Detail)
		ztesting.AssertEqual(t, "unexpected frame length.", true, len(e.Frames) > 0)
		ztesting.AssertEqual(t, "inner error mismatch.", nil, e.Inner)
	})
	t.Run("inner error without stack", func(t *testing.T) {
		e := zerrors.Definition{"c", "p", "m", "d", "e"}.NewStack(io.EOF)
		w := zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "e"}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "pkg mismatch.", w.Pkg, e.Pkg)
		ztesting.AssertEqual(t, "msg mismatch.", w.Msg, e.Msg)
		ztesting.AssertEqual(t, "ext mismatch.", w.Ext, e.Ext)
		ztesting.AssertEqual(t, "detail mismatch.", w.Detail, e.Detail)
		ztesting.AssertEqual(t, "unexpected frame length.", true, len(e.Frames) > 0)
		ztesting.AssertEqual(t, "inner error mismatch.", io.EOF, e.Inner)
	})
	t.Run("inner error with stack", func(t *testing.T) {
		inner := &zerrors.Error{Frames: []zerrors.Frame{{}, {}}}
		e := zerrors.Definition{"c", "p", "m", "d", "e"}.NewStack(inner)
		w := zerrors.Error{Code: "c", Pkg: "p", Msg: "m", Detail: "d", Ext: "e"}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "pkg mismatch.", w.Pkg, e.Pkg)
		ztesting.AssertEqual(t, "msg mismatch.", w.Msg, e.Msg)
		ztesting.AssertEqual(t, "ext mismatch.", w.Ext, e.Ext)
		ztesting.AssertEqual(t, "detail mismatch.", w.Detail, e.Detail)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "inner error mismatch.", error(inner), e.Inner)
	})
}
