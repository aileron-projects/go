package zhttp

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

var (
	// ErrCannotRewind is the error that means the body of [http.Request] cannot be rewound.
	// That means the request has a header of "Idempotency-Key" or "X-Idempotency-Key".
	// Or the [http.Request.GetBody] is nil.
	ErrCannotRewind = errors.New("zhttp: cannot rewind request body")
)

// SetupRewindBody makes the r retry-able by filling the r.GetBody field.
// It modifies the r.Body and r.GetBody if the request is rewindable.
// If the r has a header of "Idempotency-Key" or "X-Idempotency-Key",
// an [ErrCannotRewind] will be returned.
// Note that the entire body of r is read on memory.
// SetupRewindBody should be called before sending the first request.
//
// References:
//   - https://cs.opensource.google/go/go/+/refs/tags/go1.24.2:src/net/http/request.go;l=1534-1547
//   - https://cs.opensource.google/go/go/+/refs/tags/go1.24.2:src/net/http/transport.go;l=773-780
func SetupRewindBody(r *http.Request) error {
	if r.Body == nil || r.Body == http.NoBody {
		return nil // Nothing to rewind.
	}
	if len(r.Header["Idempotency-Key"]) > 0 || len(r.Header["X-Idempotency-Key"]) > 0 {
		return ErrCannotRewind
	}
	if r.GetBody == nil { // Make the body reusable.
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(body)), nil
		}
		r.Body, _ = r.GetBody()
	}
	return nil
}

// RewindBody rewinds the request body.
// It modifies the r.Body if the request is rewindable.
// If the r has a header of "Idempotency-Key" or "X-Idempotency-Key", or
// the r.GetBody is not available, an [ErrCannotRewind] will be returned.
// For making the request retry-able, call the [SetupRewindBody] before
// sending the first request.
//
// References:
//   - https://cs.opensource.google/go/go/+/refs/tags/go1.24.2:src/net/http/request.go;l=1534-1547
//   - https://cs.opensource.google/go/go/+/refs/tags/go1.24.2:src/net/http/transport.go;l=786-803
func RewindBody(r *http.Request) error {
	if r.Body == nil || r.Body == http.NoBody {
		return nil // Nothing to rewind.
	}
	if len(r.Header["Idempotency-Key"]) > 0 || len(r.Header["X-Idempotency-Key"]) > 0 {
		return ErrCannotRewind
	}
	if r.GetBody == nil { // Cannot rewind.
		return ErrCannotRewind
	}
	var err error
	r.Body, err = r.GetBody()
	return err
}

// ReadBody reads request body.
// Calling ReadBody fills r.GetBody when possible.
func ReadBody(r *http.Request) ([]byte, error) {
	if r.Body == nil || r.Body == http.NoBody {
		return nil, nil
	}
	if r.GetBody != nil {
		rc, err := r.GetBody()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(rc)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}
	r.Body, _ = r.GetBody()
	return body, nil
}
