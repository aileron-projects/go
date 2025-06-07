package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/zlog/zslog"
)

func main() {
	c := &zlog.LogicalFileConfig{
		Manager: &zlog.FileManagerConfig{
			MaxHistory: 5,
			Pattern:    "app.%i.log",
		},
		RotateBytes: 500,       // Max size of a single file.
		FileName:    "app.log", // Active file name.
	}
	f, err := zlog.NewLogicalFile(c)
	if err != nil {
		panic(err)
	}

	h := slog.NewTextHandler(f, nil) // Create slog handler with f.
	lg := zslog.New(h)
	for {
		lg.InfoContext(context.Background(), "log message", "now", time.Now())
		time.Sleep(time.Second)
	}
}
