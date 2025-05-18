package zblake2s_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zblake2s"
	"github.com/aileron-projects/go/ztesting"
)

// hashSum is the msg and sum list.
// Validation data is generated with:
//   - echo -n "test" | openssl dgst -blake2s256
var hashSum = []struct {
	name string
	sf   func(msg []byte) []byte
	eqf  func(msg, sum []byte) bool
	msg  string
	sum  string // Hex encoded sum.
}{
	{"case01", zblake2s.Sum256, zblake2s.EqualSum256, "", "69217a3079908094e11121d042354a7c1f55b6482ca1a51e1b250dfd1ed0eef9"},
	{"case02", zblake2s.Sum256, zblake2s.EqualSum256, "test", "f308fc02ce9172ad02a7d75800ecfc027109bc67987ea32aba9b8dcc7b10150e"},
}

func TestSum(t *testing.T) {
	t.Parallel()
	for _, tc := range hashSum {
		t.Run(tc.msg, func(t *testing.T) {
			got := tc.sf([]byte(tc.msg))
			ztesting.AssertEqual(t, "hash not match", tc.sum, hex.EncodeToString(got))
		})
	}
}

func TestEqualSum(t *testing.T) {
	t.Parallel()
	for _, tc := range hashSum {
		t.Run(tc.name, func(t *testing.T) {
			sum, _ := hex.DecodeString(tc.sum)
			equal := tc.eqf([]byte(tc.msg), sum)
			ztesting.AssertEqual(t, "invalid compare result", true, equal)
			slices.Reverse(sum) // Make the sum wrong.
			equal = tc.eqf([]byte(tc.msg), sum)
			ztesting.AssertEqual(t, "invalid compare result", false, equal)
		})
	}
}

// hmacSum is the msg,key and sum list.
// Validation data can be generated with:
//   - https://emn178.github.io/online-tools/blake2s/
var hmacSum = []struct {
	name     string
	sf       func(msg, key []byte) []byte
	eqf      func(msg, key, sum []byte) bool
	msg, key string // Hash
	sum      string // Hex encoded sum.
}{
	{"case01", zblake2s.HMACSum256, zblake2s.HMACEqualSum256, "test", "key", "2fc31698bc218ccc20185a737c14a42ef60eb168973300ca81f96db522a1f6c8"},
	{"case02", zblake2s.HMACSum256, zblake2s.HMACEqualSum256, "", "", "69217a3079908094e11121d042354a7c1f55b6482ca1a51e1b250dfd1ed0eef9"},
	{"case03", zblake2s.HMACSum256, zblake2s.HMACEqualSum256, "", "key", "a65f92611fdc3722a305edf1ed575947aa86209290344f817e45c3a4edfddad9"},
	{"case04", zblake2s.HMACSum256, zblake2s.HMACEqualSum256, "test", "", "f308fc02ce9172ad02a7d75800ecfc027109bc67987ea32aba9b8dcc7b10150e"},
	{"case05", zblake2s.HMACSum256, zblake2s.HMACEqualSum256, "test", "12345678901234567890123456789012", "c25c86d772d34b3e02899561ea275d345bef74623605b0f487036638e0681b38"},
	{"case06", zblake2s.HMACSum256, zblake2s.HMACEqualSum256, "test", "123456789012345678901234567890123", "bef7ed44dee3ac8408747a4def121dc9c6d604987668b1e19692ca70372a23b7"},
}

func TestHMACSum(t *testing.T) {
	t.Parallel()
	for _, tc := range hmacSum {
		t.Run(tc.name, func(t *testing.T) {
			sum := tc.sf([]byte(tc.msg), []byte(tc.key))
			got := hex.EncodeToString(sum)
			ztesting.AssertEqual(t, "invalid sum result", tc.sum, got)
		})
	}
}

func TestHMACEqualSum(t *testing.T) {
	t.Parallel()
	for _, tc := range hmacSum {
		t.Run(tc.name, func(t *testing.T) {
			sum, _ := hex.DecodeString(tc.sum)
			equal := tc.eqf([]byte(tc.msg), []byte(tc.key), sum)
			ztesting.AssertEqual(t, "invalid equal result", true, equal)
			slices.Reverse(sum) // Make the sum wrong.
			equal = tc.eqf([]byte(tc.msg), []byte(tc.key), sum)
			ztesting.AssertEqual(t, "invalid equal result", false, equal)
		})
	}
}
