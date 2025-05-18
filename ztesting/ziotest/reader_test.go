package ziotest_test

import (
	"io"
	"strings"
	"testing"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestCharsetReader(t *testing.T) {
	t.Parallel()

	t.Run("empty charset", func(t *testing.T) {
		cr := ziotest.CharsetReader("", false)
		buf := make([]byte, 3)
		n, err := cr.Read(buf)
		ztesting.AssertEqual(t, "invalid read bytes.", 0, n)
		ztesting.AssertEqual(t, "invalid returned error.", io.EOF, err)
		ztesting.AssertEqual(t, "invalid read string.", "\x00\x00\x00", string(buf))
	})

	t.Run("read chart without loop", func(t *testing.T) {
		cr := ziotest.CharsetReader("12345", false)
		var n int
		var err error

		// Read first
		buf1 := make([]byte, 3)
		n, err = cr.Read(buf1)
		ztesting.AssertEqual(t, "invalid read bytes.", 3, n)
		ztesting.AssertEqual(t, "invalid returned error.", nil, err)
		ztesting.AssertEqual(t, "invalid read string.", "123", string(buf1))

		// Read second
		buf2 := make([]byte, 3)
		n, err = cr.Read(buf2)
		ztesting.AssertEqual(t, "invalid read bytes.", 2, n)
		ztesting.AssertEqual(t, "invalid returned error.", io.EOF, err)
		ztesting.AssertEqual(t, "invalid read string.", "45\x00", string(buf2))
	})

	t.Run("read chart with loop", func(t *testing.T) {
		cr := ziotest.CharsetReader("12345", true)
		var n int
		var err error

		// Read first
		buf1 := make([]byte, 3)
		n, err = cr.Read(buf1)
		ztesting.AssertEqual(t, "invalid read bytes.", 3, n)
		ztesting.AssertEqual(t, "invalid returned error.", nil, err)
		ztesting.AssertEqual(t, "invalid read string.", "123", string(buf1))

		// Read second
		buf2 := make([]byte, 3)
		n, err = cr.Read(buf2)
		ztesting.AssertEqual(t, "invalid read bytes.", 3, n)
		ztesting.AssertEqual(t, "invalid returned error.", nil, err)
		ztesting.AssertEqual(t, "invalid read string.", "451", string(buf2))
	})
}

func TestErrReader(t *testing.T) {
	t.Parallel()

	t.Run("read -10", func(t *testing.T) {
		er := ziotest.ErrReader(strings.NewReader("123456789"), -10)
		buf := make([]byte, 3)
		n, err := er.Read(buf)
		ztesting.AssertEqual(t, "invalid read bytes.", 0, n)
		ztesting.AssertEqual(t, "invalid returned error.", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "invalid read string.", "\x00\x00\x00", string(buf))
	})

	t.Run("read 0", func(t *testing.T) {
		er := ziotest.ErrReader(strings.NewReader("123456789"), 0)
		buf := make([]byte, 3)
		n, err := er.Read(buf)
		ztesting.AssertEqual(t, "invalid read bytes.", 0, n)
		ztesting.AssertEqual(t, "invalid returned error.", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "invalid read string.", "\x00\x00\x00", string(buf))
	})

	t.Run("read more than n", func(t *testing.T) {
		er := ziotest.ErrReader(strings.NewReader("123456789"), 5)
		var n int
		var err error

		// Read first
		buf1 := make([]byte, 3)
		n, err = er.Read(buf1)
		ztesting.AssertEqual(t, "invalid read bytes.", 3, n)
		ztesting.AssertEqual(t, "invalid returned error.", nil, err)
		ztesting.AssertEqual(t, "invalid read string.", "123", string(buf1))

		// Read second
		buf2 := make([]byte, 3)
		n, err = er.Read(buf2)
		ztesting.AssertEqual(t, "invalid read bytes.", 2, n)
		ztesting.AssertEqual(t, "invalid returned error.", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "invalid read string.", "45\x00", string(buf2))
	})

	t.Run("read exactly n", func(t *testing.T) {
		er := ziotest.ErrReader(strings.NewReader("123456789"), 6)
		var n int
		var err error

		// Read first
		buf1 := make([]byte, 3)
		n, err = er.Read(buf1)
		ztesting.AssertEqual(t, "invalid read bytes.", 3, n)
		ztesting.AssertEqual(t, "invalid returned error.", nil, err)
		ztesting.AssertEqual(t, "invalid read string.", "123", string(buf1))

		// Read second
		buf2 := make([]byte, 3)
		n, err = er.Read(buf2)
		ztesting.AssertEqual(t, "invalid read bytes.", 3, n)
		ztesting.AssertEqual(t, "invalid returned error.", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "invalid read string.", "456", string(buf2))
	})
}
