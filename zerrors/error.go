package zerrors

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Attributes provides error attributes as a map.
type Attributes interface {
	// Attrs returns error attributes as map.
	Attrs() map[string]any
}

// Attrs returns error attributes as a map.
// If the given error implements the [Attributes] interface,
// is call the [Attributes.Attrs] internally.
// Attrs repeatedly unwraps the given error using [errors.Unwrap].
// Attrs returns nil map when the given error was nil.
func Attrs(err error) map[string]any {
	if err == nil {
		return nil
	}
	if a, ok := err.(Attributes); ok {
		return a.Attrs()
	}
	if errs := UnwrapErrs(err); len(errs) > 0 {
		m := make(map[string]any, len(errs))
		for i, e := range errs {
			m["cause."+strconv.Itoa(i+1)] = Attrs(e)
		}
		return m
	}
	m := map[string]any{
		"message": err.Error(),
	}
	if err = UnwrapErr(err); err != nil {
		m["cause"] = Attrs(err)
	}
	return m
}

// Error is the error type.
// Error implements [error] and [Attributer] interface.
type Error struct {
	// Cause is the error cause.
	Cause error `json:"cause,omitempty" msgpack:"cause,omitempty" xml:"cause,omitempty" yaml:"cause,omitempty"`
	// Code is the error code, name or alias for the error.
	// Code is compared in the [Errors.Is] method.
	Code Code `json:"code" msgpack:"code" xml:"code" yaml:"code"`
	// Kind is the error kind.
	Kind Kind `json:"kind" msgpack:"kind" xml:"kind" yaml:"kind"`
	// Name is the error name.
	Name string `json:"name" msgpack:"name" xml:"name" yaml:"name"`
	// Message is the error message.
	Message string `json:"message" msgpack:"message" xml:"message" yaml:"message"`
	// Detail is the error detail.
	Detail string `json:"detail,omitempty" msgpack:"detail,omitempty" xml:"detail,omitempty" yaml:"detail,omitempty"`
	// Frames is the list of stack trace frames.
	// Use [Error.WithStack] to fill this field.
	Frames []Frame `json:"frames,omitempty" msgpack:"frames,omitempty" xml:"frames,omitempty" yaml:"frames,omitempty"`
}

func (e *Error) Error() string {
	var builder strings.Builder
	builder.Grow(len(e.Code) + len(e.Kind) + len(e.Name) + len(e.Message) + len(e.Detail) + 6)
	_, _ = builder.WriteString(string(e.Code) + " ")
	_, _ = builder.WriteString(e.Name + " ")
	_, _ = builder.WriteString(string(e.Kind) + " : ")
	_, _ = builder.WriteString(e.Message)
	if e.Detail != "" {
		_, _ = builder.WriteString(" " + e.Detail)
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

// Attrs returns error attributes in map.
// Extra attributes in e.Extra is copied to the returned map.
func (e *Error) Attrs() map[string]any {
	attrs := map[string]any{
		"code":    e.Code,
		"name":    e.Name,
		"kind":    e.Kind,
		"message": e.Message,
	}
	if e.Detail != "" {
		attrs["detail"] = e.Detail
	}
	if len(e.Frames) > 0 {
		fs := make([]string, 0, len(e.Frames))
		for _, f := range e.Frames {
			fs = append(fs, f.Pkg+":"+f.File+":L"+strconv.Itoa(f.Line)+"("+f.Func+")")
		}
		attrs["frames"] = fs
	}
	if wrap := Attrs(e.Cause); wrap != nil {
		attrs["cause"] = wrap
	}
	return attrs
}

// Code is the error code type.
// For example, "E123", "E456".
type Code string

// Kind is the error kind type.
// For example, "ClientError", "ServerError".
type Kind string

// NewDefinition returns a new error definition.
// See [Definition].
func NewDefinition(code Code, kind Kind, name, message, detail string) *Definition {
	return &Definition{
		Code:    code,
		Kind:    kind,
		Name:    name,
		Message: message,
		Detail:  detail,
	}
}

// Definition is the error definition.
type Definition struct {
	// Code is the error code.
	// Code is compared in [Definition.Is].
	Code Code
	// Kind is the error kind that this error belongs to.
	// Kind is compared in [Definition.Is].
	Kind Kind
	// Name is the human readable error name.
	// Name is NOT compared in [Definition.Is].
	Name string
	// Message is the fixed error message.
	// Message is NOT compared in [Definition.Is].
	Message string
	// Detail is the text template in the format of [fmt.Sprintf].
	// Detail is NOT compared in [Definition.Is].
	Detail string
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

// New returns a new [Error] instance from the definition
// with the given error cause and the detail values.
// New does not fill the [Error.Frames].
// Use [Definition.NewStack] if stack frames are necessary.
func (d *Definition) New(cause error, values ...any) *Error {
	err := &Error{
		Cause:   cause,
		Code:    d.Code,
		Kind:    d.Kind,
		Name:    d.Name,
		Message: d.Message,
		Detail:  fmt.Sprintf(d.Detail, values...),
	}
	traceTo(nil, err)
	return err
}

// NewStack returns a new [Error] instance from the definition
// with the given error cause and the detail values.
// NewStack fills [Error.Frames] field.
// Use [Definition.New] if stack frames are not necessary.
func (d *Definition) NewStack(cause error, values ...any) *Error {
	err := &Error{
		Cause:   cause,
		Code:    d.Code,
		Kind:    d.Kind,
		Name:    d.Name,
		Message: d.Message,
		Detail:  fmt.Sprintf(d.Detail, values...),
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
