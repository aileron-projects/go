package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aileron-projects/go/zos"
)

var (
	file = flag.String("file", "env.txt", "env file path to check")
)

func main() {
	flag.Parse()
	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}
	b, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}
	kvs, err := zos.LoadEnv(b)
	if err != nil {
		panic(err)
	}
	fmt.Println("Env parse result of:", *file)
	fmt.Println("---------------------------")
	fmt.Println("Number of variables:", len(kvs))
	for k, v := range kvs {
		fmt.Println("")
		fmt.Println(">> KEY:", k)
		if strings.Contains(v, "\n") {
			fmt.Println(">> VALUE: ⬎")
			fmt.Println(v)
		} else {
			fmt.Println(">> VALUE:", v)
		}
	}
}
