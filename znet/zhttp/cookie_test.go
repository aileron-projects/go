package zhttp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aileron-projects/go/znet/zhttp"
	"github.com/aileron-projects/go/ztesting"
)

func TestGetMultiCookie(t *testing.T) {
	t.Parallel()
	t.Run("empty prefix for nil", func(t *testing.T) {
		cks, value := zhttp.GetMultiCookie("", nil)
		ztesting.AssertEqual(t, "cookies should be empty", 0, len(cks))
		ztesting.AssertEqual(t, "value should be empty", "", value)
	})
	t.Run("with prefix for nil", func(t *testing.T) {
		cks, value := zhttp.GetMultiCookie("prefix", nil)
		ztesting.AssertEqual(t, "cookies should be empty", 0, len(cks))
		ztesting.AssertEqual(t, "value should be empty", "", value)
	})
	t.Run("no index", func(t *testing.T) {
		cookies := []*http.Cookie{
			{Name: "prefix", Value: "value"},
		}
		cks, value := zhttp.GetMultiCookie("prefix", cookies)
		ztesting.AssertEqual(t, "cookies should be empty", 0, len(cks))
		ztesting.AssertEqual(t, "value should be empty", "", value)
	})
	t.Run("parse index fails", func(t *testing.T) {
		cookies := []*http.Cookie{
			{Name: "prefixNN", Value: "value"},
		}
		cks, value := zhttp.GetMultiCookie("prefix", cookies)
		ztesting.AssertEqual(t, "cookies should be empty", 0, len(cks))
		ztesting.AssertEqual(t, "value should be empty", "", value)
	})
	t.Run("index from 0 to 0", func(t *testing.T) {
		cookies := []*http.Cookie{
			{Name: "prefix0", Value: "value0"},
		}
		cks, value := zhttp.GetMultiCookie("prefix", cookies)
		ztesting.AssertEqual(t, "cookies should be obtained", 1, len(cks))
		ztesting.AssertEqual(t, "value should be obtained", "value0", value)
	})
	t.Run("index from 0 to 1", func(t *testing.T) {
		cookies := []*http.Cookie{
			{Name: "prefix0", Value: "value0"},
			{Name: "prefix1", Value: "value1"},
		}
		cks, value := zhttp.GetMultiCookie("prefix", cookies)
		ztesting.AssertEqual(t, "cookies should be obtained", 2, len(cks))
		ztesting.AssertEqual(t, "value should be obtained", "value0value1", value)
	})
	t.Run("index starts 1", func(t *testing.T) {
		cookies := []*http.Cookie{
			{Name: "prefix1", Value: "value1"},
		}
		cks, value := zhttp.GetMultiCookie("prefix", cookies)
		ztesting.AssertEqual(t, "cookies should be obtained", 1, len(cks))
		ztesting.AssertEqual(t, "value should be empty", "", value)
	})
	t.Run("ignore irrelevant cookie", func(t *testing.T) {
		cookies := []*http.Cookie{
			{Name: "foo", Value: "foo"},
			{Name: "bar", Value: "bar"},
			{Name: "prefix", Value: "value"},
			{Name: "prefix0", Value: "value0"},
			{Name: "prefix1", Value: "value1"},
		}
		cks, value := zhttp.GetMultiCookie("prefix", cookies)
		ztesting.AssertEqual(t, "cookies should be obtained", 2, len(cks))
		ztesting.AssertEqual(t, "value should be obtained", "value0value1", value)
	})
}

func TestSetMultiCookie(t *testing.T) {
	t.Parallel()
	t.Run("nil cookie", func(t *testing.T) {
		rec := httptest.NewRecorder()
		cookie := zhttp.SetMultiCookie(rec, nil)
		ztesting.AssertEqual(t, "cookie should not be set", 0, len(cookie))
	})
	t.Run("empty value", func(t *testing.T) {
		rec := httptest.NewRecorder()
		cookie := zhttp.SetMultiCookie(rec, &http.Cookie{Name: "test", Value: ""})
		ztesting.AssertEqual(t, "1 cookie should be set", 1, len(cookie))
		ztesting.AssertEqual(t, "0-th cookie not match", http.Cookie{Name: "test0", Value: ""}, *cookie[0])
	})
	t.Run("non-empty value", func(t *testing.T) {
		rec := httptest.NewRecorder()
		cookie := zhttp.SetMultiCookie(rec, &http.Cookie{Name: "test", Value: "value"})
		ztesting.AssertEqual(t, "1 cookie should be set", 1, len(cookie))
		ztesting.AssertEqual(t, "cookie not match", http.Cookie{Name: "test0", Value: "value"}, *cookie[0])
	})
	t.Run("copy attributes", func(t *testing.T) {
		rec := httptest.NewRecorder()
		cookie := zhttp.SetMultiCookie(rec, &http.Cookie{Name: "test", Value: "value", Domain: "example.com"})
		ztesting.AssertEqual(t, "1 cookie should be set", 1, len(cookie))
		ztesting.AssertEqual(t, "cookie not match", http.Cookie{Name: "test0", Value: "value", Domain: "example.com"}, *cookie[0])
	})
	t.Run("large value", func(t *testing.T) {
		rec := httptest.NewRecorder()
		cookie := zhttp.SetMultiCookie(rec, &http.Cookie{Name: "test", Value: strings.Repeat("0123456789", 400)})
		ztesting.AssertEqual(t, "2 cookie should be set", 2, len(cookie))
		ztesting.AssertEqual(t, "0-th cookie not match", http.Cookie{Name: "test0", Value: strings.Repeat("0123456789", 200)}, *cookie[0])
		ztesting.AssertEqual(t, "1-th cookie not match", http.Cookie{Name: "test1", Value: strings.Repeat("0123456789", 200)}, *cookie[1])
	})
	t.Run("max length value", func(t *testing.T) {
		rec := httptest.NewRecorder()
		cookie := zhttp.SetMultiCookie(rec, &http.Cookie{Name: "test", Value: strings.Repeat("0123456789", 396) + "01234567"})
		ztesting.AssertEqual(t, "2 cookie should be set", 2, len(cookie))
		ztesting.AssertEqual(t, "0-th cookie not match", http.Cookie{Name: "test0", Value: strings.Repeat("0123456789", 198) + "0123"}, *cookie[0])
		ztesting.AssertEqual(t, "1-th cookie not match", http.Cookie{Name: "test1", Value: "456789" + strings.Repeat("0123456789", 197) + "01234567"}, *cookie[1])
	})
}
