package zblake2b

import (
	"bytes"
	"crypto/sha512"

	"github.com/aileron-projects/go/internal/ihash"
	"golang.org/x/crypto/blake2b"
)

var (
	_ ihash.SumFunc      = Sum256
	_ ihash.SumFunc      = Sum384
	_ ihash.SumFunc      = Sum512
	_ ihash.EqualSumFunc = EqualSum256
	_ ihash.EqualSumFunc = EqualSum384
	_ ihash.EqualSumFunc = EqualSum512
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
	_ ihash.HMACSumFunc      = HMACSum256
	_ ihash.HMACSumFunc      = HMACSum384
	_ ihash.HMACSumFunc      = HMACSum512
	_ ihash.HMACEqualSumFunc = HMACEqualSum256
	_ ihash.HMACEqualSumFunc = HMACEqualSum384
	_ ihash.HMACEqualSumFunc = HMACEqualSum512
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
