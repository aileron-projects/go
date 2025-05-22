package zio_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zio"
	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestTeeReader(t *testing.T) {
	t.Parallel()
	t.Run("nil reader", func(t *testing.T) {
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReader(nil, w)
		ztesting.AssertEqual(t, "reader does not match", nil, tr)
	})

	t.Run("nil writer", func(t *testing.T) {
		r := strings.NewReader("abc")
		tr := zio.TeeReader(r, nil)
		ztesting.AssertEqual(t, "reader does not match", io.Reader(r), tr)
	})

	t.Run("no error", func(t *testing.T) {
		r := strings.NewReader("abc")
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReader(r, w)
		buf := make([]byte, 10)
		n, err := tr.Read(buf)
		ztesting.AssertEqual(t, "error is not nil", nil, err)
		ztesting.AssertEqual(t, "wrong read bytes", 3, n)
		ztesting.AssertEqual(t, "wrong content written to writer", "abc", w.String())
		ztesting.AssertEqual(t, "wrong content read", "abc", string(buf[:n]))
	})

	t.Run("read error", func(t *testing.T) {
		r := ziotest.ErrReader(strings.NewReader("abc"), 1)
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReader(r, w)
		buf := make([]byte, 10)
		n, err := tr.Read(buf)
		ztesting.AssertEqual(t, "error mismatch", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "wrong read bytes", 1, n)
		ztesting.AssertEqual(t, "wrong content written to writer", "a", w.String())
		ztesting.AssertEqual(t, "wrong content read", "a", string(buf[:n]))
	})

	t.Run("write error", func(t *testing.T) {
		r := strings.NewReader("abc")
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReader(r, ziotest.ErrWriter(w, 1))
		buf := make([]byte, 10)
		n, err := tr.Read(buf)
		ztesting.AssertEqual(t, "error mismatch", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "wrong read bytes", 1, n)
		ztesting.AssertEqual(t, "wrong content written to writer", "a", w.String())
		ztesting.AssertEqual(t, "wrong content read", "a", string(buf[:n]))
	})

	t.Run("read write error", func(t *testing.T) {
		r := ziotest.ErrReader(strings.NewReader("abc"), 1)
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReader(r, ziotest.ErrWriter(w, 1))
		buf := make([]byte, 10)
		n, err := tr.Read(buf)
		ztesting.AssertEqual(t, "error mismatch", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "wrong read bytes", 1, n)
		ztesting.AssertEqual(t, "wrong content written to writer", "a", w.String())
		ztesting.AssertEqual(t, "wrong content read", "a", string(buf[:n]))
	})
}

func TestTeeReadCloser(t *testing.T) {
	t.Parallel()
	t.Run("nil reader", func(t *testing.T) {
		w := zio.NopWriteCloser(bytes.NewBuffer(nil))
		tr := zio.TeeReadCloser(nil, w)
		ztesting.AssertEqual(t, "reader does not match", nil, tr)
	})
	t.Run("nil writer", func(t *testing.T) {
		r := zio.NopReadCloser(strings.NewReader("abc"))
		tr := zio.TeeReadCloser(r, nil)
		ztesting.AssertEqual(t, "reader does not match", io.ReadCloser(r), tr)
	})
	t.Run("no error", func(t *testing.T) {
		r := zio.NopReadCloser(strings.NewReader("abc"))
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReadCloser(r, zio.NopWriteCloser(w))
		buf := make([]byte, 10)
		n, err := tr.Read(buf)
		ztesting.AssertEqual(t, "error is not nil", nil, err)
		ztesting.AssertEqual(t, "wrong read bytes", 3, n)
		ztesting.AssertEqual(t, "wrong content written to writer", "abc", w.String())
		ztesting.AssertEqual(t, "wrong content read", "abc", string(buf[:n]))
	})
	t.Run("read error", func(t *testing.T) {
		r := zio.NopReadCloser(ziotest.ErrReader(strings.NewReader("abc"), 1))
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReadCloser(r, zio.NopWriteCloser(w))
		buf := make([]byte, 10)
		n, err := tr.Read(buf)
		ztesting.AssertEqual(t, "error mismatch", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "wrong read bytes", 1, n)
		ztesting.AssertEqual(t, "wrong content written to writer", "a", w.String())
		ztesting.AssertEqual(t, "wrong content read", "a", string(buf[:n]))
	})
	t.Run("write error", func(t *testing.T) {
		r := zio.NopReadCloser(strings.NewReader("abc"))
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReadCloser(r, zio.NopWriteCloser(ziotest.ErrWriter(w, 1)))
		buf := make([]byte, 10)
		n, err := tr.Read(buf)
		ztesting.AssertEqual(t, "error mismatch", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "wrong read bytes", 1, n)
		ztesting.AssertEqual(t, "wrong content written to writer", "a", w.String())
		ztesting.AssertEqual(t, "wrong content read", "a", string(buf[:n]))
	})
	t.Run("read write error", func(t *testing.T) {
		r := zio.NopReadCloser(ziotest.ErrReader(strings.NewReader("abc"), 1))
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReadCloser(r, zio.NopWriteCloser(ziotest.ErrWriter(w, 1)))
		buf := make([]byte, 10)
		n, err := tr.Read(buf)
		ztesting.AssertEqual(t, "error mismatch", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "wrong read bytes", 1, n)
		ztesting.AssertEqual(t, "wrong content written to writer", "a", w.String())
		ztesting.AssertEqual(t, "wrong content read", "a", string(buf[:n]))
	})
	t.Run("no close error", func(t *testing.T) {
		r := zio.NopReadCloser(strings.NewReader("abc"))
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReadCloser(r, zio.NopWriteCloser(w))
		err := tr.Close()
		ztesting.AssertEqualErr(t, "error mismatch", nil, err)
	})
	t.Run("reader close error", func(t *testing.T) {
		r := &errReadCloser{Reader: strings.NewReader("abc"), err: os.ErrClosed}
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReadCloser(r, zio.NopWriteCloser(w))
		err := tr.Close()
		ztesting.AssertEqualErr(t, "error mismatch", os.ErrClosed, err)
	})
	t.Run("writer close error", func(t *testing.T) {
		r := zio.NopReadCloser(ziotest.ErrReader(strings.NewReader("abc"), 1))
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReadCloser(r, &errWriteCloser{Writer: w, err: os.ErrClosed})
		err := tr.Close()
		ztesting.AssertEqualErr(t, "error mismatch", os.ErrClosed, err)
	})
	t.Run("reader writer close error", func(t *testing.T) {
		r := &errReadCloser{Reader: strings.NewReader("abc"), err: os.ErrPermission}
		w := bytes.NewBuffer(nil)
		tr := zio.TeeReadCloser(r, &errWriteCloser{Writer: w, err: os.ErrClosed})
		err := tr.Close()
		ztesting.AssertEqualErr(t, "error mismatch", errors.Join(os.ErrPermission, os.ErrClosed), err)
	})
}

type errReadCloser struct {
	io.Reader
	err error
}

func (r *errReadCloser) Close() error {
	return r.err
}

type errWriteCloser struct {
	io.Writer
	err error
}

func (w *errWriteCloser) Close() error {
	return w.err
}
