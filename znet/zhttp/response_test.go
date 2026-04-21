package zhttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aileron-projects/go/znet/zhttp"
	"github.com/aileron-projects/go/ztesting"
)

func TestWrapResponseWriter(t *testing.T) {
	t.Run("new wrapper", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.WriteHeader(http.StatusBadRequest)
		ztesting.AssertEqual(t, http.StatusBadRequest, w.Result().StatusCode)
	})
	t.Run("already wrapped", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww1 := zhttp.WrapResponseWriter(w)
		ww2 := zhttp.WrapResponseWriter(ww1)
		ztesting.AssertEqual(t, ww1, ww2)
		ww2.WriteHeader(http.StatusBadRequest)
		ztesting.AssertEqual(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestResponseWrapper(t *testing.T) {
	t.Run("not written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ztesting.AssertEqual(t, -1, ww.StatusCode())
		ztesting.AssertEqual(t, 0, ww.Written())
	})
	t.Run("header written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.Header().Set("Test", "value") // No affection.
		ztesting.AssertEqual(t, -1, ww.StatusCode())
		ztesting.AssertEqual(t, "value", ww.Header().Get("Test"))
		ztesting.AssertEqual(t, 0, ww.Written())
	})
	t.Run("status written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.WriteHeader(http.StatusBadRequest)
		ztesting.AssertEqual(t, http.StatusBadRequest, ww.StatusCode())
		ztesting.AssertEqual(t, 0, ww.Written())
	})
	t.Run("nil body written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.Write(nil)
		ztesting.AssertEqual(t, http.StatusOK, ww.StatusCode())
		ztesting.AssertEqual(t, 0, ww.Written())
	})
	t.Run("non-nil body written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		ww.Write([]byte("test"))
		ztesting.AssertEqual(t, http.StatusOK, ww.StatusCode())
		ztesting.AssertEqual(t, 4, ww.Written())
	})
	t.Run("hook on status written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		status := 0
		ww.StatusWritten = func(statusCode int) { status = statusCode }
		ww.WriteHeader(http.StatusBadRequest)
		ztesting.AssertEqual(t, http.StatusBadRequest, status)
	})
	t.Run("hook on body written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		status := 0
		ww.StatusWritten = func(statusCode int) { status = statusCode }
		ww.Write(nil)
		ztesting.AssertEqual(t, http.StatusOK, status)
	})
	t.Run("hijack error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ww := zhttp.WrapResponseWriter(w)
		_, _, err := ww.Hijack()
		ztesting.AssertEqual(t, true, err != nil)
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
		ztesting.AssertEqual(t, true, w.called)
	})
	t.Run("flush error", func(t *testing.T) {
		w := &testFlushErrorResponse{
			ResponseWriter: httptest.NewRecorder(),
		}
		ww := zhttp.WrapResponseWriter(w)
		ww.Flush()
		ww.Flush()
		ztesting.AssertEqual(t, true, w.called)
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
		ztesting.AssertEqual(t, true, inner.called)
	})
}
