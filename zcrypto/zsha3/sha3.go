package zsha3

import (
	"bytes"
	"crypto/hmac"

	"github.com/aileron-projects/go/zcrypto/internal"
	"golang.org/x/crypto/sha3"
)

var (
	_ internal.SumFunc      = Sum224
	_ internal.SumFunc      = Sum256
	_ internal.SumFunc      = Sum384
	_ internal.SumFunc      = Sum512
	_ internal.SumFunc      = SumShake128
	_ internal.SumFunc      = SumShake256
	_ internal.EqualSumFunc = EqualSum224
	_ internal.EqualSumFunc = EqualSum256
	_ internal.EqualSumFunc = EqualSum384
	_ internal.EqualSumFunc = EqualSum512
	_ internal.EqualSumFunc = EqualSumShake128
	_ internal.EqualSumFunc = EqualSumShake256
)

// Sum224 returns SHA3/224 hash.
// It uses [crypto/sha3.New224].
func Sum224(b []byte) []byte {
	h := sha3.New224()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, h.Size()))
}

// Sum256 returns SHA3/256 hash.
// It uses [crypto/sha3.New256].
func Sum256(b []byte) []byte {
	h := sha3.New256()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, h.Size()))
}

// Sum384 returns SHA3/384 hash.
// It uses [crypto/sha3.New384].
func Sum384(b []byte) []byte {
	h := sha3.New384()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, h.Size()))
}

// Sum512 returns SHA3/512 hash.
// It uses [crypto/sha3.New512].
func Sum512(b []byte) []byte {
	h := sha3.New512()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, h.Size()))
}

// SumShake128 returns SHAKE128 hash.
// Digest size is fixed to 256 bit.
// Use [crypto/sha3.ShakeHash] directory if variable length digest are necessary.
// It uses [crypto/sha3.NewShake128].
func SumShake128(b []byte) []byte {
	h := sha3.NewShake128()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, h.Size()))
}

// SumShake256 returns SHAKE256 hash.
// Digest size is fixed to 512 bit.
// Use [crypto/sha3.ShakeHash] directory if variable length digest are necessary.
// It uses [crypto/sha3.NewShake256].
func SumShake256(b []byte) []byte {
	h := sha3.NewShake256()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, h.Size()))
}

// EqualSum224 compares SHA3/224 hash.
// It returns if the sum matches to the hash of b.
func EqualSum224(b []byte, sum []byte) bool {
	return bytes.Equal(Sum224(b), sum)
}

// EqualSum256 compares SHA3/256 hash.
// It returns if the sum matches to the hash of b.
func EqualSum256(b []byte, sum []byte) bool {
	return bytes.Equal(Sum256(b), sum)
}

// EqualSum384 compares SHA3/384 hash.
// It returns if the sum matches to the hash of b.
func EqualSum384(b []byte, sum []byte) bool {
	return bytes.Equal(Sum384(b), sum)
}

// EqualSum512 compares SHA3/512 hash.
// It returns if the sum matches to the hash of b.
func EqualSum512(b []byte, sum []byte) bool {
	return bytes.Equal(Sum512(b), sum)
}

// EqualSumShake128 compares SHAKE128 hash.
// It returns if the sum matches to the hash of b.
func EqualSumShake128(b []byte, sum []byte) bool {
	return bytes.Equal(SumShake128(b), sum)
}

// EqualSumShake256 compares SHAKE256 hash.
// It returns if the sum matches to the hash of b.
func EqualSumShake256(b []byte, sum []byte) bool {
	return bytes.Equal(SumShake256(b), sum)
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

// HMACSum224 returns HMAC-SHA3/224 hash.
// It uses [crypto/sha3.New224].
func HMACSum224(msg, key []byte) []byte {
	mac := hmac.New(sha3.New224, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, mac.Size()))
}

// HMACSum256 returns HMAC-SHA3/256 hash.
// It uses [crypto/sha3.New256].
func HMACSum256(msg, key []byte) []byte {
	mac := hmac.New(sha3.New256, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, mac.Size()))
}

// HMACSum384 returns HMAC-SHA3/384 hash.
// It uses [crypto/sha3.New384].
func HMACSum384(msg, key []byte) []byte {
	mac := hmac.New(sha3.New384, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, mac.Size()))
}

// HMACSum512 returns HMAC-SHA3/512 hash.
// It uses [crypto/sha3.New512].
func HMACSum512(msg, key []byte) []byte {
	mac := hmac.New(sha3.New512, key)
	_, _ = mac.Write(msg)
	return mac.Sum(make([]byte, 0, mac.Size()))
}

// HMACEqualSum224 compares HMAC-SHA3/224 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum224(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum224(msg, key), sum)
}

// HMACEqualSum256 compares HMAC-SHA3/256 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum256(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum256(msg, key), sum)
}

// HMACEqualSum384 compares HMAC-SHA3/384 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum384(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum384(msg, key), sum)
}

// HMACEqualSum512 compares HMAC-SHA3/512 hash.
// It returns if the sum matches to the hash of b.
func HMACEqualSum512(msg, key, sum []byte) bool {
	return bytes.Equal(HMACSum512(msg, key), sum)
}
