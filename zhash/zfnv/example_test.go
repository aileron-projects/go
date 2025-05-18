package zfnv_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zhash/zfnv"
)

func ExampleSum32() {
	// Calculate fnv1/32 hash.
	// Validation code using python (https://pypi.org/project/fnv/):
	//    import fnv
	//    sum = fnv.hash(b"Hello Go!", algorithm=fnv.fnv, bits=32)
	//    print(f"{sum:x}")

	sum := zfnv.Sum32([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 4 5b81b1e4
}

func ExampleSum32a() {
	// Calculate fnv1a/32 hash.
	// Validation code using python (https://pypi.org/project/fnv/):
	//    import fnv
	//    sum = fnv.hash(b"Hello Go!", algorithm=fnv.fnv_1a, bits=32)
	//    print(f"{sum:x}")

	sum := zfnv.Sum32a([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 4 6c8aa582
}

func ExampleSum64() {
	// Calculate fnv1/64 hash.
	// Validation code using python (https://pypi.org/project/fnv/):
	//    import fnv
	//    sum = fnv.hash(b"Hello Go!", algorithm=fnv.fnv, bits=64)
	//    print(f"{sum:x}")

	sum := zfnv.Sum64([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 8 ffd6428808d948a4
}

func ExampleSum64a() {
	// Calculate fnv1a/64 hash.
	// Validation code using python (https://pypi.org/project/fnv/):
	//    import fnv
	//    sum = fnv.hash(b"Hello Go!", algorithm=fnv.fnv_1a, bits=64)
	//    print(f"{sum:x}")

	sum := zfnv.Sum64a([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 8 55a86077c4d631a2
}

func ExampleSum128() {
	// Calculate fnv1/128 hash.
	// Validation code using python (https://pypi.org/project/fnv/):
	//    import fnv
	//    sum = fnv.hash(b"Hello Go!", algorithm=fnv.fnv, bits=128)
	//    print(f"{sum:x}")

	sum := zfnv.Sum128([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 16 8d7b68876003b298410a7aab37c42f94
}

func ExampleSum128a() {
	// Calculate fnv1a/128 hash.
	// Validation code using python (https://pypi.org/project/fnv/):
	//    import fnv
	//    sum = fnv.hash(b"Hello Go!", algorithm=fnv.fnv_1a, bits=128)
	//    print(f"{sum:x}")

	sum := zfnv.Sum128a([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	// Output:
	// 16 7428d4b281051c165b1acb019a6c08ea
}
