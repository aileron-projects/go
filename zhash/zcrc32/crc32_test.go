package zcrc32_test

import (
	"encoding/hex"
	"testing"

	"github.com/aileron-projects/go/zhash"
	"github.com/aileron-projects/go/zhash/zcrc32"
	"github.com/aileron-projects/go/ztesting"
)

func TestAvailable(t *testing.T) {
	t.Parallel()
	ztesting.AssertEqual(t, "hash not available", true, zhash.CRC32IEEE.Available())
	ztesting.AssertEqual(t, "hash not available", true, zhash.CRC32Castagnoli.Available())
	ztesting.AssertEqual(t, "hash not available", true, zhash.CRC32Koopman.Available())
}

func TestSum(t *testing.T) {
	t.Parallel()

	// Validation data is generated with:
	// 	- https://www.sunshine2k.de/coding/javascript/crc/crc_js.html
	// 	- https://crccalc.com/?crc=Hello%20Go!&method=CRC-32&datatype=ascii&outtype=hex

	testCase := map[string]struct {
		hf   func([]byte) []byte
		data string // Hash
		want string // Hex encode of hash.
	}{
		"IEEE empty":       {zcrc32.SumIEEE, "", "00000000"},
		"IEEE":             {zcrc32.SumIEEE, "test", "d87f7e0c"},
		"Castagnoli empty": {zcrc32.SumCastagnoli, "", "00000000"},
		"Castagnoli":       {zcrc32.SumCastagnoli, "test", "86a072c0"},
		"Koopman empty":    {zcrc32.SumKoopman, "", "00000000"},
		"Koopman":          {zcrc32.SumKoopman, "test", "5c39ab1e"},
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
		"IEEE not match":       {zcrc32.EqualSumIEEE, "test", "d87f7e0b", false},
		"IEEE empty":           {zcrc32.EqualSumIEEE, "", "00000000", true},
		"IEEE":                 {zcrc32.EqualSumIEEE, "test", "d87f7e0c", true},
		"Castagnoli not match": {zcrc32.EqualSumCastagnoli, "test", "86a072c1", false},
		"Castagnoli empty":     {zcrc32.EqualSumCastagnoli, "", "00000000", true},
		"Castagnoli":           {zcrc32.EqualSumCastagnoli, "test", "86a072c0", true},
		"Koopman not match":    {zcrc32.EqualSumKoopman, "test", "5c39ab1d", false},
		"Koopman empty":        {zcrc32.EqualSumKoopman, "", "00000000", true},
		"Koopman":              {zcrc32.EqualSumKoopman, "test", "5c39ab1e", true},
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
		"IEEE":       {zhash.CRC32IEEE, "test", "d87f7e0c"},
		"Castagnoli": {zhash.CRC32Castagnoli, "test", "86a072c0"},
		"Koopman":    {zhash.CRC32Koopman, "test", "5c39ab1e"},
	}
	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			got := tc.h.Sum([]byte(tc.data))
			ztesting.AssertEqual(t, "hash not match", tc.want, hex.EncodeToString(got))
		})
	}
}
