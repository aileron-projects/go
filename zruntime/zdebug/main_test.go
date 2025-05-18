package zdebug

import (
	"os"
	"testing"
	"time"
)

var (
	ExportedDumpEnabled = dumpEnabled
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	tmpFunc := timeNow
	timeNow = func() time.Time { return time.Unix(0, 0) }
	tmpLoc := time.Local
	time.Local = time.UTC
	defer func() {
		timeNow = tmpFunc
		time.Local = tmpLoc
	}()

	return m.Run()
}
