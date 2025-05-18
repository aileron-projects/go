package zdes_test

import (
	"crypto/des"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zdes"
	"github.com/aileron-projects/go/ztesting"
)

func TestCFB(t *testing.T) {
	t.Parallel()
	t.Run("encrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		ciphertext, err := zdes.EncryptCFB(key, nil)
		ztesting.AssertEqualErr(t, "error not match", des.KeySizeError(5), err)
		ztesting.AssertEqualSlice(t, "ciphertext is not nil", nil, ciphertext)
	})
	t.Run("decrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		plaintext, err := zdes.DecryptCFB(key, []byte("12345678"))
		ztesting.AssertEqualErr(t, "error not match", des.KeySizeError(5), err)
		ztesting.AssertEqualSlice(t, "plaintext is not nil", nil, plaintext)
	})
	t.Run("decrypt invalid ciphertext length", func(t *testing.T) {
		key := []byte("12345678")
		plaintext, err := zdes.DecryptCFB(key, []byte("short"))
		ztesting.AssertEqualErr(t, "error not match", zdes.ErrCipherLength(5), err)
		ztesting.AssertEqualSlice(t, "plaintext is not nil", nil, plaintext)
	})
	t.Run("encrypt-decrypt empty", func(t *testing.T) {
		key := []byte("12345678")
		ciphertext, err := zdes.EncryptCFB(key, nil)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zdes.DecryptCFB(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", nil, plaintext)
	})
	t.Run("encrypt-decrypt", func(t *testing.T) {
		key := []byte("12345678")
		msg := []byte("test message")
		ciphertext, err := zdes.EncryptCFB(key, msg)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zdes.DecryptCFB(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", msg, plaintext)
	})
}

func TestCFB3(t *testing.T) {
	t.Parallel()
	t.Run("encrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		ciphertext, err := zdes.EncryptCFB3(key, nil)
		ztesting.AssertEqualErr(t, "error not match", des.KeySizeError(5), err)
		ztesting.AssertEqualSlice(t, "ciphertext is not nil", nil, ciphertext)
	})
	t.Run("decrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		plaintext, err := zdes.DecryptCFB3(key, []byte("12345678"))
		ztesting.AssertEqualErr(t, "error not match", des.KeySizeError(5), err)
		ztesting.AssertEqualSlice(t, "plaintext is not nil", nil, plaintext)
	})
	t.Run("decrypt invalid ciphertext length", func(t *testing.T) {
		key := []byte("123456781234567812345678")
		plaintext, err := zdes.DecryptCFB3(key, []byte("short"))
		ztesting.AssertEqualErr(t, "error not match", zdes.ErrCipherLength(5), err)
		ztesting.AssertEqualSlice(t, "plaintext is not nil", nil, plaintext)
	})
	t.Run("encrypt-decrypt empty", func(t *testing.T) {
		key := []byte("123456781234567812345678")
		ciphertext, err := zdes.EncryptCFB3(key, nil)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zdes.DecryptCFB3(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", nil, plaintext)
	})
	t.Run("encrypt-decrypt", func(t *testing.T) {
		key := []byte("123456781234567812345678")
		msg := []byte("test message")
		ciphertext, err := zdes.EncryptCFB3(key, msg)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zdes.DecryptCFB3(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", msg, plaintext)
	})
}
