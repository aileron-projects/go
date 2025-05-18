package zhttp

import (
	"cmp"
	"io"
	"net/http"
)

// WrapResponseWriter wraps the w with [ResponseWrapper].
func WrapResponseWriter(w http.ResponseWriter) *ResponseWrapper {
	if ww, ok := w.(*ResponseWrapper); ok {
		return ww // Already wrapped.
	}
	return &ResponseWrapper{
		inner: w,
	}
}

// ResponseWrapper wraps the [net/http.ResponseWriter] and
// make users accessible to the written status code
// and written body bytes.
// Use [WrapResponseWriter] to wraps a response writer with this.
type ResponseWrapper struct {
	// Body is an io writer that intercept data from
	// writing to the internal response writer.
	// Unlike [net/http.ResponseWriter.Write], writing to the Body
	// does not implicitly write [net/http.StatusOK].
	// If the Body is not nil, written bytes
	// to the [ResponseWrapper.Write] is written to the
	// Body instead of the inner ResponseWriter.
	Body io.Writer
	// inner is the wrapped response writer.
	// inner must not be nil.
	inner         http.ResponseWriter
	written       int64 // written bytes.
	status        int   // status code.
	statusWritten bool  // flag if a status code was
	flush         func()
	flushError    func() error
}

// Unwrap returns the internal response writer.
func (w *ResponseWrapper) Unwrap() http.ResponseWriter {
	return w.inner
}

// StatusCode returns a HTTP status code written to the
// response writer. If a status code has not been
// written yet, -1 is returned.
func (w *ResponseWrapper) StatusCode() int {
	if !w.statusWritten {
		return -1
	}
	return cmp.Or(w.status, http.StatusOK) // Default is 200 as the standard http package goes.
}

// WrittenBytes returns the number of bytes written
// to the response writer. It returns -1 when nothing was written.
// It returns -1 when the w.Body is not nil and no status code was written.
func (w *ResponseWrapper) WrittenBytes() int64 {
	if !w.statusWritten {
		return -1
	}
	return w.written
}

// Header returns http header.
// See [net/http.ResponseWriter.Header].
func (w *ResponseWrapper) Header() http.Header {
	return w.inner.Header()
}

// WriteHeader writes http status code.
// See [net/http.ResponseWriter.WriteHeader].
func (w *ResponseWrapper) WriteHeader(statusCode int) {
	if !w.statusWritten {
		w.status = statusCode
		w.statusWritten = true
	}
	w.inner.WriteHeader(statusCode)
}

// Write writes the data to the response writer.
// See [net/http.ResponseWriter.Write].
func (w *ResponseWrapper) Write(b []byte) (n int, err error) {
	if w.Body != nil {
		n, err = w.Body.Write(b)
		w.written += int64(n)
		return n, err
	}
	w.statusWritten = true // A [http.StatusOK] may be written.
	n, err = w.inner.Write(b)
	w.written += int64(n)
	return n, err
}

// Flush calls Flush() method of the internal response writers.
// See also the comments on the [net/http.NewResponseController].
func (w *ResponseWrapper) Flush() {
	_ = w.FlushError()
}

// FlushError calls FlushError() method of the internal response writers.
// It returns [net/http.ErrNotSupported] when the feature is not available.
// See also the comments on the [net/http.NewResponseController].
func (w *ResponseWrapper) FlushError() error {
	if w.flushError != nil {
		return w.flushError()
	}
	if w.flush != nil {
		w.flush()
		return nil
	}
	ww := w.inner
	for {
		switch t := ww.(type) {
		case interface{ FlushError() error }:
			w.flushError = t.FlushError
			return t.FlushError()
		case http.Flusher:
			w.flush = t.Flush
			t.Flush()
			return nil
		case interface{ Unwrap() http.ResponseWriter }:
			ww = t.Unwrap()
			continue
		}
		return nil // No flusher.
	}
}
