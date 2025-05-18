package zaes_test

import (
	"bytes"
	"crypto/aes"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zaes"
	"github.com/aileron-projects/go/ztesting"
)

func TestOFB(t *testing.T) {
	t.Parallel()
	t.Run("encrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		ciphertext, err := zaes.EncryptOFB(key, nil)
		ztesting.AssertEqualErr(t, "error not match", aes.KeySizeError(5), err)
		ztesting.AssertEqualSlice(t, "ciphertext is not nil", nil, ciphertext)
	})
	t.Run("decrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		plaintext, err := zaes.DecryptOFB(key, []byte("1234567890123456"))
		ztesting.AssertEqualErr(t, "error not match", aes.KeySizeError(5), err)
		ztesting.AssertEqualSlice(t, "plaintext is not nil", nil, plaintext)
	})
	t.Run("decrypt invalid ciphertext length", func(t *testing.T) {
		key := []byte("1234567890123456")
		plaintext, err := zaes.DecryptOFB(key, []byte("short"))
		ztesting.AssertEqualErr(t, "error not match", zaes.ErrCipherLength(5), err)
		ztesting.AssertEqualSlice(t, "plaintext is not nil", nil, plaintext)
	})
	t.Run("AES128: encrypt-decrypt empty", func(t *testing.T) {
		key := []byte("1234567890123456")
		ciphertext, err := zaes.EncryptOFB(key, nil)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zaes.DecryptOFB(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", nil, plaintext)
	})
	t.Run("AES128: encrypt-decrypt", func(t *testing.T) {
		key := []byte("1234567890123456")
		msg := []byte("test message")
		ciphertext, err := zaes.EncryptOFB(key, msg)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zaes.DecryptOFB(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", msg, plaintext)
	})
	t.Run("AES192: encrypt-decrypt empty", func(t *testing.T) {
		key := []byte("123456789012345678901234")
		ciphertext, err := zaes.EncryptOFB(key, nil)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zaes.DecryptOFB(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", nil, plaintext)
	})
	t.Run("AES192: encrypt-decrypt", func(t *testing.T) {
		key := []byte("123456789012345678901234")
		msg := []byte("test message")
		ciphertext, err := zaes.EncryptOFB(key, msg)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zaes.DecryptOFB(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", msg, plaintext)
	})
	t.Run("AES256: encrypt-decrypt empty", func(t *testing.T) {
		key := []byte("12345678901234567890123456789012")
		ciphertext, err := zaes.EncryptOFB(key, nil)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zaes.DecryptOFB(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", nil, plaintext)
	})
	t.Run("AES256: encrypt-decrypt", func(t *testing.T) {
		key := []byte("12345678901234567890123456789012")
		msg := []byte("test message")
		ciphertext, err := zaes.EncryptOFB(key, msg)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zaes.DecryptOFB(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", msg, plaintext)
	})
}

func TestCopyOFB(t *testing.T) {
	t.Parallel()
	t.Run("encrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		iv := []byte("1234567890123456")
		err := zaes.CopyOFB(key, iv, nil, nil)
		ztesting.AssertEqualErr(t, "error not match", aes.KeySizeError(5), err)
	})
	t.Run("encrypt-decrypt", func(t *testing.T) {
		key := []byte("1234567890123456")
		iv := []byte("1234567890123456")
		msg := "test message"
		var w, ww bytes.Buffer
		err := zaes.CopyOFB(key, iv, &w, strings.NewReader(msg))
		ztesting.AssertEqualErr(t, "error is not nil", nil, err)
		ztesting.AssertEqual(t, "length not match", len(msg), w.Len())
		ztesting.AssertEqual(t, "message unexpectedly match", false, msg == w.String())
		err = zaes.CopyOFB(key, iv, &ww, strings.NewReader(w.String()))
		ztesting.AssertEqualErr(t, "error is not nil", nil, err)
		ztesting.AssertEqual(t, "message not match", msg, ww.String())
	})
}
