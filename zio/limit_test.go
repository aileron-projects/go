package zio_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zio"
	"github.com/aileron-projects/go/ztesting"
)

func TestLimitReader(t *testing.T) {
	t.Parallel()
	t.Run("nil reader", func(t *testing.T) {
		r := zio.LimitReader(nil, 10)
		ztesting.AssertEqual(t, "non nil reader returned", nil, r)
	})

	t.Run("limit=-5", func(t *testing.T) {
		r := zio.LimitReader(strings.NewReader("1234567890"), -5)
		buf := make([]byte, 10)
		n, err := r.Read(buf)
		ztesting.AssertEqual(t, "read bytes not match", 0, n)
		ztesting.AssertEqual(t, "error not match", zio.ErrReadLimit, err)
		ztesting.AssertEqual(t, "read content invalid", make([]byte, 10), buf)
	})

	t.Run("limit=0", func(t *testing.T) {
		r := zio.LimitReader(strings.NewReader("1234567890"), 0)
		buf := make([]byte, 10)
		n, err := r.Read(buf)
		ztesting.AssertEqual(t, "read bytes not match", 0, n)
		ztesting.AssertEqual(t, "error not match", zio.ErrReadLimit, err)
		ztesting.AssertEqual(t, "read content invalid", make([]byte, 10), buf)
	})

	t.Run("limit=1", func(t *testing.T) {
		r := zio.LimitReader(strings.NewReader("1234567890"), 1)
		buf := make([]byte, 10)
		n, err := r.Read(buf)
		ztesting.AssertEqual(t, "read bytes not match", 1, n)
		ztesting.AssertEqual(t, "error not match", zio.ErrReadLimit, err)
		ztesting.AssertEqual(t, "read content invalid", []byte{'1', 0, 0, 0, 0, 0, 0, 0, 0, 0}, buf)
	})

	t.Run("limit=5", func(t *testing.T) {
		r := zio.LimitReader(strings.NewReader("1234567890"), 5)
		buf := make([]byte, 10)
		n, err := r.Read(buf)
		ztesting.AssertEqual(t, "read bytes not match", 5, n)
		ztesting.AssertEqual(t, "error not match", zio.ErrReadLimit, err)
		ztesting.AssertEqual(t, "read content invalid", []byte{'1', '2', '3', '4', '5', 0, 0, 0, 0, 0}, buf)
	})

	t.Run("limit=10", func(t *testing.T) {
		r := zio.LimitReader(strings.NewReader("1234567890"), 10)
		buf := make([]byte, 10)
		n, err := r.Read(buf)
		ztesting.AssertEqual(t, "read bytes not match", 10, n)
		ztesting.AssertEqual(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "read content invalid", []byte("1234567890"), buf)
	})

	t.Run("read multiple times", func(t *testing.T) {
		r := zio.LimitReader(strings.NewReader("1234567890"), 5)
		n1, err1 := r.Read(make([]byte, 3))
		ztesting.AssertEqual(t, "read bytes not match", 3, n1)
		ztesting.AssertEqual(t, "error not match", nil, err1)
		n2, err2 := r.Read(make([]byte, 3))
		ztesting.AssertEqual(t, "read bytes not match", 2, n2)
		ztesting.AssertEqual(t, "error not match", zio.ErrReadLimit, err2)
	})
}

func TestLimitWriter(t *testing.T) {
	t.Parallel()
	t.Run("nil reader", func(t *testing.T) {
		w := zio.LimitWriter(nil, 10)
		ztesting.AssertEqual(t, "non nil writer returned", nil, w)
	})

	t.Run("limit=-5", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		w := zio.LimitWriter(buf, -5)
		n, err := w.Write([]byte("1234567890"))
		ztesting.AssertEqual(t, "written bytes not match", 0, n)
		ztesting.AssertEqual(t, "error not match", zio.ErrWriteLimit, err)
		ztesting.AssertEqual(t, "written content invalid", "", buf.String())
	})

	t.Run("limit=0", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		w := zio.LimitWriter(buf, 0)
		n, err := w.Write([]byte("1234567890"))
		ztesting.AssertEqual(t, "written bytes not match", 0, n)
		ztesting.AssertEqual(t, "error not match", zio.ErrWriteLimit, err)
		ztesting.AssertEqual(t, "written content invalid", "", buf.String())
	})

	t.Run("limit=1", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		w := zio.LimitWriter(buf, 1)
		n, err := w.Write([]byte("1234567890"))
		ztesting.AssertEqual(t, "written bytes not match", 1, n)
		ztesting.AssertEqual(t, "error not match", zio.ErrWriteLimit, err)
		ztesting.AssertEqual(t, "written content invalid", "1", buf.String())
	})

	t.Run("limit=5", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		w := zio.LimitWriter(buf, 5)
		n, err := w.Write([]byte("1234567890"))
		ztesting.AssertEqual(t, "written bytes not match", 5, n)
		ztesting.AssertEqual(t, "error not match", zio.ErrWriteLimit, err)
		ztesting.AssertEqual(t, "written content invalid", "12345", buf.String())
	})

	t.Run("limit=10", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		w := zio.LimitWriter(buf, 10)
		n, err := w.Write([]byte("1234567890"))
		ztesting.AssertEqual(t, "written bytes not match", 10, n)
		ztesting.AssertEqual(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "written content invalid", "1234567890", buf.String())
	})

	t.Run("write multiple times", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		w := zio.LimitWriter(buf, 5)
		n1, err1 := w.Write([]byte("123"))
		ztesting.AssertEqual(t, "written bytes not match", 3, n1)
		ztesting.AssertEqual(t, "error not match", nil, err1)
		ztesting.AssertEqual(t, "written content invalid", "123", buf.String())
		n2, err2 := w.Write([]byte("456"))
		ztesting.AssertEqual(t, "written bytes not match", 2, n2)
		ztesting.AssertEqual(t, "error not match", zio.ErrWriteLimit, err2)
		ztesting.AssertEqual(t, "written content invalid", "12345", buf.String())
	})
}
