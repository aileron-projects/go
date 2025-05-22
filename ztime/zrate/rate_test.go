package zrate

import (
	"io"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestNoopToken(t *testing.T) {
	t.Parallel()
	t.Run("true", func(t *testing.T) {
		ztesting.AssertEqual(t, "incorrect status of the token", true, TokenOK.OK())
		ztesting.AssertEqual(t, "non nil error returned", nil, TokenOK.Err())
		TokenOK.Release() // Nothing happends. Just for taking coverage.
	})
	t.Run("false", func(t *testing.T) {
		ztesting.AssertEqual(t, "incorrect status of the token", false, TokenNG.OK())
		ztesting.AssertEqual(t, "non nil error returned", nil, TokenNG.Err())
		TokenNG.Release() // Nothing happends. Just for taking coverage.
	})
}

func TestToken(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		tk := &token{ok: true}
		ztesting.AssertEqual(t, "incorrect status of the token", true, tk.OK())
	})
	t.Run("ng", func(t *testing.T) {
		tk := &token{ok: false}
		ztesting.AssertEqual(t, "incorrect status of the token", false, tk.OK())
	})
	t.Run("error", func(t *testing.T) {
		tk := &token{err: io.EOF}
		ztesting.AssertEqual(t, "error mismatched", io.EOF, tk.Err())
	})
	t.Run("nil release func", func(t *testing.T) {
		tk := &token{releaseFunc: nil}
		tk.Release() // Nothing happens.
	})
	t.Run("non nil release func", func(t *testing.T) {
		var callCount int
		tk := &token{releaseFunc: func() { callCount += 1 }}
		ztesting.AssertEqual(t, "incorrect call count", 0, callCount)
		tk.Release()
		ztesting.AssertEqual(t, "release func not executed", 1, callCount)
		tk.Release()
		tk.Release()
		ztesting.AssertEqual(t, "release func called multiple times", 1, callCount)
	})
}
