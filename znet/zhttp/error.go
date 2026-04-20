package zhttp

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/zlog/zslog"
)

// ErrorHandler handles HTTP errors.
// It is intended to be used in middleware and handlers.
// Provided w, r and err should not be nil.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// NewErrorHandler returns a new error handler with the given logger.
// A default error handler checks if the err is [HTTPError] or not.
// If the err is type of [HTTPError], it returns an error response
// with status code specified in the [HTTPError.Code].
// Other errors results in an Internal Server Error.
var NewErrorHandler = func(lg zlog.Logger) ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if herr, ok := err.(*HTTPError); ok {
			if herr.Code >= 500 {
				attrs := []any{zslog.ErrorAttr(err), zslog.CallerAttr(1), zslog.FramesAttr(0)}
				lg.ErrorContext(r.Context(), "server-side error", attrs...)
			} else if lg.DebugEnabled(r.Context()) {
				attrs := []any{zslog.ErrorAttr(err), zslog.CallerAttr(1), zslog.FramesAttr(0)}
				lg.DebugContext(r.Context(), "http error handler received an error", attrs...)
			}
			http.Error(w, http.StatusText(herr.Code), herr.Code)
			return
		}
		attrs := []any{zslog.ErrorAttr(err), zslog.CallerAttr(1), zslog.FramesAttr(0)}
		lg.ErrorContext(r.Context(), "server-side error", attrs...)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// HTTPError is the HTTP error type.
type HTTPError struct {
	// Err is the internal error if any.
	Err error
	// Code is the preferred HTTP status code for the error.
	// In most case, status code should be respected but not required.
	// It must always be -1 or valid http status code such as [net/http.StatusOK].
	// StatusCode -1 indicates that response should not be written or
	// a response content is already written or the error is logging only.
	// StatusCode is compared in [HTTPError.Is].
	Code int
	// Cause is the error cause of the error.
	// It must always be non empty string.
	// Cause is compared in [HTTPError.Is].
	Cause string
	// Details is the additional information
	// of the error if any. It can be empty.
	// Cause is not compared in [HTTPError.Is].
	Detail string
}

func (e *HTTPError) Error() string {
	msg := e.Cause + " (Code:" + strconv.Itoa(e.Code) + ")."
	if e.Detail != "" {
		msg += " " + e.Detail
	}
	if e.Err != nil {
		msg += " [" + e.Err.Error() + "]"
	}
	return msg
}

// Unwrap returns the inner error if any.
func (e *HTTPError) Unwrap() error {
	return e.Err
}

// Is returns if this error is identical to the given error.
// The err is identical to the error when it has the type
// [HTTPError] and both [HTTPError.Cause] and [HTTPError.StatusCode]
// are the same.
func (e *HTTPError) Is(err error) bool {
	if err == nil || e == nil {
		return e == err
	}
	for err != nil {
		ee, ok := err.(*HTTPError)
		if ok {
			return e.Cause == ee.Cause && e.Code == ee.Code
		}
		err = errors.Unwrap(err)
	}
	return false
}
