package main

import (
	"fmt"
	"time"

	"github.com/aileron-projects/go/zlog"
)

func main() {
	c := &zlog.LogicalFileConfig{
		Manager: &zlog.FileManagerConfig{
			MaxAge:  30 * time.Second,
			Pattern: "app.%u.log",
		},
		RotateBytes: 500,       // Max size of a single file.
		FileName:    "app.log", // Active file name.
	}
	f, err := zlog.NewLogicalFile(c)
	if err != nil {
		panic(err)
	}

	initial := time.Now()
	for {
		fmt.Println("Now:", time.Now().Unix(), "\t", "30s before:", time.Now().Unix()-30)
		fmt.Fprintln(f, time.Now(), "Time past ", time.Since(initial))
		time.Sleep(time.Second)
	}
}
