package ztext_test

import (
	"testing"

	"github.com/aileron-projects/go/ztext"
)

func BenchmarkTemplate(b *testing.B) {
	template := `
nil = {{ nil }}
string = {{ string }}
int = {{ int }}
int8 = {{ int8 }}
int16 = {{ int16 }}
int32 = {{ int32 }}
int64 = {{ int64 }}
uint = {{ uint }}
uint8 = {{ uint8 }}
uint16 = {{ uint16 }}
uint32 = {{ uint32 }}
uint64 = {{ uint64 }}
float32 = {{ float32 }}
float64 = {{ float64 }}
complex64 = {{ complex64 }}
complex128 = {{ complex128 }}
struct = {{ struct }}
`
	val := map[string]any{
		"nil":        nil,
		"string":     "foo",
		"int":        123,
		"int8":       int8(123),
		"int16":      int16(123),
		"int32":      int32(123),
		"int64":      int64(123),
		"uint":       uint(123),
		"uint8":      uint8(123),
		"uint16":     uint16(123),
		"uint32":     uint32(123),
		"uint64":     uint64(123),
		"float32":    float32(1.141592653589),
		"float64":    float64(1.141592653589),
		"complex64":  complex64(123 + 456i),
		"complex128": complex128(123 + 456i),
		"struct":     struct{ a, b string }{"foo", "bar"},
	}

	tpl := ztext.NewTemplate(template, "{{", "}}")

	b.ResetTimer()
	for b.Loop() {
		tpl.Execute(val)
	}
}
