package main

import (
	"io"

	"github.com/aileron-projects/go/zerrors"
)

var (
	ErrExample1 = zerrors.NewDefinition("E001", "ErrFoo", "error foo.", map[string]string{"foo": "FOO"})
	ErrExample2 = zerrors.NewDefinition("E002", "ErrBar", "error bar.", map[string]string{"bar": "BAR"})
)

// main tests tracing errors. Use following tags when building.
//   - go build -tags zerrorstrace ./main.go
func main() {
	err1 := ErrExample1.New(io.EOF)
	err2 := ErrExample2.NewStack(err1)
	_ = err2
}
