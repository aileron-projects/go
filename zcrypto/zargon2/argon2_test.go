package zargon2_test

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zargon2"
	"github.com/aileron-projects/go/ztesting"
)

func TestNewArgon2i(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("12345678901234567890"))
		defer done()
		b, err := zargon2.NewArgon2i(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "hash not match", "31323334353637383930b73fc562d7f71262bb4099efad0d836c5e0496788c840c33cef2d90ef992e060", hex.EncodeToString(h))
	})
	t.Run("invalid param", func(t *testing.T) {
		b, err := zargon2.NewArgon2i(10, 0, 64*1024, 4, 32)
		ztesting.AssertEqual(t, "non nil hasher returned", nil, b)
		ztesting.AssertEqualErr(t, "nil error returned", errors.New("zargon2: invalid hash parameter [argon2: number of rounds too small]"), err)
	})
}

func TestNewArgon2id(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("12345678901234567890"))
		defer done()
		b, err := zargon2.NewArgon2id(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "hash not match", "3132333435363738393046a3090e30874ea4dd680f04ee9bbb3ab0a01421de4fbe845649db192db746b3", hex.EncodeToString(h))
	})
	t.Run("invalid param", func(t *testing.T) {
		b, err := zargon2.NewArgon2id(10, 0, 64*1024, 4, 32)
		ztesting.AssertEqual(t, "non nil hasher returned", nil, b)
		ztesting.AssertEqualErr(t, "nil error returned", errors.New("zargon2: invalid hash parameter [argon2: number of rounds too small]"), err)
	})
}

func TestArgon2i_Split(t *testing.T) {
	t.Run("short hash", func(t *testing.T) {
		b, err := zargon2.NewArgon2i(32, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		_, _, err = b.Split([]byte("short hash"))
		ztesting.AssertEqualErr(t, "error not returned", zargon2.ErrShortHash, err)
	})
	t.Run("split hash", func(t *testing.T) {
		b, err := zargon2.NewArgon2i(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		salt, hash, err := b.Split([]byte("1234567890abcdefghijklmnopqrstuvwxyz"))
		ztesting.AssertEqualErr(t, "non-nil error returned", nil, err)
		ztesting.AssertEqual(t, "incorrect salt", "1234567890", string(salt))
		ztesting.AssertEqual(t, "incorrect hash", "abcdefghijklmnopqrstuvwxyz", string(hash))
	})
}

func TestArgon2id_Split(t *testing.T) {
	t.Run("short hash", func(t *testing.T) {
		b, err := zargon2.NewArgon2id(32, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		_, _, err = b.Split([]byte("short hash"))
		ztesting.AssertEqualErr(t, "error not returned", zargon2.ErrShortHash, err)
	})
	t.Run("split hash", func(t *testing.T) {
		b, err := zargon2.NewArgon2id(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		salt, hash, err := b.Split([]byte("1234567890abcdefghijklmnopqrstuvwxyz"))
		ztesting.AssertEqualErr(t, "non-nil error returned", nil, err)
		ztesting.AssertEqual(t, "incorrect salt", "1234567890", string(salt))
		ztesting.AssertEqual(t, "incorrect hash", "abcdefghijklmnopqrstuvwxyz", string(hash))
	})
}

func TestArgon2i_Sum(t *testing.T) {
	// Results can be compared and validated with
	// 	- https://argon2.online/
	t.Run("success", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("12345678901234567890"))
		defer done()
		b, err := zargon2.NewArgon2i(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "hash not match", "31323334353637383930b73fc562d7f71262bb4099efad0d836c5e0496788c840c33cef2d90ef992e060", hex.EncodeToString(h))
	})
}

func TestArgon2id_Sum(t *testing.T) {
	// Results can be compared and validated with
	// 	- https://argon2.online/
	t.Run("success", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("12345678901234567890"))
		defer done()
		b, err := zargon2.NewArgon2id(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "hash not match", "3132333435363738393046a3090e30874ea4dd680f04ee9bbb3ab0a01421de4fbe845649db192db746b3", hex.EncodeToString(h))
	})
}

func TestArgon2i_Equal(t *testing.T) {
	t.Run("", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("12345678901234567890"))
		defer done()
		b, err := zargon2.NewArgon2i(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "incorrect equal result", true, b.Equal(h, []byte("test")))
		ztesting.AssertEqual(t, "incorrect equal result", false, b.Equal(h, []byte("wrong")))
	})
}

func TestArgon2id_Equal(t *testing.T) {
	t.Run("", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("12345678901234567890"))
		defer done()
		b, err := zargon2.NewArgon2id(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "incorrect equal result", true, b.Equal(h, []byte("test")))
		ztesting.AssertEqual(t, "incorrect equal result", false, b.Equal(h, []byte("wrong")))
	})
}

func TestArgon2i_Compare(t *testing.T) {
	t.Run("short hash", func(t *testing.T) {
		b, err := zargon2.NewArgon2i(32, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		err = b.Compare([]byte("short hash"), []byte("test"))
		ztesting.AssertEqualErr(t, "error not returned", zargon2.ErrShortHash, err)
	})
	t.Run("pw not match", func(t *testing.T) {
		b, err := zargon2.NewArgon2i(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, err := b.Sum([]byte("test"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		err = b.Compare(h, []byte("wrong"))
		ztesting.AssertEqualErr(t, "error not matched", zargon2.ErrNotMatch, err)
	})
	t.Run("pw match", func(t *testing.T) {
		b, err := zargon2.NewArgon2i(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, err := b.Sum([]byte("test"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		err = b.Compare(h, []byte("test"))
		ztesting.AssertEqualErr(t, "error not matched", nil, err)
	})
}

func TestArgon2id_Compare(t *testing.T) {
	t.Run("short hash", func(t *testing.T) {
		b, err := zargon2.NewArgon2id(32, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		err = b.Compare([]byte("short hash"), []byte("test"))
		ztesting.AssertEqualErr(t, "error not returned", zargon2.ErrShortHash, err)
	})
	t.Run("pw not match", func(t *testing.T) {
		b, err := zargon2.NewArgon2id(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, err := b.Sum([]byte("test"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		err = b.Compare(h, []byte("wrong"))
		ztesting.AssertEqualErr(t, "error not matched", zargon2.ErrNotMatch, err)
	})
	t.Run("pw match", func(t *testing.T) {
		b, err := zargon2.NewArgon2id(10, 1, 64*1024, 4, 32)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, err := b.Sum([]byte("test"))
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		err = b.Compare(h, []byte("test"))
		ztesting.AssertEqualErr(t, "error not matched", nil, err)
	})
}
