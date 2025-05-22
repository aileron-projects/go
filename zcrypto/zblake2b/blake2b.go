package zblake2b

import (
	"bytes"
	"crypto/sha512"

	"github.com/aileron-projects/go/zcrypto/internal"
	"golang.org/x/crypto/blake2b"
)

var (
	_ internal.SumFunc      = Sum256
	_ internal.SumFunc      = Sum384
	_ internal.SumFunc      = Sum512
	_ internal.EqualSumFunc = EqualSum256
	_ internal.EqualSumFunc = EqualSum384
	_ internal.EqualSumFunc = EqualSum512
)

// Sum256 returns BLAKE2b-256 hash.
// It uses [golang.org/x/crypto/blake2b.New256].
func Sum256(b []byte) []byte {
	h, _ := blake2b.New256(nil)
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, blake2b.Size256))
}

// Sum384 returns BLAKE2b-384 hash.
// It uses [golang.org/x/crypto/blake2b.New384].
func Sum384(b []byte) []byte {
	h, _ := blake2b.New384(nil)
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, blake2b.Size384))
}

// Sum512 returns BLAKE2b-512 hash.
// It uses [golang.org/x/crypto/blake2b.New512].
func Sum512(b []byte) []byte {
	h, _ := blake2b.New512(nil)
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, blake2b.Size))
}

// EqualSum256 compares BLAKE2b-256 hash.
// It returns if the sum matches to the hash of b.
func EqualSum256(b []byte, sum []byte) bool {
	return bytes.Equal(Sum256(b), sum)
}

// EqualSum384 compares BLAKE2b-384 hash.
// It returns if the sum matches to the hash of b.
func EqualSum384(b []byte, sum []byte) bool {
	return bytes.Equal(Sum384(b), sum)
}

// EqualSum512 compares BLAKE2b-512 hash.
// It returns if the sum matches to the hash of b.
func EqualSum512(b []byte, sum []byte) bool {
	return bytes.Equal(Sum512(b), sum)
}

var (
	_ internal.HMACSumFunc      = HMACSum256
	_ internal.HMACSumFunc      = HMACSum384
	_ internal.HMACSumFunc      = HMACSum512
	_ internal.HMACEqualSumFunc = HMACEqualSum256
	_ internal.HMACEqualSumFunc = HMACEqualSum384
	_ internal.HMACEqualSumFunc = HMACEqualSum512
)

// HMACSum256 returns HMAC-BLAKE2b-256 hash.
// If the key length is grater than 64,
// it will be shorten to 64 bytes using [sha512.Sum512].
// It uses [golang.org/x/crypto/blake2b.New256].
func HMACSum256(msg, key []byte) []byte {
	if len(key) > blake2b.Size {
		newKey := sha512.Sum512(key)
		key = newKey[:]
	}
	mac, _ := blake2b.New256(key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, blake2b.Size256))
}

// HMACSum384 returns HMAC-BLAKE2b-384 hash.
// If the key length is grater than 64,
// it will be shorten to 64 bytes using [sha512.Sum512].
// It uses [golang.org/x/crypto/blake2b.New384].
func HMACSum384(msg, key []byte) []byte {
	if len(key) > blake2b.Size {
		newKey := sha512.Sum512(key)
		key = newKey[:]
	}
	mac, _ := blake2b.New384(key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, blake2b.Size384))
}

// HMACSum512 returns HMAC-BLAKE2b-512 hash.
// If the key length is grater than 64,
// it will be shorten to 64 bytes using [sha512.Sum512].
// It uses [golang.org/x/crypto/blake2b.New512].
func HMACSum512(msg, key []byte) []byte {
	if len(key) > blake2b.Size {
		newKey := sha512.Sum512(key)
		key = newKey[:]
	}
	mac, _ := blake2b.New512(key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, blake2b.Size))
}

// HMACEqualSum256 compares HMAC-BLAKE2b-256 hash.
// It returns if the sum matches to the hash of msg.
func HMACEqualSum256(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum256(msg, key), sum)
}

// HMACEqualSum384 compares HMAC-BLAKE2b-384 hash.
// It returns if the sum matches to the hash of msg.
func HMACEqualSum384(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum384(msg, key), sum)
}

// HMACEqualSum512 compares HMAC-BLAKE2b-512 hash.
// It returns if the sum matches to the hash of msg.
func HMACEqualSum512(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum512(msg, key), sum)
}
