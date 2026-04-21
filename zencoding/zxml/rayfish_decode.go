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
func (r *RayFish) Decode(decoder *xml.Decoder) (any, error) {
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
			content, err := r.decode(decoder, t, t.End())
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

func (r *RayFish) decode(decoder *xml.Decoder, start xml.StartElement, end xml.EndElement) (map[string]any, error) {
	// Register attributes as children.
	children := make([]map[string]any, 0, len(start.Attr))
	for _, attr := range start.Attr {
		children = append(children, map[string]any{
			r.NameKey:     attrName(attr.Name, r.AttrPrefix, r.NamespaceSep),
			r.TextKey:     attr.Value,
			r.ChildrenKey: make([]map[string]any, 0), // Attribute has no child.
		})
	}

	var text strings.Builder
Loop:
	for {
		token, err := decoder.RawToken()
		if err != nil {
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			content, err := r.decode(decoder, t, t.End())
			if err != nil {
				return nil, err
			}
			children = append(children, content)
		case xml.CharData:
			trimmed := bytes.TrimSpace(t)
			if len(trimmed) == 0 {
				continue // Ignore text with only space characters.
			}
			if r.TrimSpace {
				text.WriteString(string(trimmed))
			} else {
				text.WriteString(string(t))
			}
		case xml.EndElement:
			if t == end {
				break Loop
			}
		}
	}

	// Convert XML text value to JSON value.
	var val any = text.String()
	if r.JSONValue != nil {
		v, err := r.JSONValue(text.String(), start)
		if err != nil {
			return nil, err
		}
		val = v
	} else {
		if text.String() == "" {
			val = r.emptyVal
		}
	}

	return map[string]any{
		r.NameKey:     tokenName(start.Name, r.NamespaceSep),
		r.TextKey:     val,
		r.ChildrenKey: children,
	}, nil
}
