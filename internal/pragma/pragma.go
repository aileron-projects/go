package pragma

import "sync"

// NoUnkeyedLiterals can be embedded in structs to prevent unkeyed literals.
//
// References:
//   - https://pkg.go.dev/google.golang.org/protobuf/internal/pragma#NoUnkeyedLiterals
//
// Example:
//
//	type Foo struct {
//	  NoUnkeyedLiterals
//	  name string
//	  age  int
//	}
//
//	f1 := Foo{"alice", 25} // NG. Unkeyed instantiation is not allowed.
//	f2 := Foo{name: "alice", age: 25} // OK.
type NoUnkeyedLiterals struct{}

// DoNotImplement can be embedded in interfaces to prevent trivial
// implementations of the interface.
// This is useful to prevent unauthorized implementations of an interface.
//
// References:
//   - https://pkg.go.dev/google.golang.org/protobuf/internal/pragma#DoNotImplement
//
// Example:
//
//	// Hello interface cannot be implicitly implemented.
//	type Hello interface {
//	  DoNotImplement
//	  HelloWorld()
//	}
//
//	// Foo does not implement Hello interface.
//	type Foo struct {
//	}
//	func (f *Foo) HelloWorld() {}
//
//	// Bar does implement Hello interface.
//	type Bar struct {
//	  Hello
//	}
//	func (b *Bar) HelloWorld() {}
type DoNotImplement interface{ DoNotCallMe(DoNotImplement) }

// DoNotCompare can be embedded in structs to prevent comparability.
//
// References:
//   - https://pkg.go.dev/google.golang.org/protobuf/internal/pragma#DoNotCompare
//
// Example:
//
//	// Foo should not be compared.
//	type Foo struct {
//	  DoNotCompare
//	}
//	f1 := Foo{}
//	f2 := Foo{}
//	println(f1==f2) // Error. Comparing Foo is not allowed.
type DoNotCompare [0]func()

// DoNotCopy can be embedded in structs to help prevent shallow copies.
// This does not rely on a Go language feature, but rather a special case
// within the vet checker.
//
// References:
//   - https://pkg.go.dev/google.golang.org/protobuf/internal/pragma#DoNotCopy
//   - https://golang.org/issues/8005.
//
// Example:
//
//	// Foo should not be shallow copied.
//	type Foo struct {
//	  DoNotCopy
//	}
//	f1 := Foo{}
//	f2 := f1 // Warning. Shallow copy is not allowed.
type DoNotCopy [0]sync.Mutex
