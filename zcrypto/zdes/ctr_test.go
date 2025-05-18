package zdes_test

import (
	"bytes"
	"crypto/des"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zdes"
	"github.com/aileron-projects/go/ztesting"
)

func TestCTR(t *testing.T) {
	t.Parallel()
	t.Run("encrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		ciphertext, err := zdes.EncryptCTR(key, nil)
		ztesting.AssertEqualErr(t, "error not match", des.KeySizeError(5), err)
		ztesting.AssertEqualSlice(t, "ciphertext is not nil", nil, ciphertext)
	})
	t.Run("decrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		plaintext, err := zdes.DecryptCTR(key, []byte("12345678"))
		ztesting.AssertEqualErr(t, "error not match", des.KeySizeError(5), err)
		ztesting.AssertEqualSlice(t, "plaintext is not nil", nil, plaintext)
	})
	t.Run("decrypt invalid ciphertext length", func(t *testing.T) {
		key := []byte("12345678")
		plaintext, err := zdes.DecryptCTR(key, []byte("short"))
		ztesting.AssertEqualErr(t, "error not match", zdes.ErrCipherLength(5), err)
		ztesting.AssertEqualSlice(t, "plaintext is not nil", nil, plaintext)
	})
	t.Run("encrypt-decrypt empty", func(t *testing.T) {
		key := []byte("12345678")
		ciphertext, err := zdes.EncryptCTR(key, nil)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zdes.DecryptCTR(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", nil, plaintext)
	})
	t.Run("encrypt-decrypt", func(t *testing.T) {
		key := []byte("12345678")
		msg := []byte("test message")
		ciphertext, err := zdes.EncryptCTR(key, msg)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zdes.DecryptCTR(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", msg, plaintext)
	})
}

func TestCTR3(t *testing.T) {
	t.Parallel()
	t.Run("encrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		ciphertext, err := zdes.EncryptCTR3(key, nil)
		ztesting.AssertEqualErr(t, "error not match", des.KeySizeError(5), err)
		ztesting.AssertEqualSlice(t, "ciphertext is not nil", nil, ciphertext)
	})
	t.Run("decrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		plaintext, err := zdes.DecryptCTR3(key, []byte("12345678"))
		ztesting.AssertEqualErr(t, "error not match", des.KeySizeError(5), err)
		ztesting.AssertEqualSlice(t, "plaintext is not nil", nil, plaintext)
	})
	t.Run("decrypt invalid ciphertext length", func(t *testing.T) {
		key := []byte("123456781234567812345678")
		plaintext, err := zdes.DecryptCTR3(key, []byte("short"))
		ztesting.AssertEqualErr(t, "error not match", zdes.ErrCipherLength(5), err)
		ztesting.AssertEqualSlice(t, "plaintext is not nil", nil, plaintext)
	})
	t.Run("encrypt-decrypt empty", func(t *testing.T) {
		key := []byte("123456781234567812345678")
		ciphertext, err := zdes.EncryptCTR3(key, nil)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zdes.DecryptCTR3(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", nil, plaintext)
	})
	t.Run("encrypt-decrypt", func(t *testing.T) {
		key := []byte("123456781234567812345678")
		msg := []byte("test message")
		ciphertext, err := zdes.EncryptCTR3(key, msg)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		plaintext, err := zdes.DecryptCTR3(key, ciphertext)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqualSlice(t, "plaintext does not match", msg, plaintext)
	})
}

func TestCopyCTR(t *testing.T) {
	t.Parallel()
	t.Run("encrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		iv := []byte("12345678")
		err := zdes.CopyCTR(key, iv, nil, nil)
		ztesting.AssertEqualErr(t, "error not match", des.KeySizeError(5), err)
	})
	t.Run("encrypt-decrypt", func(t *testing.T) {
		key := []byte("12345678")
		iv := []byte("12345678")
		msg := "test message"
		var w, ww bytes.Buffer
		err := zdes.CopyCTR(key, iv, &w, strings.NewReader(msg))
		ztesting.AssertEqualErr(t, "error is not nil", nil, err)
		ztesting.AssertEqual(t, "length not match", len(msg), w.Len())
		ztesting.AssertEqual(t, "message unexpectedly match", false, msg == w.String())
		err = zdes.CopyCTR(key, iv, &ww, strings.NewReader(w.String()))
		ztesting.AssertEqualErr(t, "error is not nil", nil, err)
		ztesting.AssertEqual(t, "message not match", msg, ww.String())
	})
}

func TestCopyCTR3(t *testing.T) {
	t.Parallel()
	t.Run("encrypt key invalid", func(t *testing.T) {
		key := []byte("short")
		iv := []byte("12345678")
		err := zdes.CopyCTR3(key, iv, nil, nil)
		ztesting.AssertEqualErr(t, "error not match", des.KeySizeError(5), err)
	})
	t.Run("encrypt-decrypt", func(t *testing.T) {
		key := []byte("123456781234567812345678")
		iv := []byte("12345678")
		msg := "test message"
		var w, ww bytes.Buffer
		err := zdes.CopyCTR3(key, iv, &w, strings.NewReader(msg))
		ztesting.AssertEqualErr(t, "error is not nil", nil, err)
		ztesting.AssertEqual(t, "length not match", len(msg), w.Len())
		ztesting.AssertEqual(t, "message unexpectedly match", false, msg == w.String())
		err = zdes.CopyCTR3(key, iv, &ww, strings.NewReader(w.String()))
		ztesting.AssertEqualErr(t, "error is not nil", nil, err)
		ztesting.AssertEqual(t, "message not match", msg, ww.String())
	})
}
