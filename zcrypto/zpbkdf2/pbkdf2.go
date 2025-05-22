package zpbkdf2

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"errors"
	"hash"

	"github.com/aileron-projects/go/internal/ihash"
	"github.com/aileron-projects/go/zcrypto/internal"
	"golang.org/x/crypto/pbkdf2"
)

var (
	_ ihash.PWHasher = &PBKDF2{}
)

var (
	ErrNotMatch    = errors.New("zpbkdf2: password hash value not match")
	ErrShortHash   = errors.New("zpbkdf2: password hash is too short")
	ErrUnavailable = errors.New("zpbkdf2: hash algorithm is not available")
)

// New returns a new instance of PBKDF2 hasher.
// The parameter saltLen, iter, keyLen and hash func of h are passed to the [golang.org/x/crypto/pbkdf2.Key].
// Currently following hash algorithm are supported for h.
//
//   - [crypto.MD4]
//   - [crypto.MD5]
//   - [crypto.SHA1]
//   - [crypto.SHA224]
//   - [crypto.SHA256]
//   - [crypto.SHA384]
//   - [crypto.SHA512]
//   - [crypto.SHA512_224]
//   - [crypto.SHA512_256]
//   - [crypto.RIPEMD160]
//   - [crypto.SHA3_224]
//   - [crypto.SHA3_256]
//   - [crypto.SHA3_384]
//   - [crypto.SHA3_512]
//   - [crypto.BLAKE2s_256]
//   - [crypto.BLAKE2b_256]
//   - [crypto.BLAKE2b_384]
//   - [crypto.BLAKE2b_512]
//
// Parameter ranges:
//
//   - 0 <= saltLen
//   - 2 <= iter
//   - 1 <= keyLen
//
// Recommended parameters:
//
//   - iter=310000,keyLen=32,hash=sha256 (https://datatracker.ietf.org/doc/draft-ietf-kitten-password-storage/)
//   - iter=1300000,keyLen=-,hash=sha1 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - iter=600000,keyLen=-,hash=sha256 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - iter=210000,keyLen=-,hash=sha512 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
func New(saltLen, iter, keyLen int, h crypto.Hash) (*PBKDF2, error) {
	if !h.Available() {
		return nil, ErrUnavailable
	}
	c := &PBKDF2{
		saltLen:  saltLen,
		iter:     iter,
		keyLen:   keyLen,
		hashFunc: h.New,
	}
	return c, nil
}

// PBKDF2 is the hasher type for [golang.org/x/crypto/pbkdf2].
type PBKDF2 struct {
	saltLen  int
	iter     int
	keyLen   int
	hashFunc func() hash.Hash
}

// Split splits the hashed password into salt and hash value of the password.
// If the hashedPW is too short, it returns [ErrShortHash].
func (c *PBKDF2) Split(hashedPW []byte) (salt, hash []byte, err error) {
	n := len(hashedPW) - c.saltLen
	if n < 0 {
		return nil, nil, ErrShortHash
	}
	return hashedPW[:c.saltLen], hashedPW[c.saltLen:], nil
}

// Sum returns the hash sum value of the password.
// Salt is joined at the left side of the returned sum.
// Use [PBKDF2.Split] to split salt and the sum of the password.
// [PBKDF2.Sum] always returns nil as error.
func (c *PBKDF2) Sum(password []byte) ([]byte, error) {
	salt := make([]byte, c.saltLen)
	_, err := rand.Read(salt)
	internal.MustNil(err)
	hashed := pbkdf2.Key(password, salt, c.iter, c.keyLen, c.hashFunc)
	return append(salt, hashed...), nil
}

// Equal reports if the given hashed password and the password
// are the same or not.
func (c *PBKDF2) Equal(hashedPW, pw []byte) bool {
	return c.Compare(hashedPW, pw) == nil
}

// Compare compares hashed password and the password.
// If the hashed value of the pw matched to the hashedPW,
// it returns nil. If not matched, it returns non-nil error.
func (c *PBKDF2) Compare(hashedPW, pw []byte) error {
	n := len(hashedPW) - c.saltLen
	if n < 0 {
		return ErrShortHash
	}
	salt, hashed := hashedPW[:c.saltLen], hashedPW[c.saltLen:]
	x := pbkdf2.Key(pw, salt, c.iter, c.keyLen, c.hashFunc)
	if !bytes.Equal(x, hashed) {
		return ErrNotMatch
	}
	return nil
}
