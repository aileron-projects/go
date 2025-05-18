package ztransform

import (
	"golang.org/x/text/transform"
)

// String returns a string with the result of converting s using t.
// It calls Reset on t.
// If t is nil, it returns s as-is and nil error.
// See also [golang.org/x/text/transform.String].
func String(t transform.Transformer, s string) (string, error) {
	if t == nil {
		return s, nil
	}
	ss, _, err := transform.String(t, s)
	return ss, err
}

// Bytes returns a byte slice with the result of converting b using t.
// It calls Reset on t.
// If t is nil, it returns b as-is and nil error.
// See also [golang.org/x/text/transform.Bytes].
func Bytes(t transform.Transformer, b []byte) ([]byte, error) {
	if t == nil {
		return b, nil
	}
	bb, _, err := transform.Bytes(t, b)
	return bb, err
}

// StringSlice returns a string slice with the result of
// converting all strings in ss using t.
// It calls Reset on t.
// If t is nil, it returns ss as-is and nil error.
// See also [golang.org/x/text/transform.String].
func StringSlice(t transform.Transformer, ss []string) ([]string, error) {
	if t == nil {
		return ss, nil
	}
	newSS := make([]string, len(ss))
	var err error
	for i, s := range ss {
		newSS[i], _, err = transform.String(t, s)
		if err != nil {
			return nil, err
		}
	}
	return newSS, err
}

// BytesSlice returns a byte slice with the result of
// converting all byte slices in bs using t.
// It calls Reset on t.
// If t is nil, it returns bs as-is and nil error.
// See also [golang.org/x/text/transform.Bytes].
func BytesSlice(t transform.Transformer, bs [][]byte) ([][]byte, error) {
	if t == nil {
		return bs, nil
	}
	newBS := make([][]byte, len(bs))
	var err error
	for i, b := range bs {
		newBS[i], _, err = transform.Bytes(t, b)
		if err != nil {
			return nil, err
		}
	}
	return newBS, err
}
