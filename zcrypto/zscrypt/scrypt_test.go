package zscrypt_test

import (
	"encoding/hex"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zscrypt"
	"github.com/aileron-projects/go/ztesting"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("12345678901234567890"))
		defer done()
		b, err := zscrypt.New(10, 32768, 8, 1, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "hash not match", "31323334353637383930dac13eff06858ad01e3f90af26b7a3055233ebdfaa90674bada2a0fd0e3a316a", hex.EncodeToString(h))
	})
	t.Run("param invalid", func(t *testing.T) {
		b, err := zscrypt.New(10, 32768, math.MaxInt, 1, 32)
		ztesting.AssertEqual(t, "non nil hasher returned", nil, b)
		ztesting.AssertEqualErr(t, "nil error returned", errors.New("scrypt: parameters are too large"), err)
	})
}

func TestSplit(t *testing.T) {
	t.Run("short hash", func(t *testing.T) {
		b, err := zscrypt.New(32, 32768, 8, 1, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		_, _, err = b.Split([]byte("short hash"))
		ztesting.AssertEqualErr(t, "error not returned", zscrypt.ErrShortHash, err)
	})
	t.Run("split hash", func(t *testing.T) {
		b, err := zscrypt.New(10, 32768, 8, 1, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		salt, hash, err := b.Split([]byte("1234567890abcdefghijklmnopqrstuvwxyz"))
		ztesting.AssertEqualErr(t, "non-nil error returned", nil, err)
		ztesting.AssertEqual(t, "incorrect salt", "1234567890", string(salt))
		ztesting.AssertEqual(t, "incorrect hash", "abcdefghijklmnopqrstuvwxyz", string(hash))
	})
}

func TestSum(t *testing.T) {
	// Results can be compared and validated with
	// 	- https://8gwifi.org/scrypt.jsp
	//	- https://gchq.github.io/CyberChef/#recipe=Scrypt(%7B'option':'UTF8','string':'1234567890'%7D,32768,8,1,32)&input=dGVzdA
	t.Run("success", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("12345678901234567890"))
		defer done()
		b, err := zscrypt.New(10, 32768, 8, 1, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, err := b.Sum([]byte("test"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		ztesting.AssertEqual(t, "hash not match", "31323334353637383930dac13eff06858ad01e3f90af26b7a3055233ebdfaa90674bada2a0fd0e3a316a", hex.EncodeToString(h))
	})
}

func TestEqual(t *testing.T) {
	t.Run("", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("1234567890123456789012345678901234567890"))
		defer done()
		b, err := zscrypt.New(10, 32768, 8, 1, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "incorrect equal result", true, b.Equal(h, []byte("test")))
		ztesting.AssertEqual(t, "incorrect equal result", false, b.Equal(h, []byte("wrong")))
	})
}

func TestCompare(t *testing.T) {
	t.Run("short hash", func(t *testing.T) {
		b, err := zscrypt.New(32, 32768, 8, 1, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		err = b.Compare([]byte("short hash"), []byte("test"))
		ztesting.AssertEqualErr(t, "error not returned", zscrypt.ErrShortHash, err)
	})
	t.Run("key generate error", func(t *testing.T) {
		b := &zscrypt.SCrypt{}
		err := b.Compare([]byte("short hash"), []byte("test"))
		ztesting.AssertEqualErr(t, "error not returned", errors.New("scrypt: N must be > 1 and a power of 2"), err)
	})
	t.Run("pw not match", func(t *testing.T) {
		b, err := zscrypt.New(32, 32768, 8, 1, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, err := b.Sum([]byte("test"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		err = b.Compare(h, []byte("wrong"))
		ztesting.AssertEqualErr(t, "error not matched", zscrypt.ErrNotMatch, err)
	})
	t.Run("pw match", func(t *testing.T) {
		b, err := zscrypt.New(32, 32768, 8, 1, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, err := b.Sum([]byte("test"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		err = b.Compare(h, []byte("test"))
		ztesting.AssertEqualErr(t, "error not matched", nil, err)
	})
}
