package zcrc64_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zhash/zcrc64"
)

func ExampleSumISO() {
	// Calculate CRC64 hash using ISO table.
	// Validation data can be generated with:
	// 	- https://www.sunshine2k.de/coding/javascript/crc/crc_js.html
	// 	- https://toolkitbay.com/tkb/tool/CRC-64

	sum := zcrc64.SumISO([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 8 d22707d9ac3eeee2
}

func ExampleSumECMA() {
	// Calculate CRC64 hash using ECMA table.
	// Validation data can be generated with:
	// 	- https://www.sunshine2k.de/coding/javascript/crc/crc_js.html
	// 	- https://toolkitbay.com/tkb/tool/CRC-64

	sum := zcrc64.SumECMA([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 8 85d207c9ff681c7c
}
