package ztext

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestTemplate_WithTagFunc(t *testing.T) {
	t.Parallel()

	t.Run("key found", func(t *testing.T) {
		tpl := NewTemplate("a={a} b={b}", "{", "}")
		tpl.WithTagFunc("a", func(s string) []byte { return []byte("A") })
		tpl.WithTagFunc("b", func(s string) []byte { return []byte("B") })
		got := tpl.ExecuteString(nil)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=B", got)
	})

	t.Run("key not found", func(t *testing.T) {
		tpl := NewTemplate("a={a} b={b}", "{", "}")
		tpl.WithTagFunc("a", func(s string) []byte { return []byte("A") })
		got := tpl.ExecuteString(nil)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", got)
	})

	t.Run("empty tag", func(t *testing.T) {
		tpl := NewTemplate("a={a} b={b}", "{", "}")
		tpl.WithTagFunc("a", func(s string) []byte { return []byte("A") })
		tpl.WithTagFunc("", func(s string) []byte { return []byte("B") })
		got := tpl.ExecuteString(nil)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", got)
	})

	t.Run("nil func", func(t *testing.T) {
		tpl := NewTemplate("a={a} b={b}", "{", "}")
		tpl.WithTagFunc("a", func(s string) []byte { return []byte("A") })
		tpl.WithTagFunc("b", nil)
		got := tpl.ExecuteString(nil)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", got)
	})
}

func TestTemplate_Execute(t *testing.T) {
	t.Parallel()

	tpl := NewTemplate("a={a} b={b}", "{", "}")
	tpl.WithTagFunc("a", func(s string) []byte { return []byte("A") })

	t.Run("key found", func(t *testing.T) {
		got := tpl.Execute(map[string]any{"b": "B"})
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=B", string(got))
	})

	t.Run("key not found", func(t *testing.T) {
		got := tpl.Execute(map[string]any{"c": "C"})
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", string(got))
	})

	t.Run("nil map", func(t *testing.T) {
		got := tpl.Execute(nil)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", string(got))
	})
}

func TestTemplate_ExecuteString(t *testing.T) {
	t.Parallel()

	tpl := NewTemplate("a={a} b={b}", "{", "}")
	tpl.WithTagFunc("a", func(s string) []byte { return []byte("A") })

	t.Run("key found", func(t *testing.T) {
		got := tpl.ExecuteString(map[string]any{"b": "B"})
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=B", got)
	})

	t.Run("key not found", func(t *testing.T) {
		got := tpl.ExecuteString(map[string]any{"c": "C"})
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", got)
	})

	t.Run("nil map", func(t *testing.T) {
		got := tpl.ExecuteString(nil)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", got)
	})
}

func TestTemplate_ExecuteFunc(t *testing.T) {
	t.Parallel()

	tpl := NewTemplate("a={a} b={b}", "{", "}")
	tpl.WithTagFunc("a", func(s string) []byte { return []byte("A") })

	t.Run("non-nil func", func(t *testing.T) {
		got := tpl.ExecuteFunc(func(s string) []byte { return []byte(strings.ToUpper(s)) })
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=B", string(got))
	})

	t.Run("nil func", func(t *testing.T) {
		got := tpl.ExecuteFunc(nil)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", string(got))
	})
}

func TestTemplate_ExecuteWriter(t *testing.T) {
	t.Parallel()

	tpl := NewTemplate("a={a} b={b}", "{", "}")
	tpl.WithTagFunc("a", func(s string) []byte { return []byte("A") })

	t.Run("key found", func(t *testing.T) {
		var buf bytes.Buffer
		err := tpl.ExecuteWriter(&buf, map[string]any{"b": "B"})
		ztesting.AssertEqual(t, "unexpected error returned.", nil, err)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=B", buf.String())
	})

	t.Run("key not found", func(t *testing.T) {
		var buf bytes.Buffer
		err := tpl.ExecuteWriter(&buf, map[string]any{"c": "C"})
		ztesting.AssertEqual(t, "unexpected error returned.", nil, err)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", buf.String())
	})

	t.Run("nil map", func(t *testing.T) {
		var buf bytes.Buffer
		err := tpl.ExecuteWriter(&buf, nil)
		ztesting.AssertEqual(t, "unexpected error returned.", nil, err)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", buf.String())
	})

	t.Run("write error", func(t *testing.T) {
		buf := ziotest.ErrWriter(nil, 0)
		err := tpl.ExecuteWriter(buf, map[string]any{"b": "B"})
		ztesting.AssertEqual(t, "unexpected error returned.", io.ErrClosedPipe, err)
	})
}

func TestTemplate_ExecuteWriterFunc(t *testing.T) {
	t.Parallel()

	tpl := NewTemplate("a={a} b={b}", "{", "}")
	tpl.WithTagFunc("a", func(s string) []byte { return []byte("A") })

	t.Run("non-nil func", func(t *testing.T) {
		var buf bytes.Buffer
		err := tpl.ExecuteWriterFunc(&buf, func(s string) []byte { return []byte(strings.ToUpper(s)) })
		ztesting.AssertEqual(t, "unexpected error returned.", nil, err)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=B", buf.String())
	})

	t.Run("nil func", func(t *testing.T) {
		var buf bytes.Buffer
		err := tpl.ExecuteWriterFunc(&buf, nil)
		ztesting.AssertEqual(t, "unexpected error returned.", nil, err)
		ztesting.AssertEqual(t, "wrong result of template.", "a=A b=", buf.String())
	})

	t.Run("write error", func(t *testing.T) {
		buf := ziotest.ErrWriter(nil, 0)
		err := tpl.ExecuteWriterFunc(buf, func(s string) []byte { return []byte(strings.ToUpper(s)) })
		ztesting.AssertEqual(t, "unexpected error returned.", io.ErrClosedPipe, err)
	})
}

type testStringer struct {
	s string
}

func (t *testStringer) String() string {
	return t.s
}

func TestMapValue(t *testing.T) {
	t.Parallel()
	t.Run("value", func(t *testing.T) {
		testCases := map[string]struct {
			mv     mapVal
			tag    string
			expect []byte
		}{
			"not found":  {mapVal{"key": nil}, "KEY", nil},
			"nil":        {mapVal{"key": nil}, "key", []byte("<nil>")},
			"string":     {mapVal{"key": "test"}, "key", []byte("test")},
			"[]byte":     {mapVal{"key": []byte("test")}, "key", []byte("test")},
			"bool":       {mapVal{"key": true}, "key", []byte("true")},
			"int":        {mapVal{"key": 123}, "key", []byte("123")},
			"int8":       {mapVal{"key": int8(123)}, "key", []byte("123")},
			"int16":      {mapVal{"key": int16(123)}, "key", []byte("123")},
			"int32":      {mapVal{"key": int32(123)}, "key", []byte("123")},
			"int64":      {mapVal{"key": int64(123)}, "key", []byte("123")},
			"uint":       {mapVal{"key": uint(123)}, "key", []byte("123")},
			"uint8":      {mapVal{"key": uint8(123)}, "key", []byte("123")},
			"uint16":     {mapVal{"key": uint16(123)}, "key", []byte("123")},
			"uint32":     {mapVal{"key": uint32(123)}, "key", []byte("123")},
			"uint64":     {mapVal{"key": uint64(123)}, "key", []byte("123")},
			"complex64":  {mapVal{"key": complex64(1 + 2i)}, "key", []byte("(1+2i)")},
			"complex128": {mapVal{"key": complex128(1 + 2i)}, "key", []byte("(1+2i)")},
			"stringer":   {mapVal{"key": &testStringer{"test"}}, "key", []byte("test")},
			"struct":     {mapVal{"key": struct{ x string }{"test"}}, "key", []byte("{test}")},
		}

		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				got := tc.mv.value(tc.tag)
				ztesting.AssertEqual(t, "wrong result of template.", string(tc.expect), string(got))
			})
		}
	})
}
