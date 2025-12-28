package zhttp

import (
	"net/http"
	"strconv"
	"strings"
)

// GetMultiCookie returns a single cookie from multiple.
// Cookie names must be prefix+index such as "name0", "name1", "name2", ...
// For example, if the cookies have "session0=value1", "session1=value2", "session2=value3"
// then the returned value will be "value1value2value3".
// Use [SetMultiCookie] to save a cookie value that may be larger than the value browsers allow.
// First returned map is the cookies those names are prefix+index.
// An empty string "" will be returned if there were missing index.
//
// Example:
//
//	// Let r be *http.Request
//	cookies, value := zhttp.GetMultiCookie("SESSION", r.Cookies())
//	fmt.Println("Reconstructed cookie value". value)
func GetMultiCookie(prefix string, cookies []*http.Cookie) (map[int]*http.Cookie, string) {
	cks := map[int]*http.Cookie{}
	for _, ck := range cookies {
		after, found := strings.CutPrefix(ck.Name, prefix)
		if !found {
			continue
		}
		index, err := strconv.Atoi(after)
		if err != nil {
			continue // ignore this cookie.
		}
		cks[index] = ck
	}
	var builder strings.Builder
	for i := range len(cks) {
		ck, found := cks[i]
		if !found {
			return cks, "" // i-th cookie is missing.
		}
		builder.WriteString(ck.Value)
	}
	return cks, builder.String()
}

// SetMultiCookie sets cookies which are the parts of given cookie.
// It breaks given cookie into multiple cookies and write them to the response writer.
// Each cookie value does not exceed 3968 bytes.
// Cookie attributes except for name and value are copied from the given cookie.
func SetMultiCookie(w http.ResponseWriter, cookie *http.Cookie) map[int]*http.Cookie {
	maxLen := 1<<12 - 1<<7 // Maximum length of values is 3968 bytes.
	cks := map[int]*http.Cookie{}
	if cookie == nil {
		return cks
	}
	n := len(cookie.Value) // Original value length.
	m := n/maxLen + 1      // Number of cookies.
	l := (n + m - 1) / m   // Length of each values.
	for i := range m {
		copied := *cookie
		ck := &(copied)
		ck.Name = cookie.Name + strconv.Itoa(i)
		ck.Value = cookie.Value[i*l : min((i+1)*l, n)]
		w.Header().Add("Set-Cookie", ck.String())
		cks[i] = ck
	}
	return cks
}
