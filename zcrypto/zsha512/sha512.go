package zsha512

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"

	"github.com/aileron-projects/go/zcrypto/internal"
)

var (
	_ internal.SumFunc      = Sum224
	_ internal.SumFunc      = Sum256
	_ internal.SumFunc      = Sum384
	_ internal.SumFunc      = Sum512
	_ internal.EqualSumFunc = EqualSum224
	_ internal.EqualSumFunc = EqualSum256
	_ internal.EqualSumFunc = EqualSum384
	_ internal.EqualSumFunc = EqualSum512
)

// Sum224 returns SHA512/224 hash.
// It uses [crypto/sha512.New512_224].
func Sum224(b []byte) []byte {
	h := sha512.New512_224()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, sha512.Size224))
}

// Sum256 returns SHA512/256 hash.
// It uses [crypto/sha512.New512_256].
func Sum256(b []byte) []byte {
	h := sha512.New512_256()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, sha512.Size256))
}

// Sum384 returns SHA512/384 hash.
// It uses [crypto/sha512.New384].
func Sum384(b []byte) []byte {
	h := sha512.New384()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, sha512.Size384))
}

// Sum512 returns SHA512 hash.
// It uses [crypto/sha512.New].
func Sum512(b []byte) []byte {
	h := sha512.New()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, sha512.Size))
}

// EqualSum224 compares SHA512/224 hash.
// It returns if the sum matches to the hash of b.
func EqualSum224(b []byte, sum []byte) bool {
	return bytes.Equal(Sum224(b), sum)
}

// EqualSum256 compares SHA512/256 hash.
// It returns if the sum matches to the hash of b.
func EqualSum256(b []byte, sum []byte) bool {
	return bytes.Equal(Sum256(b), sum)
}

// EqualSum384 compares SHA512/384 hash.
// It returns if the sum matches to the hash of b.
func EqualSum384(b []byte, sum []byte) bool {
	return bytes.Equal(Sum384(b), sum)
}

// EqualSum512 compares SHA512 hash.
// It returns if the sum matches to the hash of b.
func EqualSum512(b []byte, sum []byte) bool {
	return bytes.Equal(Sum512(b), sum)
}

var (
	_ internal.HMACSumFunc      = HMACSum224
	_ internal.HMACSumFunc      = HMACSum256
	_ internal.HMACSumFunc      = HMACSum384
	_ internal.HMACSumFunc      = HMACSum512
	_ internal.HMACEqualSumFunc = HMACEqualSum224
	_ internal.HMACEqualSumFunc = HMACEqualSum256
	_ internal.HMACEqualSumFunc = HMACEqualSum384
	_ internal.HMACEqualSumFunc = HMACEqualSum512
)

// HMACSum224 returns HMAC-SHA512/224 hash.
// It uses [crypto/sha512.New512_224].
func HMACSum224(msg, key []byte) []byte {
	mac := hmac.New(sha512.New512_224, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, sha512.Size224))
}

// HMACSum256 returns HMAC-SHA512/256 hash.
// It uses [crypto/sha512.New512_256].
func HMACSum256(msg, key []byte) []byte {
	mac := hmac.New(sha512.New512_256, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, sha512.Size256))
}

// HMACSum384 returns HMAC-SHA512/384 hash.
// It uses [crypto/sha512.New384].
func HMACSum384(msg, key []byte) []byte {
	mac := hmac.New(sha512.New384, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, sha512.Size384))
}

// HMACSum512 returns HMAC-SHA512 hash.
// It uses [crypto/sha512.New].
func HMACSum512(msg, key []byte) []byte {
	mac := hmac.New(sha512.New, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, sha512.Size))
}

// HMACEqualSum224 compares HMAC-SHA512/224 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum224(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum224(msg, key), sum)
}

// HMACEqualSum256 compares HMAC-SHA512/256 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum256(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum256(msg, key), sum)
}

// HMACEqualSum384 compares HMAC-SHA512/384 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum384(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum384(msg, key), sum)
}

// HMACEqualSum512 compares HMAC-SHA512 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum512(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum512(msg, key), sum)
}
