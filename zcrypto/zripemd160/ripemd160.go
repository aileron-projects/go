package zripemd160

import (
	"bytes"
	"crypto/hmac"

	"github.com/aileron-projects/go/internal/ihash"
	"golang.org/x/crypto/ripemd160"
)

var (
	_ ihash.SumFunc      = Sum
	_ ihash.EqualSumFunc = EqualSum
)

// Sum returns RIPEMD160 hash.
// It uses [crypto/ripemd160.New].
func Sum(b []byte) []byte {
	h := ripemd160.New()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, ripemd160.Size))
}

// EqualSum compares RIPEMD160 hash.
// It returns if the sum matches to the hash of b.
func EqualSum(b []byte, sum []byte) bool {
	return bytes.Equal(Sum(b), sum)
}

var (
	_ ihash.HMACSumFunc      = HMACSum
	_ ihash.HMACEqualSumFunc = HMACEqualSum
)

// HMACSum returns HMAC-RIPEMD160 hash.
// It uses [crypto/ripemd160.New].
func HMACSum(msg, key []byte) []byte {
	mac := hmac.New(ripemd160.New, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, ripemd160.Size))
}

// HMACEqualSum compares HMAC-RIPEMD160 hash.
// It returns if the sum matches to the hash of msg.
func HMACEqualSum(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum(msg, key), sum)
}
