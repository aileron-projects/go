package zsha1

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"

	"github.com/aileron-projects/go/zcrypto/internal"
)

var (
	_ internal.SumFunc      = Sum
	_ internal.EqualSumFunc = EqualSum
)

// Sum returns SHA1 hash.
// It uses [crypto/sha1.New].
func Sum(b []byte) []byte {
	h := sha1.New()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, sha1.Size))
}

// EqualSum compares SHA1 hash.
// It returns if the sum matches to the hash of b.
func EqualSum(b []byte, sum []byte) bool {
	return bytes.Equal(Sum(b), sum)
}

var (
	_ internal.HMACSumFunc      = HMACSum
	_ internal.HMACEqualSumFunc = HMACEqualSum
)

// HMACSum returns HMAC-SHA1 hash.
// It uses [crypto/sha1.New].
func HMACSum(msg, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, sha1.Size))
}

// HMACEqualSum compares HMAC-SHA1 hash.
// It returns if the sum matches to the hash of msg.
func HMACEqualSum(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum(msg, key), sum)
}
