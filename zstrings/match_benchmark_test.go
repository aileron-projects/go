package zstrings

import (
	"path"
	"testing"
)

var matchStr = `
Go is
* An open-source programming language supported by Google.
* Easy to learn and great for teams.
* Built-in concurrency and a robust standard library.
* Large ecosystem of partners, communities, and tools.
`

var matchPattern = `
Go is
\* An op*-so* p?o?r?m?i?g l*e su*orted by Google.
\* Easy to le??? and great *** teams.
\* Built\-in c*o*n*c*u*r*r*e*n*c*y and a robust stand* lib*.
\* Large ecosystem of partners, communities, and tools.*
`

func BenchmarkMatch(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Match(matchPattern, matchStr)
	}
}

func BenchmarkPathMatch(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path.Match(matchPattern, matchStr)
	}
}
