package zsha512_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zsha512"
	"github.com/aileron-projects/go/ztesting"
)

// hashSum is the msg and sum list.
// Validation data is generated with:
//   - echo -n "test" | openssl dgst -sha512-224
//   - echo -n "test" | openssl dgst -sha512-256
//   - echo -n "test" | openssl dgst -sha384
//   - echo -n "test" | openssl dgst -sha512
var hashSum = []struct {
	name string
	sf   func(msg []byte) []byte
	eqf  func(msg, sum []byte) bool
	msg  string
	sum  string // Hex encoded sum.
}{
	{"case01", zsha512.Sum224, zsha512.EqualSum224, "", "6ed0dd02806fa89e25de060c19d3ac86cabb87d6a0ddd05c333b84f4"},
	{"case02", zsha512.Sum224, zsha512.EqualSum224, "test", "06001bf08dfb17d2b54925116823be230e98b5c6c278303bc4909a8c"},
	{"case03", zsha512.Sum256, zsha512.EqualSum256, "", "c672b8d1ef56ed28ab87c3622c5114069bdd3ad7b8f9737498d0c01ecef0967a"},
	{"case04", zsha512.Sum256, zsha512.EqualSum256, "test", "3d37fe58435e0d87323dee4a2c1b339ef954de63716ee79f5747f94d974f913f"},
	{"case05", zsha512.Sum384, zsha512.EqualSum384, "", "38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b"},
	{"case06", zsha512.Sum384, zsha512.EqualSum384, "test", "768412320f7b0aa5812fce428dc4706b3cae50e02a64caa16a782249bfe8efc4b7ef1ccb126255d196047dfedf17a0a9"},
	{"case07", zsha512.Sum512, zsha512.EqualSum512, "", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
	{"case08", zsha512.Sum512, zsha512.EqualSum512, "test", "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff"},
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
//   - echo -n "test" | openssl dgst -hmac "key" -sha512-224
//   - echo -n "test" | openssl dgst -hmac "key" -sha512-256
//   - echo -n "test" | openssl dgst -hmac "key" -sha384
//   - echo -n "test" | openssl dgst -hmac "key" -sha512
var hmacSum = []struct {
	name     string
	sf       func(msg, key []byte) []byte
	eqf      func(msg, key, sum []byte) bool
	msg, key string // Hash
	sum      string // Hex encoded sum.
}{
	{"case01", zsha512.HMACSum224, zsha512.HMACEqualSum224, "test", "key", "d21afcf33ae87e09597697046d54239b7d35fa108c060d7ecfea9ca1"},
	{"case02", zsha512.HMACSum224, zsha512.HMACEqualSum224, "", "", "de43f6b96f2d08cebe1ee9c02c53d96b68c1e55b6c15d6843b410d4c"},
	{"case03", zsha512.HMACSum224, zsha512.HMACEqualSum224, "", "key", "0f57635549043abfad00d2cb62d91ac609f08c0ab27c8549c4e78f5b"},
	{"case04", zsha512.HMACSum224, zsha512.HMACEqualSum224, "test", "", "d78188da24584eabe85a12ab34f2bf8ffcaa8062c56e3dcb212b3b17"},
	{"case05", zsha512.HMACSum256, zsha512.HMACEqualSum256, "test", "key", "e73445953627a124717b15fc57ae735567a257cbae46a773616c816a4dce437b"},
	{"case06", zsha512.HMACSum256, zsha512.HMACEqualSum256, "", "", "b79c9951df595274582dc094a1ba46c33e4a36878b2d83cb8553f0fe467dcdcf"},
	{"case07", zsha512.HMACSum256, zsha512.HMACEqualSum256, "", "key", "f6a69e8f50b53a2ad52875eb41f8a4255e3f9aca453ff7d3357ae18e5464b108"},
	{"case08", zsha512.HMACSum256, zsha512.HMACEqualSum256, "test", "", "24646c883b563e6455920e3bbb187ab1f62228a3c4d64e021e717b29a1e3a51b"},
	{"case09", zsha512.HMACSum384, zsha512.HMACEqualSum384, "test", "key", "160a099ad9d6dadb46311cb4e6dfe98aca9ca519c2e0fedc8dc45da419b1173039cc131f0b5f68b2bbc2b635109b57a8"},
	{"case10", zsha512.HMACSum384, zsha512.HMACEqualSum384, "", "", "6c1f2ee938fad2e24bd91298474382ca218c75db3d83e114b3d4367776d14d3551289e75e8209cd4b792302840234adc"},
	{"case11", zsha512.HMACSum384, zsha512.HMACEqualSum384, "", "key", "99f44bb4e73c9d0ef26533596c8d8a32a5f8c10a9b997d30d89a7e35ba1ccf200b985f72431202b891fe350da410e43f"},
	{"case12", zsha512.HMACSum384, zsha512.HMACEqualSum384, "test", "", "a154ade5eb70996838bffda2b49a00df11b43b70264dc8eff989444bf6afd61064d39926b22bb8fc988089128932dcea"},
	{"case13", zsha512.HMACSum512, zsha512.HMACEqualSum512, "test", "key", "287a0fb89a7fbdfa5b5538636918e537a5b83065e4ff331268b7aaa115dde047a9b0f4fb5b828608fc0b6327f10055f7637b058e9e0dbb9e698901a3e6dd461c"},
	{"case14", zsha512.HMACSum512, zsha512.HMACEqualSum512, "", "", "b936cee86c9f87aa5d3c6f2e84cb5a4239a5fe50480a6ec66b70ab5b1f4ac6730c6c515421b327ec1d69402e53dfb49ad7381eb067b338fd7b0cb22247225d47"},
	{"case15", zsha512.HMACSum512, zsha512.HMACEqualSum512, "", "key", "84fa5aa0279bbc473267d05a53ea03310a987cecc4c1535ff29b6d76b8f1444a728df3aadb89d4a9a6709e1998f373566e8f824a8ca93b1821f0b69bc2a2f65e"},
	{"case16", zsha512.HMACSum512, zsha512.HMACEqualSum512, "test", "", "29c5fab077c009b9e6676b2f082a7ab3b0462b41acf75f075b5a7bac5619ec81c9d8bb2e25b6d33800fba279ee492ac7d05220e829464df3ca8e00298c517764"},
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
