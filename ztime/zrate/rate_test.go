package zrate

import (
	"io"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestNoopToken(t *testing.T) {
	t.Parallel()
	t.Run("true", func(t *testing.T) {
		ztesting.AssertEqual(t, true, TokenOK.OK())
		ztesting.AssertEqual(t, nil, TokenOK.Err())
		TokenOK.Release() // Nothing happens. Just for taking coverage.
	})
	t.Run("false", func(t *testing.T) {
		ztesting.AssertEqual(t, false, TokenNG.OK())
		ztesting.AssertEqual(t, nil, TokenNG.Err())
		TokenNG.Release() // Nothing happens. Just for taking coverage.
	})
}

func TestToken(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		tk := &token{ok: true}
		ztesting.AssertEqual(t, true, tk.OK())
	})
	t.Run("ng", func(t *testing.T) {
		tk := &token{ok: false}
		ztesting.AssertEqual(t, false, tk.OK())
	})
	t.Run("error", func(t *testing.T) {
		tk := &token{err: io.EOF}
		ztesting.AssertEqual(t, io.EOF, tk.Err())
	})
	t.Run("nil release func", func(t *testing.T) {
		tk := &token{releaseFunc: nil}
		tk.Release() // Nothing happens.
	})
	t.Run("non nil release func", func(t *testing.T) {
		var callCount int
		tk := &token{releaseFunc: func() { callCount += 1 }}
		ztesting.AssertEqual(t, 0, callCount)
		tk.Release()
		ztesting.AssertEqual(t, 1, callCount)
		tk.Release()
		tk.Release()
		ztesting.AssertEqual(t, 1, callCount)
	})
}
