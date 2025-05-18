package zos_test

import (
	"fmt"
	"os"

	"github.com/aileron-projects/go/zos"
)

func ExampleGetenvSlice() {
	os.Setenv("FOO", "foo1,foo2,foo3")
	os.Setenv("BAR", "bar1|bar2|bar3")

	fmt.Println(zos.GetenvSlice("FOO", ","))
	fmt.Println(zos.GetenvSlice("BAR", "|"))
	// Output:
	// [foo1 foo2 foo3]
	// [bar1 bar2 bar3]
}

func ExampleGetenvMap() {
	os.Setenv("FOO", "key1:val1,key2:val2")
	os.Setenv("BAR", "key1-alice|key1-bob|key2")

	fmt.Println(zos.GetenvMap("FOO", ",", ":"))
	fmt.Println(zos.GetenvMap("BAR", "|", "-"))
	// Output:
	// map[key1:[val1] key2:[val2]]
	// map[key1:[alice bob] key2:[]]
}

func ExampleEnvSubst() {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "FOO")

	b1, _ := zos.EnvSubst([]byte(`{FOO}=${FOO}`))
	fmt.Println(string(b1))

	b2, _ := zos.EnvSubst([]byte(`{{BAR}}=${${BAR}}`))
	fmt.Println(string(b2)) // Nested env is not supported.

	// Output:
	// {FOO}=foo
	// {{BAR}}=${FOO}
}

func ExampleEnvSubst2() {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "FOO")

	b1, _ := zos.EnvSubst2([]byte(`{FOO}=${FOO}`))
	fmt.Println(string(b1))

	b2, _ := zos.EnvSubst2([]byte(`{{BAR}}={FOO}=${${BAR}}`))
	fmt.Println(string(b2)) // Nested env is not supported.

	// Output:
	// {FOO}=foo
	// {{BAR}}={FOO}=foo
}

func ExampleEnvSubst_all() {
	os.Setenv("ABC", "abcdefg")
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "BAR")
	os.Setenv("ARR_X", "xxx")
	os.Setenv("ARR_Y", "yyy")

	txt := []byte(`
01: {parameter}                 => ${FOO}
02: {parameter:-word}           => ${BAZ:-default}
03: {parameter-word}            => ${BAZ-default}
04: {parameter:=word}           => ${BAZ:=default}
05: {parameter=word}            => ${BAZ=default}
06: {parameter:?word}           => ${BAZ:?default}
07: {parameter?word}            => ${BAZ?default}
08: {parameter:+word}           => ${BAZ:+default}
09: {parameter+word}            => ${BAZ+default}
10: {parameter:offset}          => ${ABC:3}
11: {parameter:offset:length}   => ${ABC:3:3}
12: {!prefix*}                  => ${!ARR*}
13: {!prefix@}                  => ${!ARR@}
14: {#parameter}                => ${#FOO}
15: {parameter#word}            => ${FOO#[a-z]}
16: {parameter##word}           => ${FOO##[a-z]}
17: {parameter%word}            => ${FOO%[a-z]}
18: {parameter%%word}           => ${FOO%%[a-z]}
19: {parameter/pattern/string}  => ${FOO/[a-z]/x}
20: {parameter//pattern/string} => ${FOO//[a-z]/x}
21: {parameter/#pattern/string} => ${FOO/#[a-z]/x}
22: {parameter/%pattern/string} => ${FOO/%[a-z]/x}
23: {parameter^pattern}         => ${FOO^[f]}
24: {parameter^^pattern}        => ${FOO^^[o]}
25: {parameter,pattern}         => ${BAR,[B]}
26: {parameter,,pattern}        => ${BAR,,[A]}
27: {parameter@U}               => ${FOO@U}
27: {parameter@u}               => ${FOO@u}
27: {parameter@L}               => ${BAR@L}
27: {parameter@l}               => ${BAR@l}
`)

	b, _ := zos.EnvSubst(txt)
	fmt.Println(string(b))
	// Output:
	// 01: {parameter}                 => foo
	// 02: {parameter:-word}           => default
	// 03: {parameter-word}            => default
	// 04: {parameter:=word}           => default
	// 05: {parameter=word}            => default
	// 06: {parameter:?word}           => default
	// 07: {parameter?word}            => default
	// 08: {parameter:+word}           => default
	// 09: {parameter+word}            => default
	// 10: {parameter:offset}          => defg
	// 11: {parameter:offset:length}   => def
	// 12: {!prefix*}                  => ARR_X ARR_Y
	// 13: {!prefix@}                  => ARR_X ARR_Y
	// 14: {#parameter}                => 3
	// 15: {parameter#word}            => oo
	// 16: {parameter##word}           => oo
	// 17: {parameter%word}            => fo
	// 18: {parameter%%word}           => fo
	// 19: {parameter/pattern/string}  => xoo
	// 20: {parameter//pattern/string} => xxx
	// 21: {parameter/#pattern/string} => xoo
	// 22: {parameter/%pattern/string} => fox
	// 23: {parameter^pattern}         => Foo
	// 24: {parameter^^pattern}        => fOO
	// 25: {parameter,pattern}         => bAR
	// 26: {parameter,,pattern}        => BaR
	// 27: {parameter@U}               => FOO
	// 27: {parameter@u}               => Foo
	// 27: {parameter@L}               => bar
	// 27: {parameter@l}               => bAR
}

func ExampleResolveEnv() {
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
	// Output:
	// ${FOO} ------------ foo
	// ${BAZ:-default} --- default
	// ${BAZ-default}  --- default
	// ${BAZ:=default} --- default
	// ${BAZ=default}  --- default
	// ${BAZ:?default} --- default
	// ${BAZ?default}  --- default
	// ${BAZ:+default} --- default
	// ${BAZ+default}  --- default
	// ${ABC:3} ---------- defg
	// ${ABC:3:3} -------- def
	// ${!ARR*} ---------- ARR_X ARR_Y
	// ${!ARR@} ---------- ARR_X ARR_Y
	// ${#FOO} ---------- 3
	// ${FOO#[a-z]} ----- oo
	// ${FOO##[a-z]} ---- oo
	// ${FOO%[a-z]} ----- fo
	// ${FOO%%[a-z]} ---- fo
	// ${FOO/[a-z]/x} --- xoo
	// ${FOO//[a-z]/x} -- xxx
	// ${FOO/#[a-z]/x} -- xoo
	// ${FOO/%[a-z]/x} -- fox
	// ${FOO^[f]} ------- Foo
	// ${FOO^^[o]} ------ fOO
	// ${BAR,[B]} ------- bAR
	// ${BAR,,[A]} ------ BaR
	// ${FOO@U} --------- FOO
}
