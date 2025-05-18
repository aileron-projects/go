package zerrors

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kr/pretty"
)

var (
	// writer is the default output destination.
	writer io.Writer = os.Stdout
	// timeNow returns current time.
	// timeNow should be replaced only when testing.
	timeNow func() time.Time = time.Now
)

func init() {
	if traceEnabled {
		initTrace(os.Getenv("GO_ZERRORS"), os.CreateTemp)
	}
}

// initTrace initialize error trace output.
// The argument writeTo means
//   - "stdout": Standard output.
//   - "stderr": Standard error output.
//   - "discard": Discard all output.
//   - "file": File output. The second arg createFile must be specified.
//   - Others: Ignored.
func initTrace(writeTo string, createFile func(dir string, pattern string) (*os.File, error)) {
	fmt.Println("Error tracing is enabled (-tags zerrorstrace)")
	switch strings.ToLower(writeTo) {
	case "stdout":
		writer = os.Stdout
		fmt.Println("Trace messages will be output to the standard out.")
	case "stderr":
		writer = os.Stderr
		fmt.Println("Trace messages will be output to the standard error.")
	case "discard":
		writer = io.Discard
		fmt.Println("Trace messages will be discarded.")
	case "file":
		f, err := createFile("", "zerrors-*")
		if err != nil {
			panic(err)
		}
		writer = f
		fmt.Println("Trace messages will be output to the " + f.Name())
	}
}

// traceTo prints error information.
// Use build time tags "-tags zerrorstrace" to enable output.
// The global variable writer is used if the first argument w is nil.
// The second argument e must not be nil. Otherwise panics.
func traceTo(w io.Writer, e *Error) {
	if !traceEnabled {
		return // Tracing is disabled.
	}
	f := callerFrames(2, 1)
	loc := ""
	if len(f) > 0 {
		f0 := f[0]
		loc = " Pkg:" + f0.Pkg + " File:" + f0.File + " Func:" + f0.Func + " Line:" + strconv.Itoa(f0.Line)
	}
	var buf bytes.Buffer
	_, _ = pretty.Fprintf(&buf, "%# v", Attrs(e))
	_ = buf.WriteByte('\n')

	w = cmp.Or(w, writer, io.Writer(os.Stdout))
	_, _ = w.Write([]byte(timeNow().Format(time.DateTime) + " [TRACE]" + loc + "\n"))
	_, _ = writeWithPrefix(w, buf.Bytes())
}

func writeWithPrefix(w io.Writer, p []byte) (n int, err error) {
	pp := p
	for len(pp) > 0 {
		nn, _ := w.Write([]byte("  | "))
		n += nn
		i := bytes.IndexByte(pp, '\n')
		if i < 0 {
			nn, _ := w.Write(pp)
			n += nn
			return n, nil
		} else {
			nn, _ := w.Write(pp[:i+1])
			n += nn
			pp = pp[i+1:]
		}
	}
	return n, nil
}
