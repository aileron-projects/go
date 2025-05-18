package zstrconv_test

import (
	"fmt"

	"github.com/aileron-projects/go/zerrors"
	"github.com/aileron-projects/go/zstrconv"
)

func ExampleParseNum() {
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[int]("-123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[int8]("-123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[int16]("-123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[int32]("-123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[int64]("-123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[uint]("123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[uint8]("123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[uint16]("123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[uint32]("123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[uint64]("123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[uintptr]("123")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[float32]("1.23")))
	fmt.Printf("%T\n", zerrors.Must(zstrconv.ParseNum[float64]("1.23")))
	// Output:
	// int
	// int8
	// int16
	// int32
	// int64
	// uint
	// uint8
	// uint16
	// uint32
	// uint64
	// uintptr
	// float32
	// float64
}
