package zdebug_test

import (
	"bytes"
	"fmt"

	"github.com/aileron-projects/go/zruntime/zdebug"
)

func ExampleDumpTo() {
	val := struct {
		foo int
		bar string
	}{
		foo: 123,
		bar: "bar",
	}

	var buf bytes.Buffer
	zdebug.DumpTo(&buf, val)

	// Discard the first line which contains call location to avoid test failure.
	_, output, _ := bytes.Cut(buf.Bytes(), []byte("\n"))
	fmt.Println(string(output))
	// Output:
	//   | ┌── args[0]
	//   | (struct { foo int; bar string }) {
	//   |  foo: (int) 123,
	//   |  bar: (string) (len=3) "bar"
	//   | }
}
