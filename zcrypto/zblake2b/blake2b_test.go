package zblake2b_test

import (
	"encoding/hex"
	"slices"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zblake2b"
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
	{"case01", zblake2b.Sum256, zblake2b.EqualSum256, "", "0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8"},
	{"case02", zblake2b.Sum256, zblake2b.EqualSum256, "test", "928b20366943e2afd11ebc0eae2e53a93bf177a4fcf35bcc64d503704e65e202"},
	{"case03", zblake2b.Sum384, zblake2b.EqualSum384, "", "b32811423377f52d7862286ee1a72ee540524380fda1724a6f25d7978c6fd3244a6caf0498812673c5e05ef583825100"},
	{"case04", zblake2b.Sum384, zblake2b.EqualSum384, "test", "8a84b8666c8fcfb69f2ec41f578d7c85fbdb504ea6510fb05b50fcbf7ed8153c77943bc2da73abb136834e1a0d4f22cb"},
	{"case05", zblake2b.Sum512, zblake2b.EqualSum512, "", "786a02f742015903c6c6fd852552d272912f4740e15847618a86e217f71f5419d25e1031afee585313896444934eb04b903a685b1448b755d56f701afe9be2ce"},
	{"case06", zblake2b.Sum512, zblake2b.EqualSum512, "test", "a71079d42853dea26e453004338670a53814b78137ffbed07603a41d76a483aa9bc33b582f77d30a65e6f29a896c0411f38312e1d66e0bf16386c86a89bea572"},
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
//   - https://emn178.github.io/online-tools/blake2b/
var hmacSum = []struct {
	name     string
	sf       func(msg, key []byte) []byte
	eqf      func(msg, key, sum []byte) bool
	msg, key string // Hash
	sum      string // Hex encoded sum.
}{
	{"case01", zblake2b.HMACSum256, zblake2b.HMACEqualSum256, "test", "key", "50e7edf49b2700d4cc8f2f35223b671f7823072e0659f0bc02067001965cb415"},
	{"case02", zblake2b.HMACSum256, zblake2b.HMACEqualSum256, "", "", "0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8"},
	{"case03", zblake2b.HMACSum256, zblake2b.HMACEqualSum256, "", "key", "e65edfce5a36261cd824cb0f0da736b1109dcf20d2b831d598f337bb3552a3e4"},
	{"case04", zblake2b.HMACSum256, zblake2b.HMACEqualSum256, "test", "", "928b20366943e2afd11ebc0eae2e53a93bf177a4fcf35bcc64d503704e65e202"},
	{"case05", zblake2b.HMACSum256, zblake2b.HMACEqualSum256, "test", strings.Repeat("1234567890", 6) + "1234", "0e3172c89eecbc3324b81b3aeec3f101da9f7c4afbff0f944f40380c4809f0fe"},
	{"case06", zblake2b.HMACSum256, zblake2b.HMACEqualSum256, "test", strings.Repeat("1234567890", 6) + "12345", "4952ab140a10381b511b790418d183c7c43c11ebf3c698379e36a742daa3f778"},
	{"case07", zblake2b.HMACSum384, zblake2b.HMACEqualSum384, "test", "key", "4be6f224486039d1fbaa24bf52bec4d302bd8cafd46d822b50a3c233c0e61a9ab10a6a6cc38d1545b84ab5ed3ea9a979"},
	{"case08", zblake2b.HMACSum384, zblake2b.HMACEqualSum384, "", "", "b32811423377f52d7862286ee1a72ee540524380fda1724a6f25d7978c6fd3244a6caf0498812673c5e05ef583825100"},
	{"case09", zblake2b.HMACSum384, zblake2b.HMACEqualSum384, "", "key", "be1b0f20d4fc5c8b60ef377d7134d539d696b19f6c2e142465fcbc4edb3bacfe77c4668e2372359e60ec04d7fe2daa9b"},
	{"case10", zblake2b.HMACSum384, zblake2b.HMACEqualSum384, "test", "", "8a84b8666c8fcfb69f2ec41f578d7c85fbdb504ea6510fb05b50fcbf7ed8153c77943bc2da73abb136834e1a0d4f22cb"},
	{"case11", zblake2b.HMACSum384, zblake2b.HMACEqualSum384, "test", strings.Repeat("1234567890", 6) + "1234", "9f193ea539b9ad2ccc6033162e1d4c76b664262ddbeaf3ebbbb75872aace07f2ef1876122406a900fcd94d2e7024f417"},
	{"case12", zblake2b.HMACSum384, zblake2b.HMACEqualSum384, "test", strings.Repeat("1234567890", 6) + "12345", "931a08439e726a369421802120fa8c087961597730cf0210c658f8c55f0e60e5e3872905bd7600c9e9b385a63da8a381"},
	{"case13", zblake2b.HMACSum512, zblake2b.HMACEqualSum512, "test", "key", "d7abd38eaa165b55132742b74afefcabfadf4764bd9fd8d0b90391b30e65af5eda2f92a5165de75cc9816f3e2d631fab091d89d3c39c82c4528d9bcfc901bd7a"},
	{"case14", zblake2b.HMACSum512, zblake2b.HMACEqualSum512, "", "", "786a02f742015903c6c6fd852552d272912f4740e15847618a86e217f71f5419d25e1031afee585313896444934eb04b903a685b1448b755d56f701afe9be2ce"},
	{"case15", zblake2b.HMACSum512, zblake2b.HMACEqualSum512, "", "key", "5b3cfd8f422b490b764b55eceb330b500c79cbefa9a928ad00202b8b3c5dd778a81122570434a2e3b8bfd028d105dfefd0a9576e88ed66de742ca9fbb5f8d2b6"},
	{"case16", zblake2b.HMACSum512, zblake2b.HMACEqualSum512, "test", "", "a71079d42853dea26e453004338670a53814b78137ffbed07603a41d76a483aa9bc33b582f77d30a65e6f29a896c0411f38312e1d66e0bf16386c86a89bea572"},
	{"case17", zblake2b.HMACSum512, zblake2b.HMACEqualSum512, "test", strings.Repeat("1234567890", 6) + "1234", "9907c5348f66109770e27f2b29a2da341c8f2f29e614261c066f068ff1867252fe95e8f9485df8a2a83b37e68b044b3103c753f06fa1590bae48596146efe2a9"},
	{"case18", zblake2b.HMACSum512, zblake2b.HMACEqualSum512, "test", strings.Repeat("1234567890", 6) + "12345", "500a92b0577350ba81d0d89ecee1b7f89360d78907b59d3149603cf01005cdddfa149854dd6768839887c7e88d4fd9d4299d97ffa471f37d48fbc4c265cb6604"},
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
