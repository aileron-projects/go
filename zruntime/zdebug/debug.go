package zdebug

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aileron-projects/go/zruntime"
	"github.com/davecgh/go-spew/spew"
)

var (
	// writer is the default output destination
	// used by dump and diff.
	writer io.Writer = os.Stdout
	// timeNow returns current time.
	// timeNow should be replaced only when testing.
	timeNow func() time.Time = time.Now
	// DumpConfig is the dump output format configuration.
	// See the documents on [spew.ConfigState].
	DumpConfig = spew.ConfigState{Indent: " "}
)

func init() {
	initDebug(os.Getenv("GO_ZDEBUG"), os.CreateTemp)
}

// initDebug initialize debug output.
// The argument writeTo means
//   - "stdout": Standard output.
//   - "stderr": Standard error output.
//   - "discard": Discard all output.
//   - "file": File output. The second arg createFile must be specified.
//   - Others: Ignored.
func initDebug(writeTo string, createFile func(dir string, pattern string) (*os.File, error)) {
	if dumpEnabled {
		fmt.Println("Debugging is enabled (-tags zdebugdump)")
	}
	switch strings.ToLower(writeTo) {
	case "stdout":
		writer = os.Stdout
		fmt.Println("Debug messages will be output to the standard out.")
	case "stderr":
		writer = os.Stderr
		fmt.Println("Debug messages will be output to the standard error.")
	case "discard":
		writer = io.Discard
		fmt.Println("Debug messages will be discarded.")
	case "file":
		f, err := createFile("", "zdebug-*")
		if err != nil {
			panic(err)
		}
		writer = f
		fmt.Println("Debug messages will be output to the " + f.Name())
	}
}

type prefixWriter struct {
	w       io.Writer
	newline bool
	prefix  []byte
}

func (pw *prefixWriter) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		if i := bytes.IndexByte(p, '\n'); i >= 0 {
			if pw.newline {
				nn, _ := pw.w.Write(pw.prefix)
				n += nn
			}
			nn, _ := pw.w.Write(p[:i+1])
			n += nn
			p = p[i+1:]
			pw.newline = true
		} else {
			if pw.newline {
				nn, _ := pw.w.Write(pw.prefix)
				n += nn
			}
			nn, _ := pw.w.Write(p)
			n += nn
			p = nil
			pw.newline = false
		}
	}
	return n, nil
}

var (
	// HookDumpFunc hooks [Dump] function.
	// HookDumpFunc is called when applications are run or built
	// with runtime dump enabled.
	// HookDumpFunc receives [io.Writer] that should be used for log output
	// if necessary and caller fame information in [zruntime.Frame] format.
	// In addition to them, all arguments given to the [Dump] are also given.
	// If HookDumpFunc returned true, the [Dump] function immediately returns
	// without running default procession of if.
	// Use build tag `//go:build !zdebugdump` to enable runtime dump.
	HookDumpFunc func(io.Writer, zruntime.Frame, ...any) bool = nil
)

// DumpTo writes dump results into the given writer.
// The writer will be reset after DumpTo returned.
// Unlike the [Dump], DumpTo does not requires build tag.
func DumpTo(w io.Writer, msg string, a ...any) {
	w = cmp.Or(w, writer, io.Writer(os.Stderr))
	f := zruntime.ConvertFrame(zruntime.CallerFrame(2))
	if HookDumpFunc != nil {
		if HookDumpFunc(w, f, a...) {
			return
		}
	}

	caller := "Pkg=" + f.Pkg + " File=" + f.File + " Func=" + f.Func + " Line=" + strconv.Itoa(f.Line)
	_, _ = w.Write([]byte(timeNow().Format(time.DateTime) + " [DUMP] " + msg + "\n"))
	_, _ = w.Write([]byte("  | Caller: " + caller + "\n"))

	pw := &prefixWriter{
		w:       w,
		newline: true,
		prefix:  []byte("  | "),
	}
	dc := &DumpConfig
	for i := range a {
		_, _ = pw.Write([]byte("┌── args[" + strconv.Itoa(i) + "]\n"))
		dc.Fdump(pw, a[i])
	}
}

// DumpAlways prints object dump of the given values.
// Unlike [Dump], it does not requires build tags.
func DumpAlways(msg string, a ...any) {
	DumpTo(writer, msg, a...)
}

// Dump prints object dump of the given values.
// Use build time tags "-tags zdebugdump" to enable output.
// Stdout is used as default output destination.
// If callers who wants to debug temporarily without build tags,
// use [DumpAlways] or [DumpTo] instead.
func Dump(msg string, a ...any) {
	if dumpEnabled {
		DumpTo(writer, msg, a...)
	}
}
