package zerrors

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Attributer provides error attributes as a map.
type Attributer interface {
	// Attrs returns error attributes as map.
	Attrs() map[string]any
}

// Attrs returns error attributes as a map.
// If the given error implements the [Attributer] interface,
// is call the [Attributer.Attrs] internally.
// Attrs repeatedly unwraps the given error using [errors.Unwrap].
// Attrs returns nil when nil error was nil.
func Attrs(err error) map[string]any {
	if err == nil {
		return nil
	}
	if a, ok := err.(Attributer); ok {
		return a.Attrs()
	}
	if errs := UnwrapErrs(err); len(errs) > 0 {
		m := make(map[string]any, len(errs))
		for i, e := range errs {
			m["err"+strconv.Itoa(i+1)] = Attrs(e)
		}
		return m
	}
	m := map[string]any{
		"msg": err.Error(),
	}
	if err = UnwrapErr(err); err != nil {
		m["wraps"] = Attrs(err)
	}
	return m
}

// Error is the error type.
// Error implements [error] and [Attributer] interface.
type Error struct {
	// Inner is the inner error.
	Inner error `json:"wraps,omitempty" msgpack:"wraps,omitempty" xml:"wraps,omitempty" yaml:"wraps,omitempty"`
	// Code is the error code, name or alias for the error.
	// Code is compared in the [Errors.Is] method.
	Code string `json:"code" msgpack:"code" xml:"code" yaml:"code"`
	// Pkg is the package name that this error belongs to.
	// Pkg should not be empty.
	// Pkg is compared in the [Errors.Is] method.
	Pkg string `json:"pkg" msgpack:"pkg" xml:"pkg" yaml:"pkg"`
	// Msg is the error message.
	// Msg is compared in the [Errors.Is] method.
	Msg string `json:"msg" msgpack:"msg" xml:"msg" yaml:"msg"`
	// Detail is the detail of this error.
	// It provides additional information for the error.
	// Detail is not compared in the [Errors.Is] method.
	Detail string `json:"detail,omitempty" msgpack:"detail,omitempty" xml:"detail,omitempty" yaml:"detail,omitempty"`
	// Ext is the user extensible field.
	// Detail is not compared in the [Errors.Is] method.
	Ext string `json:"ext,omitempty" msgpack:"ext,omitempty" xml:"ext,omitempty" yaml:"ext,omitempty"`
	// Frames is the list of stack trace frames.
	// Use [Error.WithStack] to fill this field.
	Frames []Frame `json:"frames,omitempty" msgpack:"frames,omitempty" xml:"frames,omitempty" yaml:"frames,omitempty"`
}

func (e *Error) Error() string {
	var builder strings.Builder
	builder.Grow(len(e.Code) + len(e.Pkg) + len(e.Msg) + len(e.Detail) + 6)
	_, _ = builder.WriteString(e.Code + ": ") // len(e.Code) + 2
	_, _ = builder.WriteString(e.Pkg + ": ")  // len(e.Pkg) + 2
	_, _ = builder.WriteString(e.Msg)         // len(e.Msg)
	if e.Detail != "" {
		_, _ = builder.WriteString(": " + e.Detail) // len(e.Detail)+2
	}
	if e.Inner != nil {
		_, _ = builder.WriteString(" [")
		_, _ = builder.WriteString(e.Inner.Error())
		_, _ = builder.WriteString("]")
	}
	return builder.String()
}

// Unwrap returns the inner error if any.
func (e *Error) Unwrap() error {
	return e.Inner
}

// Is returns if this error is identical to the given error.
// The err is identical to the error when it has the type [Error]
// and both [Error.Code], [Error.Pkg] and [Error.Msg] fields are the same.
func (e *Error) Is(err error) bool {
	if err == nil || e == nil {
		return e == err
	}
	for err != nil {
		ee, ok := err.(*Error)
		if ok {
			return e.Code == ee.Code && e.Pkg == ee.Pkg && e.Msg == ee.Msg
		}
		err = errors.Unwrap(err)
	}
	return false
}

// Attrs returns error attributes in map.
// Extra attributes in e.Extra is copied to the returned map.
func (e *Error) Attrs() map[string]any {
	attrs := map[string]any{
		"code": e.Code,
		"pkg":  e.Pkg,
		"msg":  e.Msg,
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
	if wrap := Attrs(e.Inner); wrap != nil {
		attrs["wraps"] = wrap
	}
	return attrs
}

// NewDefinition returns a new error definition.
// Only detail can be template format with [fmt.Sprintf] syntax.
func NewDefinition(code, pkg, msg, detail string) Definition {
	return Definition{code, pkg, msg, detail}
}

// Definition is the error definition type.
// Each index has the following meanings.
//
//	0: code, name or alias for the error. (e.g. "E1234")
//	1: package name that the error belongs to. (e.g. "zerrors")
//	2: error message. (e.g. "authentication failed")
//	3: template of detail in [fmt.Sprintf] syntax. (e.g. "username=%s")
//	4: custom field . Extensible by users.
type Definition [5]string

// Is returns if the target err is generated from this error definition.
// The err is identical to the definition when it has the type [Error]
// and both [Error.Code], [Error.Pkg] and [Error.Msg] fields are the same.
func (d Definition) Is(err error) bool {
	if err == nil {
		return false
	}
	for err != nil {
		ee, ok := err.(*Error)
		if ok {
			return d[0] == ee.Code && d[1] == ee.Pkg && d[2] == ee.Msg
		}
		err = errors.Unwrap(err)
	}
	return false
}

// New returns a new [Error] instance from the definition
// with the given inner error and the arguments for detail template.
// New does not fill the [Error.Frames] field.
// Use [Definition.NewStack] if stack frames are necessary.
func (d Definition) New(inner error, vs ...any) *Error {
	err := &Error{
		Inner:  inner,
		Code:   d[0],
		Pkg:    d[1],
		Msg:    d[2],
		Ext:    d[4],
		Detail: fmt.Sprintf(d[3], vs...),
	}
	traceTo(nil, err)
	return err
}

// NewStack returns a new [Error] instance from the definition
// with the given inner error and the arguments for detail.
// NewTplStack fills the [Error.Frames] for stack traces.
// Use [Definition.New] if stack frames are not necessary.
// Note that NewStack does not fill the [Error.Frames]
// if the given inner error already has stack frames.
func (d Definition) NewStack(inner error, vs ...any) *Error {
	err := &Error{
		Inner:  inner,
		Code:   d[0],
		Pkg:    d[1],
		Msg:    d[2],
		Ext:    d[4],
		Detail: fmt.Sprintf(d[3], vs...),
	}
	e := inner
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
