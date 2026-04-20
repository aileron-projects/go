package zhttp_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aileron-projects/go/znet/zhttp"
	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestSetupRewindBody(t *testing.T) {
	t.Parallel()
	t.Run("nil body", func(t *testing.T) {
		r := &http.Request{
			Body: nil,
		}
		err := zhttp.SetupRewindBody(r)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("no body", func(t *testing.T) {
		r := &http.Request{
			Body: http.NoBody,
		}
		err := zhttp.SetupRewindBody(r)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("idempotent", func(t *testing.T) {
		r := &http.Request{
			Header: http.Header{"Idempotency-Key": []string{"test"}},
			Body:   io.NopCloser(strings.NewReader("body")),
		}
		err := zhttp.SetupRewindBody(r)
		ztesting.AssertEqualErr(t, zhttp.ErrCannotRewind, err)
	})
	t.Run("no get body", func(t *testing.T) {
		r := &http.Request{
			Body: io.NopCloser(strings.NewReader("body")),
		}
		err := zhttp.SetupRewindBody(r)
		ztesting.AssertEqualErr(t, nil, err)
		gb, _ := r.GetBody()
		body, _ := io.ReadAll(gb)
		ztesting.AssertEqual(t, "body", string(body))
	})
	t.Run("get body", func(t *testing.T) {
		b := io.NopCloser(strings.NewReader("body"))
		r := &http.Request{
			GetBody: func() (io.ReadCloser, error) { return nil, nil },
			Body:    b,
		}
		err := zhttp.SetupRewindBody(r)
		ztesting.AssertEqualErr(t, nil, err)
		ztesting.AssertEqual(t, b, r.Body)
	})
	t.Run("body read error", func(t *testing.T) {
		r := &http.Request{
			Body: io.NopCloser(ziotest.ErrReader(strings.NewReader("body"), 2)),
		}
		err := zhttp.SetupRewindBody(r)
		ztesting.AssertEqualErr(t, io.ErrClosedPipe, err)
	})
}

func TestRewindBody(t *testing.T) {
	t.Parallel()
	t.Run("nil body", func(t *testing.T) {
		r := &http.Request{
			Body: nil,
		}
		err := zhttp.RewindBody(r)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("no body", func(t *testing.T) {
		r := &http.Request{
			Body: http.NoBody,
		}
		err := zhttp.RewindBody(r)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("idempotent", func(t *testing.T) {
		r := &http.Request{
			Header: http.Header{"Idempotency-Key": []string{"test"}},
			Body:   io.NopCloser(strings.NewReader("body")),
		}
		err := zhttp.RewindBody(r)
		ztesting.AssertEqualErr(t, zhttp.ErrCannotRewind, err)
	})
	t.Run("no get body", func(t *testing.T) {
		r := &http.Request{
			Body: io.NopCloser(strings.NewReader("body")),
		}
		err := zhttp.RewindBody(r)
		ztesting.AssertEqualErr(t, zhttp.ErrCannotRewind, err)
	})
	t.Run("get body", func(t *testing.T) {
		b := io.NopCloser(strings.NewReader("body"))
		r := &http.Request{
			GetBody: func() (io.ReadCloser, error) { return b, nil },
			Body:    io.NopCloser(strings.NewReader("temp")),
		}
		err := zhttp.RewindBody(r)
		ztesting.AssertEqual(t, b, r.Body)
		ztesting.AssertEqualErr(t, nil, err)
	})
}

func TestReadBody(t *testing.T) {
	t.Parallel()
	t.Run("nil body", func(t *testing.T) {
		r := &http.Request{
			Body: nil,
		}
		body, err := zhttp.ReadBody(r)
		ztesting.AssertEqual(t, 0, len(body))
		ztesting.AssertEqual(t, "", string(body))
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("no body", func(t *testing.T) {
		r := &http.Request{
			Body: http.NoBody,
		}
		body, err := zhttp.ReadBody(r)
		ztesting.AssertEqual(t, 0, len(body))
		ztesting.AssertEqual(t, "", string(body))
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("no get body", func(t *testing.T) {
		r := &http.Request{
			Body: io.NopCloser(strings.NewReader("body")),
		}
		body, err := zhttp.ReadBody(r)
		ztesting.AssertEqual(t, 4, len(body))
		ztesting.AssertEqual(t, "body", string(body))
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("get body", func(t *testing.T) {
		r := &http.Request{
			GetBody: func() (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("body")), nil
			},
			Body: io.NopCloser(strings.NewReader("temp")),
		}
		body, err := zhttp.ReadBody(r)
		ztesting.AssertEqual(t, 4, len(body))
		ztesting.AssertEqual(t, "body", string(body))
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("read multiple times", func(t *testing.T) {
		r := &http.Request{
			Body: io.NopCloser(strings.NewReader("body")),
		}
		body, err := zhttp.ReadBody(r)
		ztesting.AssertEqual(t, 4, len(body))
		ztesting.AssertEqual(t, "body", string(body))
		ztesting.AssertEqualErr(t, nil, err)
		body, err = zhttp.ReadBody(r)
		ztesting.AssertEqual(t, 4, len(body))
		ztesting.AssertEqual(t, "body", string(body))
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("get body error", func(t *testing.T) {
		r := &http.Request{
			GetBody: func() (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("body")), io.ErrUnexpectedEOF
			},
			Body: io.NopCloser(strings.NewReader("temp")),
		}
		body, err := zhttp.ReadBody(r)
		ztesting.AssertEqual(t, 0, len(body))
		ztesting.AssertEqual(t, "", string(body))
		ztesting.AssertEqualErr(t, io.ErrUnexpectedEOF, err)
	})
	t.Run("read body error", func(t *testing.T) {
		r := &http.Request{
			Body: io.NopCloser(ziotest.ErrReader(strings.NewReader("body"), 4)),
		}
		body, err := zhttp.ReadBody(r)
		ztesting.AssertEqual(t, 0, len(body))
		ztesting.AssertEqual(t, "", string(body))
		ztesting.AssertEqualErr(t, io.ErrClosedPipe, err)
	})
}
