package zruntime

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var (
	// mu protects reportTo
	mu = sync.Mutex{}
	// reportTo is the destination where
	// system errors are written to.
	reportTo io.Writer = os.Stderr
)

// ReportErr outputs the given error with stack trace
// to the [os.Stderr].
// ReportErr does nothing if the given error is nil.
// Caller can pass additional error message.
// ReportErr should be used only for error that callers
// do not know how they handle if.
func ReportErr(err error, msg string) {
	if err == nil {
		return
	}

	f := ConvertFrame(CallerFrame(1))
	loc := " Pkg:" + f.Pkg + " File:" + f.File + " Func:" + f.Func + " Line:" + strconv.Itoa(f.Line)
	stack := make([]byte, 1<<13) // Read max 8kiB.
	n := runtime.Stack(stack, false)

	mu.Lock()
	defer mu.Unlock()

	_, _ = reportTo.Write([]byte(time.Now().Format(time.DateTime) + " [REPORT]" + loc + "\n"))
	_, _ = reportTo.Write([]byte(">> Message : " + msg + "\n"))
	_, _ = reportTo.Write([]byte(">> Error   : " + err.Error() + "\n"))
	_, _ = reportTo.Write([]byte(">> Dump    : " + spew.Sdump(err)))
	const prefix = "   | "
	_, _ = reportTo.Write([]byte(">> Stack Trace:\n" + prefix))
	_, _ = reportTo.Write(bytes.ReplaceAll(stack[:n], []byte("\n"), []byte("\n"+prefix)))
	_, _ = reportTo.Write([]byte("\n"))
}
