package zripemd160_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zripemd160"
	"github.com/aileron-projects/go/ztesting"
)

// hashSum is the msg and sum list.
// Validation data is generated with:
//   - echo -n "test" | openssl dgst -ripemd160
var hashSum = []struct {
	name string
	sf   func(msg []byte) []byte
	eqf  func(msg, sum []byte) bool
	msg  string
	sum  string // Hex encoded sum.
}{
	{"case01", zripemd160.Sum, zripemd160.EqualSum, "", "9c1185a5c5e9fc54612808977ee8f548b2258d31"},
	{"case02", zripemd160.Sum, zripemd160.EqualSum, "test", "5e52fee47e6b070565f74372468cdc699de89107"},
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
//   - echo -n "test" | openssl dgst -hmac "key" -ripemd160
var hmacSum = []struct {
	name     string
	sf       func(msg, key []byte) []byte
	eqf      func(msg, key, sum []byte) bool
	msg, key string // Hash
	sum      string // Hex encoded sum.
}{
	{"case01", zripemd160.HMACSum, zripemd160.HMACEqualSum, "test", "key", "b434297d1e00df816b1e5d056b9f7fed8e1e1a19"},
	{"case02", zripemd160.HMACSum, zripemd160.HMACEqualSum, "", "", "44d86b658a3e7cbc1a2010848b53e35c917720ca"},
	{"case03", zripemd160.HMACSum, zripemd160.HMACEqualSum, "", "key", "eb123f0b89091ba4cb169cc0142520ebe3aa094e"},
	{"case04", zripemd160.HMACSum, zripemd160.HMACEqualSum, "test", "", "0ca12a8497e2b0f588a5e37dfbc3a1615c6bfc0f"},
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
