package zhttp

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/http/httpguts"
)

const (
	CausePreRoundTrip    = "znet/zhttp: PreRoundTrip hook returned an error"
	CausePostRoundTrip   = "znet/zhttp: PostRoundTrip hook returned an error"
	CauseTransport       = "znet/zhttp: proxy transport returned an error"
	CauseUpgrade         = "znet/zhttp: handling upgrade response failed"
	CauseCopyResponse    = "znet/zhttp: copying response body to client failed. may be client canceled the request"
	CauseFlushBody       = "znet/zhttp: flushing response body failed"
	CauseUpgradeMismatch = "znet/zhttp: upgrade protocol mismatch"
	CauseUpgradeSwitch   = "znet/zhttp: 101 switching protocol failed"
	CauseHijack          = "znet/zhttp: hijacking response writer failed while switching protocol"
)

// NewProxy returns a new instance of [Proxy] with the given
// proxy targets. Targets are selected by round-robin algorithm.
func NewProxy(targets ...string) (*Proxy, error) {
	ts := make([]*url.URL, 0, len(targets))
	for _, t := range targets {
		u, err := url.Parse(t)
		if err != nil {
			return nil, err
		}
		ts = append(ts, u)
	}
	var mu sync.Mutex
	var index int
	return &Proxy{
		Rewrite: func(in, out *http.Request) {
			mu.Lock()
			defer mu.Unlock()
			if index >= len(ts) {
				index = 0
			}
			target := ts[index]
			index++
			rewriteProxyURL(out.URL, target)
			SetForwardedHeaders(in, out.Header)
		},
	}, nil
}

func rewriteProxyURL(dst, src *url.URL) {
	dst.Scheme = src.Scheme
	dst.User = cmp.Or(src.User, dst.User) // Replace if non-nil.
	dst.Host = src.Host
	dst.Path = joinWithByte(src.Path, dst.Path, '/')
	dst.RawPath = joinWithByte(src.RawPath, dst.RawPath, '/')
	dst.RawQuery = joinWithByte(src.RawQuery, dst.RawQuery, '&')
	if src.Fragment != "" || src.RawFragment != "" {
		dst.Fragment = src.Fragment
		dst.RawFragment = src.RawFragment
	}
}

func joinWithByte(left, right string, b byte) string {
	if left == "" {
		return right
	}
	if right == "" {
		return left
	}
	leftHasByte := left[len(left)-1] == b
	rightHasByte := right[0] == b
	if leftHasByte && rightHasByte {
		return left + right[1:]
	}
	if !leftHasByte && !rightHasByte {
		return left + string(b) + right
	}
	return left + right
}

// Proxy is the http proxy.
type Proxy struct {
	// Rewrite modifies proxy request.
	// Client request if provided with in and the proxy request is provided with out.
	// In must not be modified and at least out.URL must be configured for proxy target.
	// Out request is created with in.Clone with the context of in.Context.
	// out.Host is always empty and must not be filled.
	// Hop-by-hop headers are removed from out before Rewrite is called.
	// See also [net/http.Request.Host].
	// Rewrite must not be nil, otherwise panics.
	Rewrite func(in, out *http.Request)

	// The transport used to perform proxy requests.
	// If nil, [net/http.DefaultTransport] is used.
	Transport http.RoundTripper

	// PreRoundTrip is called just before roundtrip.
	// Client request is provided as in and proxy request as out.
	// In and out requests should not be modified.
	// PreRoundTrip can be used for logging or hooking etc.
	// If PreRoundTrip returned non-nil error, it is passed to
	// the error handler.
	PreRoundTrip func(in, out *http.Request) error
	// PostRoundTrip is called just before roundtrip.
	// Proxy request and response are provided as in and out each.
	// PostRoundTrip can be used as response modifier.
	// If PostRoundTrip returned non-nil error, it is passed to
	// the error handler.
	PostRoundTrip func(in *http.Request, out *http.Response) error

	// ErrorHandler is the optional error handler.
	// If non-nil, any errors occurred while proxying is given.
	// If nil, a default error handler is used.
	// Given w and r is the same as tha value passed to [Proxy.ServeHTTP].
	ErrorHandler ErrorHandler[*HTTPError]
}

func (p *Proxy) handleError(w http.ResponseWriter, r *http.Request, err *HTTPError) {
	if eh := p.ErrorHandler; eh != nil {
		eh(w, r, err)
		return
	}
	if err.Code > 0 && !errors.Is(err, context.Canceled) {
		w.WriteHeader(err.Code)
		_, _ = w.Write([]byte(http.StatusText(err.Code)))
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	outReq := r.Clone(r.Context())
	outReq.Host = ""
	if r.ContentLength == 0 {
		outReq.Body = nil // Issue 16036: nil Body for http.Transport retries
	}
	if outReq.Body != nil {
		defer outReq.Body.Close()
	}
	if outReq.Header == nil {
		outReq.Header = make(http.Header, 0)
	}
	if _, ok := outReq.Header["User-Agent"]; !ok {
		outReq.Header.Set("User-Agent", "") // Don't send the default Go User-Agent.
	}

	RemoveHopByHopHeaders(outReq.Header)
	p.Rewrite(r, outReq)

	if httpguts.HeaderValuesContainsToken(r.Header["Te"], "trailers") {
		// The TE request header specifies the transfer encodings the user agent is willing to accept.
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/TE
		outReq.Header.Set("Te", "trailers")
	}
	if httpguts.HeaderValuesContainsToken(r.Header["Connection"], "Upgrade") {
		// The Upgrade header can be used to upgrade an already-established client/server connection to a different protocol.
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Upgrade
		outReq.Header.Set("Connection", "Upgrade")
		outReq.Header.Set("Upgrade", r.Header.Get("Upgrade"))
	}

	if preRT := p.PreRoundTrip; preRT != nil {
		if err := preRT(r, outReq); err != nil {
			p.handleError(w, r, &HTTPError{Code: http.StatusInternalServerError, Cause: CausePreRoundTrip})
			return
		}
	}

	transport := cmp.Or(p.Transport, http.DefaultTransport)
	outRes, err := transport.RoundTrip(outReq)
	if err != nil {
		p.handleError(w, r, &HTTPError{Err: err, Code: http.StatusBadGateway, Cause: CauseTransport})
		return
	}
	defer outRes.Body.Close()

	if postRT := p.PostRoundTrip; postRT != nil {
		if err := postRT(r, outRes); err != nil {
			p.handleError(w, r, &HTTPError{Code: http.StatusInternalServerError, Cause: CausePostRoundTrip})
			return
		}
	}

	// Deal with 101 Switching Protocols responses: (WebSocket, h2c, etc)
	if outRes.StatusCode == http.StatusSwitchingProtocols {
		if err := handleUpgradeResponse(w, outReq, outRes); err != nil {
			p.handleError(w, r, err)
		}
		return
	}

	// Copy response header.
	RemoveHopByHopHeaders(outRes.Header)
	CopyHeaders(w.Header(), outRes.Header)
	if n := len(outRes.Trailer); n > 0 {
		announcedTrailerKeys := make([]string, 0, n)
		for k := range outRes.Trailer {
			announcedTrailerKeys = append(announcedTrailerKeys, k)
		}
		w.Header().Add("Trailer", strings.Join(announcedTrailerKeys, ", "))
	}
	// Copy response status code.
	w.WriteHeader(outRes.StatusCode)
	// Copy response body.
	if err := copyResponseBody(w, outRes); err != nil {
		p.handleError(w, r, &HTTPError{Err: err, Code: -1, Cause: CauseCopyResponse}) // Client canceled?
		return
	}

	outRes.Body.Close() // Close before populating trailers.

	if len(outRes.Trailer) > 0 {
		// Force chunking if we saw a response trailer.
		if err := http.NewResponseController(w).Flush(); err != nil {
			p.handleError(w, r, &HTTPError{Err: err, Code: -1, Cause: CauseFlushBody})
			return
		}
		CopyTrailers(w.Header(), outRes.Trailer)
	}
}

// handleUpgradeResponse handles protocol upgrade.
// This method is called when [net/http.StatusSwitchingProtocols] was detected.
// See also handleUpgradeResponse function in
// https://go.dev/src/net/http/httputil/reverseproxy.go
func handleUpgradeResponse(rw http.ResponseWriter, req *http.Request, res *http.Response) *HTTPError {
	reqUpType := upgradeType(req.Header)
	resUpType := upgradeType(res.Header)
	if len(reqUpType) != len(resUpType) || !strings.EqualFold(reqUpType, resUpType) {
		return &HTTPError{
			Code:   http.StatusBadRequest,
			Cause:  CauseUpgradeMismatch,
			Detail: fmt.Sprintf("backend tried to switch protocol %q when %q was requested", resUpType, reqUpType),
		}
	}

	backConn, ok := res.Body.(io.ReadWriteCloser)
	if !ok {
		return &HTTPError{
			Code:   http.StatusInternalServerError,
			Cause:  CauseUpgrade,
			Detail: "response body is not writable (" + fmt.Sprintf("%T", res.Body) + ")",
		}
	}
	defer backConn.Close() // Ensure close.

	conn, brw, err := http.NewResponseController(rw).Hijack()
	if err != nil {
		return &HTTPError{
			Err:    err,
			Code:   http.StatusInternalServerError,
			Cause:  CauseHijack,
			Detail: "hijacking ResponseWriter failed (" + fmt.Sprintf("%T", rw) + ")",
		}
	}
	defer conn.Close() // Ensure close.

	CopyHeaders(rw.Header(), res.Header)

	resp := res               // Shallow copy response to generate a response.
	resp.Header = rw.Header() // Headers to be responded.
	resp.Body = nil           // Avoid writing body. We copy using conn and backConn.
	if err := resp.Write(brw); err != nil {
		return &HTTPError{Err: err, Code: -1, Cause: CauseCopyResponse}
	}
	if err := brw.Flush(); err != nil {
		return &HTTPError{
			Err:    err,
			Code:   -1,
			Cause:  CauseFlushBody,
			Detail: "try to flush hijack response writer",
		}
	}

	errChan := make(chan error, 1)
	go copyBuf(conn, backConn, errChan)
	go copyBuf(backConn, conn, errChan)
	if err = <-errChan; err != nil {
		return &HTTPError{Err: err, Code: -1, Cause: CauseCopyResponse}
	}
	if err = <-errChan; err != nil {
		return &HTTPError{Err: err, Code: -1, Cause: CauseCopyResponse}
	}
	return nil
}

// upgradeType returns the string of protocol to upgrade.
// "Upgrade" header must be present at most 1 in the given h.
// The header key in the given h must be canonicalized.
// i.e. "Upgrade" not "upgrade", "Connection" not "connection".
//
// Reference
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Protocol_upgrade_mechanism
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Upgrade
func upgradeType(h http.Header) string {
	if httpguts.HeaderValuesContainsToken(h["Connection"], "Upgrade") {
		return h.Get("Upgrade")
	}
	return ""
}

// pool is the buffer pool.
var pool = sync.Pool{
	New: func() any {
		buf := make([]byte, 1<<14) // 16kiB
		return &buf
	},
}

// copy copies from src to dst.
// An error, which may be nil or non-nil is sent to the errChan.
func copyBuf(dst io.Writer, src io.Reader, errChan chan<- error) {
	buf := *pool.Get().(*[]byte)
	defer pool.Put(&buf)
	_, err := io.CopyBuffer(dst, src, buf)
	errChan <- err
}

// copyResponseBody copies proxy response body.
// w should be the frontend response writer and the res
// should be the response of proxy request.
func copyResponseBody(w http.ResponseWriter, res *http.Response) error {
	var dst io.Writer
	mt, _, _ := mime.ParseMediaType(res.Header.Get("Content-Type"))
	switch {
	case res.ContentLength < 0:
		dst = withImmediateFlushWriter(w) // Any streaming type response.
	case mt == "text/event-stream":
		dst = withImmediateFlushWriter(w) // Server Sent Events (SSE)
	case httpguts.HeaderValuesContainsToken(res.Header["Transfer-Encoding"], "chunked"):
		dst = withImmediateFlushWriter(w) // Chunked response.
	default:
		dst = w
	}
	buf := *pool.Get().(*[]byte)
	defer pool.Put(&buf)
	_, err := io.CopyBuffer(dst, res.Body, buf)
	return err
}

// withImmediateFlushWriter wraps the [net/http.ResponseWriter] with
// immediateFlushWriter if the rw implements [net/http.Flusher] interface.
// If the rw does not implements [net/http.Flusher], it returns rw itself.
// withImmediateFlushWriter try to find first flusher implementation
// by recursively unwrapping the rw with interface{ Unwrap() http.ResponseWriter
func withImmediateFlushWriter(rw http.ResponseWriter) io.Writer {
	fw := rw
	for {
		if flusher, ok := fw.(http.Flusher); ok {
			return &immediateFlushWriter{
				inner:   rw, // Use rw, not inner.
				flusher: flusher,
			}
		}
		if uw, ok := fw.(interface{ Unwrap() http.ResponseWriter }); ok {
			fw = uw.Unwrap()
			continue
		}
		return rw
	}
}

// immediateFlushWriter flushes immediately after Write called.
// The inner writer and flusher must not be nil.
type immediateFlushWriter struct {
	inner   io.Writer
	flusher http.Flusher
}

func (f *immediateFlushWriter) Write(p []byte) (n int, err error) {
	defer f.flusher.Flush()
	return f.inner.Write(p)
}
