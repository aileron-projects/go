package zerrors

// UnwrapErr returns the result of calling the Unwrap method on err.
// If the given err implements Unwrap() that returns an error.
// Otherwise, UnwrapErr returns nil.
//
// UnwrapErr only calls a method of the form "Unwrap() error".
// In particular UnwrapErr does not unwrap errors returned by [errors.Join].
// See also [UnwrapErrs] and [errors.Unwrap].
func UnwrapErr(err error) error {
	u, ok := err.(interface{ Unwrap() error })
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// UnwrapErrs returns the result of calling the Unwrap method on err.
// If the given err implements Unwrap() that returns a []error.
// Otherwise, UnwrapErrs returns nil slice.
//
// UnwrapErrs only calls a method of the form "Unwrap() []error".
// In particular UnwrapErrs does not unwrap errors returned by [Wrap].
// UnwrapErrs can unwrap errors returned by [errors.Join].
// See also [UnwrapErr] and [errors.Unwrap].
func UnwrapErrs(err error) []error {
	u, ok := err.(interface{ Unwrap() []error })
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// Must panics if the given err is not nil.
// Must exits with panic(err) if the second argument err is not nil.
// If the err is nil, the value t is returned.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// MustNil panics if the given err is not nil.
// MustNil exits with panic(err) if the given err is not nil.
func MustNil(err error) {
	if err != nil {
		panic(err)
	}
}
