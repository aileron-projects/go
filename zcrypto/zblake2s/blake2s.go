package zblake2s

import (
	"bytes"
	"crypto/sha256"

	"github.com/aileron-projects/go/internal/ihash"
	"golang.org/x/crypto/blake2s"
)

var (
	_ ihash.SumFunc      = Sum256
	_ ihash.EqualSumFunc = EqualSum256
)

// Sum256 returns BLAKE2s-256 hash.
// It uses [golang.org/x/crypto/blake2s.New256].
func Sum256(b []byte) []byte {
	h, _ := blake2s.New256(nil)
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, blake2s.Size))
}

// EqualSum256 compares BLAKE2s-256 hash.
// It returns if the sum matches to the hash of b.
func EqualSum256(b []byte, sum []byte) bool {
	return bytes.Equal(Sum256(b), sum)
}

var (
	_ ihash.HMACSumFunc      = HMACSum256
	_ ihash.HMACEqualSumFunc = HMACEqualSum256
)

// HMACSum256 returns HMAC-BLAKE2s/256 hash.
// If the key length is grater than 32,
// it will be shorten to 32 bytes using [sha256.Sum256].
// It uses [golang.org/x/crypto/blake2s.New256].
func HMACSum256(msg, key []byte) []byte {
	if len(key) > blake2s.Size {
		newKey := sha256.Sum256(key)
		key = newKey[:]
	}
	mac, _ := blake2s.New256(key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, blake2s.Size))
}

// HMACEqualSum256 compares HMAC-BLAKE2s/256 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum256(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum256(msg, key), sum)
}
