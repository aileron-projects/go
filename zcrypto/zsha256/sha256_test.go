package zsha256_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zsha256"
	"github.com/aileron-projects/go/ztesting"
)

// hashSum is the msg and sum list.
// Validation data is generated with:
//   - echo -n "test" | openssl dgst -sha224
//   - echo -n "test" | openssl dgst -sha256
var hashSum = []struct {
	name string
	sf   func(msg []byte) []byte
	eqf  func(msg, sum []byte) bool
	msg  string
	sum  string // Hex encoded sum.
}{
	{"case01", zsha256.Sum224, zsha256.EqualSum224, "", "d14a028c2a3a2bc9476102bb288234c415a2b01f828ea62ac5b3e42f"},
	{"case02", zsha256.Sum224, zsha256.EqualSum224, "test", "90a3ed9e32b2aaf4c61c410eb925426119e1a9dc53d4286ade99a809"},
	{"case03", zsha256.Sum256, zsha256.EqualSum256, "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	{"case04", zsha256.Sum256, zsha256.EqualSum256, "test", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
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
//   - echo -n "test" | openssl dgst -hmac "key" -sha224
//   - echo -n "test" | openssl dgst -hmac "key" -sha256
var hmacSum = []struct {
	name     string
	sf       func(msg, key []byte) []byte
	eqf      func(msg, key, sum []byte) bool
	msg, key string // Hash
	sum      string // Hex encoded sum.
}{
	{"case01", zsha256.HMACSum224, zsha256.HMACEqualSum224, "test", "key", "76b34b643e71d7d92afd4c689c0949cbe0c5445feae907aac532a5a1"},
	{"case02", zsha256.HMACSum224, zsha256.HMACEqualSum224, "", "", "5ce14f72894662213e2748d2a6ba234b74263910cedde2f5a9271524"},
	{"case03", zsha256.HMACSum224, zsha256.HMACEqualSum224, "", "key", "5aa677c13ce1128eeb3a5c01cef7f16557cd0b76d18fd557d6ac3962"},
	{"case04", zsha256.HMACSum224, zsha256.HMACEqualSum224, "test", "", "1c28510aaa54e61cce1ec07886756f9401e6469fb080517bc76c24ab"},
	{"case05", zsha256.HMACSum256, zsha256.HMACEqualSum256, "test", "key", "02afb56304902c656fcb737cdd03de6205bb6d401da2812efd9b2d36a08af159"},
	{"case06", zsha256.HMACSum256, zsha256.HMACEqualSum256, "", "", "b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad"},
	{"case07", zsha256.HMACSum256, zsha256.HMACEqualSum256, "", "key", "5d5d139563c95b5967b9bd9a8c9b233a9dedb45072794cd232dc1b74832607d0"},
	{"case08", zsha256.HMACSum256, zsha256.HMACEqualSum256, "test", "", "43b0cef99265f9e34c10ea9d3501926d27b39f57c6d674561d8ba236e7a819fb"},
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
