package zdebug

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestInitDebug(t *testing.T) {
	origWriter := writer
	defer func() { writer = origWriter }()
	testCases := map[string]struct {
		writeTo string
		wants   io.Writer
	}{
		"empty": {
			writeTo: "",
			wants:   nil, // Nothing set.
		},
		"undefined": {
			writeTo: "undefined",
			wants:   nil, // Nothing set.
		},
		"stdout": {
			writeTo: "stdout",
			wants:   os.Stdout,
		},
		"stderr": {
			writeTo: "stderr",
			wants:   os.Stderr,
		},
		"discard": {
			writeTo: "discard",
			wants:   io.Discard,
		},
		"file": {
			writeTo: "file",
			wants:   os.Stdin, // Whichever. Use available Stdin for here.
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			writer = nil // Reset before initDebug.
			initDebug(tc.writeTo, func(dir, pattern string) (*os.File, error) {
				return os.Stdin, nil
			})
			ztesting.AssertEqual(t, "writer not matched.", tc.wants, writer)
		})
	}
}

func TestInitDebug_panics(t *testing.T) {
	t.Parallel()
	err := errors.New("file create error")
	defer func() {
		rec := recover()
		ztesting.AssertEqual(t, "error not matched.", err, rec.(error))
	}()
	initDebug("file", func(dir, pattern string) (*os.File, error) {
		return nil, err
	})
}

func TestPrefixWriter(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		write [][]byte
		wants string
	}{
		"empty": {
			write: nil,
			wants: "",
		},
		"LF only": {
			write: [][]byte{[]byte("\n")},
			wants: "\n",
		},
		"CRLF only": {
			write: [][]byte{[]byte("\r\n")},
			wants: "\r\n",
		},
		"multiple LFs": {
			write: [][]byte{[]byte("\n\n\n\n\n")},
			wants: "\n\n\n\n\n",
		},
		"multiple CRLFs": {
			write: [][]byte{[]byte("\r\n\r\n\r\n")},
			wants: "\r\n\r\n\r\n",
		},
		"1 line without LF": {
			write: [][]byte{[]byte("abc")},
			wants: "abc",
		},
		"1 line ends with LF": {
			write: [][]byte{[]byte("abc\n")},
			wants: "abc\n",
		},
		"1 line starts with LF": {
			write: [][]byte{[]byte("\nabc")},
			wants: "\nabc",
		},
		"1 line ends with CRLF": {
			write: [][]byte{[]byte("abc\r\n")},
			wants: "abc\r\n",
		},
		"1 line starts with CRLF": {
			write: [][]byte{[]byte("\r\nabc")},
			wants: "\r\nabc",
		},
		"2 lines with single LF": {
			write: [][]byte{[]byte("abc"), []byte("def\n")},
			wants: "abcdef\n",
		},
		"2 lines with LFs": {
			write: [][]byte{[]byte("abc\n"), []byte("def\n")},
			wants: "abc\ndef\n",
		},
		"2 lines with single CRLF": {
			write: [][]byte{[]byte("abc"), []byte("def\r\n")},
			wants: "abcdef\r\n",
		},
		"2 lines with CRLFs": {
			write: [][]byte{[]byte("abc\r\n"), []byte("def\r\n")},
			wants: "abc\r\ndef\r\n",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			w := &prefixWriter{w: &buf, newline: true}
			for _, v := range tc.write {
				w.Write(v)
			}
			result := buf.String()
			ztesting.AssertEqual(t, " does not match.", tc.wants, result)
		})
	}
}

func TestDump(t *testing.T) {
	var buf bytes.Buffer
	tmp := writer
	writer = &buf
	defer func() { writer = tmp }()
	Dump(int(123))
	result := buf.String()
	if dumpEnabled {
		ztesting.AssertEqual(t, "dump does not contain datetime", true, strings.Contains(result, "1970-01-01 00:00:00 [DUMP]"))
	} else {
		ztesting.AssertEqual(t, "dump should be empty", "", result)
	}
}

func TestDumpAlways(t *testing.T) {
	var buf bytes.Buffer
	tmp := writer
	writer = &buf
	defer func() { writer = tmp }()
	DumpAlways(int(123))
	result := buf.String()
	ztesting.AssertEqual(t, "dump does not contain datetime", true, strings.Contains(result, "1970-01-01 00:00:00 [DUMP]"))

}
