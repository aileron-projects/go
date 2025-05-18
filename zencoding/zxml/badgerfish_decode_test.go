package zxml

import (
	"encoding/xml"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestBadgerFish_Decode(t *testing.T) {
	t.Parallel()
	base := BadgerFish{
		TextKey:      "$",
		AttrPrefix:   "@",
		NamespaceSep: ":",
		emptyVal:     make(map[string]any, 0),
	}
	t.Run("text key", func(t *testing.T) {
		s := base
		s.TextKey = "#text"
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david">bob</alice>`)))
		want := map[string]any{"alice": map[string]any{"#text": "bob", "@charlie": "david"}}
		ztesting.AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		ztesting.AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("attr prefix", func(t *testing.T) {
		s := base
		s.AttrPrefix = "%"
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david">bob</alice>`)))
		want := map[string]any{"alice": map[string]any{"$": "bob", "%charlie": "david"}}
		ztesting.AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		ztesting.AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("namespace sep", func(t *testing.T) {
		s := base
		s.NamespaceSep = "_**_"
		d, _ := s.Decode(xml.NewDecoder(strings.NewReader(`<alice foo:charlie="david">bob</alice>`)))
		want := map[string]any{"alice": map[string]any{"$": "bob", "@foo_**_charlie": "david"}}
		ztesting.AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
	})
	t.Run("trim space", func(t *testing.T) {
		s := base
		s.TrimSpace = true
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david"> bob </alice>`)))
		want := map[string]any{"alice": map[string]any{"$": "bob", "@charlie": "david"}}
		ztesting.AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		ztesting.AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("json value", func(t *testing.T) {
		s := base
		s.JSONValue = func(s string, se xml.StartElement) (any, error) { return "value", nil }
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david">bob</alice>`)))
		want := map[string]any{"alice": map[string]any{"$": "value", "@charlie": "david"}}
		ztesting.AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		ztesting.AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("json value error", func(t *testing.T) {
		s := base
		e := errors.New("parse error")
		s.JSONValue = func(s string, se xml.StartElement) (any, error) { return nil, e }
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice charlie="david">bob</alice>`)))
		ztesting.AssertEqualErr(t, "error not match", e, err)
		ztesting.AssertEqual(t, "map is not nil", true, reflect.DeepEqual(nil, d))
	})
	t.Run("inner json value error", func(t *testing.T) {
		s := base
		e := errors.New("parse error")
		s.JSONValue = func(s string, se xml.StartElement) (any, error) { return nil, e }
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice><bob>charlie</bob></alice>`)))
		ztesting.AssertEqualErr(t, "error not match", e, err)
		ztesting.AssertEqual(t, "map is not nil", true, reflect.DeepEqual(nil, d))
	})
}

func TestBadgerFish_WithEmptyValue(t *testing.T) {
	t.Parallel()
	base := BadgerFish{
		TextKey:      "$",
		AttrPrefix:   "@",
		NamespaceSep: ":",
		emptyVal:     make(map[string]any, 0),
	}
	t.Run("nil", func(t *testing.T) {
		s := base
		s.WithEmptyValue(nil)
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice></alice>`)))
		want := map[string]any{"alice": nil}
		ztesting.AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		ztesting.AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("string", func(t *testing.T) {
		s := base
		s.WithEmptyValue("")
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice></alice>`)))
		want := map[string]any{"alice": ""}
		ztesting.AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		ztesting.AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("map", func(t *testing.T) {
		s := base
		s.WithEmptyValue(map[string]any{})
		d, err := s.Decode(xml.NewDecoder(strings.NewReader(`<alice></alice>`)))
		want := map[string]any{"alice": map[string]any{}}
		ztesting.AssertEqual(t, "map not match", true, reflect.DeepEqual(want, d))
		ztesting.AssertEqual(t, "error is not nil", nil, err)
	})
	t.Run("panic", func(t *testing.T) {
		s := base
		defer func() {
			r := recover()
			ztesting.AssertEqualErr(t, "error not match", &XMLError{Cause: CauseEmptyVal}, r.(error))
		}()
		s.WithEmptyValue(123)
	})
}
