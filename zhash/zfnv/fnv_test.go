package zfnv_test

import (
	"encoding/hex"
	"testing"

	"github.com/aileron-projects/go/zhash"
	"github.com/aileron-projects/go/zhash/zfnv"
	"github.com/aileron-projects/go/ztesting"
)

func TestAvailable(t *testing.T) {
	t.Parallel()
	ztesting.AssertEqual(t, "hash not available", true, zhash.FNV32.Available())
	ztesting.AssertEqual(t, "hash not available", true, zhash.FNV32a.Available())
	ztesting.AssertEqual(t, "hash not available", true, zhash.FNV64.Available())
	ztesting.AssertEqual(t, "hash not available", true, zhash.FNV64a.Available())
	ztesting.AssertEqual(t, "hash not available", true, zhash.FNV128.Available())
	ztesting.AssertEqual(t, "hash not available", true, zhash.FNV128a.Available())
}

func TestSum(t *testing.T) {
	t.Parallel()

	// Validation data is generated using python.
	// https://pypi.org/project/fnv/
	// Example:
	//     sum = fnv.hash(b"", algorithm=fnv.fnv, bits=32)
	//     print(f"{sum:x}")
	//     sum = fnv.hash(b"test", algorithm=fnv.fnv, bits=32)
	//     print(f"{sum:x}")

	testCase := map[string]struct {
		hf   func([]byte) []byte
		data string // Hash
		want string // Hex encode of hash.
	}{
		"32 empty":   {zfnv.Sum32, "", "811c9dc5"},
		"32":         {zfnv.Sum32, "test", "bc2c0be9"},
		"32a empty":  {zfnv.Sum32a, "", "811c9dc5"},
		"32a":        {zfnv.Sum32a, "test", "afd071e5"},
		"64 empty":   {zfnv.Sum64, "", "cbf29ce484222325"},
		"64":         {zfnv.Sum64, "test", "8c093f7e9fccbf69"},
		"64a empty":  {zfnv.Sum64a, "", "cbf29ce484222325"},
		"64a":        {zfnv.Sum64a, "test", "f9e6e6ef197c2b25"},
		"128 empty":  {zfnv.Sum128, "", "6c62272e07bb014262b821756295c58d"},
		"128":        {zfnv.Sum128, "test", "66ab2a8b6f757277b806e89c56faf339"},
		"128a empty": {zfnv.Sum128a, "", "6c62272e07bb014262b821756295c58d"},
		"128a":       {zfnv.Sum128a, "test", "69d061a9c5757277b806e99413dd99a5"},
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
		"32 not match":   {zfnv.EqualSum32, "", "811c9dc4", false},
		"32 empty":       {zfnv.EqualSum32, "", "811c9dc5", true},
		"32":             {zfnv.EqualSum32, "test", "bc2c0be9", true},
		"32a not match":  {zfnv.EqualSum32a, "", "811c9dc4", false},
		"32a empty":      {zfnv.EqualSum32a, "", "811c9dc5", true},
		"32a":            {zfnv.EqualSum32a, "test", "afd071e5", true},
		"64 not match":   {zfnv.EqualSum64, "", "cbf29ce484222324", false},
		"64 empty":       {zfnv.EqualSum64, "", "cbf29ce484222325", true},
		"64":             {zfnv.EqualSum64, "test", "8c093f7e9fccbf69", true},
		"64a not match":  {zfnv.EqualSum64a, "", "cbf29ce484222324", false},
		"64a empty":      {zfnv.EqualSum64a, "", "cbf29ce484222325", true},
		"64a":            {zfnv.EqualSum64a, "test", "f9e6e6ef197c2b25", true},
		"128 not match":  {zfnv.EqualSum128, "", "6c62272e07bb014262b821756295c58c", false},
		"128 empty":      {zfnv.EqualSum128, "", "6c62272e07bb014262b821756295c58d", true},
		"128":            {zfnv.EqualSum128, "test", "66ab2a8b6f757277b806e89c56faf339", true},
		"128a not match": {zfnv.EqualSum128a, "", "6c62272e07bb014262b821756295c58c", false},
		"128a empty":     {zfnv.EqualSum128a, "", "6c62272e07bb014262b821756295c58d", true},
		"128a":           {zfnv.EqualSum128a, "test", "69d061a9c5757277b806e99413dd99a5", true},
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
		"32":   {zhash.FNV32, "test", "bc2c0be9"},
		"32a":  {zhash.FNV32a, "test", "afd071e5"},
		"64":   {zhash.FNV64, "test", "8c093f7e9fccbf69"},
		"64a":  {zhash.FNV64a, "test", "f9e6e6ef197c2b25"},
		"128":  {zhash.FNV128, "test", "66ab2a8b6f757277b806e89c56faf339"},
		"128a": {zhash.FNV128a, "test", "69d061a9c5757277b806e99413dd99a5"},
	}
	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			got := tc.h.Sum([]byte(tc.data))
			ztesting.AssertEqual(t, "hash not match", tc.want, hex.EncodeToString(got))
		})
	}
}
