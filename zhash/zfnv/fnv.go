package zfnv

import (
	"bytes"
	"hash"
	"hash/fnv"

	"github.com/aileron-projects/go/zhash"
)

func init() {
	zhash.RegisterHash(zhash.FNV32, func() hash.Hash { return fnv.New32() })
	zhash.RegisterHash(zhash.FNV32a, func() hash.Hash { return fnv.New32a() })
	zhash.RegisterHash(zhash.FNV64, func() hash.Hash { return fnv.New64() })
	zhash.RegisterHash(zhash.FNV64a, func() hash.Hash { return fnv.New64a() })
	zhash.RegisterHash(zhash.FNV128, func() hash.Hash { return fnv.New128() })
	zhash.RegisterHash(zhash.FNV128a, func() hash.Hash { return fnv.New128a() })
}

// Sum32 returns FNV1/32 hash.
// It uses [hash/fnv.New32].
func Sum32(b []byte) []byte {
	h := fnv.New32()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, 4))
}

// Sum32a returns FNV1a/32 hash.
// It uses [hash/fnv.New32a].
func Sum32a(b []byte) []byte {
	h := fnv.New32a()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, 4))
}

// Sum64 returns FNV1/64 hash.
// It uses [hash/fnv.New64].
func Sum64(b []byte) []byte {
	h := fnv.New64()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, 8))
}

// Sum64a returns FNV1a/64 hash.
// It uses [hash/fnv.New64a].
func Sum64a(b []byte) []byte {
	h := fnv.New64a()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, 8))
}

// Sum128 returns FNV1/128 hash.
// It uses [hash/fnv.New128].
func Sum128(b []byte) []byte {
	h := fnv.New128()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, 16))
}

// Sum128a returns FNV1a/128 hash.
// It uses [hash/fnv.New128a].
func Sum128a(b []byte) []byte {
	h := fnv.New128a()
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, 16))
}

// EqualSum32 compares FNV1/32 hash.
// It returns if the sum matches to the hash of b.
func EqualSum32(b []byte, sum []byte) bool {
	return bytes.Equal(Sum32(b), sum)
}

// EqualSum32a compares FNV1a/32 hash.
// It returns if the sum matches to the hash of b.
func EqualSum32a(b []byte, sum []byte) bool {
	return bytes.Equal(Sum32a(b), sum)
}

// EqualSum64 compares FNV1/64 hash.
// It returns if the sum matches to the hash of b.
func EqualSum64(b []byte, sum []byte) bool {
	return bytes.Equal(Sum64(b), sum)
}

// EqualSum64s compares FNV1s/64 hash.
// It returns if the sum matches to the hash of b.
func EqualSum64a(b []byte, sum []byte) bool {
	return bytes.Equal(Sum64a(b), sum)
}

// EqualSum128 compares FNV1/128 hash.
// It returns if the sum matches to the hash of b.
func EqualSum128(b []byte, sum []byte) bool {
	return bytes.Equal(Sum128(b), sum)
}

// EqualSum128a compares FNV1a/128 hash.
// It returns if the sum matches to the hash of b.
func EqualSum128a(b []byte, sum []byte) bool {
	return bytes.Equal(Sum128a(b), sum)
}
