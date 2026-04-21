package ziotest_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestErrWriter(t *testing.T) {
	t.Parallel()

	t.Run("n is -10", func(t *testing.T) {
		var buf bytes.Buffer
		ew := ziotest.ErrWriter(&buf, -10)
		n, err := ew.Write([]byte("1234"))
		ztesting.AssertEqual(t, 0, n)
		ztesting.AssertEqual(t, io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "", buf.String())
	})

	t.Run("n is 0", func(t *testing.T) {
		var buf bytes.Buffer
		ew := ziotest.ErrWriter(&buf, 0)
		n, err := ew.Write([]byte("1234"))
		ztesting.AssertEqual(t, 0, n)
		ztesting.AssertEqual(t, io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "", buf.String())
	})

	t.Run("write more than n", func(t *testing.T) {
		var buf bytes.Buffer
		ew := ziotest.ErrWriter(&buf, 9)
		var n int
		var err error

		// First write.
		n, err = ew.Write([]byte("1234"))
		ztesting.AssertEqual(t, 4, n)
		ztesting.AssertEqual(t, nil, err)
		ztesting.AssertEqual(t, "1234", buf.String())

		// Second write.
		n, err = ew.Write([]byte("5678"))
		ztesting.AssertEqual(t, 4, n)
		ztesting.AssertEqual(t, nil, err)
		ztesting.AssertEqual(t, "12345678", buf.String())

		// Third write.
		n, err = ew.Write([]byte("9012"))
		ztesting.AssertEqual(t, 1, n)
		ztesting.AssertEqual(t, io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "123456789", buf.String())
	})

	t.Run("write exactly 8", func(t *testing.T) {
		var buf bytes.Buffer
		ew := ziotest.ErrWriter(&buf, 8)
		var n int
		var err error

		// First write.
		n, err = ew.Write([]byte("1234"))
		ztesting.AssertEqual(t, 4, n)
		ztesting.AssertEqual(t, nil, err)
		ztesting.AssertEqual(t, "1234", buf.String())

		// Second write.
		n, err = ew.Write([]byte("5678"))
		ztesting.AssertEqual(t, 4, n)
		ztesting.AssertEqual(t, io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "12345678", buf.String())
	})

	t.Run("inner writer error", func(t *testing.T) {
		var buf bytes.Buffer
		ew := ziotest.ErrWriter(ziotest.ErrWriter(&buf, 2), 4)
		n, err := ew.Write([]byte("123"))
		ztesting.AssertEqual(t, 2, n)
		ztesting.AssertEqual(t, io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "12", buf.String())
	})
}
