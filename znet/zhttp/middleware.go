package zhttp

import (
	"net/http"
	"slices"
)

var (
	_ ServerMiddleware = &ServerMiddlewareChain{}
	_ ClientMiddleware = &ClientMiddlewareChain{}
)

// HandlerFunc type is an adapter to allow the
// use of ordinary function as HTTP handlers.
// If f is a function with the appropriate signature,
// HandlerFunc(f) is [http.Handler] that calls f.
// HandlerFunc is the same as [http.HandlerFunc].
//
// Example:
//
//	var h http.Handler
//	h = HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte("ok"))
//	})
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

// RoundTripperFunc type is an adapter to allow the
// use of ordinary function as HTTP round tripper.
// If f is a function with the appropriate signature,
// RoundTripperFunc(f) is [http.RoundTripper] that calls f.
//
// Example:
//
//	var r http.RoundTripper
//	r = RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
//		// Some process.
//		return http.DefaultTransport.RoundTrip(r)
//	})
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// ServerMiddlewareFunc type is an adapter to allow the
// use of ordinary function as server middleware.
// If f is a function with the appropriate signature,
// ServerMiddlewareFunc(f) is [ServerMiddleware] that calls f.
//
// Example:
//
//	var m ServerMiddleware
//	m = ServerMiddlewareFunc(func(next http.Handler) http.Handler {
//		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			// Some process.
//			next.ServeHTTP(w, r)
//		})
//	})
type ServerMiddlewareFunc func(next http.Handler) http.Handler

func (f ServerMiddlewareFunc) ServerMiddleware(next http.Handler) http.Handler {
	return f(next)
}

// ClientMiddlewareFunc type is an adapter to allow the
// use of ordinary function as client middleware.
// If f is a function with the appropriate signature,
// ClientMiddlewareFunc(f) is [ClientMiddleware] that calls f.
//
// Example:
//
//	var m ClientMiddleware
//	m = ClientMiddlewareFunc(func(next http.RoundTripper) http.RoundTripper {
//		return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
//			// Some process.
//			return next.RoundTrip(r)
//		})
//	})
type ClientMiddlewareFunc func(next http.RoundTripper) http.RoundTripper

func (f ClientMiddlewareFunc) ClientMiddleware(next http.RoundTripper) http.RoundTripper {
	return f(next)
}

// ServerMiddleware is the server side middleware.
type ServerMiddleware interface {
	ServerMiddleware(next http.Handler) http.Handler
}

// ClientMiddleware is the client side middleware.
type ClientMiddleware interface {
	ClientMiddleware(next http.RoundTripper) http.RoundTripper
}

// NewHandler returns a http handler with server middleware.
// NewHandler is short for ServerMiddlewareChain(ms).Handler(h).
func NewHandler(h http.Handler, ms ...ServerMiddleware) http.Handler {
	return ServerMiddlewareChain(ms).Handler(h)
}

// ServerMiddlewareChain is the middleware chain of [ServerMiddleware].
// ServerMiddlewareChain implements the [ServerMiddleware] interface itself.
type ServerMiddlewareChain []ServerMiddleware

// ServerMiddleware is the implementation of [ServerMiddleware.ServerMiddleware].
func (c ServerMiddlewareChain) ServerMiddleware(next http.Handler) http.Handler {
	return c.Handler(next)
}

// Handler returns a handler with middleware.
// Additional middleware is applied at the end of the chain
// but they are not registered to this middleware chain.
func (c ServerMiddlewareChain) Handler(h http.Handler, ms ...ServerMiddleware) http.Handler {
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i].ServerMiddleware(h)
	}
	for i := len(c) - 1; i >= 0; i-- {
		h = c[i].ServerMiddleware(h)
	}
	return h
}

// Add adds given middleware to the chain.
// Middleware are added at the end of the chain keeping the order.
// Note that nil middleware will cause an panic.
func (c *ServerMiddlewareChain) Add(ms ...ServerMiddleware) {
	*c = append(*c, ms...)
}

// Insert inserts the given middleware at the position of index.
// If index<=0, ms are added at the beginning of the chain.
// If index>=len(c),  ms are added at the end of the chain.
// Note that nil middleware will cause an panic.
//
//	      ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	chain │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	      ├──────┼──────┼──────┼──────┼─────┼────────┤
//	      ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index 0      1      2      3      4    n-1       n
func (c *ServerMiddlewareChain) Insert(index int, ms ...ServerMiddleware) {
	index = min(max(index, 0), len(*c))
	*c = slices.Insert(*c, index, ms...)
}

// InsertAll inserts the ms between each middleware.
// That is, the given ms is inserted at the index from 0 to n.
// Inserting to an empty chain will be ms itself.
// Note that nil middleware will cause an panic.
//
//	      ○      ○      ○      ○      ○      ○       ○
//	      ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	chain │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	      ├──────┼──────┼──────┼──────┼─────┼────────┤
//	      ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index 0      1      2      3      4    n-1       n
func (c *ServerMiddlewareChain) InsertAll(ms ...ServerMiddleware) {
	c.BeforeAll(ms...)
	c.Add(ms...)
}

// BeforeAll inserts the ms before all existing middleware.
// That is, the given ms is inserted at the index from 0 to n-1.
// Note that nil middleware will cause an panic.
//
//	      ○      ○      ○      ○      ○      ○       ×
//	      ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	chain │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	      ├──────┼──────┼──────┼──────┼─────┼────────┤
//	      ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index 0      1      2      3      4    n-1       n
func (c *ServerMiddlewareChain) BeforeAll(ms ...ServerMiddleware) {
	k := 1 + len(ms)
	newChain := make([]ServerMiddleware, k*len(*c))
	for i, cc := range *c {
		for j, m := range ms {
			newChain[k*i+j] = m
		}
		newChain[k*i+k-1] = cc
	}
	*c = newChain
}

// AfterAll inserts the ms after all existing middleware.
// That is, the given ms is inserted at the index from 1 to n.
// Note that nil middleware will cause an panic.
//
//	      ×      ○      ○      ○      ○      ○       ○
//	      ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	chain │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	      ├──────┼──────┼──────┼──────┼─────┼────────┤
//	      ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index 0      1      2      3      4    n-1       n
func (c *ServerMiddlewareChain) AfterAll(ms ...ServerMiddleware) {
	k := 1 + len(ms)
	newChain := make([]ServerMiddleware, k*len(*c))
	for i, cc := range *c {
		newChain[k*i] = cc
		for j, m := range ms {
			newChain[k*i+j+1] = m
		}
	}
	*c = newChain
}

// NewRoundTripper returns a http round tripper with client middleware.
// NewRoundTripper is short for ClientMiddlewareChain(ms).RoundTripper(rt).
func NewRoundTripper(rt http.RoundTripper, ms ...ClientMiddleware) http.RoundTripper {
	return ClientMiddlewareChain(ms).RoundTripper(rt)
}

// ClientMiddlewareChain is the middleware chain of [ClientMiddleware].
// ClientMiddlewareChain implements the [ClientMiddleware] interface itself.
type ClientMiddlewareChain []ClientMiddleware

// ClientMiddleware is the implementation of [ClientMiddleware.ClientMiddleware].
func (c ClientMiddlewareChain) ClientMiddleware(next http.RoundTripper) http.RoundTripper {
	return c.RoundTripper(next)
}

// RoundTripper returns a round tripper with middleware.
// Additional middleware is applied at the end of the chain
// but they are not registered to this middleware chain.
func (c ClientMiddlewareChain) RoundTripper(rt http.RoundTripper, ms ...ClientMiddleware) http.RoundTripper {
	for i := len(ms) - 1; i >= 0; i-- {
		rt = ms[i].ClientMiddleware(rt)
	}
	for i := len(c) - 1; i >= 0; i-- {
		rt = c[i].ClientMiddleware(rt)
	}
	return rt
}

// Add adds given middleware to the chain.
// Middleware are added at the end of the chain keeping the order.
// Note that nil value will cause an panic.
func (c *ClientMiddlewareChain) Add(ms ...ClientMiddleware) {
	*c = append(*c, ms...)
}

// Insert inserts the given middleware at the position of index.
// If index<=0, ms are added at the beginning of the chain.
// If index>=len(c),  ms are added at the end of the chain.
// Note that nil middleware will cause an panic.
//
//	      ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	chain │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	      ├──────┼──────┼──────┼──────┼─────┼────────┤
//	      ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index 0      1      2      3      4    n-1       n
func (c *ClientMiddlewareChain) Insert(index int, ms ...ClientMiddleware) {
	index = min(max(index, 0), len(*c))
	*c = slices.Insert(*c, index, ms...)
}

// InsertAll inserts the ms between each middleware.
// That is, the given ms is inserted at the index from 0 to n.
// Inserting to an empty chain will be ms itself.
// Note that nil middleware will cause an panic.
//
//	      ○      ○      ○      ○      ○      ○       ○
//	      ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	chain │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	      ├──────┼──────┼──────┼──────┼─────┼────────┤
//	      ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index 0      1      2      3      4    n-1       n
func (c *ClientMiddlewareChain) InsertAll(ms ...ClientMiddleware) {
	c.BeforeAll(ms...)
	c.Add(ms...)
}

// BeforeAll inserts the ms before all existing middleware.
// That is, the given ms is inserted at the index from 0 to n-1.
// Note that nil middleware will cause an panic.
//
//	      ○      ○      ○      ○      ○      ○       ×
//	      ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	chain │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	      ├──────┼──────┼──────┼──────┼─────┼────────┤
//	      ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index 0      1      2      3      4    n-1       n
func (c *ClientMiddlewareChain) BeforeAll(ms ...ClientMiddleware) {
	k := 1 + len(ms)
	newChain := make([]ClientMiddleware, k*len(*c))
	for i, cc := range *c {
		for j, m := range ms {
			newChain[k*i+j] = m
		}
		newChain[k*i+k-1] = cc
	}
	*c = newChain
}

// AfterAll inserts the ms after all existing middleware.
// That is, the given ms is inserted at the index from 1 to n.
// Note that nil middleware will cause an panic.
//
//	      ×      ○      ○      ○      ○      ○       ○
//	      ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	chain │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	      ├──────┼──────┼──────┼──────┼─────┼────────┤
//	      ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index 0      1      2      3      4    n-1       n
func (c *ClientMiddlewareChain) AfterAll(ms ...ClientMiddleware) {
	k := 1 + len(ms)
	newChain := make([]ClientMiddleware, k*len(*c))
	for i, cc := range *c {
		newChain[k*i] = cc
		for j, m := range ms {
			newChain[k*i+j+1] = m
		}
	}
	*c = newChain
}
