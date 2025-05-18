package zhttp

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestHTTPHeader(t *testing.T) {
	t.Parallel()
	t.Run("add", func(t *testing.T) {
		h := httpHeader{}
		h.Add("Test", "v1")
		h.Add("Test", "v2")
		ztesting.AssertEqualSlice(t, "values not match", []string{"v1", "v2"}, h["Test"])
	})
	t.Run("del", func(t *testing.T) {
		h := httpHeader{"Test": {"value"}}
		h.Del("Test")
		ztesting.AssertEqualSlice(t, "values not match", []string{}, h["Test"])
	})
	t.Run("get exist", func(t *testing.T) {
		h := httpHeader{"Test": {"v1", "v2"}}
		ztesting.AssertEqual(t, "values not match", "v1", h.Get("Test"))
	})
	t.Run("get non exist", func(t *testing.T) {
		h := httpHeader{}
		ztesting.AssertEqual(t, "values not match", "", h.Get("Test"))
	})
	t.Run("set", func(t *testing.T) {
		h := httpHeader{"Test": {"v1", "v2"}}
		h.Set("Test", "v3")
		ztesting.AssertEqualSlice(t, "values not match", []string{"v3"}, h["Test"])
	})
	t.Run("values", func(t *testing.T) {
		h := httpHeader{"Test": {"v1", "v2"}}
		ztesting.AssertEqualSlice(t, "values not match", []string{"v1", "v2"}, h.Values("Test"))
	})
}

func TestSetForwardedHeaders(t *testing.T) {
	t.Parallel()
	t.Run("valid remote addr", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: "127.0.0.1:1234",
			Host:       "test.com",
		}
		h := http.Header{}
		SetForwardedHeaders(r, h)
		ztesting.AssertEqual(t, "Forwarded not match", `for="127.0.0.1"; host="test.com"; proto=http`, h.Get("Forwarded"))
		ztesting.AssertEqual(t, "X-Forwarded-For not match", "127.0.0.1", h.Get("X-Forwarded-For"))
		ztesting.AssertEqual(t, "X-Forwarded-Port not match", "1234", h.Get("X-Forwarded-Port"))
		ztesting.AssertEqual(t, "X-Forwarded-Host not match", "test.com", h.Get("X-Forwarded-Host"))
		ztesting.AssertEqual(t, "X-Forwarded-Proto not match", "http", h.Get("X-Forwarded-Proto"))
	})
	t.Run("invalid remote addr", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: "INVALID IP",
			Host:       "test.com",
		}
		h := http.Header{}
		SetForwardedHeaders(r, h)
		ztesting.AssertEqual(t, "Forwarded not match", "", h.Get("Forwarded"))
		ztesting.AssertEqual(t, "X-Forwarded-For not match", "", h.Get("X-Forwarded-For"))
		ztesting.AssertEqual(t, "X-Forwarded-Port not match", "", h.Get("X-Forwarded-Port"))
		ztesting.AssertEqual(t, "X-Forwarded-Host not match", "test.com", h.Get("X-Forwarded-Host"))
		ztesting.AssertEqual(t, "X-Forwarded-Proto not match", "http", h.Get("X-Forwarded-Proto"))
	})
	t.Run("https", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: "127.0.0.1:1234",
			Host:       "test.com",
			TLS:        &tls.ConnectionState{},
		}
		h := http.Header{}
		SetForwardedHeaders(r, h)
		ztesting.AssertEqual(t, "Forwarded not match", `for="127.0.0.1"; host="test.com"; proto=https`, h.Get("Forwarded"))
		ztesting.AssertEqual(t, "X-Forwarded-For not match", "127.0.0.1", h.Get("X-Forwarded-For"))
		ztesting.AssertEqual(t, "X-Forwarded-Port not match", "1234", h.Get("X-Forwarded-Port"))
		ztesting.AssertEqual(t, "X-Forwarded-Host not match", "test.com", h.Get("X-Forwarded-Host"))
		ztesting.AssertEqual(t, "X-Forwarded-Proto not match", "https", h.Get("X-Forwarded-Proto"))
	})
	t.Run("Forwarded exists", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: "127.0.0.1:1234",
			Host:       "test.com",
			Header: http.Header{
				"Forwarded":         []string{`for="192.168.0.1"`},
				"X-Forwarded-For":   []string{"192.168.0.1"},
				"X-Forwarded-Port":  []string{"5678"},
				"X-Forwarded-Host":  []string{"prior.com"},
				"X-Forwarded-Proto": []string{"https"},
			},
		}
		h := http.Header{}
		SetForwardedHeaders(r, h)
		ztesting.AssertEqual(t, "Forwarded not match", `for="192.168.0.1", for="127.0.0.1"; host="test.com"; proto=http`, h.Get("Forwarded"))
		ztesting.AssertEqual(t, "X-Forwarded-For not match", "192.168.0.1, 127.0.0.1", h.Get("X-Forwarded-For"))
		ztesting.AssertEqual(t, "X-Forwarded-Port not match", "1234", h.Get("X-Forwarded-Port"))
		ztesting.AssertEqual(t, "X-Forwarded-Host not match", "test.com", h.Get("X-Forwarded-Host"))
		ztesting.AssertEqual(t, "X-Forwarded-Proto not match", "http", h.Get("X-Forwarded-Proto"))
	})
	t.Run("don't set Forwarded", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: "127.0.0.1:1234",
			Host:       "test.com",
			Header:     http.Header{"Forwarded": nil},
		}
		h := http.Header{}
		SetForwardedHeaders(r, h)
		_, found := h["Forwarded"]
		ztesting.AssertEqual(t, "found not match", false, found)
	})
	t.Run("don't set X-Forwarded-For", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: "127.0.0.1:1234",
			Host:       "test.com",
			Header:     http.Header{"X-Forwarded-For": nil},
		}
		h := http.Header{}
		SetForwardedHeaders(r, h)
		_, found := h["X-Forwarded-For"]
		ztesting.AssertEqual(t, "found not match", false, found)
	})
	t.Run("don't set X-Forwarded-Port", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: "127.0.0.1:1234",
			Host:       "test.com",
			Header:     http.Header{"X-Forwarded-Port": nil},
		}
		h := http.Header{}
		SetForwardedHeaders(r, h)
		_, found := h["X-Forwarded-Port"]
		ztesting.AssertEqual(t, "found not match", false, found)
	})
	t.Run("don't set X-Forwarded-Host", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: "127.0.0.1:1234",
			Host:       "test.com",
			Header:     http.Header{"X-Forwarded-Host": nil},
		}
		h := http.Header{}
		SetForwardedHeaders(r, h)
		_, found := h["X-Forwarded-Host"]
		ztesting.AssertEqual(t, "found not match", false, found)
	})
	t.Run("don't set X-Forwarded-Proto", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: "127.0.0.1:1234",
			Host:       "test.com",
			Header:     http.Header{"X-Forwarded-Proto": nil},
		}
		h := http.Header{}
		SetForwardedHeaders(r, h)
		_, found := h["X-Forwarded-Proto"]
		ztesting.AssertEqual(t, "found not match", false, found)
	})
}

func TestRemoveHopByHopHeaders(t *testing.T) {
	t.Parallel()
	h := http.Header{
		"Connection":          []string{"foo, bar"},
		"Keep-Alive":          []string{"Value-Keep-Alive"},
		"Proxy-Authenticate":  []string{"Value-Proxy-Authenticate"},
		"Proxy-Authorization": []string{"Value-Proxy-Authorization"},
		"Te":                  []string{"Value-Te"},
		"Trailer":             []string{"Value-Trailer"},
		"Transfer-Encoding":   []string{"Value-Transfer-Encoding"},
		"Upgrade":             []string{"Value-Upgrade"},
		"Proxy-Connection":    []string{"Value-Proxy-Connection"},
		"Foo":                 []string{"Value-Foo"},
		"Bar":                 []string{"Value-Bar"},
		"Baz":                 []string{"Value-Baz"},
	}
	RemoveHopByHopHeaders(h)
	ztesting.AssertEqual(t, "header not deleted", "", h.Get("Keep-Alive"))
	ztesting.AssertEqual(t, "header not deleted", "", h.Get("Proxy-Authenticate"))
	ztesting.AssertEqual(t, "header not deleted", "", h.Get("Proxy-Authorization"))
	ztesting.AssertEqual(t, "header not deleted", "", h.Get("Te"))
	ztesting.AssertEqual(t, "header not deleted", "", h.Get("Trailer"))
	ztesting.AssertEqual(t, "header not deleted", "", h.Get("Transfer-Encoding"))
	ztesting.AssertEqual(t, "header not deleted", "", h.Get("Upgrade"))
	ztesting.AssertEqual(t, "header not deleted", "", h.Get("Proxy-Connection"))
	ztesting.AssertEqual(t, "header not deleted", "", h.Get("Foo"))
	ztesting.AssertEqual(t, "header not deleted", "", h.Get("Bar"))
	ztesting.AssertEqual(t, "header was deleted", "Value-Baz", h.Get("Baz"))
}

func TestCopyHeaders(t *testing.T) {
	t.Parallel()
	dst := http.Header{
		"Foo": []string{"foo"},
		"Bar": []string{"bar1"},
	}
	src := http.Header{
		"Bar": []string{"bar2"},
		"Baz": []string{"baz"},
	}
	CopyHeaders(dst, src)
	ztesting.AssertEqualSlice(t, "value not match", []string{"foo"}, dst.Values("Foo"))
	ztesting.AssertEqualSlice(t, "value not match", []string{"bar1", "bar2"}, dst.Values("Bar"))
	ztesting.AssertEqualSlice(t, "value not match", []string{"baz"}, dst.Values("Baz"))
}

func TestCopyTrailers(t *testing.T) {
	t.Parallel()
	dst := http.Header{
		"Foo": []string{"foo"},
		"Bar": []string{"bar1"},
	}
	src := http.Header{
		"Bar": []string{"bar2"},
		"Baz": []string{"baz"},
	}
	CopyTrailers(dst, src)
	ztesting.AssertEqualSlice(t, "value not match", []string{"foo"}, dst.Values("Foo"))
	ztesting.AssertEqualSlice(t, "value not match", []string{"bar1"}, dst.Values("Bar"))
	ztesting.AssertEqualSlice(t, "value not match", []string{"bar2"}, dst.Values(http.TrailerPrefix+"Bar"))
	ztesting.AssertEqualSlice(t, "value not match", []string{"baz"}, dst.Values(http.TrailerPrefix+"Baz"))
}

func TestParseQualifiedHeader(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		header string
		values []string
		params []map[string]string
	}{
		"empty1":             {"", []string{}, []map[string]string{}},
		"empty2":             {";,", []string{}, []map[string]string{}},
		"default":            {"foo", []string{"foo"}, []map[string]string{{}}},
		"zero":               {"foo; q=0.000", []string{}, []map[string]string{}},
		"minimum":            {"foo; q=0.001", []string{"foo"}, []map[string]string{{"q": "0.001"}}},
		"maximum":            {"foo; q=1.000", []string{"foo"}, []map[string]string{{"q": "1.000"}}},
		"too much precision": {"foo; q=0.1234", []string{"foo"}, []map[string]string{{"q": "0.1234"}}},
		"more than 1.0":      {"foo; q=1.1", []string{}, []map[string]string{}},
		"less than 0.0":      {"foo; q=-0.1", []string{}, []map[string]string{}},
		"not a number":       {"foo; q=xxx", []string{}, []map[string]string{}},
		"quoted":             {"foo; q=\"0.1\"", []string{"foo"}, []map[string]string{{"q": "0.1"}}},
		"case01":             {"foo, bar", []string{"foo", "bar"}, []map[string]string{{}, {}}},
		"case02":             {"foo;q=0.1, bar", []string{"bar", "foo"}, []map[string]string{{}, {"q": "0.1"}}},
		"case03":             {"foo;q=0.1, bar;q=0.2, baz", []string{"baz", "bar", "foo"}, []map[string]string{{}, {"q": "0.2"}, {"q": "0.1"}}},
		"case04":             {"foo;q=0.1, bar;q=0.0, baz", []string{"baz", "foo"}, []map[string]string{{}, {"q": "0.1"}}},
		"case05":             {"foo;q=0.1, , baz", []string{"baz", "foo"}, []map[string]string{{}, {"q": "0.1"}}},
		"case06":             {"foo;p=0.1, bar;q=0.2", []string{"foo", "bar"}, []map[string]string{{"p": "0.1"}, {"q": "0.2"}}},
		"rfc example01": {
			header: "audio/*; q=0.2, audio/basic", // Example in RFC9110
			values: []string{"audio/basic", "audio/*"},
			params: []map[string]string{{}, {"q": "0.2"}},
		},
		"rfc example02": {
			header: "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c", // Example in RFC9110
			values: []string{"text/html", "text/x-c", "text/x-dvi", "text/plain"},
			params: []map[string]string{{}, {}, {"q": "0.8"}, {"q": "0.5"}},
		},
		"rfc example03": {
			header: "text/*;q=0.3, text/plain;q=0.7, text/plain;format=flowed, text/plain;format=fixed;q=0.4, */*;q=0.5", // Example in RFC9110
			values: []string{"text/plain", "text/plain", "*/*", "text/plain", "text/*"},
			params: []map[string]string{{"format": "flowed"}, {"q": "0.7"}, {"q": "0.5"}, {"format": "fixed", "q": "0.4"}, {"q": "0.3"}},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			values, params := ParseQualifiedHeader(tc.header)
			ztesting.AssertEqualSlice(t, "values not match", tc.values, values)
			for i := range tc.params {
				ztesting.AssertEqualMap(t, "params not match", tc.params[i], params[i])
			}
		})
	}
}

func TestParseHeader(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		header string
		values []string
		params []map[string]string
	}{
		"case01": {"", []string{}, []map[string]string{}},
		"case02": {",", []string{}, []map[string]string{}},
		"case03": {" , ", []string{}, []map[string]string{}},
		"case04": {" , ; ", []string{}, []map[string]string{}},
		"case05": {" , ; , ", []string{}, []map[string]string{}},
		"case06": {"foo", []string{"foo"}, []map[string]string{{}}},
		"case07": {"foo; p=x", []string{"foo"}, []map[string]string{{"p": "x"}}},
		"case08": {"foo; p=x; q=y", []string{"foo"}, []map[string]string{{"p": "x", "q": "y"}}},
		"case09": {"foo, bar", []string{"foo", "bar"}, []map[string]string{{}, {}}},
		"case10": {"foo, bar", []string{"foo", "bar"}, []map[string]string{{}, {}}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			values, params := ParseHeader(tc.header)
			ztesting.AssertEqualSlice(t, "values not match", tc.values, values)
			for i := range tc.params {
				ztesting.AssertEqualMap(t, "params not match", tc.params[i], params[i])
			}
		})
	}
}

func TestMatchMediaType(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		mt   string
		list []string
		want int
	}{
		"case01": {"", []string{}, -1},
		"case02": {"", []string{""}, -1},
		"case03": {"", []string{"text/plain"}, -1},
		"case04": {"", []string{"text/*"}, -1},
		"case05": {"", []string{"*/*"}, -1},
		"case06": {"text/plain", []string{}, -1},
		"case07": {"text/plain", []string{""}, -1},
		"case08": {"text/plain", []string{"text/plain"}, 0},
		"case09": {"text/plain", []string{"text/*"}, 0},
		"case10": {"text/plain", []string{"*/plain"}, 0},
		"case11": {"text/plain", []string{"*/*"}, 0},
		"case12": {"text/plain", []string{"application/json"}, -1},
		"case13": {"text/plain", []string{"application/*"}, -1},
		"case14": {"text/plain", []string{"*/json"}, -1},
		"case15": {"text/plain", []string{"application/json", "text/plain"}, 1},
		"case16": {"text/plain", []string{"application/json", "text/*"}, 1},
		"case17": {"text/plain", []string{"application/json", "*/plain"}, 1},
		"case18": {"text/plain", []string{"application/json", "*/*"}, 1},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := MatchMediaType(tc.mt, tc.list)
			ztesting.AssertEqual(t, "returned type not match", tc.want, got)
		})
	}
}

func TestScanElement(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		elements string
		want     []string
	}{
		"case01": {"", []string{}},                                          // Example in RFC9110
		"case02": {",", []string{}},                                         // Example in RFC9110
		"case03": {",   ,", []string{}},                                     // Example in RFC9110
		"case04": {"foo,bar", []string{"foo", "bar"}},                       // Example in RFC9110
		"case05": {"foo ,bar,", []string{"foo", "bar"}},                     // Example in RFC9110
		"case06": {"foo , ,bar,charlie", []string{"foo", "bar", "charlie"}}, // Example in RFC9110
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var got []string
			s := tc.elements
			for s != "" {
				var e string
				e, s = ScanElement(s)
				if e != "" {
					got = append(got, e)
				}
			}
			ztesting.AssertEqualSlice(t, "elems not match", tc.want, got)
		})
	}
}
