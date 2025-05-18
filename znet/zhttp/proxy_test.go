package zhttp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestRewriteProxyURL(t *testing.T) {
	t.Parallel()
	parse := func(s string) *url.URL {
		u, err := url.Parse(s)
		if err != nil {
			panic(err)
		}
		return u
	}
	testCases := map[string]struct {
		dst, src *url.URL
		want     *url.URL
	}{
		"case1":            {&url.URL{}, parse("example.com"), &url.URL{Path: "example.com"}},
		"case2":            {&url.URL{}, parse("http://example.com"), &url.URL{Scheme: "http", Host: "example.com"}},
		"case3":            {&url.URL{}, parse("http://example.com/test"), &url.URL{Scheme: "http", Host: "example.com", Path: "/test"}},
		"rewrite schema":   {parse("https://example.com"), parse("http://example.com"), &url.URL{Scheme: "http", Host: "example.com"}},
		"rewrite host":     {parse("http://foo.com"), parse("http://bar.com"), &url.URL{Scheme: "http", Host: "bar.com"}},
		"rewrite path":     {parse("http://example.com/foo"), parse("http://example.com/bar"), &url.URL{Scheme: "http", Host: "example.com", Path: "/bar/foo"}},
		"rewrite query":    {parse("http://example.com?foo=alice"), parse("http://example.com?bar=bob"), &url.URL{Scheme: "http", Host: "example.com", RawQuery: "bar=bob&foo=alice"}},
		"rewrite fragment": {parse("http://example.com#foo"), parse("http://example.com#bar"), &url.URL{Scheme: "http", Host: "example.com", Fragment: "bar"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rewriteProxyURL(tc.dst, tc.src)
			matched := reflect.DeepEqual(tc.want, tc.dst)
			t.Logf("dst: %#v\n", tc.dst)
			t.Logf("want: %#v\n", tc.want)
			ztesting.AssertEqual(t, "url not match", true, matched)
		})
	}
}

func TestJoinWithByte(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		left, right string
		want        string
	}{
		"case01": {"", "", ""},
		"case02": {"foo", "", "foo"},
		"case03": {"/foo", "", "/foo"},
		"case04": {"foo/", "", "foo/"},
		"case05": {"/foo/", "", "/foo/"},
		"case06": {"", "bar", "bar"},
		"case07": {"foo", "bar", "foo/bar"},
		"case08": {"foo", "/bar", "foo/bar"},
		"case09": {"foo", "bar/", "foo/bar/"},
		"case10": {"foo", "/bar/", "foo/bar/"},
		"case11": {"/foo", "bar", "/foo/bar"},
		"case12": {"/foo", "/bar", "/foo/bar"},
		"case13": {"/foo", "bar/", "/foo/bar/"},
		"case14": {"/foo", "/bar/", "/foo/bar/"},
		"case15": {"foo/", "bar", "foo/bar"},
		"case16": {"foo/", "/bar", "foo/bar"},
		"case17": {"foo/", "bar/", "foo/bar/"},
		"case18": {"foo/", "/bar/", "foo/bar/"},
		"case19": {"/foo/", "bar", "/foo/bar"},
		"case20": {"/foo/", "/bar", "/foo/bar"},
		"case21": {"/foo/", "bar/", "/foo/bar/"},
		"case22": {"/foo/", "/bar/", "/foo/bar/"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := joinWithByte(tc.left, tc.right, '/')
			ztesting.AssertEqual(t, "join not match", tc.want, got)
		})
	}
}

type nopTransport struct {
	r *http.Request
}

func (t *nopTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.r = r
	return nil, errors.New("roundtrip error")
}

func TestNewProxy(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		targets []string
		err     error
	}{
		"single":   {targets: []string{"http://test.com"}, err: nil},
		"multiple": {targets: []string{"http://test1.com", "http://test2.com"}, err: nil},
		"error":    {targets: []string{"%http://test.com"}, err: &url.Error{Op: "parse", URL: "%http://test.com"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			proxy, err := NewProxy(tc.targets...)
			if tc.err != nil {
				ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
				return
			}
			ztesting.AssertEqualErr(t, "error should be nil", nil, err)
			nt := &nopTransport{}
			proxy.Transport = nt
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			w := httptest.NewRecorder()
			for i := range 2 * len(tc.targets) {
				target := tc.targets[i%len(tc.targets)]
				proxy.ServeHTTP(w, r)
				ztesting.AssertEqual(t, "url not match", target, nt.r.URL.String())
			}
		})
	}
}

type testTransport struct {
	req  *http.Request
	resp *http.Response
	err  error
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.req = req
	return t.resp, t.err
}

func TestProxy(t *testing.T) {
	t.Parallel()
	t.Run("simple request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		res := &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     http.Header{"Test": {"foo"}},
			Body:       io.NopCloser(strings.NewReader("resp body")),
		}
		tp := &testTransport{resp: res}
		proxy := &Proxy{
			Rewrite:   func(in, out *http.Request) {},
			Transport: tp,
		}
		resp := httptest.NewRecorder()
		proxy.ServeHTTP(resp, req)
		ztesting.AssertEqual(t, "method not match", http.MethodGet, tp.req.Method)
		ztesting.AssertEqual(t, "target not match", "http://test.com", tp.req.URL.String())
		ztesting.AssertEqual(t, "response status not match", http.StatusBadRequest, resp.Result().StatusCode)
		ztesting.AssertEqual(t, "response header not match", "foo", resp.Header().Get("Test"))
		ztesting.AssertEqual(t, "response body not match", "resp body", resp.Body.String())
	})
	t.Run("simple request with body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", strings.NewReader("req body"))
		res := &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     http.Header{"Test": {"foo"}},
			Body:       io.NopCloser(strings.NewReader("resp body")),
		}
		tp := &testTransport{resp: res}
		proxy := &Proxy{
			Rewrite:   func(in, out *http.Request) {},
			Transport: tp,
		}
		resp := httptest.NewRecorder()
		proxy.ServeHTTP(resp, req)
		body, _ := io.ReadAll(tp.req.Body)
		ztesting.AssertEqual(t, "method not match", http.MethodGet, tp.req.Method)
		ztesting.AssertEqual(t, "target not match", "http://test.com", tp.req.URL.String())
		ztesting.AssertEqual(t, "request body not match", "req body", string(body))
		ztesting.AssertEqual(t, "response status not match", http.StatusBadRequest, resp.Result().StatusCode)
		ztesting.AssertEqual(t, "response header not match", "foo", resp.Header().Get("Test"))
		ztesting.AssertEqual(t, "response body not match", "resp body", resp.Body.String())
	})
	t.Run("nil request header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		req.Header = nil
		res := &http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody}
		tp := &testTransport{resp: res}
		proxy := &Proxy{
			Rewrite:   func(in, out *http.Request) {},
			Transport: tp,
		}
		resp := httptest.NewRecorder()
		proxy.ServeHTTP(resp, req)
		t.Log(tp.req.Header)
		ztesting.AssertEqual(t, "header not match", true, reflect.DeepEqual(http.Header{"User-Agent": []string{""}}, tp.req.Header))
	})
	t.Run("request with trailer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		req.Header.Set("Te", "trailers")
		res := &http.Response{StatusCode: http.StatusBadRequest, Header: http.Header{}, Body: http.NoBody}
		tp := &testTransport{resp: res}
		proxy := &Proxy{
			Rewrite:   func(in, out *http.Request) {},
			Transport: tp,
		}
		resp := httptest.NewRecorder()
		proxy.ServeHTTP(resp, req)
		ztesting.AssertEqual(t, "trailers not found", "trailers", tp.req.Header.Get("Te"))
	})
	t.Run("response with trailer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		res := &http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody, Trailer: http.Header{"Test": {"foo"}}}
		proxy := &Proxy{
			Rewrite:   func(in, out *http.Request) {},
			Transport: &testTransport{resp: res},
		}
		rec := httptest.NewRecorder()
		resp := &testFlushResponse{ResponseWriter: rec}
		proxy.ServeHTTP(resp, req)
		ztesting.AssertEqual(t, "trailers not match", "foo", rec.Result().Trailer.Get("Test"))
		ztesting.AssertEqual(t, "response not flushed", true, resp.flushed)
	})
	t.Run("protocol switch", func(t *testing.T) {
		respRW := bufio.NewReadWriter(nil, bufio.NewWriter(&bytes.Buffer{}))
		conn := &testConn{Reader: strings.NewReader("foo")}
		rw := &testHijackResponse{ResponseWriter: httptest.NewRecorder(), conn: conn, rw: respRW}

		req := &http.Request{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}}
		var buf bytes.Buffer
		body := &readWriteCloser{Reader: strings.NewReader("bar"), Writer: &buf}
		res := &http.Response{
			StatusCode: http.StatusSwitchingProtocols,
			Header:     http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}},
			Body:       body,
		}

		proxy := &Proxy{
			Rewrite:   func(in, out *http.Request) {},
			Transport: &testTransport{resp: res},
		}
		proxy.ServeHTTP(rw, req)
		ztesting.AssertEqual(t, "copied content not match", "foo", buf.String())
		ztesting.AssertEqual(t, "copied content not match", "bar", conn.content.String())
	})
	t.Run("pre roundtrip error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		res := &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}
		var gotErr *HTTPError
		proxy := &Proxy{
			Rewrite:      func(in, out *http.Request) {},
			Transport:    &testTransport{resp: res},
			PreRoundTrip: func(in, out *http.Request) error { return errors.New("pre error") },
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err *HTTPError) { gotErr = err },
		}
		resp := httptest.NewRecorder()
		proxy.ServeHTTP(resp, req)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CausePreRoundTrip, Code: http.StatusInternalServerError}, gotErr)
	})
	t.Run("post roundtrip error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		res := &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}
		var gotErr *HTTPError
		proxy := &Proxy{
			Rewrite:       func(in, out *http.Request) {},
			Transport:     &testTransport{resp: res},
			PostRoundTrip: func(in *http.Request, out *http.Response) error { return errors.New("pre error") },
			ErrorHandler:  func(w http.ResponseWriter, r *http.Request, err *HTTPError) { gotErr = err },
		}
		resp := httptest.NewRecorder()
		proxy.ServeHTTP(resp, req)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CausePostRoundTrip, Code: http.StatusInternalServerError}, gotErr)
	})
	t.Run("upgrade error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		res := &http.Response{StatusCode: http.StatusSwitchingProtocols, Body: http.NoBody}
		var gotErr *HTTPError
		proxy := &Proxy{
			Rewrite:      func(in, out *http.Request) {},
			Transport:    &testTransport{resp: res},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err *HTTPError) { gotErr = err },
		}
		resp := httptest.NewRecorder()
		proxy.ServeHTTP(resp, req)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseUpgrade, Code: http.StatusInternalServerError}, gotErr)
	})
	t.Run("copy response body error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		body := io.NopCloser(ziotest.ErrReader(strings.NewReader("test response"), 4))
		res := &http.Response{StatusCode: http.StatusOK, Body: body}
		var gotErr *HTTPError
		proxy := &Proxy{
			Rewrite:      func(in, out *http.Request) {},
			Transport:    &testTransport{resp: res},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err *HTTPError) { gotErr = err },
		}
		resp := httptest.NewRecorder()
		proxy.ServeHTTP(resp, req)
		ztesting.AssertEqual(t, "copied body not match", "test", resp.Body.String())
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseCopyResponse, Code: -1}, gotErr)
	})
	t.Run("copy trailer error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
		res := &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Trailer: http.Header{"Test": {"foo"}}}
		var gotErr *HTTPError
		proxy := &Proxy{
			Rewrite:      func(in, out *http.Request) {},
			Transport:    &testTransport{resp: res},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err *HTTPError) { gotErr = err },
		}
		resp := struct{ http.ResponseWriter }{httptest.NewRecorder()} // Resp does not implement flusher.
		proxy.ServeHTTP(resp, req)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseFlushBody, Code: -1}, gotErr)
	})
}

type testConn struct {
	io.Reader
	net.Conn
	closeErr error

	closed  bool
	content bytes.Buffer
}

func (c *testConn) Read(b []byte) (n int, err error) {
	return c.Reader.Read(b)
}

func (c *testConn) Write(b []byte) (n int, err error) {
	return c.content.Write(b)
}

func (c *testConn) Close() error {
	c.closed = true
	return c.closeErr
}

type testHijackResponse struct {
	http.ResponseWriter
	conn      net.Conn
	rw        *bufio.ReadWriter
	hijackErr error
}

func (r *testHijackResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.conn, r.rw, r.hijackErr
}

type readWriteCloser struct {
	io.Reader
	io.Writer
}

func (rwc *readWriteCloser) Close() error {
	return nil
}

func TestHandleUpgradeResponse(t *testing.T) {
	t.Parallel()
	t.Run("upgrade success", func(t *testing.T) {
		respRW := bufio.NewReadWriter(nil, bufio.NewWriter(&bytes.Buffer{}))
		conn := &testConn{Reader: strings.NewReader("foo")}
		rw := &testHijackResponse{ResponseWriter: httptest.NewRecorder(), conn: conn, rw: respRW}

		req := &http.Request{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}}
		var buf bytes.Buffer
		body := &readWriteCloser{Reader: strings.NewReader("bar"), Writer: &buf}
		res := &http.Response{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}, Body: body}

		err := handleUpgradeResponse(rw, req, res)
		ztesting.AssertEqualErr(t, "error not match", (*HTTPError)(nil), err)
		ztesting.AssertEqual(t, "copied content not match", "foo", buf.String())
		ztesting.AssertEqual(t, "copied content not match", "bar", conn.content.String())
	})
	t.Run("upgrade not match", func(t *testing.T) {
		respRW := bufio.NewReadWriter(nil, bufio.NewWriter(&bytes.Buffer{}))
		conn := &testConn{Reader: strings.NewReader("foo")}
		rw := &testHijackResponse{ResponseWriter: httptest.NewRecorder(), conn: conn, rw: respRW}

		req := &http.Request{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test1"}}}
		var buf bytes.Buffer
		body := &readWriteCloser{Reader: strings.NewReader("bar"), Writer: &buf}
		res := &http.Response{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test2"}}, Body: body}

		err := handleUpgradeResponse(rw, req, res)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseUpgradeMismatch, Code: http.StatusBadRequest}, err)
		ztesting.AssertEqual(t, "copied content not match", "", buf.String())
		ztesting.AssertEqual(t, "copied content not match", "", conn.content.String())
	})
	t.Run("no connection header", func(t *testing.T) {
		respRW := bufio.NewReadWriter(nil, bufio.NewWriter(&bytes.Buffer{}))
		conn := &testConn{Reader: strings.NewReader("foo")}
		rw := &testHijackResponse{ResponseWriter: httptest.NewRecorder(), conn: conn, rw: respRW}

		req := &http.Request{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}}
		var buf bytes.Buffer
		body := &readWriteCloser{Reader: strings.NewReader("bar"), Writer: &buf}
		res := &http.Response{Header: http.Header{"Connection": {""}, "Upgrade": {"test2"}}, Body: body} // No connection header.

		err := handleUpgradeResponse(rw, req, res)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseUpgradeMismatch, Code: http.StatusBadRequest}, err)
		ztesting.AssertEqual(t, "copied content not match", "", buf.String())
		ztesting.AssertEqual(t, "copied content not match", "", conn.content.String())
	})
	t.Run("non ReadWriteCloser body", func(t *testing.T) {
		respRW := bufio.NewReadWriter(nil, bufio.NewWriter(&bytes.Buffer{}))
		conn := &testConn{Reader: strings.NewReader("foo")}
		rw := &testHijackResponse{ResponseWriter: httptest.NewRecorder(), conn: conn, rw: respRW}

		req := &http.Request{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}}
		var buf bytes.Buffer
		res := &http.Response{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}, Body: nil} // Body is nil.

		err := handleUpgradeResponse(rw, req, res)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseUpgrade, Code: http.StatusInternalServerError}, err)
		ztesting.AssertEqual(t, "copied content not match", "", buf.String())
		ztesting.AssertEqual(t, "copied content not match", "", conn.content.String())
	})
	t.Run("hijack failed", func(t *testing.T) {
		respRW := bufio.NewReadWriter(nil, bufio.NewWriter(&bytes.Buffer{}))
		conn := &testConn{Reader: strings.NewReader("foo")}
		hijackErr := errors.New("hijack failed")
		rw := &testHijackResponse{ResponseWriter: httptest.NewRecorder(), conn: conn, rw: respRW, hijackErr: hijackErr} // Hijack error.

		req := &http.Request{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}}
		var buf bytes.Buffer
		body := &readWriteCloser{Reader: strings.NewReader("bar"), Writer: &buf}
		res := &http.Response{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}, Body: body}

		err := handleUpgradeResponse(rw, req, res)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseHijack, Code: http.StatusInternalServerError}, err)
		ztesting.AssertEqual(t, "copied content not match", "", buf.String())
		ztesting.AssertEqual(t, "copied content not match", "", conn.content.String())
	})
	t.Run("flush error", func(t *testing.T) {
		ew := ziotest.ErrWriter(&bytes.Buffer{}, 0) // Error writer results in flush error.
		respRW := bufio.NewReadWriter(nil, bufio.NewWriter(ew))
		conn := &testConn{Reader: strings.NewReader("foo")}
		rw := &testHijackResponse{ResponseWriter: httptest.NewRecorder(), conn: conn, rw: respRW}

		req := &http.Request{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}}
		var buf bytes.Buffer
		body := &readWriteCloser{Reader: strings.NewReader("bar"), Writer: &buf}
		res := &http.Response{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}, Body: body}

		err := handleUpgradeResponse(rw, req, res)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseFlushBody, Code: -1}, err)
		ztesting.AssertEqual(t, "copied content not match", "", buf.String())
		ztesting.AssertEqual(t, "copied content not match", "", conn.content.String())
	})
	t.Run("response header write error", func(t *testing.T) {
		respRW := bufio.NewReadWriter(nil, bufio.NewWriterSize(ziotest.ErrWriter(&bytes.Buffer{}, 0), 1)) // Error writer.
		conn := &testConn{Reader: strings.NewReader("foo")}
		rw := &testHijackResponse{ResponseWriter: httptest.NewRecorder(), conn: conn, rw: respRW}

		req := &http.Request{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}}
		var buf bytes.Buffer
		body := &readWriteCloser{Reader: strings.NewReader("bar"), Writer: &buf}
		res := &http.Response{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}, Body: body}

		err := handleUpgradeResponse(rw, req, res)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseCopyResponse, Code: -1}, err)
		ztesting.AssertEqual(t, "copied content not match", "", buf.String())
		ztesting.AssertEqual(t, "copied content not match", "", conn.content.String())
	})
	t.Run("copy front to back error", func(t *testing.T) {
		respRW := bufio.NewReadWriter(nil, bufio.NewWriter(&bytes.Buffer{}))
		conn := &testConn{Reader: ziotest.ErrReader(strings.NewReader("foo"), 1)} // Error reader.
		rw := &testHijackResponse{ResponseWriter: httptest.NewRecorder(), conn: conn, rw: respRW}

		req := &http.Request{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}}
		var buf bytes.Buffer
		body := &readWriteCloser{Reader: strings.NewReader("bar"), Writer: &buf}
		res := &http.Response{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}, Body: body}

		err := handleUpgradeResponse(rw, req, res)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseCopyResponse, Code: -1}, err)
		ztesting.AssertEqual(t, "copied content not match", "f", buf.String())
		// ztesting.AssertEqual(t, "copied content not match", "bar", conn.content.String()) // Do not check.
	})
	t.Run("copy back to front error", func(t *testing.T) {
		respRW := bufio.NewReadWriter(nil, bufio.NewWriter(&bytes.Buffer{}))
		conn := &testConn{Reader: strings.NewReader("foo")}
		rw := &testHijackResponse{ResponseWriter: httptest.NewRecorder(), conn: conn, rw: respRW}

		req := &http.Request{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}}
		var buf bytes.Buffer
		body := &readWriteCloser{Reader: ziotest.ErrReader(strings.NewReader("bar"), 1), Writer: &buf} // Error reader.
		res := &http.Response{Header: http.Header{"Connection": {"Upgrade"}, "Upgrade": {"test"}}, Body: body}

		err := handleUpgradeResponse(rw, req, res)
		ztesting.AssertEqualErr(t, "error not match", &HTTPError{Cause: CauseCopyResponse, Code: -1}, err)
		// ztesting.AssertEqual(t, "copied content not match", "foo", buf.String()) // Do not check.
		ztesting.AssertEqual(t, "copied content not match", "b", conn.content.String())
	})
}

type testFlushResponse struct {
	http.ResponseWriter
	body    bytes.Buffer
	flushed bool
}

func (r *testFlushResponse) Write(p []byte) (int, error) {
	r.body.Write(p)
	return r.ResponseWriter.Write(p)
}

func (r *testFlushResponse) Flush() {
	r.flushed = true
}

func TestCopyResponseBody(t *testing.T) {
	t.Parallel()
	t.Run("stream response", func(t *testing.T) {
		rec := httptest.NewRecorder()
		w := &testFlushResponse{ResponseWriter: rec}
		res := &http.Response{ContentLength: -1, Body: io.NopCloser(strings.NewReader("test"))}
		copyResponseBody(w, res)
		ztesting.AssertEqual(t, "flush not called", true, w.flushed)
		ztesting.AssertEqual(t, "written body not match not", "test", w.body.String())
	})
	t.Run("sse", func(t *testing.T) {
		rec := httptest.NewRecorder()
		w := &testFlushResponse{ResponseWriter: rec}
		res := &http.Response{ContentLength: 0,
			Header: http.Header{"Content-Type": {"text/event-stream"}},
			Body:   io.NopCloser(strings.NewReader("test"))}
		copyResponseBody(w, res)
		ztesting.AssertEqual(t, "flush not called", true, w.flushed)
		ztesting.AssertEqual(t, "written body not match not", "test", w.body.String())
	})
	t.Run("chunked", func(t *testing.T) {
		rec := httptest.NewRecorder()
		w := &testFlushResponse{ResponseWriter: rec}
		res := &http.Response{ContentLength: 0,
			Header: http.Header{"Transfer-Encoding": {"chunked"}},
			Body:   io.NopCloser(strings.NewReader("test"))}
		copyResponseBody(w, res)
		ztesting.AssertEqual(t, "flush not called", true, w.flushed)
		ztesting.AssertEqual(t, "written body not match not", "test", w.body.String())
	})
	t.Run("default", func(t *testing.T) {
		rec := httptest.NewRecorder()
		w := &testFlushResponse{ResponseWriter: rec}
		res := &http.Response{Body: io.NopCloser(strings.NewReader("test"))}
		copyResponseBody(w, res)
		ztesting.AssertEqual(t, "flush should not be called", false, w.flushed)
		ztesting.AssertEqual(t, "written body not match not", "test", w.body.String())
	})
}

type testResponseWrapper struct {
	http.ResponseWriter
	body bytes.Buffer
}

func (r *testResponseWrapper) Write(p []byte) (int, error) {
	r.body.Write(p)
	return r.ResponseWriter.Write(p)
}

func (r *testResponseWrapper) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func TestWithImmediateFlushWriter(t *testing.T) {
	t.Parallel()
	t.Run("no flush", func(t *testing.T) {
		rec := httptest.NewRecorder()
		w := withImmediateFlushWriter(rec)
		w.Write([]byte("test"))
		ztesting.AssertEqual(t, "body not written", "test", rec.Body.String())
	})
	t.Run("unwrap flush", func(t *testing.T) {
		rec := httptest.NewRecorder()
		tr := &testResponseWrapper{ResponseWriter: rec}
		w := withImmediateFlushWriter(tr)
		w.Write([]byte("test"))
		ztesting.AssertEqual(t, "body not written", "test", tr.body.String())
		ztesting.AssertEqual(t, "body not written", "test", rec.Body.String())
	})
	t.Run("no flusher", func(t *testing.T) {
		tr := &testResponseWrapper{ResponseWriter: nil}
		w := withImmediateFlushWriter(tr)
		ztesting.AssertEqual(t, "writer not match", io.Writer(tr), w)
	})
}
