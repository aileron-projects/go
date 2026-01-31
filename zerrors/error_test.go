package zerrors_test

import (
	"fmt"
	"io"
	"log/slog"
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
	t.Run("interface", func(t *testing.T) {
		def := zerrors.NewDefinition("c", "k", "m", nil)
		err := def.NewStack(nil)
		m := zerrors.ToMap(err)
		ztesting.AssertEqual(t, "code mismatch", "c", m["code"])
		ztesting.AssertEqual(t, "kind mismatch", "k", m["kind"])
		ztesting.AssertEqual(t, "message mismatch", "m", m["message"])
		ztesting.AssertEqual(t, "empty frame", true, len(m["frames"].([]string)) > 0)
	})
}

func TestToSlogAttrs(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		m := zerrors.ToSlogAttrs(nil)
		ztesting.AssertEqual(t, "attrs mismatch", nil, m)
	})
	t.Run("primitive error", func(t *testing.T) {
		m := zerrors.ToSlogAttrs(io.EOF)
		want := []slog.Attr{slog.String("message", "EOF")}
		ztesting.AssertEqual(t, "attrs mismatch", want, m)
	})
	t.Run("wrapped error", func(t *testing.T) {
		err := fmt.Errorf("outer error [%w]", io.EOF)
		m := zerrors.ToSlogAttrs(err)
		want := []slog.Attr{
			slog.String("message", "outer error [EOF]"),
			slog.GroupAttrs("cause", slog.String("message", "EOF")),
		}
		ztesting.AssertEqual(t, "attrs mismatch", want, m)
	})
	t.Run("interface", func(t *testing.T) {
		def := zerrors.NewDefinition("c", "k", "m", nil)
		err := def.New(nil)
		m := zerrors.ToSlogAttrs(err)
		want := []slog.Attr{
			slog.String("code", "c"),
			slog.String("kind", "k"),
			slog.String("message", "m"),
		}
		ztesting.AssertEqual(t, "attrs mismatch", want, m)
	})
}

func TestError_Unwrap(t *testing.T) {
	t.Parallel()
	e := &zerrors.Error{Cause: io.EOF}
	u := e.Unwrap()
	ztesting.AssertEqual(t, "unwrapped error is incorrect.", io.EOF, u)
}

func TestError_Error(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		err  *zerrors.Error
		want string
	}{
		"code": {
			err:  &zerrors.Error{Code: "c"},
			want: "c  :",
		},
		"kind": {
			err:  &zerrors.Error{Kind: "k"},
			want: " k :",
		},
		"message": {
			err:  &zerrors.Error{Message: "m"},
			want: "  :m",
		},
		"code kind": {
			err:  &zerrors.Error{Code: "c", Kind: "k"},
			want: "c k :",
		},
		"code kind message": {
			err:  &zerrors.Error{Code: "c", Kind: "k", Message: "m"},
			want: "c k :m",
		},
		"attrs": {
			err:  &zerrors.Error{Code: "c", Kind: "k", Message: "m", Attrs: map[string]string{"foo": "bar"}},
			want: "c k :m (foo=bar)",
		},
		"cause": {
			err:  &zerrors.Error{Code: "c", Kind: "k", Message: "m", Cause: io.EOF},
			want: "c k :m [EOF]",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.err.Error()
			ztesting.AssertEqual(t, "error string not match", tc.want, got)
		})
	}
}

func TestError_Is(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		use    *zerrors.Error
		target error
		same   bool
	}{
		"nil": {
			use:    nil,
			target: nil,
			same:   false,
		},
		"nil pointer": {
			use:    nil,
			target: (*zerrors.Error)(nil),
			same:   true,
		},
		"nil target": {
			use:    &zerrors.Error{Cause: io.EOF, Code: "c", Kind: "k"},
			target: nil,
			same:   false,
		},
		"equal": {
			use:    &zerrors.Error{Code: "c", Kind: "k"},
			target: &zerrors.Error{Code: "c", Kind: "k"},
			same:   true,
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

func TestError_Map(t *testing.T) {
	t.Parallel()
	t.Run("minimum", func(t *testing.T) {
		e := &zerrors.Error{Code: "c", Kind: "k", Message: "m"}
		got := e.Map()
		want := map[string]any{
			"code":    "c",
			"kind":    "k",
			"message": "m",
		}
		ztesting.AssertEqual(t, "maps not matched", want, got)
	})
	t.Run("cause", func(t *testing.T) {
		e := &zerrors.Error{Cause: io.EOF}
		got := e.Map()
		want := map[string]any{
			"message": "EOF",
		}
		ztesting.AssertEqual(t, "maps not matched", want, got["cause"].(map[string]any))
	})
	t.Run("attrs", func(t *testing.T) {
		e := &zerrors.Error{Attrs: map[string]string{"foo": "bar"}}
		got := e.Map()
		want := map[string]string{
			"foo": "bar",
		}
		ztesting.AssertEqual(t, "maps not matched", want, got["attrs"].(map[string]string))
	})
	t.Run("frames", func(t *testing.T) {
		def := zerrors.NewDefinition("c", "k", "m", nil)
		e := def.NewStack(nil)
		got := e.Map()
		ztesting.AssertEqual(t, "empty frame", true, len(got["frames"].([]string)) > 0)
	})
}

func TestError_SlogAttrs(t *testing.T) {
	t.Parallel()
	t.Run("minimum", func(t *testing.T) {
		e := &zerrors.Error{Code: "c", Kind: "k", Message: "m"}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("code", "c"),
			slog.String("kind", "k"),
			slog.String("message", "m"),
		}
		ztesting.AssertEqual(t, "attrs not matched", want, got)
	})
	t.Run("cause", func(t *testing.T) {
		e := &zerrors.Error{Cause: io.EOF}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("code", ""),
			slog.String("kind", ""),
			slog.String("message", ""),
			slog.GroupAttrs("cause", slog.String("message", "EOF")),
		}
		ztesting.AssertEqual(t, "attrs not matched", want, got)
	})
	t.Run("attrs", func(t *testing.T) {
		e := &zerrors.Error{Attrs: map[string]string{"foo": "bar"}}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("code", ""),
			slog.String("kind", ""),
			slog.String("message", ""),
			slog.GroupAttrs("attrs", slog.String("foo", "bar")),
		}
		ztesting.AssertEqual(t, "attrs not matched", want, got)
	})
	t.Run("frames", func(t *testing.T) {
		def := zerrors.NewDefinition("c", "k", "m", nil)
		e := def.NewStack(nil)
		got := e.SlogAttrs()
		for _, a := range got {
			if a.Key == "frames" {
				ztesting.AssertEqual(t, "empty frame", true, len(a.Value.Any().([]string)) > 0)
				return
			}
		}
		t.Error("frame does not exist.")
	})
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
