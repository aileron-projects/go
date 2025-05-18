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

// Wrap returns a new error that wraps the given inner error.
// Wrap returns nil if the given inner error is nil.
func Wrap(inner error, msg string) *NestErr {
	if inner == nil {
		return nil
	}
	return &NestErr{
		Err: inner,
		Msg: msg,
	}
}

// NestErr is the error type that wraps an error
// with another error.
type NestErr struct {
	// Err is the inner, or wrapped error.
	// Err should not be nil.
	// Otherwise, [NestErr.Error] will panics.
	Err error
	// Msg is outer error message.
	Msg string
}

func (e *NestErr) Error() string {
	return e.Msg + " [" + e.Err.Error() + "]"
}

func (e *NestErr) Unwrap() error {
	return e.Err
}
