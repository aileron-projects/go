package zio_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zio"
	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestCopy(t *testing.T) {
	t.Parallel()

	t.Run("read write success", func(t *testing.T) {
		r := strings.NewReader("1234567890")
		w := bytes.NewBuffer(nil)
		written, err := zio.Copy(w, r)
		ztesting.AssertEqual(t, "invalid written bytes", 10, written)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})

	t.Run("reader EOF", func(t *testing.T) {
		r := ziotest.ErrReaderWith(strings.NewReader("1234567890"), 5, io.EOF)
		w := bytes.NewBuffer(nil)
		written, err := zio.Copy(w, r)
		ztesting.AssertEqual(t, "invalid written bytes", 5, written)
		ztesting.AssertEqual(t, "error not matched", nil, err)
	})

	t.Run("reader UnexpectedEOF", func(t *testing.T) {
		r := ziotest.ErrReaderWith(strings.NewReader("1234567890"), 5, io.ErrUnexpectedEOF)
		w := bytes.NewBuffer(nil)
		written, err := zio.Copy(w, r)
		ztesting.AssertEqual(t, "invalid written bytes", 5, written)
		ztesting.AssertEqual(t, "error not matched", io.ErrUnexpectedEOF, err)
	})

	t.Run("write error", func(t *testing.T) {
		r := strings.NewReader("1234567890")
		w := ziotest.ErrWriterWith(bytes.NewBuffer(nil), 5, io.ErrClosedPipe)
		written, err := zio.Copy(w, r)
		ztesting.AssertEqual(t, "invalid written bytes", 5, written)
		ztesting.AssertEqual(t, "error not matched", io.ErrClosedPipe, err)
	})

	t.Run("read write error", func(t *testing.T) {
		r := ziotest.ErrReaderWith(strings.NewReader("1234567890"), 5, io.ErrUnexpectedEOF)
		w := ziotest.ErrWriterWith(bytes.NewBuffer(nil), 4, io.ErrClosedPipe)
		written, err := zio.Copy(w, r)
		ztesting.AssertEqual(t, "invalid written bytes", 4, written) // Write error is returned.
		ztesting.AssertEqual(t, "error not matched", io.ErrClosedPipe, err)
	})

	t.Run("nRead != nWrite", func(t *testing.T) {
		r := strings.NewReader("1234567890")
		w := ziotest.ShortWriter(bytes.NewBuffer(nil), 4)
		written, err := zio.Copy(w, r)
		ztesting.AssertEqual(t, "invalid written bytes", 4, written)
		ztesting.AssertEqual(t, "error not matched", io.ErrShortWrite, err)
	})
}
