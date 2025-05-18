package zruntime

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestReportErr(t *testing.T) {
	// No parallel because tests replace a global variable.

	tmp := reportTo
	defer func() {
		reportTo = tmp
	}()

	t.Run("report nil", func(t *testing.T) {
		var buf bytes.Buffer
		reportTo = &buf
		ReportErr(nil, "description")
		ztesting.AssertEqual(t, "report is not empty", "", buf.String())
	})
	t.Run("report error", func(t *testing.T) {
		var buf bytes.Buffer
		reportTo = &buf
		ReportErr(io.EOF, "description")
		ztesting.AssertEqual(t, "report does not contain error", true, strings.Contains(buf.String(), ">> Error   : EOF"))
	})
}
