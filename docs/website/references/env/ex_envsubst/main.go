package main

import (
	"fmt"
	"os"

	"github.com/aileron-projects/go/zos"
)

func main() {
	os.Setenv("ABC", "abcdefg")
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "BAR")
	os.Setenv("ARR_X", "xxx")
	os.Setenv("ARR_Y", "yyy")

	must := func(b []byte, err error) string {
		if err != nil {
			panic(err)
		}
		return string(b)
	}
	fmt.Println("${FOO} ------------", must(zos.ResolveEnv([]byte("${FOO}"))))
	fmt.Println("${BAZ:-default} ---", must(zos.ResolveEnv([]byte("${BAZ:-default}"))))
	fmt.Println("${BAZ-default}  ---", must(zos.ResolveEnv([]byte("${BAZ-default}"))))
	fmt.Println("${BAZ:=default} ---", must(zos.ResolveEnv([]byte("${BAZ:=default}"))))
	fmt.Println("${BAZ=default}  ---", must(zos.ResolveEnv([]byte("${BAZ=default}"))))
	fmt.Println("${BAZ:?default} ---", must(zos.ResolveEnv([]byte("${BAZ:?default}"))))
	fmt.Println("${BAZ?default}  ---", must(zos.ResolveEnv([]byte("${BAZ?default}"))))
	fmt.Println("${BAZ:+default} ---", must(zos.ResolveEnv([]byte("${BAZ:+default}"))))
	fmt.Println("${BAZ+default}  ---", must(zos.ResolveEnv([]byte("${BAZ+default}"))))
	fmt.Println("${ABC:3} ----------", must(zos.ResolveEnv([]byte("${ABC:3}"))))
	fmt.Println("${ABC:3:3} --------", must(zos.ResolveEnv([]byte("${ABC:3:3}"))))
	fmt.Println("${!ARR*} ----------", must(zos.ResolveEnv([]byte("${!ARR*}"))))
	fmt.Println("${!ARR@} ----------", must(zos.ResolveEnv([]byte("${!ARR@}"))))
	fmt.Println("${#FOO} ----------", must(zos.ResolveEnv([]byte("${#FOO}"))))
	fmt.Println("${FOO#[a-z]} -----", must(zos.ResolveEnv([]byte("${FOO#[a-z]}"))))
	fmt.Println("${FOO##[a-z]} ----", must(zos.ResolveEnv([]byte("${FOO##[a-z]}"))))
	fmt.Println("${FOO%[a-z]} -----", must(zos.ResolveEnv([]byte("${FOO%[a-z]}"))))
	fmt.Println("${FOO%%[a-z]} ----", must(zos.ResolveEnv([]byte("${FOO%%[a-z]}"))))
	fmt.Println("${FOO/[a-z]/x} ---", must(zos.ResolveEnv([]byte("${FOO/[a-z]/x}"))))
	fmt.Println("${FOO//[a-z]/x} --", must(zos.ResolveEnv([]byte("${FOO//[a-z]/x}"))))
	fmt.Println("${FOO/#[a-z]/x} --", must(zos.ResolveEnv([]byte("${FOO/#[a-z]/x}"))))
	fmt.Println("${FOO/%[a-z]/x} --", must(zos.ResolveEnv([]byte("${FOO/%[a-z]/x}"))))
	fmt.Println("${FOO^[f]} -------", must(zos.ResolveEnv([]byte("${FOO^[f]}"))))
	fmt.Println("${FOO^^[o]} ------", must(zos.ResolveEnv([]byte("${FOO^^[o]}"))))
	fmt.Println("${BAR,[B]} -------", must(zos.ResolveEnv([]byte("${BAR,[B]}"))))
	fmt.Println("${BAR,,[A]} ------", must(zos.ResolveEnv([]byte("${BAR,,[A]}"))))
	fmt.Println("${FOO@U} ---------", must(zos.ResolveEnv([]byte("${FOO@U}"))))
}
