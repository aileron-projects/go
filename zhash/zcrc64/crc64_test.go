package zcrc64_test

import (
	"encoding/hex"
	"testing"

	"github.com/aileron-projects/go/zhash"
	"github.com/aileron-projects/go/zhash/zcrc64"
	"github.com/aileron-projects/go/ztesting"
)

func TestAvailable(t *testing.T) {
	t.Parallel()
	ztesting.AssertEqual(t, "hash not available", true, zhash.CRC64ISO.Available())
	ztesting.AssertEqual(t, "hash not available", true, zhash.CRC64ECMA.Available())
}

func TestSum(t *testing.T) {
	t.Parallel()

	// Validation data is generated with:
	// 	- https://www.sunshine2k.de/coding/javascript/crc/crc_js.html
	// 	- https://toolkitbay.com/tkb/tool/CRC-64

	testCase := map[string]struct {
		hf   func([]byte) []byte
		data string // Hash
		want string // Hex encode of hash.
	}{
		"ISO empty":  {zcrc64.SumISO, "", "0000000000000000"},
		"ISO":        {zcrc64.SumISO, "test", "287c72c850000000"},
		"ECMA empty": {zcrc64.SumECMA, "", "0000000000000000"},
		"ECMA":       {zcrc64.SumECMA, "test", "fa15fda7c10c75a5"},
	}
	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			got := tc.hf([]byte(tc.data))
			ztesting.AssertEqual(t, "hash not match", tc.want, hex.EncodeToString(got))
		})
	}
}

func TestEqualSum(t *testing.T) {
	t.Parallel()
	testCase := map[string]struct {
		cf   func([]byte, []byte) bool
		data string // Hash
		sum  string // Hex encode of hash.
		want bool
	}{
		"ISO not match":  {zcrc64.EqualSumISO, "test", "287c72c850000001", false},
		"ISO empty":      {zcrc64.EqualSumISO, "", "0000000000000000", true},
		"ISO":            {zcrc64.EqualSumISO, "test", "287c72c850000000", true},
		"ECMA not match": {zcrc64.EqualSumECMA, "test", "fa15fda7c10c75a4", false},
		"ECMA empty":     {zcrc64.EqualSumECMA, "", "0000000000000000", true},
		"ECMA":           {zcrc64.EqualSumECMA, "test", "fa15fda7c10c75a5", true},
	}
	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			d, _ := hex.DecodeString(tc.sum)
			got := tc.cf([]byte(tc.data), d)
			ztesting.AssertEqual(t, "invalid compare result", tc.want, got)
		})
	}
}

func TestHashSum(t *testing.T) {
	t.Parallel()
	testCase := map[string]struct {
		h    zhash.Hash
		data string // Hash
		want string // Hex encode of hash.
	}{
		"ISO":  {zhash.CRC64ISO, "test", "287c72c850000000"},
		"ECMA": {zhash.CRC64ECMA, "test", "fa15fda7c10c75a5"},
	}
	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			got := tc.h.Sum([]byte(tc.data))
			ztesting.AssertEqual(t, "hash not match", tc.want, hex.EncodeToString(got))
		})
	}
}
