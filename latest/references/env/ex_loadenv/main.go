package main

import (
	"github.com/aileron-projects/go/zos"
)

func main() {
	if err := zos.LoadEnv("env.txt"); err != nil {
		panic(err)
	}
}
