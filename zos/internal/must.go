package internal

// MustNil panics if err is not nil.
func MustNil(err error) {
	if err != nil {
		panic(err)
	}
}
