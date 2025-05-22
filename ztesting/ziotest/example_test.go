package ziotest_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing/iotest"

	"github.com/aileron-projects/go/ztesting/ziotest"
)

func ExampleCharsetReader() {
	// r should returns "abc" repeatedly.
	r := ziotest.CharsetReader("abc", true)

	buf := make([]byte, 9)
	fmt.Println(r.Read(buf))
	fmt.Println(string(buf))
	// Output:
	// 9 <nil>
	// abcabcabc
}

func ExampleErrReader() {
	// er should return the first 10 characters from the r.
	r := strings.NewReader("abcdefghijklmnopqrstuvwxyz")
	er := ziotest.ErrReader(r, 10)

	buf := make([]byte, 4)
	fmt.Println(er.Read(buf))
	fmt.Println(er.Read(buf))
	fmt.Println(er.Read(buf))
	// Output:
	// 4 <nil>
	// 4 <nil>
	// 2 io: read/write on closed pipe
}

func ExampleErrReader_innerError() {
	// The internal reader r returns error.
	r := iotest.ErrReader(io.ErrUnexpectedEOF)
	er := ziotest.ErrReader(r, 10)

	buf := make([]byte, 4)
	fmt.Println(er.Read(buf))
	// Output:
	// 0 unexpected EOF
}

func ExampleErrReaderWith() {
	// er should return the EOF error after read 10 characters from the r.
	r := strings.NewReader("abcdefghijklmnopqrstuvwxyz")
	er := ziotest.ErrReaderWith(r, 10, io.EOF)

	buf := make([]byte, 4)
	fmt.Println(er.Read(buf))
	fmt.Println(er.Read(buf))
	fmt.Println(er.Read(buf))
	// Output:
	// 4 <nil>
	// 4 <nil>
	// 2 EOF
}

func ExampleErrWriter() {
	// ew should return an error after 10 characters written.
	ew := ziotest.ErrWriter(nil, 10)

	fmt.Println(ew.Write([]byte("abcd")))
	fmt.Println(ew.Write([]byte("abcd")))
	fmt.Println(ew.Write([]byte("abcd")))
	fmt.Println(ew.Write([]byte("abcd")))
	// Output:
	// 4 <nil>
	// 4 <nil>
	// 2 io: read/write on closed pipe
	// 0 io: read/write on closed pipe
}

func ExampleErrWriterWith() {
	// ew should return an EOF error after 10 characters written.
	ew := ziotest.ErrWriterWith(nil, 10, io.EOF)

	fmt.Println(ew.Write([]byte("abcd")))
	fmt.Println(ew.Write([]byte("abcd")))
	fmt.Println(ew.Write([]byte("abcd")))
	fmt.Println(ew.Write([]byte("abcd")))
	// Output:
	// 4 <nil>
	// 4 <nil>
	// 2 EOF
	// 0 EOF
}

func ExampleShortReader() {
	r := strings.NewReader("abcdefg")
	sr := ziotest.ShortReader(r, 3)

	buf := make([]byte, 10)
	n, err := sr.Read(buf)
	fmt.Println(n, string(buf[:n]), err)
	// Output:
	// 3 abc <nil>
}

func ExampleShortWriter() {
	w := bytes.NewBuffer(nil)
	sw := ziotest.ShortWriter(w, 3)

	n, err := sw.Write([]byte("abcdefg"))
	fmt.Println(n, w.String(), err)
	// Output:
	// 3 abc <nil>
}
