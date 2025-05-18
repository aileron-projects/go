package zsha3_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zsha3"
	"github.com/aileron-projects/go/ztesting"
)

// hashSum is the msg and sum list.
// Validation data is generated with:
//   - echo -n "test" | openssl dgst -sha3-224
//   - echo -n "test" | openssl dgst -sha3-256
//   - echo -n "test" | openssl dgst -sha3-384
//   - echo -n "test" | openssl dgst -sha3-512
//   - echo -n "test" | openssl dgst -shake128 -xoflen 32
//   - echo -n "test" | openssl dgst -shake256 -xoflen 64
var hashSum = []struct {
	name string
	sf   func(msg []byte) []byte
	eqf  func(msg, sum []byte) bool
	msg  string
	sum  string // Hex encoded sum.
}{
	{"case01", zsha3.Sum224, zsha3.EqualSum224, "", "6b4e03423667dbb73b6e15454f0eb1abd4597f9a1b078e3f5b5a6bc7"},
	{"case02", zsha3.Sum224, zsha3.EqualSum224, "test", "3797bf0afbbfca4a7bbba7602a2b552746876517a7f9b7ce2db0ae7b"},
	{"case03", zsha3.Sum256, zsha3.EqualSum256, "", "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"},
	{"case04", zsha3.Sum256, zsha3.EqualSum256, "test", "36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"},
	{"case05", zsha3.Sum384, zsha3.EqualSum384, "", "0c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f004"},
	{"case06", zsha3.Sum384, zsha3.EqualSum384, "test", "e516dabb23b6e30026863543282780a3ae0dccf05551cf0295178d7ff0f1b41eecb9db3ff219007c4e097260d58621bd"},
	{"case07", zsha3.Sum512, zsha3.EqualSum512, "", "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"},
	{"case08", zsha3.Sum512, zsha3.EqualSum512, "test", "9ece086e9bac491fac5c1d1046ca11d737b92a2b2ebd93f005d7b710110c0a678288166e7fbe796883a4f2e9b3ca9f484f521d0ce464345cc1aec96779149c14"},
	{"case09", zsha3.SumShake128, zsha3.EqualSumShake128, "", "7f9c2ba4e88f827d616045507605853ed73b8093f6efbc88eb1a6eacfa66ef26"},
	{"case10", zsha3.SumShake128, zsha3.EqualSumShake128, "test", "d3b0aa9cd8b7255622cebc631e867d4093d6f6010191a53973c45fec9b07c774"},
	{"case11", zsha3.SumShake256, zsha3.EqualSumShake256, "", "46b9dd2b0ba88d13233b3feb743eeb243fcd52ea62b81b82b50c27646ed5762fd75dc4ddd8c0f200cb05019d67b592f6fc821c49479ab48640292eacb3b7c4be"},
	{"case12", zsha3.SumShake256, zsha3.EqualSumShake256, "test", "b54ff7255705a71ee2925e4a3e30e41aed489a579d5595e0df13e32e1e4dd202a7c7f68b31d6418d9845eb4d757adda6ab189e1bb340db818e5b3bc725d992fa"},
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
//   - echo -n "test" | openssl dgst -hmac "key" -sha3-224
//   - echo -n "test" | openssl dgst -hmac "key" -sha3-256
//   - echo -n "test" | openssl dgst -hmac "key" -sha3-384
//   - echo -n "test" | openssl dgst -hmac "key" -sha3-512
var hmacSum = []struct {
	name     string
	sf       func(msg, key []byte) []byte
	eqf      func(msg, key, sum []byte) bool
	msg, key string // Hash
	sum      string // Hex encoded sum.
}{
	{"case01", zsha3.HMACSum224, zsha3.HMACEqualSum224, "test", "key", "c1cb62bf2208eaec3d9aa235aeed737d5b5fc2382c37c16b3d656495"},
	{"case02", zsha3.HMACSum224, zsha3.HMACEqualSum224, "", "", "1b9044e0d5bb4ef944bc00f1b26c483ac3e222f4640935d089a49083"},
	{"case03", zsha3.HMACSum224, zsha3.HMACEqualSum224, "", "key", "8f481e10aa1ab054f9862d9b2c2ec2be515ec8355e60c452eff83efc"},
	{"case04", zsha3.HMACSum224, zsha3.HMACEqualSum224, "test", "", "e79f17d1d4946b347fdeff6fd6517cfe13f347dcd0ba6f658439c08e"},
	{"case05", zsha3.HMACSum256, zsha3.HMACEqualSum256, "test", "key", "28e9ff660abb162f7653415efac4e3a0b0a40395f0e6b45fc67afbed15cbeb41"},
	{"case06", zsha3.HMACSum256, zsha3.HMACEqualSum256, "", "", "e841c164e5b4f10c9f3985587962af72fd607a951196fc92fb3a5251941784ea"},
	{"case07", zsha3.HMACSum256, zsha3.HMACEqualSum256, "", "key", "74f3c030ecc36a1835d04a333ebb7fce2688c0c78fb0bcf9592213331c884c75"},
	{"case08", zsha3.HMACSum256, zsha3.HMACEqualSum256, "test", "", "e7cedcb60ab24ac069ca1fa3f4a4757cf50903931a268da62fbb0d2ef36f6193"},
	{"case09", zsha3.HMACSum384, zsha3.HMACEqualSum384, "test", "key", "a229a41408e064854f94d3736814dd90d267df4b912adb21a886fc3e8047c1bd4fd57a6b463bea036fd87706a8986fad"},
	{"case10", zsha3.HMACSum384, zsha3.HMACEqualSum384, "", "", "adca89f07bbfbeaf58880c1572379ea2416568fd3b66542bd42599c57c4567e6ae086299ea216c6f3e7aef90b6191d24"},
	{"case11", zsha3.HMACSum384, zsha3.HMACEqualSum384, "", "key", "9139ba623c8c521d0a103bcf868041c73fa30a9e89d2a5fca9102a748be86dc15853b6b50cce3a24c008bce88182006d"},
	{"case12", zsha3.HMACSum384, zsha3.HMACEqualSum384, "test", "", "f877bfee0712e2fdedc88a308afb7bda297ec9ee4a9f4268977a480097e292a8513f2ccf46f4f9411c9c3369e7621a07"},
	{"case13", zsha3.HMACSum512, zsha3.HMACEqualSum512, "test", "key", "ddf457f9f8dcf91c459f2d36c79ad2625d86109ebf6835df5e235cb5039f7906b1e04ddbb589ffa590825c4312137fc01dd8dc29bb6b8a0d117d6953283ec2b2"},
	{"case14", zsha3.HMACSum512, zsha3.HMACEqualSum512, "", "", "cbcf45540782d4bc7387fbbf7d30b3681d6d66cc435cafd82546b0fce96b367ea79662918436fba442e81a01d0f9592dfcd30f7a7a8f1475693d30be4150ca84"},
	{"case15", zsha3.HMACSum512, zsha3.HMACEqualSum512, "", "key", "7539119b6367aa902bdc6f558d20c906d6acbd4aba3fd344eb08b0200144a1fa453ff6e7919962358be53f6db2a320d1852c52a3dea3e907070775f7a91f1282"},
	{"case16", zsha3.HMACSum512, zsha3.HMACEqualSum512, "test", "", "66a67fae72b6c23efa5f2aad2a06354c6679d19d706efa5b070e5996534748de5dc61607b46a12efeaa9c1d124f140c9c6429a5153ab2a9884a1c9a37a341143"},
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
