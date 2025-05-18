package zhttp_test

import (
	"fmt"

	"github.com/aileron-projects/go/znet/zhttp"
)

func ExampleMatchMediaType() {
	fmt.Println(zhttp.MatchMediaType("text/plain", []string{"application/json", "application/*"}))
	fmt.Println(zhttp.MatchMediaType("text/plain", []string{"application/json", "text/plain"}))
	fmt.Println(zhttp.MatchMediaType("text/plain", []string{"application/json", "text/*"}))
	fmt.Println(zhttp.MatchMediaType("text/plain", []string{"application/json", "*/*"}))
	fmt.Println(zhttp.MatchMediaType("text-plain", []string{"application/json", "*/*"}))
	// Output:
	// -1
	// 1
	// 1
	// 1
	// -1
}

func ExampleParseQualifiedHeader() {
	fmt.Println(zhttp.ParseQualifiedHeader("foo"))
	fmt.Println(zhttp.ParseQualifiedHeader("foo; q=0.5"))
	fmt.Println(zhttp.ParseQualifiedHeader("foo; q=0.123"))
	fmt.Println(zhttp.ParseQualifiedHeader("foo; q=0.1234")) // Parsed as 0.123
	fmt.Println(zhttp.ParseQualifiedHeader("foo, bar"))
	fmt.Println(zhttp.ParseQualifiedHeader("foo; q=0.3, bar; q=0.5"))
	fmt.Println(zhttp.ParseQualifiedHeader("foo; q=0.0, bar; q=0.5")) // foo is ignored
	fmt.Println(zhttp.ParseQualifiedHeader("foo; q=1.2, bar; q=0.5")) // foo is ignored
	// Output:
	// [foo] [map[]]
	// [foo] [map[q:0.5]]
	// [foo] [map[q:0.123]]
	// [foo] [map[q:0.1234]]
	// [foo bar] [map[] map[]]
	// [bar foo] [map[q:0.5] map[q:0.3]]
	// [bar] [map[q:0.5]]
	// [bar] [map[q:0.5]]
}

func ExampleParseHeader() {
	fmt.Println(zhttp.ParseHeader("foo"))
	fmt.Println(zhttp.ParseHeader("foo; p=x"))
	fmt.Println(zhttp.ParseHeader("foo; p=x; q=y"))
	fmt.Println(zhttp.ParseHeader("foo; p=x; q=y, bar"))
	// Output:
	// [foo] [map[]]
	// [foo] [map[p:x]]
	// [foo] [map[p:x q:y]]
	// [foo bar] [map[p:x q:y] map[]]
}

func ExampleScanElement() {
	println := func(s1, s2 string) {
		fmt.Printf("`%s` | `%s`\n", s1, s2)
	}
	println(zhttp.ScanElement("foo"))
	println(zhttp.ScanElement("foo; p=x"))
	println(zhttp.ScanElement("foo; p=x; q=y"))
	println(zhttp.ScanElement("foo; p=x; q=y, bar"))
	// Output:
	// `foo` | ``
	// `foo; p=x` | ``
	// `foo; p=x; q=y` | ``
	// `foo; p=x; q=y` | ` bar`
}
