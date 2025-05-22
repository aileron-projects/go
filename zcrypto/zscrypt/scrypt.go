package zscrypt

import (
	"bytes"
	"crypto/rand"
	"errors"

	"github.com/aileron-projects/go/zcrypto/internal"
	"golang.org/x/crypto/scrypt"
)

var (
	_ internal.PWHasher = &SCrypt{}
)

var (
	ErrNotMatch  = errors.New("zscrypt: password hash value not match")
	ErrShortHash = errors.New("zscrypt: password hash is too short")
)

// New returns a new instance of SCrypt hasher with given parameters.
// The parameters n, r, p and keyLen are passed to the [golang.org/x/crypto/scrypt.Key].
//
// Parameter ranges:
//
//   - 0 <= saltLen
//   - 1 < n. n must be power of 2
//   - n <= [math.MaxInt]/128/r
//   - r*p < 1<<30 (1,073,741,824)
//   - r <= [math.MaxInt]/128/p
//   - r <= [math.MaxInt]/256
//
// Recommended parameters:
//
//   - n=32768,r=8,p=1,keyLen=32 (https://datatracker.ietf.org/doc/draft-ietf-kitten-password-storage/)
//   - n=131072,r=8,p=1,keyLen=32 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - n=65536,r=8,p=2,keyLen=32 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - n=32768,r=8,p=3,keyLen=32 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - n=16384,r=8,p=5,keyLen=32 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - n=8192,r=8,p=10,keyLen=32 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
func New(saltLen, n, r, p, keyLen int) (*SCrypt, error) {
	c := &SCrypt{
		saltLen: saltLen,
		n:       n,
		r:       r,
		p:       p,
		keyLen:  keyLen,
	}
	if _, err := c.Sum([]byte("")); err != nil { // Validate parameters.
		return nil, err
	}
	return c, nil
}

// SCrypt is the hasher type for [golang.org/x/crypto/scrypt].
type SCrypt struct {
	saltLen int
	n       int
	r       int
	p       int
	keyLen  int
}

// Split splits the hashed password into salt and hash value of the password.
// If the hashedPW is too short, it returns [ErrShortHash].
func (c *SCrypt) Split(hashedPW []byte) (salt, hash []byte, err error) {
	n := len(hashedPW) - c.saltLen
	if n < 0 {
		return nil, nil, ErrShortHash
	}
	return hashedPW[:c.saltLen], hashedPW[c.saltLen:], nil
}

// Sum returns the hash sum value of the password.
// Salt is joined at the left side of the returned sum.
// Use [SCrypt.Split] to split salt and the sum of the password.
func (c *SCrypt) Sum(password []byte) ([]byte, error) {
	salt := make([]byte, c.saltLen)
	_, err := rand.Read(salt)
	internal.MustNil(err)
	hashed, err := scrypt.Key(password, salt, c.n, c.r, c.p, c.keyLen)
	if err != nil {
		return nil, err
	}
	return append(salt, hashed...), nil
}

// Equal reports if the given hashed password and the password
// are the same or not.
func (c *SCrypt) Equal(hashedPW, pw []byte) bool {
	return c.Compare(hashedPW, pw) == nil
}

// Compare compares hashed password and the password.
// If the hashed value of the pw matched to the hashedPW,
// it returns nil. If not matched, it returns non-nil error.
func (c *SCrypt) Compare(hashedPW, pw []byte) error {
	n := len(hashedPW) - c.saltLen
	if n < 0 {
		return ErrShortHash
	}
	salt, hashed := hashedPW[:c.saltLen], hashedPW[c.saltLen:]
	x, err := scrypt.Key(pw, salt, c.n, c.r, c.p, c.keyLen)
	if err != nil {
		return err
	}
	if !bytes.Equal(x, hashed) {
		return ErrNotMatch
	}
	return nil
}
