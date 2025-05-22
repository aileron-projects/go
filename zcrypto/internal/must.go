package internal

// MustNil panics if err is not nil.
// MustNil is a defensive assertion used primarily in security-sensitive areas.
// It ensures that an error, which "should never" occur, is truly absent.
// This acts as a safeguard against future changes in upstream APIs.
//
// This is intended to be used with [crypto/rand.Read].
// As noted in the documentation for [crypto/rand.Read], it does not return
// a non-nil error. While we believe this interface will never change,
// we include this check just to be absolutely sure — for security reasons.
//
// Example:
//
//	buf := make([]byte, 10)
//	_, err := rand.Read(buf) // Never returns a non-nil error.
//	internal.MustNil(err)    // But we check anyway, just in case.
func MustNil(err error) {
	if err != nil {
		panic(err)
	}
}
