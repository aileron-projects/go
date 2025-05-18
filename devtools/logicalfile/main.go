package main

import (
	"fmt"
	"time"

	"github.com/aileron-projects/go/zlog"
)

func main() {
	c := &zlog.LogicalFileConfig{
		Manager: &zlog.FileManagerConfig{
			GzipLv:     5,
			MaxHistory: 5,
			// MaxAge:     30 * time.Second,
			MaxTotalBytes: 1000,
			Pattern:       "app.%u.log",
			SrcDir:        "./src",
			DstDir:        "./dst",
		},
		RotateBytes: 500,
		FileName:    "app.log",
	}
	f, err := zlog.NewLogicalFile(c)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Fprintln(f, time.Now(), ": TEST LOG ", "1234567890123456789012345678901234567890")
		time.Sleep(500 * time.Millisecond)
	}
}
