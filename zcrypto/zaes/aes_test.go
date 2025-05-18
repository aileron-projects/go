package zaes

import (
	"crypto/aes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/internal"
	"github.com/aileron-projects/go/ztesting"
)

func TestErrCipherLength(t *testing.T) {
	t.Parallel()
	e := ErrCipherLength(1)
	ztesting.AssertEqual(t, "message not match", "zaes: incorrect ciphertext length. got:1", e.Error())
}

func TestNewAES(t *testing.T) {
	t.Parallel()
	t.Run("invalid key", func(t *testing.T) {
		cb, iv, err := newAES([]byte("short"))
		ztesting.AssertEqualErr(t, "unexpected error", aes.KeySizeError(5), err)
		ztesting.AssertEqual(t, "non nil cipher block returned", nil, cb)
		ztesting.AssertEqualSlice(t, "non empty iv returned", nil, iv)
	})
	t.Run("16 bytes key", func(t *testing.T) {
		_, iv, err := newAES([]byte("1234567890123456"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		ztesting.AssertEqual(t, "iv length invalid", 16, len(iv))
	})
	t.Run("24 bytes key", func(t *testing.T) {
		_, iv, err := newAES([]byte("123456789012345678901234"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		ztesting.AssertEqual(t, "iv length invalid", 16, len(iv))
	})
	t.Run("32 bytes key", func(t *testing.T) {
		_, iv, err := newAES([]byte("12345678901234567890123456789012"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		ztesting.AssertEqual(t, "iv length invalid", 16, len(iv))
	})
}

// validation data is generated with:
//   - https://emn178.github.io/online-tools/aes/encrypt/
var testCases = []struct {
	name   string
	enc    internal.EncryptFunc
	dec    internal.EncryptFunc
	key    string
	plain  string // UTF-8
	cipher string // Hex encoded
}{
	{"case01", EncryptECB, DecryptECB, "1234567890123456", "", "050187a0cde5a9872cbab091ab73e553"},
	{"case02", EncryptCFB, DecryptCFB, "1234567890123456", "", "6162636465666768696a6b6c6d6e6f70"}, // iv only.
	{"case03", EncryptCTR, DecryptCTR, "1234567890123456", "", "6162636465666768696a6b6c6d6e6f70"}, // iv only.
	{"case04", EncryptOFB, DecryptOFB, "1234567890123456", "", "6162636465666768696a6b6c6d6e6f70"}, // iv only.
	{"case05", EncryptCBC, DecryptCBC, "1234567890123456", "", "6162636465666768696a6b6c6d6e6f70221b5712b4b4f99a05d61856652965f0"},
	{"case06", EncryptGCM, DecryptGCM, "1234567890123456", "", "6162636465666768696a6b6c744e626989cd58f39d0f49841b14e5b1"},
	{"case07", EncryptECB, DecryptECB, "1234567890123456", "test", "ddfbda2e0e480e5bdeb30b97ce155073"},
	{"case08", EncryptCFB, DecryptCFB, "1234567890123456", "test", "6162636465666768696a6b6c6d6e6f7088c8022f"},
	{"case09", EncryptCTR, DecryptCTR, "1234567890123456", "test", "6162636465666768696a6b6c6d6e6f7088c8022f"},
	{"case10", EncryptOFB, DecryptOFB, "1234567890123456", "test", "6162636465666768696a6b6c6d6e6f7088c8022f"},
	{"case11", EncryptCBC, DecryptCBC, "1234567890123456", "test", "6162636465666768696a6b6c6d6e6f7047044d19c39025966715ce8e6c8bf2e9"},
	{"case12", EncryptGCM, DecryptGCM, "1234567890123456", "test", "6162636465666768696a6b6cd47e67d8519b0ac61061e03dbf6ae74b7022d16c"},
	{"case13", EncryptECB, DecryptECB, "1234567890123456", "test longer message", "d4293ee4f789456df439056ef2b4dd509d942dff49f6860590ff61671a15d7cb"},
	{"case14", EncryptCFB, DecryptCFB, "1234567890123456", "test longer message", "6162636465666768696a6b6c6d6e6f7088c8022ff75733de2feaf62f56c80bfa76776a"},
	{"case15", EncryptCTR, DecryptCTR, "1234567890123456", "test longer message", "6162636465666768696a6b6c6d6e6f7088c8022ff75733de2feaf62f56c80bfa68c56a"},
	{"case16", EncryptOFB, DecryptOFB, "1234567890123456", "test longer message", "6162636465666768696a6b6c6d6e6f7088c8022ff75733de2feaf62f56c80bfa8787e2"},
	{"case17", EncryptCBC, DecryptCBC, "1234567890123456", "test longer message", "6162636465666768696a6b6c6d6e6f702e73120a0f66746a0abbd90aca464690cef1a42600bf73c42a88b960adf2b2e3"},
	{"case18", EncryptGCM, DecryptGCM, "1234567890123456", "test longer message", "6162636465666768696a6b6cd47e67d832c6010ccf3638de775609ce4096af5406c8a186cbfbcd9d48e1a615f313e1"},
}

func TestEncryptDecrypt(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_encrypt", func(t *testing.T) {
			done := ztesting.ReplaceRandReader(strings.NewReader("abcdefghijklmnop"))
			defer done()
			ciphertext, err := tc.enc([]byte(tc.key), []byte(tc.plain))
			ztesting.AssertEqualErr(t, "error is not nil", nil, err)
			ztesting.AssertEqual(t, "ciphertext not match", tc.cipher, hex.EncodeToString(ciphertext))
		})
		t.Run(tc.name+"_decrypt", func(t *testing.T) {
			cipher, _ := hex.DecodeString(tc.cipher)
			plain, err := tc.dec([]byte(tc.key), cipher)
			ztesting.AssertEqualErr(t, "error is not nil", nil, err)
			ztesting.AssertEqual(t, "plaintext not match", tc.plain, string(plain))
		})
	}
}
