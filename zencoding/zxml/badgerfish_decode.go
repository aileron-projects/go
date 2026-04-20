package zxml

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

// Decode decodes XML document read by the decoder into Go values.
//
// Returned value can be:
//   - map[string]any
//   - []map[string]any
//
// Decode:
//   - Ignores [encoding/xml.Comment]
//   - Ignores [encoding/xml.ProcInst]
//   - Ignores [encoding/xml.Directive]
//   - Does not identify CDATA (Limitation of [encoding/xml])
func (b *BadgerFish) Decode(decoder *xml.Decoder) (any, error) {
	objs := make([]map[string]any, 0)
	var token xml.Token
	var err error
	for {
		if token, err = decoder.RawToken(); err != nil {
			if err != io.EOF {
				return nil, &XMLError{Err: err, Cause: CauseXMLDecoder}
			}
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			content, err := b.decode(decoder, t, t.End())
			if err != nil {
				return nil, &XMLError{Err: err, Cause: CauseXMLDecoder}
			}
			objs = append(objs, content)
		}
	}
	if err == io.EOF {
		err = nil
	}
	switch len(objs) {
	case 0:
		return map[string]any{}, err
	case 1:
		return objs[0], err
	default:
		return objs, err
	}
}

func (b *BadgerFish) decode(decoder *xml.Decoder, start xml.StartElement, end xml.EndElement) (map[string]any, error) {
	// Convert attributes into map object.
	attrs := b.attrsToMap(start.Attr)

	var text strings.Builder
	var keys []string
	var children []any

Loop:
	for {
		token, err := decoder.RawToken()
		if err != nil {
			return nil, err
		}
		switch t := token.(type) {
		case xml.CharData:
			trimmed := bytes.TrimSpace(t)
			if len(trimmed) == 0 {
				continue // Ignore text with only space characters.
			}
			if b.TrimSpace {
				text.WriteString(string(trimmed))
			} else {
				text.WriteString(string(t))
			}
		case xml.StartElement:
			content, err := b.decode(decoder, t, t.End())
			if err != nil {
				return nil, err
			}
			for k, v := range content {
				keys = append(keys, k)
				children = append(children, v)
			}
		case xml.EndElement:
			if t == end {
				break Loop
			}
		}
	}

	// Convert XML text value to JSON value.
	var val any
	if b.JSONValue != nil {
		v, err := b.JSONValue(text.String(), start)
		if err != nil {
			return nil, err
		}
		val = v
	} else {
		if text.String() != "" {
			val = text.String()
		}
	}

	if val != nil {
		keys = append(keys, b.TextKey)
		children = append(children, val)
	}
	// Merge attributes and child elements.
	// Key-value pairs of keys-children are stored
	// in the attrs variable.
	mergeChildren(attrs, keys, children)

	if len(attrs) == 0 {
		return map[string]any{
			tokenName(start.Name, b.NamespaceSep): b.emptyVal,
		}, nil
	}

	return map[string]any{
		tokenName(start.Name, b.NamespaceSep): attrs,
	}, nil
}

// attrsToMap converts attributes to map.
// It uses namespaces given by the second argument ns.
// For example, if the namespace attributes are given,
// returned map would like be
//
//	"@xmlns": {
//		"$": "http://abc.com/",
//		"ns": "http://xyz.com/"
//	}
func (b *BadgerFish) attrsToMap(attrs []xml.Attr) map[string]any {
	m := make(map[string]any, 0)
	for _, attr := range attrs {
		name := attr.Name
		var key string
		switch name.Space {
		case "": // Format <elem foo="bar"> or <elem xmlns="http://abc.com/">
			if name.Local != "xmlns" {
				m[attrName(attr.Name, b.AttrPrefix, b.NamespaceSep)] = attr.Value
				continue
			}
			key = b.TextKey
		case "xmlns": // Format <elem xmlns:foo="http://abc.com/">
			key = name.Local
		default: // Format <elem foo:bar="baz">
			m[attrName(attr.Name, b.AttrPrefix, b.NamespaceSep)] = attr.Value
			continue
		}
		if v, ok := m[b.AttrPrefix+"xmlns"]; ok {
			v.(map[string]any)[key] = attr.Value
			continue
		} else {
			m[b.AttrPrefix+"xmlns"] = map[string]any{key: attr.Value}
		}
	}
	return m
}
