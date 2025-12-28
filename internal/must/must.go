package must

// Nil ensures the given error is nil.
// MustNil is a defensive assertion used primarily in security-sensitive areas.
// It ensures that an error, which "should never" occur, is truly absent.
// This acts as a safeguard against future changes in upstream APIs.
//
// Example:
//
//	buf := make([]byte, 10)
//	_, err := rand.Read(buf) // Never returns a non-nil error.
//	internal.MustNil(err)    // But we check anyway, just in case.
func Nil(err error) {
	if err != nil {
		panic(err)
	}
}
