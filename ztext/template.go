package ztext

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"sync/atomic"
)

const (
	tplValue = iota
	tplTag
)

// Template returns a new instance of the Template.
// Allowed tag name pattern is `[0-9a-zA-Z_\-\.]+`.
func NewTemplate(tpl string, start, end string) *Template {
	exp := []byte(tpl)
	escapeStart := regexp.QuoteMeta(start)
	escapeEnd := regexp.QuoteMeta(end)
	reg := regexp.MustCompile(escapeStart + ` *[0-9a-zA-Z_\-\.]+ *` + escapeEnd)
	indexes := reg.FindAllIndex(exp, -1)

	types := []int{}
	values := [][]byte{}

	pos := 0
	for _, ids := range indexes {
		if ids[0] > pos {
			val := exp[pos:ids[0]]
			val = bytes.ReplaceAll(val, []byte(escapeStart), []byte(start))
			val = bytes.ReplaceAll(val, []byte(escapeEnd), []byte(end))
			types = append(types, tplValue)
			values = append(values, val)
		}
		types = append(types, tplTag)
		values = append(values, bytes.Trim(exp[ids[0]+len(start):ids[1]-len(end)], " "))
		pos = ids[1]
	}
	if pos < len(exp) {
		types = append(types, 0)
		values = append(values, exp[pos:])
	}

	fast := &Template{
		valTypes: slices.Clip(types),
		values:   slices.Clip(values),
	}
	fast.bufSize.Store(int64(len(tpl)))
	return fast
}

// Template is more simple but fast template engine than [text.Template].
// Use [NewTemplate] to instantiate Template.
// Template supports formatting primitive types with the listed way.
// [fmt.Sprint] is used to fallback.
//
// - nil : "<nil>"
// - string : v
// - []byte : string(v)
// - bool : strconv.FormatBool(v)
// - int,int8,int16,int32,int64 : strconv.FormatInt(int64(v), 10)
// - uint,uint8,uint16,uint32,uint64 : strconv.FormatUint(uint64(v), 10)
// - float32 : strconv.FormatFloat(float64(v), 'g', -1, 32)
// - float64 : strconv.FormatFloat(float64(v), 'g', -1, 64)
// - complex64 : strconv.FormatComplex(complex128(v), 'g', -1, 64)
// - complex128 : strconv.FormatComplex(complex128(v), 'g', -1, 128)
// - fmt.Stringer : v.String()
// - others : fmt.Sprint(v)
type Template struct {
	// value types represents if template sections are value or tag.
	valTypes []int
	// value holds value or tag value.
	// Values are written to writers as it is and tag values are written
	// to writers after evaluation.
	values [][]byte
	// bufSize is the initial buffer size.
	bufSize atomic.Int64
	// tagFuncs are tag evaluation functions.
	tagFuncs map[string]func(string) []byte
}

// WithTagFunc registers a tag function that returns a value
// corresponding to tags.
// Registered tag function is used for evaluating tag values.
// The given tagFunc replaces existing tagFunc if there already
// were tagFunc which has the same tag value.
// WithTagFunc ignores the given tagFunc if the given argument tag is empty
// or the tagFunc is nil itself.
// WithTagFunc is not safe for concurrent call.
func (t *Template) WithTagFunc(tag string, tagFunc func(string) []byte) {
	if tag == "" || tagFunc == nil {
		return
	}
	if t.tagFuncs == nil {
		t.tagFuncs = map[string]func(string) []byte{}
	}
	t.tagFuncs[tag] = tagFunc
}

// Execute executes template and return resulting []byte.
// See the comment on [Template] for supported data types.
func (t *Template) Execute(m map[string]any) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize.Load()))
	_ = t.ExecuteWriter(buf, m) // Write error may not be occurred.
	return buf.Bytes()
}

// ExecuteString executes template and returns resulting string.
// See the comment on [Template] for supported data types.
func (t *Template) ExecuteString(m map[string]any) string {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize.Load()))
	_ = t.ExecuteWriter(buf, m) // Write error may not be occurred.
	return buf.String()
}

// ExecuteFunc executes the template with given tag function.
func (t *Template) ExecuteFunc(f func(string) []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, t.bufSize.Load()))
	_ = t.ExecuteWriterFunc(buf, f) // Write error may not be occurred.
	return buf.Bytes()
}

// ExecuteWriter executes template and writes results into the
// given writer. [w.Write] will be called multiple times.
// See the comment on [Template] for supported data types.
func (t *Template) ExecuteWriter(w io.Writer, m map[string]any) error {
	var err error
	mv := mapVal(m)
	for i := range t.valTypes {
		switch t.valTypes[i] {
		case tplValue:
			_, err = w.Write(t.values[i])
		case tplTag:
			if t.tagFuncs != nil {
				tf, ok := t.tagFuncs[string(t.values[i])]
				if ok {
					_, err = w.Write(tf(string(t.values[i])))
					break
				}
			}
			_, err = w.Write(mv.value(string(t.values[i])))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// ExecuteWriterFunc executes the template by evaluating tags with given
// tag function and write results into the given writer.
// Note that the registered tag func is prior to the given function.
// [io.Write] of the given w will be called multiple times.
func (t *Template) ExecuteWriterFunc(w io.Writer, f func(string) []byte) error {
	var err error
	for i := range t.valTypes {
		switch t.valTypes[i] {
		case tplValue:
			_, err = w.Write(t.values[i])
		case tplTag:
			if t.tagFuncs != nil {
				tf, ok := t.tagFuncs[string(t.values[i])]
				if ok {
					_, err = w.Write(tf(string(t.values[i])))
					break
				}
			}
			fallthrough // To default.
		default:
			if f != nil {
				_, err = w.Write(f(string(t.values[i])))
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// mapVal is the map type that is used in the [Template].
// Following
type mapVal map[string]any

func (m mapVal) value(tag string) []byte {
	val, ok := m[tag]
	if !ok {
		return nil
	}
	switch v := val.(type) {
	case nil:
		return []byte("<nil>")
	case string:
		return []byte(v)
	case []byte:
		return v
	case bool:
		return []byte(strconv.FormatBool(v))
	case int:
		return []byte(strconv.FormatInt(int64(v), 10))
	case int8:
		return []byte(strconv.FormatInt(int64(v), 10))
	case int16:
		return []byte(strconv.FormatInt(int64(v), 10))
	case int32:
		return []byte(strconv.FormatInt(int64(v), 10))
	case int64:
		return []byte(strconv.FormatInt(v, 10))
	case float32:
		return []byte(strconv.FormatFloat(float64(v), 'g', -1, 32))
	case float64:
		return []byte(strconv.FormatFloat(float64(v), 'g', -1, 64))
	case uint:
		return []byte(strconv.FormatUint(uint64(v), 10))
	case uint8:
		return []byte(strconv.FormatUint(uint64(v), 10))
	case uint16:
		return []byte(strconv.FormatUint(uint64(v), 10))
	case uint32:
		return []byte(strconv.FormatUint(uint64(v), 10))
	case uint64:
		return []byte(strconv.FormatUint(v, 10))
	case complex64:
		return []byte(strconv.FormatComplex(complex128(v), 'g', -1, 64))
	case complex128:
		return []byte(strconv.FormatComplex(complex128(v), 'g', -1, 128))
	case fmt.Stringer:
		return []byte(v.String())
	default:
		return fmt.Append(nil, v) // Fallback to "%+v"
	}
}
