package internal_test

import (
	"io"
	"testing"

	"github.com/aileron-projects/go/zcrypto/internal"
)

func TestMustNil(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		defer func() {
			r := recover()
			if r != nil {
				t.Error("recovered value is not nil")
			}
		}()
		internal.MustNil(nil)
	})
	t.Run("non nil", func(t *testing.T) {
		defer func() {
			r := recover()
			if r.(error) != io.EOF {
				t.Error("recovered error is not io.EOF")
			}
		}()
		internal.MustNil(io.EOF)
	})
}
