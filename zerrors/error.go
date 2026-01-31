package zerrors

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"strconv"
	"strings"
)

// ToMap returns the given error in map.
// If the given error implements the interface { Map() map[string]any },
// is call the method Map().
// ToMap repeatedly unwraps the given error by interfaxe{ Unwrap() error }
// or interface{ Unwrap []error }.
// ToMap returns nil map when the given error was nil.
func ToMap(err error) map[string]any {
	if err == nil {
		return nil
	}
	if a, ok := err.(interface{ Map() map[string]any }); ok {
		return a.Map()
	}
	m := map[string]any{"message": err.Error()}
	if e := UnwrapErr(err); e != nil {
		m["cause"] = ToMap(e)
		return m
	}
	if errs := UnwrapErrs(err); len(errs) > 0 {
		s := make([]map[string]any, len(errs))
		for i, e := range errs {
			s[i] = ToMap(e)
		}
		m["cause"] = s
	}
	return m
}

// ToSlogAttrs returns the error in [slog.Attr] format.
// If the given error implements the interface { SlogAttrs() []slog.Attr },
// is call the method SlogAttrs().
// ToSlogAttrs repeatedly unwraps the given error by interfaxe{ Unwrap() error }
// or interface{ Unwrap []error }.
// ToSlogAttrs returns nil slice when the given error was nil.
func ToSlogAttrs(err error) []slog.Attr {
	if err == nil {
		return nil
	}
	if a, ok := err.(interface{ SlogAttrs() []slog.Attr }); ok {
		return a.SlogAttrs()
	}
	s := []slog.Attr{slog.String("message", err.Error())}
	if e := UnwrapErr(err); e != nil {
		s = append(s, slog.GroupAttrs("cause", ToSlogAttrs(e)...))
		return s
	}
	if errs := UnwrapErrs(err); len(errs) > 0 {
		for i, e := range errs {
			s = append(s, slog.GroupAttrs("cause."+strconv.Itoa(i+1), ToSlogAttrs(e)...))
		}
	}
	return s
}

// Error is the general error type.
type Error struct {
	// Cause is the error cause.
	Cause error `json:"cause,omitempty" msgpack:"cause,omitempty" xml:"cause,omitempty" yaml:"cause,omitempty"`
	// Code is the error code, name or alias for the error.
	// Code is compared in the [Errors.Is] method.
	Code string `json:"code" msgpack:"code" xml:"code" yaml:"code"`
	// Kind is the error kind.
	Kind string `json:"kind" msgpack:"kind" xml:"kind" yaml:"kind"`
	// Message is the error message.
	Message string `json:"message" msgpack:"message" xml:"message" yaml:"message"`
	// Attrs are the attribution, or extra information, to this error.
	Attrs map[string]string `json:"attrs" msgpack:"attrs" xml:"attrs" yaml:"attrs"`
	// Frames is the list of stack trace frames.
	// Use [Error.WithStack] to fill this field.
	Frames []Frame `json:"frames,omitempty" msgpack:"frames,omitempty" xml:"frames,omitempty" yaml:"frames,omitempty"`
}

// Error implements [error] interface.
func (e *Error) Error() string {
	var builder strings.Builder
	builder.Grow(len(e.Code) + len(e.Kind) + len(e.Message) + 3)
	_, _ = builder.WriteString(e.Code + " ")
	_, _ = builder.WriteString(e.Kind + " :")
	_, _ = builder.WriteString(e.Message)
	if len(e.Attrs) > 0 {
		kvs := make([]string, 0, len(e.Attrs))
		for k, v := range e.Attrs {
			kvs = append(kvs, k+"="+v)
		}
		_, _ = builder.WriteString(" (" + strings.Join(kvs, ",") + ")")
	}
	if e.Cause != nil {
		_, _ = builder.WriteString(" [")
		_, _ = builder.WriteString(e.Cause.Error())
		_, _ = builder.WriteString("]")
	}
	return builder.String()
}

// Unwrap returns the inner error if any.
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is returns if this error is identical to the given error.
// This can be used with [errors.Is].
func (e *Error) Is(err error) bool {
	if err == nil || e == nil {
		return e == err
	}
	for err != nil {
		ee, ok := err.(*Error)
		if ok {
			return e.Code == ee.Code && e.Kind == ee.Kind
		}
		err = errors.Unwrap(err)
	}
	return false
}

// Map returns error information in map.
func (e *Error) Map() map[string]any {
	m := map[string]any{
		"code":    e.Code,
		"kind":    e.Kind,
		"message": e.Message,
	}
	if e.Attrs != nil {
		m["attrs"] = maps.Clone(e.Attrs)
	}
	if len(e.Frames) > 0 {
		fs := make([]string, 0, len(e.Frames))
		for _, f := range e.Frames {
			fs = append(fs, f.Pkg+":"+f.File+":L"+strconv.Itoa(f.Line)+"("+f.Func+")")
		}
		m["frames"] = fs
	}
	if cause := ToMap(e.Cause); cause != nil {
		m["cause"] = cause
	}
	return m
}

// Map returns error information in slice of [slog.Attr].
//
// Example:
//
//	lg := slog.New(slog.NewJSONHandler(os.Stdout, nil))
//	def := zerrors.NewDefinition("E123", "KindXXX", "example. foo=%s", map[string]string{"tag": "val"})
//	err := def.New(io.EOF, "bar")
//
//	lg.InfoContext(context.Background(), "message.", "error", err.SlogAttrs())
//	// JSON logger output >>
//	// {"level":"INFO","msg":"message.","error":{"code":"E123","kind":"KindXXX","message":"example. foo=bar","attrs":{"tag":"val"},"cause":{"message":"EOF"}}}
//	// Text logger output >>
//	// level=INFO msg="message." error.code=E123 error.kind=KindXXX error.message="example. foo=bar" error.attrs.tag=val error.cause.message=EOF
func (e *Error) SlogAttrs() []slog.Attr {
	a := []slog.Attr{
		slog.String("code", e.Code),
		slog.String("kind", e.Kind),
		slog.String("message", e.Message),
	}
	if e.Attrs != nil {
		aa := []slog.Attr{}
		for k, v := range e.Attrs {
			aa = append(aa, slog.String(k, v))
		}
		a = append(a, slog.GroupAttrs("attrs", aa...))
	}
	if len(e.Frames) > 0 {
		fs := make([]string, 0, len(e.Frames))
		for _, f := range e.Frames {
			fs = append(fs, f.Pkg+":"+f.File+":L"+strconv.Itoa(f.Line)+"("+f.Func+")")
		}
		a = append(a, slog.Any("frames", fs))
	}
	if causes := ToSlogAttrs(e.Cause); len(causes) > 0 {
		a = append(a, slog.GroupAttrs("cause", causes...))
	}
	return a
}

// NewDefinition returns a new error definition.
// See [Definition].
func NewDefinition(code, kind, message string, attrs map[string]string) *Definition {
	return &Definition{
		Code:    code,
		Kind:    kind,
		Message: message,
		Attrs:   maps.Clone(attrs),
	}
}

// Definition is the error definition.
type Definition struct {
	// Code is the error code.
	// Code is compared in [Definition.Is].
	Code string
	// Kind is the error kind that this error belongs to.
	// Kind is compared in [Definition.Is].
	Kind string
	// Message is the error message.
	// Format string for fmt.Sprintf can be used.
	// Message is NOT compared in [Definition.Is].
	Message string
	// Attrs are the attribution, or extra information, to this kind.
	// Attrs are NOT compared in [Definition.Is].
	Attrs map[string]string
}

// Is returns if the target err is the same as this definition.
func (d *Definition) Is(err error) bool {
	for err != nil {
		e, ok := err.(*Error)
		if ok {
			return d.Code == e.Code && d.Kind == e.Kind
		}
		err = errors.Unwrap(err)
	}
	return false
}

// New returns a new [Error] instance from the definition.
// New does not fill the [Error.Frames].
// Use [Definition.NewStack] when stack frames are necessary.
func (d *Definition) New(cause error, values ...any) *Error {
	err := &Error{
		Cause:   cause,
		Code:    d.Code,
		Kind:    d.Kind,
		Message: fmt.Sprintf(d.Message, values...),
		Attrs:   maps.Clone(d.Attrs),
	}
	traceTo(nil, err)
	return err
}

// NewStack returns a new [Error] instance from the definition.
// NewStack fills [Error.Frames] field.
// Use [Definition.New] when stack frames are not necessary.
func (d *Definition) NewStack(cause error, values ...any) *Error {
	err := &Error{
		Cause:   cause,
		Code:    d.Code,
		Kind:    d.Kind,
		Message: fmt.Sprintf(d.Message, values...),
		Attrs:   maps.Clone(d.Attrs),
	}
	e := cause
	for e != nil {
		ee, ok := e.(*Error)
		if ok && len(ee.Frames) > 0 {
			return err // The inner error already has frames.
		}
		e = errors.Unwrap(e)
	}
	err.Frames = callerFrames(1, 64) // Max 64 frames.
	traceTo(nil, err)
	return err
}
