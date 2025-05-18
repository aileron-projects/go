package zhttp_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-projects/go/znet/zhttp"
	"github.com/aileron-projects/go/ztesting"
)

func TestWrapResponseWriter(t *testing.T) {
	t.Run("new wrapper", func(t *testing.T) {
		w := http.ResponseWriter(httptest.NewRecorder())
		ww := zhttp.WrapResponseWriter(w)
		ztesting.AssertEqual(t, "internal writer not match", w, ww.Unwrap())
	})
	t.Run("already wrapped", func(t *testing.T) {
		w := http.ResponseWriter(httptest.NewRecorder())
		ww1 := zhttp.WrapResponseWriter(w)
		ww2 := zhttp.WrapResponseWriter(ww1)
		ztesting.AssertEqual(t, "internal writer not match", w, ww2.Unwrap())
	})
}

func TestResponseWrapper(t *testing.T) {
	t.Run("not written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ztesting.AssertEqual(t, "status code not match", -1, ww.StatusCode())
		ztesting.AssertEqual(t, "written bytes not match", -1, ww.WrittenBytes())
	})
	t.Run("header written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.Header().Set("Test", "value") // No affection.
		ztesting.AssertEqual(t, "status code not match", -1, ww.StatusCode())
		ztesting.AssertEqual(t, "header not match", "value", ww.Header().Get("Test"))
		ztesting.AssertEqual(t, "written bytes not match", -1, ww.WrittenBytes())
	})
	t.Run("status written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.WriteHeader(http.StatusBadRequest)
		ztesting.AssertEqual(t, "status code not match", http.StatusBadRequest, ww.StatusCode())
		ztesting.AssertEqual(t, "written bytes not match", 0, ww.WrittenBytes())
	})
	t.Run("nil body written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.Write(nil)
		ztesting.AssertEqual(t, "status code not match", http.StatusOK, ww.StatusCode())
		ztesting.AssertEqual(t, "written bytes not match", 0, ww.WrittenBytes())
	})
	t.Run("non-nil body written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.Write([]byte("test"))
		ztesting.AssertEqual(t, "status code not match", http.StatusOK, ww.StatusCode())
		ztesting.AssertEqual(t, "written bytes not match", 4, ww.WrittenBytes())
	})
	t.Run("body intercept pattern1", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.Body = bytes.NewBuffer(nil)
		ww.Write([]byte("test"))
		ztesting.AssertEqual(t, "status code not match", -1, ww.StatusCode())
		ztesting.AssertEqual(t, "written bytes not match", -1, ww.WrittenBytes())
	})
	t.Run("body intercept pattern2", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.Body = bytes.NewBuffer(nil)
		ww.WriteHeader(http.StatusBadRequest)
		ww.Write([]byte("test"))
		ztesting.AssertEqual(t, "status code not match", http.StatusBadRequest, ww.StatusCode())
		ztesting.AssertEqual(t, "written bytes not match", 4, ww.WrittenBytes())
	})
}

type testFlushResponse struct {
	http.ResponseWriter
	called bool
}

func (w *testFlushResponse) Flush() {
	w.called = true
}

type testFlushErrorResponse struct {
	http.ResponseWriter
	called bool
}

func (w *testFlushErrorResponse) FlushError() error {
	w.called = true
	return nil
}

type testUnwrapResponse struct {
	http.ResponseWriter
}

func (w *testUnwrapResponse) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func TestResponseWrapper_Flush(t *testing.T) {
	t.Run("no flusher", func(t *testing.T) {
		w := &struct{ http.ResponseWriter }{
			ResponseWriter: httptest.NewRecorder(),
		}
		ww := zhttp.WrapResponseWriter(w)
		ww.Flush()
		ww.Flush()
	})
	t.Run("flush", func(t *testing.T) {
		w := &testFlushResponse{
			ResponseWriter: httptest.NewRecorder(),
		}
		ww := zhttp.WrapResponseWriter(w)
		ww.Flush()
		ww.Flush()
		ztesting.AssertEqual(t, "flush was not called", true, w.called)
	})
	t.Run("flush error", func(t *testing.T) {
		w := &testFlushErrorResponse{
			ResponseWriter: httptest.NewRecorder(),
		}
		ww := zhttp.WrapResponseWriter(w)
		ww.Flush()
		ww.Flush()
		ztesting.AssertEqual(t, "flush was not called", true, w.called)
	})
	t.Run("flush after unwrap", func(t *testing.T) {
		inner := &testFlushResponse{
			ResponseWriter: httptest.NewRecorder(),
		}
		w := &testUnwrapResponse{
			ResponseWriter: inner,
		}
		ww := zhttp.WrapResponseWriter(w)
		ww.Flush()
		ww.Flush()
		ztesting.AssertEqual(t, "flush was not called", true, inner.called)
	})
}
