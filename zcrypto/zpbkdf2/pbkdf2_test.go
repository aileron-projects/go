package zpbkdf2_test

import (
	"crypto"
	_ "crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zpbkdf2"
	"github.com/aileron-projects/go/ztesting"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("12345678901234567890"))
		defer done()
		b, err := zpbkdf2.New(10, 210000, 32, crypto.SHA256)
		ztesting.AssertEqualErr(t, nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "313233343536373839306d94d84234ee7201fb39aff04ffd9d7bbc5339bb9553a3e47146c45bf0aece7d", hex.EncodeToString(h))
	})
	t.Run("hash not available", func(t *testing.T) {
		b, err := zpbkdf2.New(10, 210000, 32, crypto.Hash(999))
		ztesting.AssertEqual(t, nil, b)
		ztesting.AssertEqualErr(t, zpbkdf2.ErrUnavailable, err)
	})
}

func TestSplit(t *testing.T) {
	t.Run("short hash", func(t *testing.T) {
		b, err := zpbkdf2.New(32, 210000, 32, crypto.SHA256)
		ztesting.AssertEqualErr(t, nil, err)
		_, _, err = b.Split([]byte("short hash"))
		ztesting.AssertEqualErr(t, zpbkdf2.ErrShortHash, err)
	})
	t.Run("split hash", func(t *testing.T) {
		b, err := zpbkdf2.New(10, 210000, 32, crypto.SHA256)
		ztesting.AssertEqualErr(t, nil, err)
		salt, hash, err := b.Split([]byte("1234567890abcdefghijklmnopqrstuvwxyz"))
		ztesting.AssertEqualErr(t, nil, err)
		ztesting.AssertEqual(t, "1234567890", string(salt))
		ztesting.AssertEqual(t, "abcdefghijklmnopqrstuvwxyz", string(hash))
	})
}

func TestSum(t *testing.T) {
	// Results can be compared and validated with
	// 	- https://bcrypt-tools.vercel.app/pbkdf2
	//	- https://neurotechnics.com/tools/pbkdf2-test
	t.Run("success", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("12345678901234567890"))
		defer done()
		b, err := zpbkdf2.New(10, 210000, 32, crypto.SHA256)
		ztesting.AssertEqualErr(t, nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "313233343536373839306d94d84234ee7201fb39aff04ffd9d7bbc5339bb9553a3e47146c45bf0aece7d", hex.EncodeToString(h))
	})
}

func TestEqual(t *testing.T) {
	t.Run("", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("1234567890123456789012345678901234567890"))
		defer done()
		b, err := zpbkdf2.New(10, 210000, 32, crypto.SHA256)
		ztesting.AssertEqualErr(t, nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, true, b.Equal(h, []byte("test")))
		ztesting.AssertEqual(t, false, b.Equal(h, []byte("wrong")))
	})
}

func TestCompare(t *testing.T) {
	t.Run("short hash", func(t *testing.T) {
		b, err := zpbkdf2.New(32, 210000, 32, crypto.SHA256)
		ztesting.AssertEqualErr(t, nil, err)
		err = b.Compare([]byte("short hash"), []byte("test"))
		ztesting.AssertEqualErr(t, zpbkdf2.ErrShortHash, err)
	})
	t.Run("pw not match", func(t *testing.T) {
		b, err := zpbkdf2.New(32, 210000, 32, crypto.SHA256)
		ztesting.AssertEqualErr(t, nil, err)
		h, err := b.Sum([]byte("test"))
		ztesting.AssertEqualErr(t, nil, err)
		err = b.Compare(h, []byte("wrong"))
		ztesting.AssertEqualErr(t, zpbkdf2.ErrNotMatch, err)
	})
	t.Run("pw match", func(t *testing.T) {
		b, err := zpbkdf2.New(32, 210000, 32, crypto.SHA256)
		ztesting.AssertEqualErr(t, nil, err)
		h, err := b.Sum([]byte("test"))
		ztesting.AssertEqualErr(t, nil, err)
		err = b.Compare(h, []byte("test"))
		ztesting.AssertEqualErr(t, nil, err)
	})
}
