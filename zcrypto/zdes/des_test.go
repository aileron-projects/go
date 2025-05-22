package zdes

import (
	"crypto/des"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/internal"
	"github.com/aileron-projects/go/ztesting"
)

func TestErrCipherLength(t *testing.T) {
	t.Parallel()
	e := ErrCipherLength(1)
	ztesting.AssertEqual(t, "message not match", "zdes: incorrect ciphertext length. got:1", e.Error())
}

func TestNewDES(t *testing.T) {
	t.Parallel()
	t.Run("invalid key", func(t *testing.T) {
		cb, iv, err := newDES([]byte("short"))
		ztesting.AssertEqualErr(t, "unexpected error", des.KeySizeError(5), err)
		ztesting.AssertEqual(t, "non nil cipher block returned", nil, cb)
		ztesting.AssertEqual(t, "non empty iv returned", nil, iv)
	})
	t.Run("valid key", func(t *testing.T) {
		_, iv, err := newDES([]byte("12345678"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		ztesting.AssertEqual(t, "iv length invalid", 8, len(iv))
	})
}

func TestNew3DES(t *testing.T) {
	t.Parallel()
	t.Run("invalid key", func(t *testing.T) {
		cb, iv, err := new3DES([]byte("short"))
		ztesting.AssertEqualErr(t, "unexpected error", des.KeySizeError(5), err)
		ztesting.AssertEqual(t, "non nil cipher block returned", nil, cb)
		ztesting.AssertEqual(t, "non empty iv returned", nil, iv)
	})
	t.Run("valid key", func(t *testing.T) {
		_, iv, err := new3DES([]byte("123456781234567812345678"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		ztesting.AssertEqual(t, "iv length invalid", 8, len(iv))
	})
}

// validation data is generated with:
//   - https://emn178.github.io/online-tools/des/encrypt/
//   - https://emn178.github.io/online-tools/triple-des/encrypt/
var testCases = []struct {
	name   string
	enc    internal.EncryptFunc
	dec    internal.EncryptFunc
	key    string
	plain  string // UTF-8
	cipher string // Hex encoded
}{
	{"case01", EncryptECB, DecryptECB, "12345678", "", "feb959b7d4642fcb"},
	{"case02", EncryptCFB, DecryptCFB, "12345678", "", "6162636431323334"}, // iv only.
	{"case03", EncryptCTR, DecryptCTR, "12345678", "", "6162636431323334"}, // iv only.
	{"case04", EncryptOFB, DecryptOFB, "12345678", "", "6162636431323334"}, // iv only.
	{"case05", EncryptCBC, DecryptCBC, "12345678", "", "61626364313233346ac4b68ea609ad4a"},
	{"case06", EncryptECB3, DecryptECB3, "123456789012345678901234", "", "9335cc2fc785c26b"},
	{"case07", EncryptCFB3, DecryptCFB3, "123456789012345678901234", "", "6162636431323334"}, // iv only
	{"case08", EncryptCTR3, DecryptCTR3, "123456789012345678901234", "", "6162636431323334"}, // iv only
	{"case09", EncryptOFB3, DecryptOFB3, "123456789012345678901234", "", "6162636431323334"}, // iv only
	{"case10", EncryptCBC3, DecryptCBC3, "123456789012345678901234", "", "61626364313233348520af2b86686531"},
	{"case11", EncryptECB, DecryptECB, "12345678", "test", "a04b686b118af67b"},
	{"case12", EncryptCFB, DecryptCFB, "12345678", "test", "6162636431323334c612c982"},
	{"case13", EncryptCTR, DecryptCTR, "12345678", "test", "6162636431323334c612c982"},
	{"case14", EncryptOFB, DecryptOFB, "12345678", "test", "6162636431323334c612c982"},
	{"case15", EncryptCBC, DecryptCBC, "12345678", "test", "6162636431323334cc6f189c03a31e87"},
	{"case16", EncryptECB3, DecryptECB3, "123456789012345678901234", "test", "362d818a1de70945"},
	{"case17", EncryptCFB3, DecryptCFB3, "123456789012345678901234", "test", "616263643132333456840595"},
	{"case18", EncryptCTR3, DecryptCTR3, "123456789012345678901234", "test", "616263643132333456840595"},
	{"case19", EncryptOFB3, DecryptOFB3, "123456789012345678901234", "test", "616263643132333456840595"},
	{"case20", EncryptCBC3, DecryptCBC3, "123456789012345678901234", "test", "61626364313233349ad357ed51b1d3c7"},
	{"case21", EncryptECB, DecryptECB, "12345678", "test longer message", "14880cc1c4aade8c0c89fdc8e47904e01f6d080d85bdc006"},
	{"case22", EncryptCFB, DecryptCFB, "12345678", "test longer message", "6162636431323334c612c9829a3d662bcf6ca69a90fd26a9cc54ca"},
	{"case23", EncryptCTR, DecryptCTR, "12345678", "test longer message", "6162636431323334c612c9829a3d662bf6d72490c441eb5f51e642"},
	{"case24", EncryptOFB, DecryptOFB, "12345678", "test longer message", "6162636431323334c612c9829a3d662bb6ae5e8e68387108a71731"},
	{"case25", EncryptCBC, DecryptCBC, "12345678", "test longer message", "6162636431323334811d81698190d4219cf12a213e011148fcd334bf55f5b3a4"},
	{"case26", EncryptECB3, DecryptECB3, "123456789012345678901234", "test longer message", "6bc9c110ecd2117b1c5d729f8d4debf29d0ff1e24811702d"},
	{"case27", EncryptCFB3, DecryptCFB3, "123456789012345678901234", "test longer message", "616263643132333456840595bd4569092cce00cc9aa10fb388c6d6"},
	{"case28", EncryptCTR3, DecryptCTR3, "123456789012345678901234", "test longer message", "616263643132333456840595bd4569096f220ae4e06e75b7eee73b"},
	{"case29", EncryptOFB3, DecryptOFB3, "123456789012345678901234", "test longer message", "616263643132333456840595bd4569098c3117ee4be67429e7be5f"},
	{"case30", EncryptCBC3, DecryptCBC3, "123456789012345678901234", "test longer message", "61626364313233347e0cd96dd3c76058b5a785120f4c558f14d84898402d8f69"},
}

func TestEncryptDecrypt(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name+"_encrypt", func(t *testing.T) {
			done := ztesting.ReplaceRandReader(strings.NewReader("abcd1234"))
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
