package zerrors

import (
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestCallerFrames(t *testing.T) {
	t.Parallel()
	t.Run("skip=0", func(t *testing.T) {
		fs := callerFrames(0, 1)
		ztesting.AssertEqual(t, "length not match", 1, len(fs))
	})
	t.Run("skip=-999", func(t *testing.T) {
		fs := callerFrames(-999, 1)
		ztesting.AssertEqual(t, "length not match", 1, len(fs))
	})
	t.Run("skip=999", func(t *testing.T) {
		fs := callerFrames(999, 1)
		ztesting.AssertEqual(t, "length not match", 0, len(fs))
	})
	t.Run("size=0", func(t *testing.T) {
		fs := callerFrames(0, 0)
		ztesting.AssertEqual(t, "length not match", 0, len(fs))
	})
	t.Run("size=2", func(t *testing.T) {
		fs := callerFrames(0, 2)
		ztesting.AssertEqual(t, "length not match", 2, len(fs))
	})
}
