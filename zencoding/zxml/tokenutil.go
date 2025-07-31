package zxml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	CauseDataType    = "xmlconv: invalid data type."
	CauseJSONStruct  = "xmlconv: invalid json structure."
	CauseJSONKey     = "xmlconv: invalid json key."
	CauseEmptyVal    = "xmlconv: invalid empty value."
	CauseJSONEncoder = "xmlconv: json encoder return with an error."
	CauseJSONDecoder = "xmlconv: json decoder return with an error."
	CauseXMLEncoder  = "xmlconv: xml encoder return with an error."
	CauseXMLDecoder  = "xmlconv: xml decoder return with an error."
)

// XMLError reports cause and detail of an error
// with an inner error if any.
type XMLError struct {
	Err    error
	Cause  string
	Detail string
}

func (e *XMLError) Error() string {
	msg := "zxml: " + e.Cause
	if e.Detail != "" {
		msg += " " + e.Detail
	}
	if e.Err != nil {
		msg += " [" + e.Err.Error() + "]"
	}
	return msg
}

func (e *XMLError) Unwrap() error {
	return e.Err
}

// Is returns if this error is identical to the given error.
// The err is identical to the error when it has the type [XMLError]
// and [XMLError.Cause] field is the same.
func (e *XMLError) Is(err error) bool {
	if err == nil || e == nil {
		return e == err
	}
	for err != nil {
		ee, ok := err.(*XMLError)
		if ok {
			return e.Cause == ee.Cause
		}
		err = errors.Unwrap(err)
	}
	return false
}

// attrName converts XML attribute token name to string.
//
// Examples (prefix="@", sep=":"):
//
//	foo="bar" --> @foo
//	ns:foo="bar" --> @ns:foo
//	xmlns:foo="bar" --> @xmlns:foo
//	xmlns="bar" --> @xmlns
func attrName(name xml.Name, prefix, sep string) string {
	if name.Space == "" {
		return prefix + name.Local
	}
	return prefix + name.Space + sep + name.Local
}

// tokenName returns name with namespace.
func tokenName(name xml.Name, sep string) string {
	if name.Space == "" {
		return name.Local
	}
	return name.Space + sep + name.Local
}

// restoreNamespace returns key string with namespace when possible.
// It split the key string with nsSep int 2 parts and returns a string
// joining the 2 parts with a ":".
// If nsSep is not contained or retrieved namespace is empty,
// it returns the key string as-is.
//
// Example:
//
//	":", "ns:foo" --> "ns:foo"
//	":", "foo" --> "foo"
//	"_", "_foo" --> "_foo"
//	"_", "ns_foo" --> "ns:foo"
//	"", "ns_foo" --> "ns_foo"
func restoreNamespace(nsSep string, key string) string {
	if nsSep == ":" || nsSep == "" {
		return key
	}
	before, after, found := strings.Cut(key, nsSep)
	if found && before != "" {
		return before + ":" + after
	}
	return key
}

// jsonValueToToken returns a [encoding/xml.Token] from the value
// parsed from JSON document.
// Supported types follow the [encoding/json.Token] which listed below.
// Unsupported types results in an error.
//
// Supported Types:
//   - nil
//   - string
//   - bool
//   - float64
//   - json.Number
func jsonValueToToken(trimSpace bool, value any) (xml.Token, error) {
	switch v := value.(type) {
	case string:
		if trimSpace {
			v = strings.TrimSpace(v)
		}
		return xml.CharData(v), nil
	case bool:
		return xml.CharData(strconv.FormatBool(v)), nil
	case float64:
		return xml.CharData(strconv.FormatFloat(v, 'g', -1, 64)), nil
	case json.Number:
		return xml.CharData(v), nil
	case nil: // Keep empty
		return xml.CharData(nil), nil
	default:
		return nil, &XMLError{
			Cause:  CauseDataType,
			Detail: fmt.Sprintf("json value must be null,string,bool or number. got %T:%+v", value, value),
		}
	}
}
