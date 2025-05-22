package zmd4

import (
	"bytes"
	"crypto/hmac"

	"github.com/aileron-projects/go/zcrypto/internal"
	"golang.org/x/crypto/md4"
)

var (
	_ internal.SumFunc      = Sum
	_ internal.EqualSumFunc = EqualSum
)

// Sum returns MD4 hash.
// It uses [golang.org/x/crypto/md4.New].
func Sum(b []byte) []byte {
	h := md4.New()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, md4.Size))
}

// EqualSum compares MD4 hash.
// It returns if the sum matches to the hash of b.
func EqualSum(b []byte, sum []byte) bool {
	return bytes.Equal(Sum(b), sum)
}

var (
	_ internal.HMACSumFunc      = HMACSum
	_ internal.HMACEqualSumFunc = HMACEqualSum
)

// HMACSum returns HMAC-MD4 hash.
// It uses [crypto/md4.New].
func HMACSum(msg, key []byte) []byte {
	mac := hmac.New(md4.New, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, md4.Size))
}

// HMACEqualSum compares HMAC-MD4 hash.
// It returns if the sum matches to the hash of msg.
func HMACEqualSum(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum(msg, key), sum)
}
