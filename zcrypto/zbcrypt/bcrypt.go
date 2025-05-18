package zbcrypt

import (
	"github.com/aileron-projects/go/internal/ihash"
	"golang.org/x/crypto/bcrypt"
)

var (
	_ ihash.PWHasher = &BCrypt{}
)

// New returns a new instance of BCrypt hasher with given parameters.
// If the cost is less than [bcrypt.MinCost]=4, then the [bcrypt.DefaultCost]=10 is used.
// If the cost is larger than [bcrypt.MaxCost]=31, an error is returned.
//
// Parameter ranges:
//
//   - 4 <= cost <= 31. It is [bcrypt.MinCost] <= cost <= [bcrypt.MaxCost].
//
// Recommended parameters:
//
//   - cost>=12 (https://datatracker.ietf.org/doc/draft-ietf-kitten-password-storage/)
//   - cost>=10 (https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
func New(cost int) (*BCrypt, error) {
	c := &BCrypt{
		cost: cost,
	}
	if _, err := c.Sum([]byte("")); err != nil { // Validate parameters.
		return nil, err
	}
	return c, nil
}

// BCrypt is the hasher type for [golang.org/x/crypto/bcrypt].
type BCrypt struct {
	cost int
}

// Sum returns the hash sum value of the password.
func (c *BCrypt) Sum(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, c.cost)
}

// Equal reports if the given hashed password and the password
// are the same or not.
func (c *BCrypt) Equal(hashedPW, pw []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPW, pw) == nil
}

// Compare compares hashed password and the password.
// If the hashed value of the pw matched to the hashedPW,
// it returns nil. If not matched, it returns non-nil error.
func (c *BCrypt) Compare(hashedPW, pw []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPW, pw)
}
