package main

import (
	"io"

	"github.com/aileron-projects/go/zerrors"
)

var (
	ErrExample1 = zerrors.NewDefinition("E001", "main", "example error 01", "detail")
	ErrExample2 = zerrors.NewDefinition("E002", "main", "example error 02", "")
)

// main tests tracing errors. Use following tags when building.
//   - go build -tags zerrorstrace ./main.go
func main() {
	err1 := ErrExample1.New(io.EOF)
	err2 := ErrExample2.NewStack(err1)
	_ = err2
}
