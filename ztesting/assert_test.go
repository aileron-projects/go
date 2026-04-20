package ztesting

import (
	"fmt"
	"io"
	"testing"
)

type testT struct {
	helperCalled bool
	gotArgs      []any
}

func (t *testT) Helper()           { t.helperCalled = true }
func (t *testT) Error(args ...any) { t.gotArgs = args }

func TestAssertEqual(t *testing.T) {
	t.Cleanup(func() { mocT = nil })
	t.Run("equal", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqual(t, "foo", "foo")
		mocT = nil // Reset moc
		AssertEqual(t, true, tt.helperCalled)
		AssertEqual(t, 0, len(tt.gotArgs))
	})
	t.Run("not equal", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqual(t, "foo", "bar")
		mocT = nil // Reset moc
		AssertEqual(t, true, tt.helperCalled)
		AssertEqual(t, 1, len(tt.gotArgs))
	})
}

func TestAssertEqualErr(t *testing.T) {
	t.Cleanup(func() { mocT = nil })
	t.Run("equal pointer", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqualErr(t, io.EOF, io.EOF)
		mocT = nil // Reset moc
		AssertEqual(t, true, tt.helperCalled)
		AssertEqual(t, 0, len(tt.gotArgs))
	})
	t.Run("equal by is", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqualErr(t, io.EOF, fmt.Errorf("wrap [%w]", io.EOF))
		mocT = nil // Reset moc
		AssertEqual(t, true, tt.helperCalled)
		AssertEqual(t, 0, len(tt.gotArgs))
	})
	t.Run("equal by message", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqualErr(t, io.EOF, fmt.Errorf("EOF"))
		mocT = nil // Reset moc
		AssertEqual(t, true, tt.helperCalled)
		AssertEqual(t, 0, len(tt.gotArgs))
	})
	t.Run("not equal", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqualErr(t, io.EOF, io.ErrUnexpectedEOF)
		mocT = nil // Reset moc
		AssertEqual(t, true, tt.helperCalled)
		AssertEqual(t, 1, len(tt.gotArgs))
	})
}
