package zhttp

import (
	"net"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// httpHeader is the http header type.
// Unlike [http.Header], operations for the header
// does not check if the header is nil or not
// and does not apply [textproto.CanonicalMIMEHeaderKey] for keys.
// Caller must ensure that the given keys are in canonical form.
// Use this where performance is important.
type httpHeader http.Header

func (h httpHeader) Add(key string, value string) {
	h[key] = append(h[key], value)
}

func (h httpHeader) Del(key string) {
	delete(h, key)
}

func (h httpHeader) Get(key string) string {
	v := h[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

func (h httpHeader) Set(key string, value string) {
	h[key] = []string{value}
}

func (h httpHeader) Values(key string) []string {
	return h[key]
}

// SetForwardedHeaders sets the Forwarded, X-Forwarded-For, X-Forwarded-Host,
// and X-Forwarded-Proto headers to the given h.
// Argument r and h must not be nil.
// Forwarded header is defined in RFC7239.
//
// References:
//   - https://go.dev/src/net/http/httputil/reverseproxy.go
//   - https://datatracker.ietf.org/doc/rfc7239/
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Forwarded
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-Host
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-Proto
func SetForwardedHeaders(r *http.Request, h http.Header) {
	in := httpHeader(r.Header) // For performance.
	out := httpHeader(h)       // For performance.
	var validIP bool
	var forwarded string
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		validIP = true
		forwarded += "for=\"" + ip + "\""
		if prior := in.Values("X-Forwarded-For"); len(prior) > 0 {
			ip = strings.Join(prior, ", ") + ", " + ip
		}
		setHeader(in, out, "X-Forwarded-For", ip)
		setHeader(in, out, "X-Forwarded-Port", port)
	} else {
		out.Del("X-Forwarded-For")
		out.Del("X-Forwarded-Port")
	}

	forwarded += "; host=\"" + r.Host + "\""
	setHeader(in, out, "X-Forwarded-Host", r.Host)

	if r.TLS == nil {
		forwarded += "; proto=http"
		setHeader(in, out, "X-Forwarded-Proto", "http")
	} else {
		forwarded += "; proto=https"
		setHeader(in, out, "X-Forwarded-Proto", "https")
	}

	if validIP {
		if prior := in.Values("Forwarded"); len(prior) > 0 {
			forwarded = strings.Join(prior, ", ") + ", " + forwarded
		}
		setHeader(in, out, "Forwarded", forwarded)
	}
}

func setHeader(in, out httpHeader, key, value string) {
	prior, ok := in[key]
	omit := ok && prior == nil // nil now means don't populate the header
	if !omit {
		out.Set(key, value)
	}
}

// RemoveHopByHopHeaders removes hop-by-hop headers.
// See the references about removed headers.
//
// References:
//   - https://go.dev/src/net/http/httputil/reverseproxy.go
//   - https://datatracker.ietf.org/doc/rfc7230/
//   - https://datatracker.ietf.org/doc/rfc2616/
//
// Removed headers:
//   - Connection
//   - Keep-Alive
//   - Proxy-Authenticate
//   - Proxy-Authorization
//   - Te
//   - Trailer
//   - Transfer-Encoding
//   - Upgrade
//   - Proxy-Connection
//   - Headers in "Connection"
func RemoveHopByHopHeaders(h http.Header) {
	hh := httpHeader(h) // For performance.
	// RFC 7230, section 6.1: Remove headers listed in the "Connection" header.
	for _, conn := range hh.Values("Connection") {
		var elem string
		for conn != "" {
			elem, conn = ScanElement(conn)
			hh.Del(http.CanonicalHeaderKey(elem))
		}
	}
	// RFC 2616, section 13.5.1: Remove a set of known hop-by-hop headers.
	// This behavior is superseded by the RFC 7230 Connection header, but
	// preserve it for backwards compatibility.
	hh.Del("Connection")          // RFC2616 HTTP/1.1
	hh.Del("Keep-Alive")          // RFC2616 HTTP/1.1
	hh.Del("Proxy-Authenticate")  // RFC2616 HTTP/1.1
	hh.Del("Proxy-Authorization") // RFC2616 HTTP/1.1
	hh.Del("Te")                  // RFC2616 HTTP/1.1
	hh.Del("Trailer")             // RFC2616 HTTP/1.1
	hh.Del("Transfer-Encoding")   // RFC2616 HTTP/1.1
	hh.Del("Upgrade")             // RFC2616 HTTP/1.1
	hh.Del("Proxy-Connection")    // non-standard
}

// CopyHeaders copies headers from src to dst.
// It does not delete existing values but append to it.
func CopyHeaders(dst, src http.Header) {
	for k, v := range src {
		k = http.CanonicalHeaderKey(k)
		dst[k] = append(dst[k], v...)
	}
}

// CopyTrailers copies trailers from src to dst.
// Header keys in src are copied with the prefix
// defined as [http.TrailerPrefix].
// It does not delete existing trailers in dst but append to it.
func CopyTrailers(dst, src http.Header) {
	for k, v := range src {
		k = http.CanonicalHeaderKey(k)
		if !strings.HasPrefix(k, http.TrailerPrefix) {
			k = http.TrailerPrefix + k
		}
		dst[k] = append(dst[k], v...)
	}
}

// MatchMediaType returns the first index of the list
// that matched to the mt. It returns -1 when no matching
// media type was found. Given media type should be in the
// form of "<BaseType>/<SubType>" as defined in RFC 9110.
// An wildcard character "*" can be used for <BaseType> or/and <SubType>.
// But it cannot contain any parameters such as "charset".
// For example "text/plain" matches to the "text/plain", "*/plain", "text/*", "*/*".
// Given strings are not validated whether they are valid media types or not.
//
// References:
//   - https://datatracker.ietf.org/doc/html/rfc9110
//   - https://www.iana.org/assignments/media-types/media-types.xhtml
func MatchMediaType(mt string, list []string) int {
	base, sub, found := strings.Cut(mt, "/")
	if !found {
		return -1
	}
	for i, target := range list {
		b, s, found := strings.Cut(target, "/")
		if !found {
			continue
		}
		if base != "*" && b != "*" && b != base {
			continue
		}
		if sub != "*" && s != "*" && s != sub {
			continue
		}
		return i
	}
	return -1
}

// ParseQualifiedHeader parses header value with qualifier, or q-value.
// It returns sorted values and parameters with ascending order with q-value.
// Q-value, as defined in RFC9110, should be a floating value ranges 0.000 <= q <= 1.000.
// If a q-value has more than 3 digits after decimal point, only first 3 digits are valid and used.
// For example, q-value "0.1234" will be parsed as "0.123".
//
// The entries with the following q-value are excluded from the returned values as it is invalid.
//   - The entry with q=0.
//   - The entry with q<0 or q>1.
//   - The entry with q that have non-parsable number.
//
// References:
//   - https://datatracker.ietf.org/doc/rfc9110/
//   - https://datatracker.ietf.org/doc/rfc8941/
func ParseQualifiedHeader(s string) (values []string, params []map[string]string) {
	values, params = ParseHeader(s)
	qs := make([]float64, 0, len(values))
	for i := 0; i < len(values); i++ { // len(values) must be evaluated every loop.
		qStr, ok := params[i]["q"]
		if !ok {
			qs = append(qs, 1.0)
			continue
		}
		if len(qStr) > len("N.NNN") {
			qStr = qStr[:len("N.NNN")]
		}
		q, err := strconv.ParseFloat(qStr, 64)
		if err != nil || q <= 0.0 || 1.0 < q {
			values = slices.Delete(values, i, i+1)
			params = slices.Delete(params, i, i+1)
			i--
			continue
		}
		qs = append(qs, q)
	}
	sort.SliceStable(qs, func(i, j int) bool {
		if qs[i] > qs[j] {
			values[i], values[j] = values[j], values[i]
			params[i], params[j] = params[j], params[i]
			return true
		}
		return false
	})
	return values, params
}

// ParseHeader parses header value with parameters.
// It returns parsed values and their accompanying parameters.
// Parameters are defined as the following format in
// RFC9110 - 5.6.6. Parameters.
// Note that even white spaces around the "=" is not allowed in RFC9110
// this function accept it and trims white spaces for convenience.
//
//	parameters      = *( OWS ";" OWS [ parameter ] )
//	parameter       = parameter-name "=" parameter-value
//	parameter-name  = token
//	parameter-value = ( token / quoted-string )
//
// References:
//   - https://datatracker.ietf.org/doc/rfc9110/
//   - https://datatracker.ietf.org/doc/rfc7230/
func ParseHeader(s string) (values []string, params []map[string]string) {
	values = []string{}
	params = []map[string]string{}
	for s != "" {
		var v, p string
		v, s = ScanElement(s)
		if v == "" {
			break
		}
		i := strings.IndexByte(v, ';')
		if i < 0 {
			values = append(values, v)
			params = append(params, make(map[string]string, 0))
			continue
		}
		v, p, _ = strings.Cut(v, ";")
		v = trimPrefixOWS(trimSuffixOWS(v))
		if v == "" {
			continue // Ignore empty value.
		}
		ps := make(map[string]string, 2)
		values = append(values, v)
		params = append(params, ps)
		for p != "" {
			var pp string
			pp, p, _ = strings.Cut(p, ";")
			key, val, _ := strings.Cut(pp, "=")
			key = trimPrefixOWS(trimSuffixOWS(key))
			val = trimDQUOTE(trimPrefixOWS(trimSuffixOWS(val)))
			ps[key] = val
		}
	}
	return values, params
}

// ScanElement scans the s which is expected to be in list-based field values
// defined in RFC9110 and returns the first found element.
// It splits s by ',' to parse the first element.
// OWS, or optional white spaces, are trimmed and empty elements
// are ignored. The returned elem will be an empty string
// when no valid elements was found in the s.
//
//	RFC9110 5.6. Common Rules for Defining Field Values
//	Empty elements do not contribute to the count of elements present.
//	A recipient MUST parse and ignore a reasonable number of empty list elements:
//	enough to handle common mistakes by senders that merge values, but not so much
//	that they could be used as a denial-of-service mechanism. In other words,
//	a recipient MUST accept lists that satisfy the following syntax:
//
//	#element => [ element ] *( OWS "," OWS [ element ] )
//
//	Valid example:
//	  "foo,bar"
//	  "foo ,bar,"
//	  "foo , ,bar,charlie"
//
// References:
//   - https://datatracker.ietf.org/doc/rfc9110/
//   - https://datatracker.ietf.org/doc/rfc7230/
func ScanElement(s string) (elem, rest string) {
	for {
		i := strings.IndexByte(s, ',')
		if i < 0 {
			break
		}
		var elem string
		elem, s = s[:i], s[i+1:]
		elem = trimPrefixOWS(trimSuffixOWS(elem)) // Trim optional white spaces.
		if elem == "" {                           // Ignore empty element.
			continue
		}
		return elem, s
	}
	return trimPrefixOWS(trimSuffixOWS(s)), ""
}

// trimPrefixOWS removes leading ' ' and '\t'.
func trimPrefixOWS(s string) string {
	for len(s) > 0 {
		if s[0] == ' ' || s[0] == '\t' {
			s = s[1:]
			continue
		}
		break
	}
	return s
}

// trimSuffixOWS removes trailing ' ' and '\t'.
func trimSuffixOWS(s string) string {
	for len(s) > 0 {
		n := len(s)
		if s[n-1] == ' ' || s[n-1] == '\t' {
			s = s[:n-1]
			continue
		}
		break
	}
	return s
}

// trimDQUOTE removes leading and trailing double quotations '"' from s.
func trimDQUOTE(s string) string {
	n := len(s)
	if n < 2 {
		return s
	}
	if s[0] == '"' && s[n-1] == '"' {
		return s[1 : n-1]
	}
	return s
}
