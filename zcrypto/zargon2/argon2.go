package zargon2

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/internal"
	"golang.org/x/crypto/argon2"
)

var (
	_ internal.PWHasher = &Argon2i{}
	_ internal.PWHasher = &Argon2id{}
)

var (
	ErrNotMatch  = errors.New("zargon2: password hash value not match")
	ErrShortHash = errors.New("zargon2: password hash is too short")
)

// NewArgon2i returns a new instance of Argon2i hasher.
// Note that [NewArgon2id] is recommended from the security consideration than [NewArgon2i].
// The parameter saltLen, time, memory, threads and keyLen are passed to the [golang.org/x/crypto/argon2.Key].
//
// Parameter ranges:
//
//   - 0 <= saltLen
//   - 1 <= time
//   - 1 <= threads
//
// Recommended parameters:
//
//   - time=1,memory=2*1024*1024,threads=4,keyLen=32 (https://datatracker.ietf.org/doc/draft-ietf-kitten-password-storage/)
//   - time=1,memory=2*1024*1024,threads=4,keyLen=32 (https://datatracker.ietf.org/doc/rfc9106/)
//   - time=1,memory=64*1024,threads=4,keyLen=32 (https://datatracker.ietf.org/doc/rfc9106/)
//   - time=3,memory=12*1024,threads=1 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - time=4,memory=9*1024,threads=1 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - time=5,memory=7*1024,threads=1 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
func NewArgon2i(saltLen int, time, memory uint32, threads uint8, keyLen uint32) (a *Argon2i, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("zargon2: invalid hash parameter [%v]", r)
		}
	}()
	c := &Argon2i{
		saltLen: saltLen,
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
	}
	_, _ = c.Sum([]byte("")) // Validate parameters.
	return c, nil
}

// Argon2i is the hasher type for [golang.org/x/crypto/argon2].
type Argon2i struct {
	saltLen int
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// Split splits the hashed password into salt and hash value of the password.
// If the hashedPW is too short, it returns [ErrShortHash].
func (c *Argon2i) Split(hashedPW []byte) (salt, hash []byte, err error) {
	n := len(hashedPW) - c.saltLen
	if n < 0 {
		return nil, nil, ErrShortHash
	}
	return hashedPW[:c.saltLen], hashedPW[c.saltLen:], nil
}

// Sum returns the hash sum value of the password.
// Salt is joined at the left side of the returned sum.
// Use [Argon2i.Split] to split salt and the sum of the password.
func (c *Argon2i) Sum(password []byte) ([]byte, error) {
	salt := make([]byte, c.saltLen)
	_, err := rand.Read(salt)
	internal.MustNil(err)
	hashed := argon2.Key(password, salt, c.time, c.memory, c.threads, c.keyLen)
	return append(salt, hashed...), nil
}

// Equal reports if the given hashed password and the password
// are the same or not.
func (c *Argon2i) Equal(hashedPW, pw []byte) bool {
	return c.Compare(hashedPW, pw) == nil
}

// Compare compares hashed password and the password.
// If the hashed value of the pw matched to the hashedPW,
// it returns nil. If not matched, it returns non-nil error.
func (c *Argon2i) Compare(hashedPW, pw []byte) error {
	n := len(hashedPW) - c.saltLen
	if n < 0 {
		return ErrShortHash
	}
	salt, hashed := hashedPW[:c.saltLen], hashedPW[c.saltLen:]
	x := argon2.Key(pw, salt, c.time, c.memory, c.threads, c.keyLen)
	if !bytes.Equal(x, hashed) {
		return ErrNotMatch
	}
	return nil
}

// NewArgon2i returns a new instance of Argon2i hasher.
// Note that [NewArgon2id] is recommended from the security consideration than [NewArgon2i].
// The parameter saltLen, time, memory, threads and keyLen are passed to the [golang.org/x/crypto/argon2.IDKey].
//
// Parameter ranges:
//
//   - 0 <= saltLen
//   - 1 <= time
//   - 1 <= threads
//
// Recommended parameters:
//
//   - time=1,memory=2*1024*1024,threads=4,keyLen=32 (https://datatracker.ietf.org/doc/draft-ietf-kitten-password-storage/)
//   - time=1,memory=2*1024*1024,threads=4,keyLen=32 (https://datatracker.ietf.org/doc/rfc9106/)
//   - time=1,memory=64*1024,threads=4,keyLen=32 (https://datatracker.ietf.org/doc/rfc9106/)
//   - time=1,memory=46*1024,threads=1 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - time=2,memory=19*1024,threads=1 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - time=3,memory=12*1024,threads=1 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - time=4,memory=9*1024,threads=1 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
//   - time=5,memory=7*1024,threads=1 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
func NewArgon2id(saltLen int, time, memory uint32, threads uint8, keyLen uint32) (a *Argon2id, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("zargon2: invalid hash parameter [%v]", r)
		}
	}()
	c := &Argon2id{
		saltLen: saltLen,
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
	}
	_, _ = c.Sum([]byte("")) // Validate parameters.
	return c, nil
}

// Argon2id is the hasher type for [golang.org/x/crypto/argon2].
type Argon2id struct {
	saltLen int
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// Split splits the hashed password into salt and hash value of the password.
// If the hashedPW is too short, it returns [ErrShortHash].
func (c *Argon2id) Split(hashedPW []byte) (salt, hash []byte, err error) {
	n := len(hashedPW) - c.saltLen
	if n < 0 {
		return nil, nil, ErrShortHash
	}
	return hashedPW[:c.saltLen], hashedPW[c.saltLen:], nil
}

// Sum returns the hash sum value of the password.
// Salt is joined at the left side of the returned sum.
// Use [Argon2id.Split] to split salt and the sum of the password.
// [Argon2id.Sum] always returns nil as error.
func (c *Argon2id) Sum(password []byte) ([]byte, error) {
	salt := make([]byte, c.saltLen)
	_, err := rand.Read(salt)
	internal.MustNil(err)
	hashed := argon2.IDKey(password, salt, c.time, c.memory, c.threads, c.keyLen)
	return append(salt, hashed...), nil
}

// Equal reports if the given hashed password and the password
// are the same or not.
func (c *Argon2id) Equal(hashedPW, pw []byte) bool {
	return c.Compare(hashedPW, pw) == nil
}

// Compare compares hashed password and the password.
// If the hashed value of the pw matched to the hashedPW,
// it returns nil. If not matched, it returns non-nil error.
func (c *Argon2id) Compare(hashedPW, pw []byte) error {
	n := len(hashedPW) - c.saltLen
	if n < 0 {
		return ErrShortHash
	}
	salt, hashed := hashedPW[:c.saltLen], hashedPW[c.saltLen:]
	x := argon2.IDKey(pw, salt, c.time, c.memory, c.threads, c.keyLen)
	if !bytes.Equal(x, hashed) {
		return ErrNotMatch
	}
	return nil
}
