package zerrors_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/aileron-projects/go/zerrors"
)

func ExampleToMap() {
	e1 := errors.New("example1")
	e2 := errors.New("example2")

	fmt.Println(zerrors.ToMap(e1))
	fmt.Println(zerrors.ToMap(fmt.Errorf("example3 [%w]", e1)))
	fmt.Println(zerrors.ToMap(errors.Join(e1, e2)))
	// Output:
	// map[message:example1]
	// map[cause:map[message:example1] message:example3 [example1]]
	// map[cause:[map[message:example1] map[message:example2]] message:example1
	// example2]
}

func ExampleToSlogAttrs() {
	e1 := errors.New("example1")
	e2 := errors.New("example2")

	fmt.Println(zerrors.ToSlogAttrs(e1))
	fmt.Println(zerrors.ToSlogAttrs(fmt.Errorf("example3 [%w]", e1)))
	fmt.Println(zerrors.ToSlogAttrs(errors.Join(e1, e2)))
	// Output:
	// [message=example1]
	// [message=example3 [example1] cause=[message=example1]]
	// [message=example1
	// example2 cause.1=[message=example1] cause.2=[message=example2]]
}

func ExampleDefinition_New() {
	// Define an error with detail template.
	def := zerrors.NewDefinition("E123", "KindXXX", "example error. foo=%s bar=%s.", map[string]string{"foo": "bar"})

	fmt.Println(def.New(nil, "FOO", "BAR").Error())        // With arguments.
	fmt.Println(def.New(nil).Error())                      // No arguments.
	fmt.Println(def.New(nil, "FOO").Error())               // Insufficient arguments.
	fmt.Println(def.New(nil, "FOO", "BAR", "BAZ").Error()) // Too many arguments.
	fmt.Println(def.New(io.EOF, "FOO", "BAR").Error())     // With inner error.
	// Output:
	// E123 KindXXX :example error. foo=FOO bar=BAR. (foo=bar)
	// E123 KindXXX :example error. foo=%!s(MISSING) bar=%!s(MISSING). (foo=bar)
	// E123 KindXXX :example error. foo=FOO bar=%!s(MISSING). (foo=bar)
	// E123 KindXXX :example error. foo=FOO bar=BAR.%!(EXTRA string=BAZ) (foo=bar)
	// E123 KindXXX :example error. foo=FOO bar=BAR. (foo=bar) [EOF]
}

func ExampleError_slog() {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 { // Remove time to avoid test failure.
				return slog.Attr{}
			}
			return a
		},
	}
	lgJSON := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	lgText := slog.New(slog.NewTextHandler(os.Stdout, opts))

	def := zerrors.NewDefinition("E123", "KindXXX", "example. foo=%s", map[string]string{"tag": "val"})
	err := def.New(io.EOF, "bar")

	lgJSON.InfoContext(context.Background(), "message.", "error", err.SlogAttrs())
	lgText.InfoContext(context.Background(), "message.", "error", err.SlogAttrs())
	// Output:
	// {"level":"INFO","msg":"message.","error":{"code":"E123","kind":"KindXXX","message":"example. foo=bar","attrs":{"tag":"val"},"cause":{"message":"EOF"}}}
	// level=INFO msg=message. error.code=E123 error.kind=KindXXX error.message="example. foo=bar" error.attrs.tag=val error.cause.message=EOF
}

func ExampleErr_slog() {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 { // Remove time to avoid test failure.
				return slog.Attr{}
			}
			return a
		},
	}
	lgJSON := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	lgText := slog.New(slog.NewTextHandler(os.Stdout, opts))

	err := zerrors.NewErr(io.EOF, "example", "foo=%s", "bar")

	lgJSON.InfoContext(context.Background(), "message.", "error", err.SlogAttrs())
	lgText.InfoContext(context.Background(), "message.", "error", err.SlogAttrs())
	// Output:
	// {"level":"INFO","msg":"message.","error":{"message":"example","detail":"foo=bar","cause":{"message":"EOF"}}}
	// level=INFO msg=message. error.message=example error.detail="foo=bar" error.cause.message=EOF
}
