package zhttp_test

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zlog/zslog"
	"github.com/aileron-projects/go/znet/zhttp"
	"github.com/aileron-projects/go/ztesting"
)

func TestNewErrorHandler(t *testing.T) {
	t.Parallel()
	t.Run("server-side error", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
		eh := zhttp.NewErrorHandler(lg)
		r := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		w := httptest.NewRecorder()
		eh(w, r, &zhttp.HTTPError{Code: http.StatusBadGateway})
		ztesting.AssertEqual(t, true, strings.Contains(buf.String(), `level=ERROR`))
		ztesting.AssertEqual(t, http.StatusBadGateway, w.Result().StatusCode)
	})
	t.Run("client-side error with info logger", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
		eh := zhttp.NewErrorHandler(lg)
		r := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		w := httptest.NewRecorder()
		eh(w, r, &zhttp.HTTPError{Code: http.StatusBadRequest})
		ztesting.AssertEqual(t, "", buf.String())
		ztesting.AssertEqual(t, http.StatusBadRequest, w.Result().StatusCode)
	})
	t.Run("client-side error with debug logger", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		eh := zhttp.NewErrorHandler(lg)
		r := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		w := httptest.NewRecorder()
		eh(w, r, &zhttp.HTTPError{Code: http.StatusBadRequest})
		ztesting.AssertEqual(t, true, strings.Contains(buf.String(), `level=DEBUG`))
		ztesting.AssertEqual(t, http.StatusBadRequest, w.Result().StatusCode)
	})
	t.Run("non http error", func(t *testing.T) {
		var buf bytes.Buffer
		lg := zslog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
		eh := zhttp.NewErrorHandler(lg)
		r := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		w := httptest.NewRecorder()
		eh(w, r, io.EOF)
		ztesting.AssertEqual(t, true, strings.Contains(buf.String(), `level=ERROR`))
		ztesting.AssertEqual(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestHTTPError_Error(t *testing.T) {
	t.Parallel()
	t.Run("nil inner error", func(t *testing.T) {
		err := &zhttp.HTTPError{
			Err:    nil,
			Code:   http.StatusOK,
			Cause:  "test cause",
			Detail: "test detail",
		}
		got := err.Error()
		want := "test cause (Code:200). test detail"
		ztesting.AssertEqual(t, want, got)
	})
	t.Run("empty detail", func(t *testing.T) {
		err := &zhttp.HTTPError{
			Err:    nil,
			Code:   http.StatusOK,
			Cause:  "test cause",
			Detail: "",
		}
		got := err.Error()
		want := "test cause (Code:200)."
		ztesting.AssertEqual(t, want, got)
	})
	t.Run("non-nil inner error", func(t *testing.T) {
		err := &zhttp.HTTPError{
			Err:   io.EOF,
			Code:  http.StatusOK,
			Cause: "test cause",
		}
		got := err.Error()
		want := "test cause (Code:200). [EOF]"
		ztesting.AssertEqual(t, want, got)
	})
}

func TestHTTPError_Is(t *testing.T) {
	t.Parallel()
	t.Run("same error", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		ztesting.AssertEqual(t, true, err1.Is(err2))
	})
	t.Run("same after unwrap", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err3 := fmt.Errorf("outer error [%w]", err2)
		ztesting.AssertEqual(t, true, err1.Is(err3))
	})
	t.Run("not match nil", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		ztesting.AssertEqual(t, false, err1.Is(nil))
	})
	t.Run("not match after unwrap", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := fmt.Errorf("outer error [%w]", io.EOF)
		ztesting.AssertEqual(t, false, err1.Is(err2))
	})
	t.Run("unwrap error", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo", Err: io.EOF}
		ztesting.AssertEqual(t, false, err1.Is(err2.Unwrap()))
	})
	t.Run("status not match", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := &zhttp.HTTPError{Code: http.StatusBadGateway, Cause: "foo"}
		ztesting.AssertEqual(t, false, err1.Is(err2))
	})
	t.Run("cause not match", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "bar"}
		ztesting.AssertEqual(t, false, err1.Is(err2))
	})
}
