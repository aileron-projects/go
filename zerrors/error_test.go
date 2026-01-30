package zerrors_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/aileron-projects/go/zerrors"
	"github.com/aileron-projects/go/ztesting"
)

func TestToMap(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		m := zerrors.ToMap(nil)
		ztesting.AssertEqual(t, "returned map mismatch", nil, m)
	})
	t.Run("primitive error", func(t *testing.T) {
		m := zerrors.ToMap(io.EOF)
		want := map[string]any{"message": "EOF"}
		ztesting.AssertEqual(t, "returned map mismatch", want, m)
	})
	t.Run("wrapped error", func(t *testing.T) {
		err := fmt.Errorf("outer error [%w]", io.EOF)
		m := zerrors.ToMap(err)
		want := map[string]any{
			"message": "outer error [EOF]",
			"cause": map[string]any{
				"message": "EOF",
			},
		}
		ztesting.AssertEqual(t, "msg mismatch", want, m)
	})
}

func TestError_Unwrap(t *testing.T) {
	t.Parallel()

	e := &zerrors.Error{Cause: io.EOF}
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
			use:    &zerrors.Error{Cause: io.EOF, Code: "c", Kind: "k"},
			target: nil,
			same:   false,
		},
		"not equal": {
			use:    &zerrors.Error{Cause: io.EOF, Code: "c", Kind: "k"},
			target: io.EOF,
			same:   false,
		},
		"code mismatch": {
			use:    &zerrors.Error{Code: "c", Kind: "k"},
			target: &zerrors.Error{Code: "C", Kind: "k"},
			same:   false,
		},
		"kind mismatch": {
			use:    &zerrors.Error{Code: "c", Kind: "k"},
			target: &zerrors.Error{Code: "c", Kind: "K"},
			same:   false,
		},
		"message mismatch": {
			use:    &zerrors.Error{Code: "c", Kind: "k", Message: "m"},
			target: &zerrors.Error{Code: "c", Kind: "k", Message: "M"},
			same:   true,
		},
		"attrs mismatch": {
			use:    &zerrors.Error{Code: "c", Kind: "k", Attrs: map[string]string{"foo": "bar"}},
			target: &zerrors.Error{Code: "c", Kind: "k", Attrs: map[string]string{"FOO": "Bar"}},
			same:   true,
		},
		"same after unwrap": {
			use:    &zerrors.Error{Code: "c", Kind: "k"},
			target: fmt.Errorf("outer error [%w]", &zerrors.Error{Code: "c", Kind: "k"}),
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
			def:    zerrors.Definition{Code: "c", Kind: "k"},
			target: nil,
			same:   false,
		},
		"not equal": {
			def:    zerrors.Definition{Code: "c", Kind: "k"},
			target: io.EOF,
			same:   false,
		},
		"code mismatch": {
			def:    zerrors.Definition{Code: "c", Kind: "k"},
			target: &zerrors.Error{Code: "C", Kind: "k"},
			same:   false,
		},
		"kind mismatch": {
			def:    zerrors.Definition{Code: "c", Kind: "k"},
			target: &zerrors.Error{Code: "c", Kind: "K"},
			same:   false,
		},
		"message mismatch": {
			def:    zerrors.Definition{Code: "c", Kind: "k", Message: "m"},
			target: &zerrors.Error{Code: "c", Kind: "k", Message: "M"},
			same:   true,
		},
		"attrs mismatch": {
			def:    zerrors.Definition{Code: "c", Kind: "k", Attrs: map[string]string{"foo": "bar"}},
			target: &zerrors.Error{Code: "c", Kind: "k", Attrs: map[string]string{"FOO": "Bar"}},
			same:   true,
		},
		"same after unwrap": {
			def:    zerrors.Definition{Code: "c", Kind: "k"},
			target: fmt.Errorf("outer error [%w]", &zerrors.Error{Code: "c", Kind: "k"}),
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
		ztesting.AssertEqual(t, "kind mismatch.", w.Kind, e.Kind)
		ztesting.AssertEqual(t, "message mismatch.", w.Message, e.Message)
		ztesting.AssertEqual(t, "attrs mismatch.", w.Attrs, e.Attrs)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "cause mismatch.", nil, e.Cause)
	})
	t.Run("non inner error", func(t *testing.T) {
		e := (&zerrors.Definition{Code: "c", Kind: "k", Message: "m"}).New(nil)
		w := zerrors.Error{Code: "c", Kind: "k", Message: "m"}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "kind mismatch.", w.Kind, e.Kind)
		ztesting.AssertEqual(t, "message mismatch.", w.Message, e.Message)
		ztesting.AssertEqual(t, "attrs mismatch.", w.Attrs, e.Attrs)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "cause mismatch.", nil, e.Cause)
	})
	t.Run("inner error", func(t *testing.T) {
		e := (&zerrors.Definition{Code: "c", Kind: "k", Message: "m"}).New(io.EOF)
		w := zerrors.Error{Code: "c", Kind: "k", Message: "m"}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "kind mismatch.", w.Kind, e.Kind)
		ztesting.AssertEqual(t, "message mismatch.", w.Message, e.Message)
		ztesting.AssertEqual(t, "attrs mismatch.", w.Attrs, e.Attrs)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "cause mismatch.", io.EOF, e.Cause)
	})
	t.Run("attrs error", func(t *testing.T) {
		e := (&zerrors.Definition{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}).New(nil)
		w := zerrors.Error{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "kind mismatch.", w.Kind, e.Kind)
		ztesting.AssertEqual(t, "message mismatch.", w.Message, e.Message)
		ztesting.AssertEqual(t, "attrs mismatch.", w.Attrs, e.Attrs)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "cause mismatch.", nil, e.Cause)
	})
	t.Run("format message", func(t *testing.T) {
		e := (&zerrors.Definition{Code: "c", Kind: "k", Message: "foo=%s"}).New(nil, "xxx")
		w := zerrors.Error{Code: "c", Kind: "k", Message: "foo=xxx"}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "kind mismatch.", w.Kind, e.Kind)
		ztesting.AssertEqual(t, "message mismatch.", w.Message, e.Message)
		ztesting.AssertEqual(t, "attrs mismatch.", w.Attrs, e.Attrs)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "cause mismatch.", nil, e.Cause)
	})
}

func TestDefinition_NewStack(t *testing.T) {
	t.Parallel()
	t.Run("zero value", func(t *testing.T) {
		var ed zerrors.Definition
		e := ed.NewStack(nil)
		w := zerrors.Error{}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "kind mismatch.", w.Kind, e.Kind)
		ztesting.AssertEqual(t, "message mismatch.", w.Message, e.Message)
		ztesting.AssertEqual(t, "attrs mismatch.", w.Attrs, e.Attrs)
		ztesting.AssertEqual(t, "unexpected frame length.", true, len(e.Frames) > 0)
		ztesting.AssertEqual(t, "cause mismatch.", nil, e.Cause)
	})
	t.Run("cause without stack", func(t *testing.T) {
		e := (&zerrors.Definition{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}).NewStack(io.EOF)
		w := zerrors.Error{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "kind mismatch.", w.Kind, e.Kind)
		ztesting.AssertEqual(t, "message mismatch.", w.Message, e.Message)
		ztesting.AssertEqual(t, "attrs mismatch.", w.Attrs, e.Attrs)
		ztesting.AssertEqual(t, "unexpected frame length.", true, len(e.Frames) > 0)
		ztesting.AssertEqual(t, "cause mismatch.", io.EOF, e.Cause)
	})
	t.Run("inner error with stack", func(t *testing.T) {
		inner := &zerrors.Error{Frames: []zerrors.Frame{{}, {}}}
		e := (&zerrors.Definition{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}).NewStack(inner)
		w := zerrors.Error{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}}
		ztesting.AssertEqual(t, "code mismatch.", w.Code, e.Code)
		ztesting.AssertEqual(t, "kind mismatch.", w.Kind, e.Kind)
		ztesting.AssertEqual(t, "message mismatch.", w.Message, e.Message)
		ztesting.AssertEqual(t, "attrs mismatch.", w.Attrs, e.Attrs)
		ztesting.AssertEqual(t, "unexpected frame length.", 0, len(e.Frames))
		ztesting.AssertEqual(t, "cause mismatch.", error(inner), e.Cause)
	})
}
