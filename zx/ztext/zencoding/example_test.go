package zencoding_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/aileron-projects/go/zx/ztext/zencoding"
	"golang.org/x/text/encoding/japanese"
)

func ExampleNewDecodeReader() {
	d := japanese.ShiftJIS.NewDecoder()
	r := strings.NewReader("\xe0\x9f") // ShiftJIS string.

	rr := zencoding.NewDecodeReader(d, r)
	b, _ := io.ReadAll(rr) // UTF-8 string obtained.
	fmt.Println(string(b))
	// Output:
	// 燹
}

func ExampleNewDecodeWriter() {
	d := japanese.ShiftJIS.NewDecoder()
	var w bytes.Buffer

	ww := zencoding.NewDecodeWriter(d, &w)
	ww.Write([]byte("\xe0\x9f")) // Write ShiftJIS string.
	fmt.Println(w.String())      // UTF-8 string obtained.
	// Output:
	// 燹
}

func ExampleNewEncodeReader() {
	e := japanese.ShiftJIS.NewEncoder()
	r := strings.NewReader("燹") // UTF-8 string.

	rr := zencoding.NewEncodeReader(e, r)
	b, _ := io.ReadAll(rr) // ShiftJIS bytes.
	fmt.Printf("%2x", b)
	// Output:
	// e09f
}

func ExampleNewEncodeWriter() {
	e := japanese.ShiftJIS.NewEncoder()
	var w bytes.Buffer

	ww := zencoding.NewEncodeWriter(e, &w)
	ww.Write([]byte("燹"))        // UTF-8 string.
	fmt.Printf("%2x", w.Bytes()) // ShiftJIS bytes.
	// Output:
	// e09f
}
