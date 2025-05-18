package zrc4_test

import (
	"bytes"
	"crypto/rc4"
	"encoding/hex"
	"io"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zrc4"
	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestNewStreamReader(t *testing.T) {
	t.Parallel()
	// Validation data is generated with
	// 	- https://emn178.github.io/online-tools/rc4/encrypt/
	key := []byte("secret-key")
	t.Run("invalid key", func(t *testing.T) {
		_, err := zrc4.NewStreamReader(nil, nil)
		ztesting.AssertEqualErr(t, "error not match", rc4.KeySizeError(0), err)
	})
	t.Run("encrypt", func(t *testing.T) {
		r, err := zrc4.NewStreamReader(key, strings.NewReader("test"))
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		buf := make([]byte, 100)
		n, _ := r.Read(buf)
		ztesting.AssertEqual(t, "ciphertext not match", "c46d3def", hex.EncodeToString(buf[:n]))
	})
	t.Run("decrypt", func(t *testing.T) {
		c, _ := hex.DecodeString("c46d3def")
		r, err := zrc4.NewStreamReader(key, bytes.NewReader(c))
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		buf := make([]byte, 100)
		n, _ := r.Read(buf)
		ztesting.AssertEqual(t, "plaintext not match", "test", string(buf[:n]))
	})
}

func TestNewStreamWriter(t *testing.T) {
	t.Parallel()
	// Validation data is generated with
	// 	- https://emn178.github.io/online-tools/rc4/encrypt/
	key := []byte("secret-key")
	t.Run("invalid key", func(t *testing.T) {
		_, err := zrc4.NewStreamWriter(nil, nil)
		ztesting.AssertEqualErr(t, "error not match", rc4.KeySizeError(0), err)
	})
	t.Run("encrypt", func(t *testing.T) {
		var buf bytes.Buffer
		w, err := zrc4.NewStreamWriter(key, &buf)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		w.Write([]byte("test"))
		ztesting.AssertEqual(t, "ciphertext not match", "c46d3def", hex.EncodeToString(buf.Bytes()))
	})
	t.Run("decrypt", func(t *testing.T) {
		c, _ := hex.DecodeString("c46d3def")
		var buf bytes.Buffer
		w, err := zrc4.NewStreamWriter(key, &buf)
		w.Write(c)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "plaintext not match", "test", buf.String())
	})
}

func TestCopy(t *testing.T) {
	t.Parallel()
	key := []byte("secret-key")
	t.Run("invalid key", func(t *testing.T) {
		err := zrc4.Copy(nil, nil, nil)
		ztesting.AssertEqualErr(t, "error not match", rc4.KeySizeError(0), err)
	})
	t.Run("encrypt", func(t *testing.T) {
		var buf bytes.Buffer
		err := zrc4.Copy(key, &buf, strings.NewReader("test"))
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "ciphertext not match", "c46d3def", hex.EncodeToString(buf.Bytes()))
	})
	t.Run("decrypt", func(t *testing.T) {
		c, _ := hex.DecodeString("c46d3def")
		var buf bytes.Buffer
		err := zrc4.Copy(key, &buf, bytes.NewReader(c))
		ztesting.AssertEqualErr(t, "error not match", nil, err)
		ztesting.AssertEqual(t, "plaintext not match", "test", buf.String())
	})
	t.Run("read error", func(t *testing.T) {
		c, _ := hex.DecodeString("c46d3def")
		var buf bytes.Buffer
		err := zrc4.Copy(key, &buf, ziotest.ErrReader(bytes.NewReader(c), 3))
		ztesting.AssertEqualErr(t, "error not match", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "plaintext not match", "tes", buf.String())
	})
	t.Run("write error", func(t *testing.T) {
		c, _ := hex.DecodeString("c46d3def")
		var buf bytes.Buffer
		err := zrc4.Copy(key, ziotest.ErrWriter(&buf, 3), bytes.NewReader(c))
		ztesting.AssertEqualErr(t, "error not match", io.ErrClosedPipe, err)
		ztesting.AssertEqual(t, "plaintext not match", "tes", buf.String())
	})
	t.Run("short write", func(t *testing.T) {
		c, _ := hex.DecodeString("c46d3def")
		var buf bytes.Buffer
		err := zrc4.Copy(key, ziotest.ErrWriterWith(&buf, 3, nil), bytes.NewReader(c))
		ztesting.AssertEqualErr(t, "error not match", io.ErrShortWrite, err)
		ztesting.AssertEqual(t, "plaintext not match", "tes", buf.String())
	})
}
