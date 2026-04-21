package zerrors_test

import (
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/aileron-projects/go/zerrors"
	"github.com/aileron-projects/go/ztesting"
)

func TestErr_Unwrap(t *testing.T) {
	t.Parallel()
	e := &zerrors.Err{Cause: io.EOF}
	u := e.Unwrap()
	ztesting.AssertEqual(t, io.EOF, u)
}

func TestErr_Error(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		err  *zerrors.Err
		want string
	}{
		"message": {
			err:  &zerrors.Err{Message: "m"},
			want: "m",
		},
		"detail": {
			err:  &zerrors.Err{Detail: "d"},
			want: " d",
		},
		"cause": {
			err:  &zerrors.Err{Cause: io.EOF},
			want: " [EOF]",
		},
		"all": {
			err:  &zerrors.Err{Cause: io.EOF, Message: "m", Detail: "d"},
			want: "m d [EOF]",
		},
		"message detail": {
			err:  &zerrors.Err{Message: "m", Detail: "d"},
			want: "m d",
		},
		"message cause": {
			err:  &zerrors.Err{Message: "m", Cause: io.EOF},
			want: "m [EOF]",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.err.Error()
			ztesting.AssertEqual(t, tc.want, got)
		})
	}
}

func TestErr_Is(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		use    *zerrors.Err
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
			target: (*zerrors.Err)(nil),
			same:   true,
		},
		"nil target": {
			use:    &zerrors.Err{Cause: io.EOF, Message: "m", Detail: "d"},
			target: nil,
			same:   false,
		},
		"equal": {
			use:    &zerrors.Err{Message: "m", Detail: "d"},
			target: &zerrors.Err{Message: "m", Detail: "d"},
			same:   true,
		},
		"not equal": {
			use:    &zerrors.Err{Cause: io.EOF, Message: "m", Detail: "d"},
			target: io.EOF,
			same:   false,
		},
		"message mismatch": {
			use:    &zerrors.Err{Message: "m", Detail: "d"},
			target: &zerrors.Err{Message: "M", Detail: "d"},
			same:   false,
		},
		"detail mismatch": {
			use:    &zerrors.Err{Message: "m", Detail: "d"},
			target: &zerrors.Err{Message: "m", Detail: "D"},
			same:   true,
		},
		"same after unwrap": {
			use:    &zerrors.Err{Message: "m", Detail: "d"},
			target: fmt.Errorf("outer error [%w]", &zerrors.Err{Message: "m", Detail: "d"}),
			same:   true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			is := tc.use.Is(tc.target)
			ztesting.AssertEqual(t, tc.same, is)
		})
	}
}

func TestErr_Map(t *testing.T) {
	t.Parallel()
	t.Run("minimum", func(t *testing.T) {
		e := &zerrors.Err{Message: "m", Detail: "d"}
		got := e.Map()
		want := map[string]any{
			"message": "m",
			"detail":  "d",
		}
		ztesting.AssertEqual(t, want, got)
	})
	t.Run("cause", func(t *testing.T) {
		e := &zerrors.Err{Cause: io.EOF}
		got := e.Map()
		want := map[string]any{
			"message": "EOF",
		}
		ztesting.AssertEqual(t, want, got["cause"].(map[string]any))
	})
}

func TestErr_SlogAttrs(t *testing.T) {
	t.Parallel()
	t.Run("minimum", func(t *testing.T) {
		e := &zerrors.Err{Message: "m", Detail: "d"}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("message", "m"),
			slog.String("detail", "d"),
		}
		ztesting.AssertEqual(t, want, got)
	})
	t.Run("cause", func(t *testing.T) {
		e := &zerrors.Err{Cause: io.EOF}
		got := e.SlogAttrs()
		want := []slog.Attr{
			slog.String("message", ""),
			slog.String("detail", ""),
			slog.GroupAttrs("cause", slog.String("message", "EOF")),
		}
		ztesting.AssertEqual(t, want, got)
	})
}
