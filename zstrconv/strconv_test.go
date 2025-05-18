package zstrconv_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/aileron-projects/go/zstrconv"
	"github.com/aileron-projects/go/ztesting"
)

func TestParseNum(t *testing.T) {
	t.Parallel()
	t.Run("int", func(t *testing.T) {
		v, err := zstrconv.ParseNum[int]("123")
		ztesting.AssertEqual(t, "incorrect parse result", int(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[int]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", int(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("int8", func(t *testing.T) {
		v, err := zstrconv.ParseNum[int8]("123")
		ztesting.AssertEqual(t, "incorrect parse result", int8(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[int8]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", int8(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("int16", func(t *testing.T) {
		v, err := zstrconv.ParseNum[int16]("123")
		ztesting.AssertEqual(t, "incorrect parse result", int16(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[int16]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", int16(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("int16", func(t *testing.T) {
		v, err := zstrconv.ParseNum[int16]("123")
		ztesting.AssertEqual(t, "incorrect parse result", int16(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[int16]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", int16(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("int32", func(t *testing.T) {
		v, err := zstrconv.ParseNum[int32]("123")
		ztesting.AssertEqual(t, "incorrect parse result", int32(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[int32]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", int32(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("int64", func(t *testing.T) {
		v, err := zstrconv.ParseNum[int64]("123")
		ztesting.AssertEqual(t, "incorrect parse result", int64(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[int64]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", int64(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("uint", func(t *testing.T) {
		v, err := zstrconv.ParseNum[uint]("123")
		ztesting.AssertEqual(t, "incorrect parse result", uint(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[uint]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", uint(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("uint8", func(t *testing.T) {
		v, err := zstrconv.ParseNum[uint8]("123")
		ztesting.AssertEqual(t, "incorrect parse result", uint8(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[uint8]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", uint8(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("uint16", func(t *testing.T) {
		v, err := zstrconv.ParseNum[uint16]("123")
		ztesting.AssertEqual(t, "incorrect parse result", uint16(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[uint16]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", uint16(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("uint32", func(t *testing.T) {
		v, err := zstrconv.ParseNum[uint32]("123")
		ztesting.AssertEqual(t, "incorrect parse result", uint32(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[uint32]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", uint32(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("uint64", func(t *testing.T) {
		v, err := zstrconv.ParseNum[uint64]("123")
		ztesting.AssertEqual(t, "incorrect parse result", uint64(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[uint64]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", uint64(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("uintptr", func(t *testing.T) {
		v, err := zstrconv.ParseNum[uintptr]("123")
		ztesting.AssertEqual(t, "incorrect parse result", uintptr(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[uintptr]("1234567890123456789012345678901234567890")
		ztesting.AssertEqual(t, "incorrect parse result", uintptr(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("float32", func(t *testing.T) {
		v, err := zstrconv.ParseNum[float32]("123")
		ztesting.AssertEqual(t, "incorrect parse result", float32(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[float32]("1.23e+10000")
		ztesting.AssertEqual(t, "incorrect parse result", float32(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
	t.Run("float64", func(t *testing.T) {
		v, err := zstrconv.ParseNum[float64]("123")
		ztesting.AssertEqual(t, "incorrect parse result", float64(123), v)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		vv, err := zstrconv.ParseNum[float64]("1.23e+10000")
		ztesting.AssertEqual(t, "incorrect parse result", float64(0), vv)
		ztesting.AssertEqualErr(t, "non nil error returned", strconv.ErrRange, err)
	})
}

type testStringer struct {
	s string
}

func (t *testStringer) String() string {
	return t.s
}

func TestFormatAny(t *testing.T) {
	t.Parallel()
	testCase := map[string]struct {
		a    any
		want string
	}{
		"nil":           {nil, "<nil>"},
		"string":        {"foo", "foo"},
		"[]byte":        {[]byte("foo"), "foo"},
		"bool":          {true, "true"},
		"int":           {int(math.MaxInt), strconv.FormatInt(math.MaxInt, 10)},
		"int8":          {int8(math.MaxInt8), "127"},
		"int16":         {int16(math.MaxInt16), "32767"},
		"int32":         {int32(math.MaxInt32), "2147483647"},
		"int64":         {int64(math.MaxInt64), "9223372036854775807"},
		"uint":          {uint(math.MaxUint), strconv.FormatUint(math.MaxUint, 10)},
		"uint8":         {uint8(math.MaxUint8), "255"},
		"uint16":        {uint16(math.MaxUint16), "65535"},
		"uint32":        {uint32(math.MaxUint32), "4294967295"},
		"uint64":        {uint64(math.MaxUint64), "18446744073709551615"},
		"float32":       {float32(16777215), "1.6777215e+07"},
		"float32_upper": {float32(16777216), "1.6777216e+07"},
		"float37_over":  {float32(16777217), "1.6777216e+07"},
		"float32_-max":  {-math.MaxFloat32, "-3.4028234663852886e+38"},
		"float32_+max":  {math.MaxFloat32, "3.4028234663852886e+38"},
		"float64":       {float64(9007199254740991), "9.007199254740991e+15"},
		"float64_upper": {float64(9007199254740992), "9.007199254740992e+15"},
		"float64_over":  {float64(9007199254740993), "9.007199254740992e+15"},
		"float64_-max":  {-math.MaxFloat64, "-1.7976931348623157e+308"},
		"float64_+max":  {math.MaxFloat64, "1.7976931348623157e+308"},
		"complex64":     {complex64(math.MaxFloat32 + math.MaxFloat32*1i), "(3.4028235e+38+3.4028235e+38i)"},
		"complex128":    {complex128(math.MaxFloat64 + math.MaxFloat64*1i), "(1.7976931348623157e+308+1.7976931348623157e+308i)"},
		"stringer":      {&testStringer{"foo"}, "foo"},
		"others":        {struct{ a string }{"A"}, "{A}"},
	}
	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			got := zstrconv.FormatAny(tc.a)
			ztesting.AssertEqual(t, "wrong string expression", tc.want, got)
		})
	}
}
