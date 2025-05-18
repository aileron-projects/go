package zhttp_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/aileron-projects/go/znet/zhttp"
	"github.com/aileron-projects/go/ztesting"
)

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
		ztesting.AssertEqual(t, "error message not match", want, got)
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
		ztesting.AssertEqual(t, "error message not match", want, got)
	})
	t.Run("non-nil inner error", func(t *testing.T) {
		err := &zhttp.HTTPError{
			Err:   io.EOF,
			Code:  http.StatusOK,
			Cause: "test cause",
		}
		got := err.Error()
		want := "test cause (Code:200). [EOF]"
		ztesting.AssertEqual(t, "error message not match", want, got)
	})
}

func TestHTTPError_Is(t *testing.T) {
	t.Parallel()
	t.Run("same error", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		ztesting.AssertEqual(t, "errors not match", true, err1.Is(err2))
	})
	t.Run("same after unwrap", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err3 := fmt.Errorf("outer error [%w]", err2)
		ztesting.AssertEqual(t, "errors not match", true, err1.Is(err3))
	})
	t.Run("not match nil", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		ztesting.AssertEqual(t, "errors unexpectedly matched", false, err1.Is(nil))
	})
	t.Run("not match after unwrap", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := fmt.Errorf("outer error [%w]", io.EOF)
		ztesting.AssertEqual(t, "errors unexpectedly matched", false, err1.Is(err2))
	})
	t.Run("unwrap error", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo", Err: io.EOF}
		ztesting.AssertEqual(t, "errors unexpectedly matched", false, err1.Is(err2.Unwrap()))
	})
	t.Run("status not match", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := &zhttp.HTTPError{Code: http.StatusBadGateway, Cause: "foo"}
		ztesting.AssertEqual(t, "errors unexpectedly matched", false, err1.Is(err2))
	})
	t.Run("cause not match", func(t *testing.T) {
		err1 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "foo"}
		err2 := &zhttp.HTTPError{Code: http.StatusOK, Cause: "bar"}
		ztesting.AssertEqual(t, "errors unexpectedly matched", false, err1.Is(err2))
	})
}
