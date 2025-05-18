package zhmac_test

import (
	"crypto"
	_ "crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zhmac"
	"github.com/aileron-projects/go/ztesting"
)

func TestSum(t *testing.T) {
	t.Parallel()
	t.Run("not available", func(t *testing.T) {
		defer func() {
			r := recover()
			ztesting.AssertEqualErr(t, "error not match", zhmac.ErrNotAvailable, r.(error))
		}()
		zhmac.Sum(crypto.Hash(999), nil, nil)
	})
	t.Run("available", func(t *testing.T) {
		sum := zhmac.Sum(crypto.SHA256, []byte("test"), []byte("key"))
		ztesting.AssertEqual(t, "invalid hash", "02afb56304902c656fcb737cdd03de6205bb6d401da2812efd9b2d36a08af159", hex.EncodeToString(sum))
	})
}

func TestEqual(t *testing.T) {
	t.Parallel()
	t.Run("match", func(t *testing.T) {
		sum, _ := hex.DecodeString("02afb56304902c656fcb737cdd03de6205bb6d401da2812efd9b2d36a08af159")
		equal := zhmac.Equal(crypto.SHA256, []byte("test"), []byte("key"), sum)
		ztesting.AssertEqual(t, "incorrect comparison result", true, equal)
	})
	t.Run("mismatch", func(t *testing.T) {
		sum, _ := hex.DecodeString("02afb56304902c656fcb737cdd03de6205bb6d401da2812efd9b2d36a08af159")
		equal := zhmac.Equal(crypto.SHA256, []byte("wrong"), []byte("key"), sum)
		ztesting.AssertEqual(t, "incorrect comparison result", false, equal)
	})
}

func TestCompare(t *testing.T) {
	t.Parallel()
	t.Run("hash unavailable", func(t *testing.T) {
		sum, _ := hex.DecodeString("02afb56304902c656fcb737cdd03de6205bb6d401da2812efd9b2d36a08af159")
		err := zhmac.Compare(crypto.Hash(999), []byte("test"), []byte("key"), sum)
		ztesting.AssertEqualErr(t, "incorrect comparison result", zhmac.ErrNotAvailable, err)
	})
	t.Run("match", func(t *testing.T) {
		sum, _ := hex.DecodeString("02afb56304902c656fcb737cdd03de6205bb6d401da2812efd9b2d36a08af159")
		err := zhmac.Compare(crypto.SHA256, []byte("test"), []byte("key"), sum)
		ztesting.AssertEqualErr(t, "incorrect comparison result", nil, err)
	})
	t.Run("mismatch", func(t *testing.T) {
		sum, _ := hex.DecodeString("02afb56304902c656fcb737cdd03de6205bb6d401da2812efd9b2d36a08af159")
		err := zhmac.Compare(crypto.SHA256, []byte("wrong"), []byte("key"), sum)
		ztesting.AssertEqualErr(t, "incorrect comparison result", zhmac.ErrNotMatch, err)
	})
}
