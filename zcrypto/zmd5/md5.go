package zmd5

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"

	"github.com/aileron-projects/go/internal/ihash"
)

var (
	_ ihash.SumFunc      = Sum
	_ ihash.EqualSumFunc = EqualSum
)

// Sum returns MD5 hash.
// It uses [crypto/md5.New].
func Sum(b []byte) []byte {
	h := md5.New()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, md5.Size))
}

// EqualSum compares MD5 hash.
// It returns if the sum matches to the hash of b.
func EqualSum(b []byte, sum []byte) bool {
	return bytes.Equal(Sum(b), sum)
}

var (
	_ ihash.HMACSumFunc      = HMACSum
	_ ihash.HMACEqualSumFunc = HMACEqualSum
)

// HMACSum returns HMAC-MD5 hash.
// It uses [crypto/md5.New].
func HMACSum(msg, key []byte) []byte {
	mac := hmac.New(md5.New, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, md5.Size))
}

// HMACEqualSum compares HMAC-MD5 hash.
// It returns if the sum matches to the hash of msg.
func HMACEqualSum(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum(msg, key), sum)
}
