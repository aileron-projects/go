package zstrconv

import (
	"reflect"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestParseNum(t *testing.T) {
	t.Parallel()
	t.Run("unsupported type", func(t *testing.T) {
		vv, err := paseNum("true", reflect.Bool)
		ztesting.AssertEqual(t, "incorrect parse result", nil, vv)
		ztesting.AssertEqualErr(t, "non nil error returned", ErrTypeSupported, err)
	})
}
