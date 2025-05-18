package zchacha20_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zchacha20"
	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

var (
	errWrongKey = errors.New("chacha20: wrong key size")
)

func TestNewStreamReader(t *testing.T) {
	t.Parallel()
	key := []byte("32 bytes secret key for chacha20")
	nonce := []byte("123456789012")
	t.Run("invalid key", func(t *testing.T) {
		_, err := zchacha20.NewStreamReader(nil, nonce, nil)
		ztesting.AssertEqualErr(t, "error not match", errWrongKey, err)
	})
	t.Run("encrypt", func(t *testing.T) {
		r, err := zchacha20.NewStreamReader(key, nonce, strings.NewReader("test"))
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		buf := make([]byte, 100)
		n, _ := r.Read(buf)
		ztesting.AssertEqual(t, "ciphertext not match", "3f416920", hex.EncodeToString(buf[:n]))
	})
	t.Run("decrypt", func(t *testing.T) {
		c, _ := hex.DecodeString("3f416920")
		r, err := zchacha20.NewStreamReader(key, nonce, bytes.NewReader(c))
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		buf := make([]byte, 100)
		n, _ := r.Read(buf)
		ztesting.AssertEqual(t, "plaintext not match", "test", string(buf[:n]))
	})
}

func TestNewStreamWriter(t *testing.T) {
	t.Parallel()
	key := []byte("32 bytes secret key for chacha20")
	nonce := []byte("123456789012")
	t.Run("invalid key", func(t *testing.T) {
		_, err := zchacha20.NewStreamWriter(nil, nonce, nil)
		ztesting.AssertEqualErr(t, "error not match", errWrongKey, err)
	})
	t.Run("encrypt", func(t *testing.T) {
		var buf bytes.Buffer
		w, err := zchacha20.NewStreamWriter(key, nonce, &buf)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		w.Write([]byte("test"))
		ztesting.AssertEqual(t, "ciphertext not match", "3f416920", hex.EncodeToString(buf.Bytes()))
	})
	t.Run("decrypt", func(t *testing.T) {
		c, _ := hex.DecodeString("3f416920")
		var buf bytes.Buffer
		w, err := zchacha20.NewStreamWriter(key, nonce, &buf)
		w.Write(c)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "plaintext not match", "test", buf.String())
	})
}

func TestCopy(t *testing.T) {
	t.Parallel()
	key := []byte("32 bytes secret key for chacha20")
	nonce := []byte("123456789012")
	t.Run("invalid key", func(t *testing.T) {
		err := zchacha20.Copy(nil, nonce, nil, nil)
		ztesting.AssertEqualErr(t, "error not match", errWrongKey, err)
	})
	t.Run("encrypt", func(t *testing.T) {
		var buf bytes.Buffer
		err := zchacha20.Copy(key, nonce, &buf, strings.NewReader("test"))
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "ciphertext not match", "3f416920", hex.EncodeToString(buf.Bytes()))
	})
	t.Run("decrypt", func(t *testing.T) {
		c, _ := hex.DecodeString("3f416920")
		var buf bytes.Buffer
		err := zchacha20.Copy(key, nonce, &buf, bytes.NewReader(c))
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "plaintext not match", "test", buf.String())
	})
	t.Run("read error", func(t *testing.T) {
		c, _ := hex.DecodeString("3f416920")
		var buf bytes.Buffer
		err := zchacha20.Copy(key, nonce, &buf, ziotest.ErrReader(bytes.NewReader(c), 3))
		ztesting.AssertEqualErr(t, "error not match", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "plaintext not match", "tes", buf.String())
	})
	t.Run("write error", func(t *testing.T) {
		c, _ := hex.DecodeString("3f416920")
		var buf bytes.Buffer
		err := zchacha20.Copy(key, nonce, ziotest.ErrWriter(&buf, 3), bytes.NewReader(c))
		ztesting.AssertEqualErr(t, "error not match", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "plaintext not match", "tes", buf.String())
	})
	t.Run("short write", func(t *testing.T) {
		c, _ := hex.DecodeString("3f416920")
		var buf bytes.Buffer
		err := zchacha20.Copy(key, nonce, ziotest.ErrWriterWith(&buf, 3, nil), bytes.NewReader(c))
		ztesting.AssertEqualErr(t, "error not match", io.ErrShortWrite, err)
		ztesting.AssertEqual(t, "plaintext not match", "tes", buf.String())
	})
}
