package main

import (
	"io"

	"github.com/aileron-projects/go/zerrors"
)

var (
	ErrExample1 = zerrors.NewDefinition("E001", "KindSystemError", "ErrFoo", "error foo.", "detail foo.")
	ErrExample2 = zerrors.NewDefinition("E002", "KindSystemError", "ErrBar", "error bar.", "detail bar.")
)

// main tests tracing errors. Use following tags when building.
//   - go build -tags zerrorstrace ./main.go
func main() {
	err1 := ErrExample1.New(io.EOF)
	err2 := ErrExample2.NewStack(err1)
	_ = err2
}
