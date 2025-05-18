package zhash

import (
	"bytes"
	"crypto"
	"errors"
	"hash"
	"strconv"

	"github.com/aileron-projects/go/internal/ihash"
)

var (
	_ ihash.SumFunc        = Hash(0).Sum
	_ ihash.EqualSumFunc   = Hash(0).Equal
	_ ihash.CompareSumFunc = Hash(0).Compare
)

var (
	ErrUnknown      = errors.New("zhash: unknown hash algorithm")
	ErrNotAvailable = errors.New("zhash: hash function not available")
	ErrNotMatch     = errors.New("zhash: hashed value not match")
)

// Hash identifies a 'Non Cryptographic' hash function
// that is implemented in another package.
// Note that this identifies NON CRYPTOGRAPHIC hash.
// Use [crypto.Hash] for cryptographic hash.
type Hash uint

const (
	// Hash algorithms supported by the standard [crypto] package.
	_      = crypto.MD4
	_      = crypto.MD5
	_      = crypto.SHA1
	_      = crypto.SHA224
	_      = crypto.SHA256
	_      = crypto.SHA384
	_      = crypto.SHA512
	_      = crypto.MD5SHA1
	_      = crypto.RIPEMD160
	_      = crypto.SHA3_224
	_      = crypto.SHA3_256
	_      = crypto.SHA3_384
	_      = crypto.SHA3_512
	_      = crypto.SHA512_224
	_      = crypto.SHA512_256
	_      = crypto.BLAKE2s_256
	_      = crypto.BLAKE2b_256
	_      = crypto.BLAKE2b_384
	_      = crypto.BLAKE2b_512
	offset = 1 + iota
)

const (
	FNV32           Hash = offset + iota // Non cryptographic. import github.com/aileron-projects/go/zhash/zfnv
	FNV32a                               // Non cryptographic. import github.com/aileron-projects/go/zhash/zfnv
	FNV64                                // Non cryptographic. import github.com/aileron-projects/go/zhash/zfnv
	FNV64a                               // Non cryptographic. import github.com/aileron-projects/go/zhash/zfnv
	FNV128                               // Non cryptographic. import github.com/aileron-projects/go/zhash/zfnv
	FNV128a                              // Non cryptographic. import github.com/aileron-projects/go/zhash/zfnv
	CRC32IEEE                            // Non cryptographic. import github.com/aileron-projects/go/zhash/zcrc32
	CRC32Castagnoli                      // Non cryptographic. import github.com/aileron-projects/go/zhash/zcrc32
	CRC32Koopman                         // Non cryptographic. import github.com/aileron-projects/go/zhash/zcrc32
	CRC64ISO                             // Non cryptographic. import github.com/aileron-projects/go/zhash/zcrc64
	CRC64ECMA                            // Non cryptographic. import github.com/aileron-projects/go/zhash/zcrc64
	maxHash
)

var digestSizes = [...]uint8{
	FNV32 - offset:           4,
	FNV32a - offset:          4,
	FNV64 - offset:           8,
	FNV64a - offset:          8,
	FNV128 - offset:          16,
	FNV128a - offset:         16,
	CRC32IEEE - offset:       4,
	CRC32Castagnoli - offset: 4,
	CRC32Koopman - offset:    4,
	CRC64ISO - offset:        8,
	CRC64ECMA - offset:       8,
}

// hashes contains functions that creates a new hash instance.
var hashes = make([]func() hash.Hash, maxHash-offset)

// RegisterHash registers a hash function.
// It panics [ErrUnknown] when unknown hash was given.
func RegisterHash(h Hash, f func() hash.Hash) {
	if h < offset || h >= maxHash {
		panic(ErrUnknown)
	}
	hashes[h-offset] = f
}

func (h Hash) String() string {
	switch h {
	case FNV32:
		return "FNV1/32"
	case FNV32a:
		return "FNV1a/32"
	case FNV64:
		return "FNV1/64"
	case FNV64a:
		return "FNV1a/64"
	case FNV128:
		return "FNV1/128"
	case FNV128a:
		return "FNV1a/128"
	case CRC32IEEE:
		return "CRC32-IEEE"
	case CRC32Castagnoli:
		return "CRC32-Castagnoli"
	case CRC32Koopman:
		return "CRC32-Koopman"
	case CRC64ISO:
		return "CRC64-ISO"
	case CRC64ECMA:
		return "CRC64-ECMA"
	default:
		return "unknown hash value " + strconv.Itoa(int(h))
	}
}

// Size returns the digest size of the hash.
// It panics with [ErrUnknown] when unknown hash was given.
func (h Hash) Size() int {
	if h < offset || h >= maxHash {
		panic(ErrUnknown)
	}
	return int(digestSizes[h-offset])
}

// Available reports whether the given hash algorithm is linked into the binary.
func (h Hash) Available() bool {
	return h >= offset && h < maxHash && hashes[h-offset] != nil
}

// New returns a new hash instance.
// It panics with [ErrNotAvailable] if hash function is not available.
func (h Hash) New() hash.Hash {
	if !h.Available() {
		panic(ErrNotAvailable)
	}
	return hashes[h-offset]()
}

// NewFunc returns a function to generate new hash instances.
// It panics with [ErrNotAvailable] if hash function is not available.
func (h Hash) NewFunc() func() hash.Hash {
	if !h.Available() {
		panic(ErrNotAvailable)
	}
	return hashes[h-offset]
}

// Sum returns hash sum of the given value.
// It panics with [ErrNotAvailable] if hash function is not available.
func (h Hash) Sum(value []byte) []byte {
	if !h.Available() {
		panic(ErrNotAvailable)
	}
	hf := hashes[h-offset]()
	_, _ = hf.Write(value)
	return hf.Sum(make([]byte, 0, h.Size()))
}

// Equal returns if the two hashed value and non-hashed value are equal or not.
// To know why hashes are not matched, use [Compare] instead.
func (h Hash) Equal(hashed, value []byte) bool {
	return h.Compare(hashed, value) == nil
}

// Compare compares the hashed value and the non-hashed value.
// Compare may return following errors:
//   - [ErrNotMatch]: hashed value does not match.
//   - [ErrNotAvailable]: hash function specified by h is not available.
//   - Any errors returned from the hash function.
func (h Hash) Compare(hashed, value []byte) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err, _ = r.(error)
		}
	}()
	if !bytes.Equal(hashed, h.Sum(value)) {
		return ErrNotMatch
	}
	return nil
}
