package internal_test

import (
	"io"
	"testing"

	"github.com/aileron-projects/go/znet/internal"
	"github.com/aileron-projects/go/ztesting"
)

type testCloser struct {
	name     string
	closed   int
	closeErr error
}

func (c *testCloser) Close() error {
	c.closed++
	return c.closeErr
}

func TestCloserStore(t *testing.T) {
	t.Parallel()
	t.Run("store single value", func(t *testing.T) {
		s := internal.CloserStore[*testCloser]{}
		foo := &testCloser{name: "foo"}
		s.Store(foo)
		ztesting.AssertEqual(t, "length not match", 1, s.Length())
		s.Delete(foo)
		ztesting.AssertEqual(t, "length not match", 0, s.Length())
	})
	t.Run("store different value", func(t *testing.T) {
		s := internal.CloserStore[*testCloser]{}
		foo := &testCloser{name: "foo"}
		bar := &testCloser{name: "bar"}
		s.Store(foo)
		ztesting.AssertEqual(t, "length not match", 1, s.Length())
		s.Store(bar)
		ztesting.AssertEqual(t, "length not match", 2, s.Length())
		s.Delete(foo)
		s.Delete(bar)
		ztesting.AssertEqual(t, "length not match", 0, s.Length())
	})
	t.Run("close", func(t *testing.T) {
		s := internal.CloserStore[*testCloser]{}
		foo := &testCloser{name: "foo"}
		bar := &testCloser{name: "bar"}
		s.Store(foo)
		s.Store(bar)
		ztesting.AssertEqual(t, "length not match", 2, s.Length())
		err := s.CloseAll()
		ztesting.AssertEqualErr(t, "close error not match", nil, err)
		ztesting.AssertEqual(t, "length not match", 0, s.Length())
		ztesting.AssertEqual(t, "foo is not closed", 1, foo.closed)
		ztesting.AssertEqual(t, "bar is not closed", 1, bar.closed)
	})
	t.Run("close error", func(t *testing.T) {
		s := internal.CloserStore[*testCloser]{}
		foo := &testCloser{name: "foo", closeErr: io.EOF}              // Return dummy error
		bar := &testCloser{name: "bar", closeErr: io.ErrUnexpectedEOF} // Return dummy error
		s.Store(foo)
		s.Store(bar)
		err := s.CloseAll()
		errs := err.(interface{ Unwrap() []error }).Unwrap()
		ztesting.AssertEqual(t, "length of errors not match", 2, len(errs))
		ztesting.AssertEqual(t, "foo is not closed", 1, foo.closed)
		ztesting.AssertEqual(t, "bar is not closed", 1, bar.closed)
	})
}
