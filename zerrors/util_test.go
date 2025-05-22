package zerrors_test

import (
	"errors"
	"io"
	"testing"

	"github.com/aileron-projects/go/zerrors"
	"github.com/aileron-projects/go/ztesting"
)

func TestUnwrapErr(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		err := zerrors.UnwrapErr(nil)
		ztesting.AssertEqual(t, "unexpected error returned.", nil, err)
	})

	t.Run("has interface", func(t *testing.T) {
		err := zerrors.UnwrapErr(&zerrors.NestErr{
			Err: io.EOF,
			Msg: "outer error",
		})
		ztesting.AssertEqual(t, "unexpected error returned.", io.EOF, err)
	})

	t.Run("has no interface", func(t *testing.T) {
		err := zerrors.UnwrapErr(errors.Join(io.EOF, io.EOF)) // Implements interface{ Unwrap() []error }
		ztesting.AssertEqual(t, "unexpected error returned.", nil, err)
	})
}

func TestUnwrapErrs(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		err := zerrors.UnwrapErrs(nil)
		ztesting.AssertEqual(t, "unexpected error returned.", nil, err)
	})

	t.Run("has no interface", func(t *testing.T) {
		err := zerrors.UnwrapErrs(&zerrors.NestErr{
			Err: io.EOF,
			Msg: "outer error",
		})
		ztesting.AssertEqual(t, "unexpected error returned.", nil, err)
	})

	t.Run("has interface", func(t *testing.T) {
		errs := zerrors.UnwrapErrs(errors.Join(io.EOF, io.ErrUnexpectedEOF)) // Implements interface{ Unwrap() []error }
		ztesting.AssertEqual(t, "wrong number of errors.", 2, len(errs))
		ztesting.AssertEqual(t, "wrong unwrapped error returned.", []error{io.EOF, io.ErrUnexpectedEOF}, errs)
	})
}

func TestMust(t *testing.T) {
	t.Parallel()
	t.Run("nil error", func(t *testing.T) {
		val := zerrors.Must("dummy", nil)
		if val != "dummy" {
			t.Errorf("unexpected value returned. want:\"dummy\" got:%s", val)
		}
	})

	t.Run("non-nil error", func(t *testing.T) {
		defer func() {
			rec := recover()
			if rec != io.EOF {
				t.Errorf("unexpected panic value returned. want:%#v got:%#v", io.EOF, rec)
			}
		}()
		_ = zerrors.Must("dummy", io.EOF)
	})
}

func TestMustNil(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		defer func() {
			rec := recover()
			if rec != nil {
				t.Errorf("unexpected panic value returned. want:%#v got:%#v", nil, rec)
			}
		}()
		zerrors.MustNil(nil)
	})

	t.Run("non-nil", func(t *testing.T) {
		defer func() {
			rec := recover()
			if rec != io.EOF {
				t.Errorf("unexpected panic value returned. want:%#v got:%#v", io.EOF, rec)
			}
		}()
		zerrors.MustNil(io.EOF)
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		err := zerrors.Wrap(nil, "test")
		if err != nil {
			t.Errorf("unexpected error returned. want:%#v got:%#v", nil, err)
		}
	})

	t.Run("non-nil", func(t *testing.T) {
		err := zerrors.Wrap(io.EOF, "test")
		if err == nil {
			t.Errorf("returned error is nil. got:%#v", err)
		} else {
			if err.Err != io.EOF {
				t.Errorf("inner error is incorrect. want:%#v got:%#v", io.EOF, err.Err)
			}
			if err.Msg != "test" {
				t.Errorf("outer error message is incorrect. want:'%#v' got:'%#v'", "test", err.Msg)
			}
		}
	})
}
