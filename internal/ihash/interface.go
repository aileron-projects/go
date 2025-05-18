package ihash

// PWHasher is the interface for password hash.
type PWHasher interface {
	Sum(password []byte) ([]byte, error)
	Equal(hashedPW, pw []byte) bool
	Compare(hashedPW, pw []byte) error
}

// SumFunc returns hash sum of the b.
type SumFunc func(b []byte) []byte

// EqualSumFunc returns if the sum and the hash of b matches or not.
// Unlike [CompareSumFunc], it returns true or false.
type EqualSumFunc func(b, sum []byte) bool

// CompareSumFunc returns if the sum and the hash of b matches or not.
// If matches, it returns nil error.
// If not matches, it returns non-nil error.
type CompareSumFunc func(b, sum []byte) error

// HMACSumFunc returns HMAC of the msg.
type HMACSumFunc func(msg, key []byte) []byte

// HMACEqualSumFunc returns if the sum and the HMAC of msg matches or not.
// Unlike [HMACCompareSumFunc], it returns true or false.
type HMACEqualSumFunc func(msg, key, sum []byte) bool

// HMACCompareSumFunc returns if the sum and the HMAC of msg matches or not.
// If matches, it returns nil error.
// If not matches, it returns non-nil error.
type HMACCompareSumFunc func(msg, key, sum []byte) error
