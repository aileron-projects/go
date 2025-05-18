package zhttp_test

import (
	"net/http"
	"testing"

	"github.com/aileron-projects/go/znet/zhttp"
	"github.com/aileron-projects/go/ztesting"
)

var (
	nopHandler      = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	nopRoundTripper = zhttp.RoundTripperFunc(func(r *http.Request) (*http.Response, error) { return nil, nil })
)

func TestServerMiddlewareFunc(t *testing.T) {
	t.Parallel()
	called := false
	m := zhttp.ServerMiddlewareFunc(func(next http.Handler) http.Handler {
		return zhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	})
	m.ServerMiddleware(nopHandler).ServeHTTP(nil, nil)
	ztesting.AssertEqual(t, "middleware was not called", true, called)
}

func TestClientMiddlewareFunc(t *testing.T) {
	t.Parallel()
	called := false
	m := zhttp.ClientMiddlewareFunc(func(next http.RoundTripper) http.RoundTripper {
		return zhttp.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			called = true
			return next.RoundTrip(r)
		})
	})
	m.ClientMiddleware(nopRoundTripper).RoundTrip(nil)
	ztesting.AssertEqual(t, "middleware was not called", true, called)
}

type testMiddleware struct {
	name string
	list *[]string
}

func (m *testMiddleware) ServerMiddleware(next http.Handler) http.Handler {
	return zhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*m.list = append(*m.list, m.name)
		next.ServeHTTP(w, r)
	})
}

func (m *testMiddleware) ClientMiddleware(next http.RoundTripper) http.RoundTripper {
	return zhttp.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		*m.list = append(*m.list, m.name)
		return next.RoundTrip(r)
	})
}

func TestNewHandler(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ms   []string
		want []string
	}{
		"case01": {[]string{}, []string{}},
		"case02": {[]string{"m1"}, []string{"m1"}},
		"case03": {[]string{"m1", "m2"}, []string{"m1", "m2"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			ms := []zhttp.ServerMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			zhttp.NewHandler(nopHandler, ms...).ServeHTTP(nil, nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestServerMiddlewareChain_Handler(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		chain []string
		ms    []string
		want  []string
	}{
		"0 chain add 0": {[]string{}, []string{}, []string{}},
		"0 chain add 1": {[]string{}, []string{"mA"}, []string{"mA"}},
		"0 chain add 2": {[]string{}, []string{"mA", "mB"}, []string{"mA", "mB"}},
		"1 chain add 0": {[]string{"m1"}, []string{}, []string{"m1"}},
		"1 chain add 1": {[]string{"m1"}, []string{"mA"}, []string{"m1", "mA"}},
		"1 chain add 2": {[]string{"m1"}, []string{"mA", "mB"}, []string{"m1", "mA", "mB"}},
		"2 chain add 0": {[]string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain add 1": {[]string{"m1", "m2"}, []string{"mA"}, []string{"m1", "m2", "mA"}},
		"2 chain add 2": {[]string{"m1", "m2"}, []string{"mA", "mB"}, []string{"m1", "m2", "mA", "mB"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			var chain zhttp.ServerMiddlewareChain
			for _, n := range tc.chain {
				chain.Add(&testMiddleware{name: n, list: &list})
			}
			ms := []zhttp.ServerMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			chain.Handler(nopHandler, ms...).ServeHTTP(nil, nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestServerMiddlewareChain_Insert(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		index int
		chain []string
		ms    []string
		want  []string
	}{
		"0 chain insert 0 at 0": {0, []string{}, []string{}, []string{}},
		"0 chain insert 1 at 0": {0, []string{}, []string{"mA"}, []string{"mA"}},
		"0 chain insert 2 at 0": {0, []string{}, []string{"mA", "mB"}, []string{"mA", "mB"}},
		"0 chain insert 0 at 1": {1, []string{}, []string{}, []string{}},
		"0 chain insert 1 at 1": {1, []string{}, []string{"mA"}, []string{"mA"}},
		"0 chain insert 2 at 1": {1, []string{}, []string{"mA", "mB"}, []string{"mA", "mB"}},
		"1 chain insert 0 at 0": {0, []string{"m1"}, []string{}, []string{"m1"}},
		"1 chain insert 1 at 0": {0, []string{"m1"}, []string{"mA"}, []string{"mA", "m1"}},
		"1 chain insert 2 at 0": {0, []string{"m1"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1"}},
		"1 chain insert 0 at 1": {1, []string{"m1"}, []string{}, []string{"m1"}},
		"1 chain insert 1 at 1": {1, []string{"m1"}, []string{"mA"}, []string{"m1", "mA"}},
		"1 chain insert 2 at 1": {1, []string{"m1"}, []string{"mA", "mB"}, []string{"m1", "mA", "mB"}},
		"2 chain insert 0 at 0": {0, []string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain insert 1 at 0": {0, []string{"m1", "m2"}, []string{"mA"}, []string{"mA", "m1", "m2"}},
		"2 chain insert 2 at 0": {0, []string{"m1", "m2"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1", "m2"}},
		"2 chain insert 0 at 1": {1, []string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain insert 1 at 1": {1, []string{"m1", "m2"}, []string{"mA"}, []string{"m1", "mA", "m2"}},
		"2 chain insert 2 at 1": {1, []string{"m1", "m2"}, []string{"mA", "mB"}, []string{"m1", "mA", "mB", "m2"}},
		"2 chain insert 0 at 2": {2, []string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain insert 1 at 2": {2, []string{"m1", "m2"}, []string{"mA"}, []string{"m1", "m2", "mA"}},
		"2 chain insert 2 at 2": {2, []string{"m1", "m2"}, []string{"mA", "mB"}, []string{"m1", "m2", "mA", "mB"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			var chain zhttp.ServerMiddlewareChain
			for _, n := range tc.chain {
				chain.Add(&testMiddleware{name: n, list: &list})
			}
			ms := []zhttp.ServerMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			chain.Insert(tc.index, ms...)
			chain.ServerMiddleware(nopHandler).ServeHTTP(nil, nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestServerMiddlewareChain_InsertAll(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		chain []string
		ms    []string
		want  []string
	}{
		"0 chain add 0": {[]string{}, []string{}, []string{}},
		"0 chain add 1": {[]string{}, []string{"mA"}, []string{"mA"}},
		"0 chain add 2": {[]string{}, []string{"mA", "mB"}, []string{"mA", "mB"}},
		"1 chain add 0": {[]string{"m1"}, []string{}, []string{"m1"}},
		"1 chain add 1": {[]string{"m1"}, []string{"mA"}, []string{"mA", "m1", "mA"}},
		"1 chain add 2": {[]string{"m1"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1", "mA", "mB"}},
		"2 chain add 0": {[]string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain add 1": {[]string{"m1", "m2"}, []string{"mA"}, []string{"mA", "m1", "mA", "m2", "mA"}},
		"2 chain add 2": {[]string{"m1", "m2"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1", "mA", "mB", "m2", "mA", "mB"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			var chain zhttp.ServerMiddlewareChain
			for _, n := range tc.chain {
				chain.Add(&testMiddleware{name: n, list: &list})
			}
			ms := []zhttp.ServerMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			chain.InsertAll(ms...)
			chain.ServerMiddleware(nopHandler).ServeHTTP(nil, nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestServerMiddlewareChain_BeforeAll(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		chain []string
		ms    []string
		want  []string
	}{
		"0 chain add 0": {[]string{}, []string{}, []string{}},
		"0 chain add 1": {[]string{}, []string{"mA"}, []string{}},
		"0 chain add 2": {[]string{}, []string{"mA", "mB"}, []string{}},
		"1 chain add 0": {[]string{"m1"}, []string{}, []string{"m1"}},
		"1 chain add 1": {[]string{"m1"}, []string{"mA"}, []string{"mA", "m1"}},
		"1 chain add 2": {[]string{"m1"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1"}},
		"2 chain add 0": {[]string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain add 1": {[]string{"m1", "m2"}, []string{"mA"}, []string{"mA", "m1", "mA", "m2"}},
		"2 chain add 2": {[]string{"m1", "m2"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1", "mA", "mB", "m2"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			var chain zhttp.ServerMiddlewareChain
			for _, n := range tc.chain {
				chain.Add(&testMiddleware{name: n, list: &list})
			}
			ms := []zhttp.ServerMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			chain.BeforeAll(ms...)
			chain.ServerMiddleware(nopHandler).ServeHTTP(nil, nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestServerMiddlewareChain_AfterAll(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		chain []string
		ms    []string
		want  []string
	}{
		"0 chain add 0": {[]string{}, []string{}, []string{}},
		"0 chain add 1": {[]string{}, []string{"mA"}, []string{}},
		"0 chain add 2": {[]string{}, []string{"mA", "mB"}, []string{}},
		"1 chain add 0": {[]string{"m1"}, []string{}, []string{"m1"}},
		"1 chain add 1": {[]string{"m1"}, []string{"mA"}, []string{"m1", "mA"}},
		"1 chain add 2": {[]string{"m1"}, []string{"mA", "mB"}, []string{"m1", "mA", "mB"}},
		"2 chain add 0": {[]string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain add 1": {[]string{"m1", "m2"}, []string{"mA"}, []string{"m1", "mA", "m2", "mA"}},
		"2 chain add 2": {[]string{"m1", "m2"}, []string{"mA", "mB"}, []string{"m1", "mA", "mB", "m2", "mA", "mB"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			var chain zhttp.ServerMiddlewareChain
			for _, n := range tc.chain {
				chain.Add(&testMiddleware{name: n, list: &list})
			}
			ms := []zhttp.ServerMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			chain.AfterAll(ms...)
			chain.ServerMiddleware(nopHandler).ServeHTTP(nil, nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestNewRoundTripper(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ms   []string
		want []string
	}{
		"case01": {[]string{}, []string{}},
		"case02": {[]string{"m1"}, []string{"m1"}},
		"case03": {[]string{"m1", "m2"}, []string{"m1", "m2"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			ms := []zhttp.ClientMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			zhttp.NewRoundTripper(nopRoundTripper, ms...).RoundTrip(nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestClientMiddlewareChain_RoundTripper(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		chain []string
		ms    []string
		want  []string
	}{
		"0 chain add 0": {[]string{}, []string{}, []string{}},
		"0 chain add 1": {[]string{}, []string{"mA"}, []string{"mA"}},
		"0 chain add 2": {[]string{}, []string{"mA", "mB"}, []string{"mA", "mB"}},
		"1 chain add 0": {[]string{"m1"}, []string{}, []string{"m1"}},
		"1 chain add 1": {[]string{"m1"}, []string{"mA"}, []string{"m1", "mA"}},
		"1 chain add 2": {[]string{"m1"}, []string{"mA", "mB"}, []string{"m1", "mA", "mB"}},
		"2 chain add 0": {[]string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain add 1": {[]string{"m1", "m2"}, []string{"mA"}, []string{"m1", "m2", "mA"}},
		"2 chain add 2": {[]string{"m1", "m2"}, []string{"mA", "mB"}, []string{"m1", "m2", "mA", "mB"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			var chain zhttp.ClientMiddlewareChain
			for _, n := range tc.chain {
				chain.Add(&testMiddleware{name: n, list: &list})
			}
			ms := []zhttp.ClientMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			chain.RoundTripper(nopRoundTripper, ms...).RoundTrip(nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestClientMiddlewareChain_Insert(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		index int
		chain []string
		ms    []string
		want  []string
	}{
		"0 chain insert 0 at 0": {0, []string{}, []string{}, []string{}},
		"0 chain insert 1 at 0": {0, []string{}, []string{"mA"}, []string{"mA"}},
		"0 chain insert 2 at 0": {0, []string{}, []string{"mA", "mB"}, []string{"mA", "mB"}},
		"0 chain insert 0 at 1": {1, []string{}, []string{}, []string{}},
		"0 chain insert 1 at 1": {1, []string{}, []string{"mA"}, []string{"mA"}},
		"0 chain insert 2 at 1": {1, []string{}, []string{"mA", "mB"}, []string{"mA", "mB"}},
		"1 chain insert 0 at 0": {0, []string{"m1"}, []string{}, []string{"m1"}},
		"1 chain insert 1 at 0": {0, []string{"m1"}, []string{"mA"}, []string{"mA", "m1"}},
		"1 chain insert 2 at 0": {0, []string{"m1"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1"}},
		"1 chain insert 0 at 1": {1, []string{"m1"}, []string{}, []string{"m1"}},
		"1 chain insert 1 at 1": {1, []string{"m1"}, []string{"mA"}, []string{"m1", "mA"}},
		"1 chain insert 2 at 1": {1, []string{"m1"}, []string{"mA", "mB"}, []string{"m1", "mA", "mB"}},
		"2 chain insert 0 at 0": {0, []string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain insert 1 at 0": {0, []string{"m1", "m2"}, []string{"mA"}, []string{"mA", "m1", "m2"}},
		"2 chain insert 2 at 0": {0, []string{"m1", "m2"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1", "m2"}},
		"2 chain insert 0 at 1": {1, []string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain insert 1 at 1": {1, []string{"m1", "m2"}, []string{"mA"}, []string{"m1", "mA", "m2"}},
		"2 chain insert 2 at 1": {1, []string{"m1", "m2"}, []string{"mA", "mB"}, []string{"m1", "mA", "mB", "m2"}},
		"2 chain insert 0 at 2": {2, []string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain insert 1 at 2": {2, []string{"m1", "m2"}, []string{"mA"}, []string{"m1", "m2", "mA"}},
		"2 chain insert 2 at 2": {2, []string{"m1", "m2"}, []string{"mA", "mB"}, []string{"m1", "m2", "mA", "mB"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			var chain zhttp.ClientMiddlewareChain
			for _, n := range tc.chain {
				chain.Add(&testMiddleware{name: n, list: &list})
			}
			ms := []zhttp.ClientMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			chain.Insert(tc.index, ms...)
			chain.ClientMiddleware(nopRoundTripper).RoundTrip(nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestClientMiddlewareChain_InsertAll(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		chain []string
		ms    []string
		want  []string
	}{
		"0 chain add 0": {[]string{}, []string{}, []string{}},
		"0 chain add 1": {[]string{}, []string{"mA"}, []string{"mA"}},
		"0 chain add 2": {[]string{}, []string{"mA", "mB"}, []string{"mA", "mB"}},
		"1 chain add 0": {[]string{"m1"}, []string{}, []string{"m1"}},
		"1 chain add 1": {[]string{"m1"}, []string{"mA"}, []string{"mA", "m1", "mA"}},
		"1 chain add 2": {[]string{"m1"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1", "mA", "mB"}},
		"2 chain add 0": {[]string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain add 1": {[]string{"m1", "m2"}, []string{"mA"}, []string{"mA", "m1", "mA", "m2", "mA"}},
		"2 chain add 2": {[]string{"m1", "m2"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1", "mA", "mB", "m2", "mA", "mB"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			var chain zhttp.ClientMiddlewareChain
			for _, n := range tc.chain {
				chain.Add(&testMiddleware{name: n, list: &list})
			}
			ms := []zhttp.ClientMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			chain.InsertAll(ms...)
			chain.ClientMiddleware(nopRoundTripper).RoundTrip(nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestClientMiddlewareChain_BeforeAll(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		chain []string
		ms    []string
		want  []string
	}{
		"0 chain add 0": {[]string{}, []string{}, []string{}},
		"0 chain add 1": {[]string{}, []string{"mA"}, []string{}},
		"0 chain add 2": {[]string{}, []string{"mA", "mB"}, []string{}},
		"1 chain add 0": {[]string{"m1"}, []string{}, []string{"m1"}},
		"1 chain add 1": {[]string{"m1"}, []string{"mA"}, []string{"mA", "m1"}},
		"1 chain add 2": {[]string{"m1"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1"}},
		"2 chain add 0": {[]string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain add 1": {[]string{"m1", "m2"}, []string{"mA"}, []string{"mA", "m1", "mA", "m2"}},
		"2 chain add 2": {[]string{"m1", "m2"}, []string{"mA", "mB"}, []string{"mA", "mB", "m1", "mA", "mB", "m2"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			var chain zhttp.ClientMiddlewareChain
			for _, n := range tc.chain {
				chain.Add(&testMiddleware{name: n, list: &list})
			}
			ms := []zhttp.ClientMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			chain.BeforeAll(ms...)
			chain.ClientMiddleware(nopRoundTripper).RoundTrip(nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}

func TestClientMiddlewareChain_AfterAll(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		chain []string
		ms    []string
		want  []string
	}{
		"0 chain add 0": {[]string{}, []string{}, []string{}},
		"0 chain add 1": {[]string{}, []string{"mA"}, []string{}},
		"0 chain add 2": {[]string{}, []string{"mA", "mB"}, []string{}},
		"1 chain add 0": {[]string{"m1"}, []string{}, []string{"m1"}},
		"1 chain add 1": {[]string{"m1"}, []string{"mA"}, []string{"m1", "mA"}},
		"1 chain add 2": {[]string{"m1"}, []string{"mA", "mB"}, []string{"m1", "mA", "mB"}},
		"2 chain add 0": {[]string{"m1", "m2"}, []string{}, []string{"m1", "m2"}},
		"2 chain add 1": {[]string{"m1", "m2"}, []string{"mA"}, []string{"m1", "mA", "m2", "mA"}},
		"2 chain add 2": {[]string{"m1", "m2"}, []string{"mA", "mB"}, []string{"m1", "mA", "mB", "m2", "mA", "mB"}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			list := []string{}
			var chain zhttp.ClientMiddlewareChain
			for _, n := range tc.chain {
				chain.Add(&testMiddleware{name: n, list: &list})
			}
			ms := []zhttp.ClientMiddleware{}
			for _, n := range tc.ms {
				ms = append(ms, &testMiddleware{name: n, list: &list})
			}
			chain.AfterAll(ms...)
			chain.ClientMiddleware(nopRoundTripper).RoundTrip(nil)
			ztesting.AssertEqualSlice(t, "middleware not match", tc.want, list)
		})
	}
}
