package zhttp

import (
	"bufio"
	"net"
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
// make written status code and written body bytes accessible.
// Use [WrapResponseWriter] to wraps a response writer.
type ResponseWrapper struct {
	// StatusWritten is the callback function that will
	// be called when a status code was written to the
	// internal ResponseWriter.
	StatusWritten func(statusCode int)

	inner      http.ResponseWriter
	written    int64
	flush      func()
	flushError func() error

	statusWritten bool
	status        int
}

func (w *ResponseWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(w.inner).Hijack()
}

// StatusCode returns a HTTP status code written to the
// response writer. If status code has not been written,
// it returns -1.
func (w *ResponseWrapper) StatusCode() int {
	if !w.statusWritten {
		return -1
	}
	return w.status
}

// Written returns the number of bytes written to the internal
// response writer. If nothing written, it returns 0.
func (w *ResponseWrapper) Written() int64 {
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
	if sw := w.StatusWritten; sw != nil {
		sw(statusCode)
	}
	w.statusWritten = true
	w.status = statusCode
	w.inner.WriteHeader(statusCode)
}

// Write writes the data to the response writer.
// See [net/http.ResponseWriter.Write].
func (w *ResponseWrapper) Write(b []byte) (n int, err error) {
	if !w.statusWritten {
		w.WriteHeader(http.StatusOK)
	}
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
			if !w.statusWritten {
				w.WriteHeader(http.StatusOK)
			}
			w.flushError = t.FlushError
			return t.FlushError()
		case http.Flusher:
			if !w.statusWritten {
				w.WriteHeader(http.StatusOK)
			}
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
