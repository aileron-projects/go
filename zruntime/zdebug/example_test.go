package zdebug_test

import (
	"bytes"
	"fmt"
	"strings"

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
	zdebug.DumpTo(&buf, "", val)

	output := buf.String()
	_, output, _ = strings.Cut(output, "\n") // Discard first line.
	_, output, _ = strings.Cut(output, "\n") // Discard second line.
	fmt.Println(string(output))
	// Output:
	//   | ┌── args[0]
	//   | (struct { foo int; bar string }) {
	//   |  foo: (int) 123,
	//   |  bar: (string) (len=3) "bar"
	//   | }
}
