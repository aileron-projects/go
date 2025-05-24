package zerrors_test

import (
	"errors"
	"fmt"
	"io"

	"github.com/aileron-projects/go/zerrors"
)

func ExampleAttrs() {
	e1 := errors.New("example1")
	e2 := errors.New("example2")

	fmt.Println(zerrors.Attrs(e1))
	fmt.Println(zerrors.Attrs(fmt.Errorf("example3 [%w]", e1)))
	fmt.Println(zerrors.Attrs(errors.Join(e1, e2)))
	// Output:
	// map[msg:example1]
	// map[msg:example3 [example1] wraps:map[msg:example1]]
	// map[err1:map[msg:example1] err2:map[msg:example2]]
}

func ExampleDefinition_New() {
	// Define an error with detail template.
	def := zerrors.NewDefinition("E123", "pkgX", "example error", "foo=%s bar=%s")

	fmt.Println(def.New(nil, "FOO", "BAR").Error())        // With arguments.
	fmt.Println(def.New(nil).Error())                      // No arguments.
	fmt.Println(def.New(nil, "FOO").Error())               // Insufficient arguments.
	fmt.Println(def.New(nil, "FOO", "BAR", "BAZ").Error()) // Too many arguments.
	fmt.Println(def.New(io.EOF, "FOO", "BAR").Error())     // With inner error.
	// Output:
	// E123: pkgX: example error: foo=FOO bar=BAR
	// E123: pkgX: example error: foo=%!s(MISSING) bar=%!s(MISSING)
	// E123: pkgX: example error: foo=FOO bar=%!s(MISSING)
	// E123: pkgX: example error: foo=FOO bar=BAR%!(EXTRA string=BAZ)
	// E123: pkgX: example error: foo=FOO bar=BAR [EOF]
}
