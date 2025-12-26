package main

import (
	"fmt"

	"github.com/aileron-projects/go/zos"
)

var env = `
FOO=foo # This is a comment.
BAR=bar
export BAZ=baz
URL=http://example.com
USERNAME=foo
PASSWORD=bar
SECRET_URL=http://${USERNAME}:${PASSWORD}@example.com

MULTILINE_A="
one
two
"

MULTILINE_B="
one\n
two
"

QUOTE_SINGLE='single quoted. " can be used.'
QUOTE_DOUBLE="double quoted. ' can be used."
QUOTE_SINGLE_ESCAPE='single quotation \'escaped\'.'
QUOTE_DOUBLE_ESCAPE="double quotation \"escaped\"."
`

func main() {
	m, err := zos.ParseEnv([]byte(env))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", m)
}
