package ztesting

import (
	"cmp"
	"errors"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type test interface {
	Helper()
	Error(args ...any)
}

// mocT replaces [testing.T] for testing.
var mocT test

// AssertEqual checks if the given two values are the same.
// AssertEqual compares want and got using [reflect.DeepEqual].
// See https://go.dev/wiki/TestComments
func AssertEqual[T any](t *testing.T, msg string, want, got T) {
	tt := cmp.Or(mocT, test(t)) // Use mocT to test this func.
	tt.Helper()
	if reflect.DeepEqual(want, got) {
		return
	}
	msg += "\n"
	msg += "-want: " + spew.Sdump(want)
	msg += "+got: " + spew.Sdump(got)
	tt.Error(msg)
}

// AssertEqualErr checks if the given two errors are the same.
// Errors are checked by following order and considered the same
// when one of them returned true.
//
//   - Compare pointer: want == got
//   - Compare error: errors.Is(got, want)
//   - Compare message: want.Error() == got.Error()
//
// See https://go.dev/wiki/TestComments
func AssertEqualErr(t *testing.T, errReason string, want, got error) {
	tt := cmp.Or(mocT, test(t)) // Use mocT to test this func itself.
	tt.Helper()
	if want == got {
		return // nil == nil is also here.
	}
	if errors.Is(got, want) {
		return
	}
	if want != nil && got != nil {
		if want.Error() == got.Error() {
			return
		}
	}
	errReason += "\n"
	errReason += "-want: " + spew.Sdump(want)
	errReason += "+got: " + spew.Sdump(got)
	tt.Error(errReason)
}
