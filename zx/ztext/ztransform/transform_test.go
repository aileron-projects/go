package ztransform_test

import (
	"io"
	"testing"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/zx/ztext/ztransform"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var et = &errTransformer{
	NopResetter: &transform.NopResetter{},
	err:         io.ErrShortWrite,
}

type errTransformer struct {
	*transform.NopResetter
	err error
}

func (t *errTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	return 0, 0, t.err
}

func TestString(t *testing.T) {
	t.Parallel()
	sjisEnc := japanese.ShiftJIS.NewEncoder()
	sjisDec := japanese.ShiftJIS.NewDecoder()
	testCases := map[string]struct {
		t           transform.Transformer
		input, want string
		err         error
	}{
		"nil t 01": {nil, "", "", nil},
		"nil t 02": {nil, "test", "test", nil},
		"sjis encode": {sjisEnc,
			"月日は百代の過客にして、行かふ年も又旅人也。",
			"\x8c\x8e\x93\xfa\x82\xcd\x95\x53\x91\xe3\x82\xcc\x89\xdf\x8b\x71" +
				"\x82\xc9\x82\xb5\x82\xc4\x81\x41\x8d\x73\x82\xa9\x82\xd3\x94\x4e" +
				"\x82\xe0\x96\x94\x97\xb7\x90\x6c\x96\xe7\x81\x42",
			nil},
		"sjis decode": {sjisDec,
			"\x8c\x8e\x93\xfa\x82\xcd\x95\x53\x91\xe3\x82\xcc\x89\xdf\x8b\x71" +
				"\x82\xc9\x82\xb5\x82\xc4\x81\x41\x8d\x73\x82\xa9\x82\xd3\x94\x4e" +
				"\x82\xe0\x96\x94\x97\xb7\x90\x6c\x96\xe7\x81\x42",
			"月日は百代の過客にして、行かふ年も又旅人也。",
			nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := ztransform.String(tc.t, tc.input)
			ztesting.AssertEqual(t, "string not match", tc.want, got)
			ztesting.AssertEqual(t, "error not match", tc.err, err)
		})
	}
}

func TestBytes(t *testing.T) {
	t.Parallel()
	sjisEnc := japanese.ShiftJIS.NewEncoder()
	sjisDec := japanese.ShiftJIS.NewDecoder()
	testCases := map[string]struct {
		t           transform.Transformer
		input, want string
		err         error
	}{
		"nil t 01": {nil, "", "", nil},
		"nil t 02": {nil, "test", "test", nil},
		"sjis encode": {sjisEnc,
			"月日は百代の過客にして、行かふ年も又旅人也。",
			"\x8c\x8e\x93\xfa\x82\xcd\x95\x53\x91\xe3\x82\xcc\x89\xdf\x8b\x71" +
				"\x82\xc9\x82\xb5\x82\xc4\x81\x41\x8d\x73\x82\xa9\x82\xd3\x94\x4e" +
				"\x82\xe0\x96\x94\x97\xb7\x90\x6c\x96\xe7\x81\x42",
			nil},
		"sjis decode": {sjisDec,
			"\x8c\x8e\x93\xfa\x82\xcd\x95\x53\x91\xe3\x82\xcc\x89\xdf\x8b\x71" +
				"\x82\xc9\x82\xb5\x82\xc4\x81\x41\x8d\x73\x82\xa9\x82\xd3\x94\x4e" +
				"\x82\xe0\x96\x94\x97\xb7\x90\x6c\x96\xe7\x81\x42",
			"月日は百代の過客にして、行かふ年も又旅人也。",
			nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := ztransform.Bytes(tc.t, []byte(tc.input))
			ztesting.AssertEqual(t, "bytes not match", tc.want, string(got))
			ztesting.AssertEqual(t, "error not match", tc.err, err)
		})
	}
}

func TestStringSlice(t *testing.T) {
	t.Parallel()
	sjisEnc := japanese.ShiftJIS.NewEncoder()
	sjisDec := japanese.ShiftJIS.NewDecoder()
	testCases := map[string]struct {
		t           transform.Transformer
		input, want []string
		err         error
	}{
		"error":    {et, []string{""}, nil, io.ErrShortWrite},
		"nil t 01": {nil, []string{""}, []string{""}, nil},
		"nil t 02": {nil, []string{"test"}, []string{"test"}, nil},
		"sjis encode": {sjisEnc,
			[]string{"月日", "百代", "過客"},
			[]string{"\x8c\x8e\x93\xfa", "\x95\x53\x91\xe3", "\x89\xdf\x8b\x71"},
			nil},
		"sjis decode": {sjisDec,
			[]string{"\x8c\x8e\x93\xfa", "\x95\x53\x91\xe3", "\x89\xdf\x8b\x71"},
			[]string{"月日", "百代", "過客"},
			nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := ztransform.StringSlice(tc.t, tc.input)
			ztesting.AssertEqual(t, "error not match", tc.err, err)
			if err != nil {
				ztesting.AssertEqual(t, "non zeo length slice returned", 0, len(got))
				return
			}
			ztesting.AssertEqualSlice(t, "string slice not match", tc.want, got)
		})
	}
}

func TestBytesSlice(t *testing.T) {
	t.Parallel()
	sjisEnc := japanese.ShiftJIS.NewEncoder()
	sjisDec := japanese.ShiftJIS.NewDecoder()
	testCases := map[string]struct {
		t           transform.Transformer
		input, want [][]byte
		err         error
	}{
		"error":    {et, [][]byte{[]byte("")}, [][]byte{[]byte("")}, io.ErrShortWrite},
		"nil t 01": {nil, [][]byte{[]byte("")}, [][]byte{[]byte("")}, nil},
		"nil t 02": {nil, [][]byte{[]byte("test")}, [][]byte{[]byte("test")}, nil},
		"sjis encode": {sjisEnc,
			[][]byte{[]byte("月日"), []byte("百代"), []byte("過客")},
			[][]byte{[]byte("\x8c\x8e\x93\xfa"), []byte("\x95\x53\x91\xe3"), []byte("\x89\xdf\x8b\x71")},
			nil},
		"sjis decode": {sjisDec,
			[][]byte{[]byte("\x8c\x8e\x93\xfa"), []byte("\x95\x53\x91\xe3"), []byte("\x89\xdf\x8b\x71")},
			[][]byte{[]byte("月日"), []byte("百代"), []byte("過客")},
			nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := ztransform.BytesSlice(tc.t, tc.input)
			ztesting.AssertEqual(t, "error not match", tc.err, err)
			if err != nil {
				ztesting.AssertEqual(t, "non zeo length slice returned", 0, len(got))
				return
			}
			for i := range tc.want {
				ztesting.AssertEqualSlice(t, "bytes slice not match", tc.want[i], got[i])
			}
		})
	}
}
