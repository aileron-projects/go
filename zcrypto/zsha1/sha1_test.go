package zsha1_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zsha1"
	"github.com/aileron-projects/go/ztesting"
)

// hashSum is the msg and sum list.
// Validation data is generated with:
//   - echo -n "test" | openssl dgst -sha1
var hashSum = []struct {
	name string
	sf   func(msg []byte) []byte
	eqf  func(msg, sum []byte) bool
	msg  string
	sum  string // Hex encoded sum.
}{
	{"case01", zsha1.Sum, zsha1.EqualSum, "", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
	{"case02", zsha1.Sum, zsha1.EqualSum, "test", "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"},
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
//   - echo -n "test" | openssl dgst -hmac "key" -sha1
var hmacSum = []struct {
	name     string
	sf       func(msg, key []byte) []byte
	eqf      func(msg, key, sum []byte) bool
	msg, key string // Hash
	sum      string // Hex encoded sum.
}{
	{"case01", zsha1.HMACSum, zsha1.HMACEqualSum, "test", "key", "671f54ce0c540f78ffe1e26dcf9c2a047aea4fda"},
	{"case02", zsha1.HMACSum, zsha1.HMACEqualSum, "", "", "fbdb1d1b18aa6c08324b7d64b71fb76370690e1d"},
	{"case03", zsha1.HMACSum, zsha1.HMACEqualSum, "", "key", "f42bb0eeb018ebbd4597ae7213711ec60760843f"},
	{"case04", zsha1.HMACSum, zsha1.HMACEqualSum, "test", "", "72d224d342638722965fa0df997f91ab2e9ed94d"},
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
