package zcrc32_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zhash/zcrc32"
)

func ExampleSumIEEE() {
	// Calculate CRC32 hash using [crc32.IEEE] table.
	// Validation data can be generated with:
	// 	- https://www.sunshine2k.de/coding/javascript/crc/crc_js.html
	// 	- https://crccalc.com/?crc=Hello%20Go!&method=CRC-32&datatype=ascii&outtype=hex

	sum := zcrc32.SumIEEE([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 4 48e1410f
}

func ExampleSumCastagnoli() {
	// Calculate CRC32 hash using [crc32.Castagnoli] table.
	// Validation data can be generated with:
	// 	- https://www.sunshine2k.de/coding/javascript/crc/crc_js.html
	// 	- https://crccalc.com/?crc=Hello%20Go!&method=CRC-32&datatype=ascii&outtype=hex

	sum := zcrc32.SumCastagnoli([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 4 c83c2422
}

func ExampleSumKoopman() {
	sum := zcrc32.SumKoopman([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 4 1f83620a
}
