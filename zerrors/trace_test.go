package zerrors

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestInitTrace(t *testing.T) {
	if !traceEnabled {
		return
	}
	defer func() {
		writer = os.Stdout // Reset
	}()

	testCases := map[string]struct {
		writeTo string
		wants   io.Writer
	}{
		"empty": {
			writeTo: "",
			wants:   os.Stdout,
		},
		"undefined": {
			writeTo: "undefined",
			wants:   os.Stdout,
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
			writer = os.Stdout // Reset to Stdout before call initDebug.
			initTrace(tc.writeTo, func(dir, pattern string) (*os.File, error) {
				return os.Stdin, nil
			})
			ztesting.AssertEqual(t, "writer not matched.", tc.wants, writer)
		})
	}
}

func TestInitTrace_panics(t *testing.T) {
	if !traceEnabled {
		return
	}
	t.Parallel()
	err := errors.New("file create error")
	defer func() {
		rec := recover()
		if rec != err {
			t.Errorf("error not matched. want:%#v got:%#v", err, rec)
		}
	}()
	initTrace("file", func(dir, pattern string) (*os.File, error) {
		return nil, err
	})
}

func TestWriteWithPrefix(t *testing.T) {
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
			wants: "  | \n",
		},
		"CRLF only": {
			write: [][]byte{[]byte("\r\n")},
			wants: "  | \r\n",
		},
		"multiple LFs": {
			write: [][]byte{[]byte("\n\n\n\n\n")},
			wants: "  | \n  | \n  | \n  | \n  | \n",
		},
		"multiple CRLFs": {
			write: [][]byte{[]byte("\r\n\r\n\r\n")},
			wants: "  | \r\n  | \r\n  | \r\n",
		},
		"1 line without LF": {
			write: [][]byte{[]byte("abc")},
			wants: "  | abc",
		},
		"1 line ends with LF": {
			write: [][]byte{[]byte("abc\n")},
			wants: "  | abc\n",
		},
		"1 line starts with LF": {
			write: [][]byte{[]byte("\nabc")},
			wants: "  | \n  | abc",
		},
		"1 line ends with CRLF": {
			write: [][]byte{[]byte("abc\r\n")},
			wants: "  | abc\r\n",
		},
		"1 line starts with CRLF": {
			write: [][]byte{[]byte("\r\nabc")},
			wants: "  | \r\n  | abc",
		},
		"2 lines with single LF": {
			write: [][]byte{[]byte("abc"), []byte("def\n")},
			wants: "  | abc  | def\n",
		},
		"2 lines with LFs": {
			write: [][]byte{[]byte("abc\n"), []byte("def\n")},
			wants: "  | abc\n  | def\n",
		},
		"2 lines with single CRLF": {
			write: [][]byte{[]byte("abc"), []byte("def\r\n")},
			wants: "  | abc  | def\r\n",
		},
		"2 lines with CRLFs": {
			write: [][]byte{[]byte("abc\r\n"), []byte("def\r\n")},
			wants: "  | abc\r\n  | def\r\n",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			for _, v := range tc.write {
				writeWithPrefix(&buf, v)
			}
			result := buf.String()
			ztesting.AssertEqual(t, "written string does not match.", tc.wants, result)
		})
	}
}

func TestTraceTo(t *testing.T) {
	if !traceEnabled {
		return
	}
	t.Parallel()
	testCases := map[string]struct {
		err   *Error
		wants []string
	}{
		"empty": {
			err:   &Error{},
			wants: []string{`1970-01-01 00:00:00 [TRACE]`, `"msg":  ""`, `"pkg":  ""`},
		},
		"with pkg": {
			err:   &Error{Pkg: "test"},
			wants: []string{`1970-01-01 00:00:00 [TRACE]`, `"pkg":  "test"`},
		},
		"with msg": {
			err:   &Error{Msg: "test"},
			wants: []string{`1970-01-01 00:00:00 [TRACE]`, `"msg":  "test"`},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			traceTo(&buf, tc.err)
			result := buf.String()
			for _, w := range tc.wants {
				if !strings.Contains(result, "1970-01-01 00:00:00 [TRACE]") {
					t.Error("dump result does not contain date and time. got:\n" + result)
				}
				if !strings.Contains(result, w) {
					t.Error("expected dump result not output. want: `" + w + "` got:\n" + result)
				}
			}
		})
	}
}
