package zsha256

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"

	"github.com/aileron-projects/go/internal/ihash"
)

var (
	_ ihash.SumFunc      = Sum224
	_ ihash.SumFunc      = Sum256
	_ ihash.EqualSumFunc = EqualSum224
	_ ihash.EqualSumFunc = EqualSum256
)

// Sum224 returns SHA256/224 hash.
// It uses [crypto/sha256.New224].
func Sum224(b []byte) []byte {
	h := sha256.New224()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, sha256.Size224))
}

// Sum256 returns SHA256 hash.
// It uses [crypto/sha256.New].
func Sum256(b []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, sha256.Size))
}

// EqualSum224 compares SHA256/224 hash.
// It returns if the sum matches to the hash of b.
func EqualSum224(b []byte, sum []byte) bool {
	return bytes.Equal(Sum224(b), sum)
}

// EqualSum256 compares SHA256 hash.
// It returns if the sum matches to the hash of b.
func EqualSum256(b []byte, sum []byte) bool {
	return bytes.Equal(Sum256(b), sum)
}

var (
	_ ihash.HMACSumFunc      = HMACSum224
	_ ihash.HMACSumFunc      = HMACSum256
	_ ihash.HMACEqualSumFunc = HMACEqualSum224
	_ ihash.HMACEqualSumFunc = HMACEqualSum256
)

// HMACSum224 returns HMAC-SHA256/224 hash.
// It uses [crypto/sha256.New224].
func HMACSum224(msg, key []byte) []byte {
	mac := hmac.New(sha256.New224, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, sha256.Size224))
}

// HMACSum256 returns HMAC-SHA256 hash.
// It uses [crypto/sha256.New].
func HMACSum256(msg, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, sha256.Size))
}

// HMACEqualSum224 compares HMAC-SHA256/224 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum224(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum224(msg, key), sum)
}

// HMACEqualSum256 compares HMAC-SHA256 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum256(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum256(msg, key), sum)
}
