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
		AssertEqual(t, "error message", "foo", "foo")
		mocT = nil // Reset moc
		AssertEqual(t, "helper is not called", true, tt.helperCalled)
		AssertEqual(t, "length of args not match", 0, len(tt.gotArgs))
	})
	t.Run("not equal", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqual(t, "error message", "foo", "bar")
		mocT = nil // Reset moc
		AssertEqual(t, "helper is not called", true, tt.helperCalled)
		AssertEqual(t, "length of args not match", 1, len(tt.gotArgs))
	})
}

func TestAssertEqualErr(t *testing.T) {
	t.Cleanup(func() { mocT = nil })
	t.Run("equal pointer", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqualErr(t, "error message", io.EOF, io.EOF)
		mocT = nil // Reset moc
		AssertEqual(t, "helper is not called", true, tt.helperCalled)
		AssertEqual(t, "length of args not match", 0, len(tt.gotArgs))
	})
	t.Run("equal by is", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqualErr(t, "error message", io.EOF, fmt.Errorf("wrap [%w]", io.EOF))
		mocT = nil // Reset moc
		AssertEqual(t, "helper is not called", true, tt.helperCalled)
		AssertEqual(t, "length of args not match", 0, len(tt.gotArgs))
	})
	t.Run("equal by message", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqualErr(t, "error message", io.EOF, fmt.Errorf("EOF"))
		mocT = nil // Reset moc
		AssertEqual(t, "helper is not called", true, tt.helperCalled)
		AssertEqual(t, "length of args not match", 0, len(tt.gotArgs))
	})
	t.Run("not equal", func(t *testing.T) {
		tt := &testT{}
		mocT = tt
		AssertEqualErr(t, "error message", io.EOF, io.ErrUnexpectedEOF)
		mocT = nil // Reset moc
		AssertEqual(t, "helper is not called", true, tt.helperCalled)
		AssertEqual(t, "length of args not match", 1, len(tt.gotArgs))
	})
}
