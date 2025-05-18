package zstrconv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrTypeSupported = errors.New("zstrconv: unsupported type")
)

// Number is a constraint that permits any numeric type.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// ParseNum parses primitive numbers and return the parsed value.
// If some error had occurred while parsing s, it is returned with
// the zero value of T.
// ParseNum should not be used when performance is important.
func ParseNum[T Number](s string) (T, error) {
	var t T
	v, err := paseNum(strings.TrimSpace(s), reflect.TypeOf(t).Kind())
	if err != nil {
		return t, fmt.Errorf("zstrconv: parse failed. value=`%s` [%w]", s, err)
	}
	return v.(T), nil
}

func paseNum(v string, kind reflect.Kind) (any, error) {
	switch kind {
	case reflect.Int:
		return strconv.Atoi(v)
	case reflect.Int8:
		vv, err := strconv.ParseInt(v, 10, 8)
		return int8(vv), err
	case reflect.Int16:
		vv, err := strconv.ParseInt(v, 10, 16)
		return int16(vv), err
	case reflect.Int32:
		vv, err := strconv.ParseInt(v, 10, 32)
		return int32(vv), err
	case reflect.Int64:
		vv, err := strconv.ParseInt(v, 10, 64)
		return vv, err
	case reflect.Uint:
		vv, err := strconv.ParseUint(v, 10, 32)
		return uint(vv), err
	case reflect.Uint8:
		vv, err := strconv.ParseUint(v, 10, 8)
		return uint8(vv), err
	case reflect.Uint16:
		vv, err := strconv.ParseUint(v, 10, 16)
		return uint16(vv), err
	case reflect.Uint32:
		vv, err := strconv.ParseUint(v, 10, 32)
		return uint32(vv), err
	case reflect.Uint64:
		vv, err := strconv.ParseUint(v, 10, 64)
		return vv, err
	case reflect.Uintptr:
		vv, err := strconv.ParseUint(v, 10, 64)
		return uintptr(vv), err
	case reflect.Float32:
		vv, err := strconv.ParseFloat(v, 32)
		return float32(vv), err
	case reflect.Float64:
		vv, err := strconv.ParseFloat(v, 64)
		return vv, err
	default:
		return nil, ErrTypeSupported
	}
}

// FormatAny returns formatted string of any primitive data type.
// For given value v=a.(type), following format is applied.
//   - nil : "<nil>"
//   - string : v
//   - []byte : string(v)
//   - bool : strconv.FormatBool(v)
//   - int,int8,int16,int32,int64 : strconv.FormatInt(int64(v), 10)
//   - uint,uint8,uint16,uint32,uint64 : strconv.FormatUint(uint64(v), 10)
//   - float32 : strconv.FormatFloat(float64(v), 'g', -1, 32)
//   - float64 : strconv.FormatFloat(float64(v), 'g', -1, 64)
//   - complex64 : strconv.FormatComplex(complex128(v), 'g', -1, 64)
//   - complex128 : strconv.FormatComplex(complex128(v), 'g', -1, 128)
//   - fmt.Stringer : v.String()
//   - others : fmt.Sprint(v)
func FormatAny(a any) string {
	switch v := a.(type) {
	case nil:
		return "<nil>"
	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(float64(v), 'g', -1, 64)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case complex64:
		return strconv.FormatComplex(complex128(v), 'g', -1, 64)
	case complex128:
		return strconv.FormatComplex(complex128(v), 'g', -1, 128)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprint(v) // Fallback to "%+v"
	}
}
