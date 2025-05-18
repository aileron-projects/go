package internal_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/internal"
	"github.com/aileron-projects/go/zcrypto/zaes"
	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestCopy(t *testing.T) {
	t.Parallel()
	t.Run("copy success", func(t *testing.T) {
		key := []byte("1234567890123456")
		iv := []byte("1234567890123456")
		c, _ := aes.NewCipher(key)
		s := cipher.NewCTR(c, iv)
		var w, ww bytes.Buffer
		err := internal.Copy(s, &w, strings.NewReader("test"))
		ztesting.AssertEqualErr(t, "error is not nil", nil, err)
		ztesting.AssertEqual(t, "message unexpectedly match", false, w.String() == "test")
		err = zaes.CopyOFB(key, iv, &ww, strings.NewReader(w.String()))
		ztesting.AssertEqualErr(t, "error is not nil", nil, err)
		ztesting.AssertEqual(t, "message not match", "test", ww.String())
	})
	t.Run("read error", func(t *testing.T) {
		key := []byte("1234567890123456")
		iv := []byte("1234567890123456")
		c, _ := aes.NewCipher(key)
		s := cipher.NewCTR(c, iv)
		var w bytes.Buffer
		err := internal.Copy(s, &w, ziotest.ErrReader(strings.NewReader("test"), 3))
		ztesting.AssertEqualErr(t, "error not match", io.ErrClosedPipe, err)
	})
	t.Run("write error", func(t *testing.T) {
		key := []byte("1234567890123456")
		iv := []byte("1234567890123456")
		c, _ := aes.NewCipher(key)
		s := cipher.NewCTR(c, iv)
		var w bytes.Buffer
		err := internal.Copy(s, ziotest.ErrWriter(&w, 3), strings.NewReader("test"))
		ztesting.AssertEqualErr(t, "error not match", io.ErrClosedPipe, err)
	})
	t.Run("short write error", func(t *testing.T) {
		key := []byte("1234567890123456")
		iv := []byte("1234567890123456")
		c, _ := aes.NewCipher(key)
		s := cipher.NewCTR(c, iv)
		var w bytes.Buffer
		err := internal.Copy(s, ziotest.ShortWriter(&w, 3), strings.NewReader("test"))
		ztesting.AssertEqualErr(t, "error not match", io.ErrShortWrite, err)
	})
}
