package zripemd160_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zripemd160"
)

func ExampleSum() {
	// Calculate RIPEMD160 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -ripemd160

	sum := zripemd160.Sum([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zripemd160.EqualSum([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zripemd160.EqualSum([]byte("Bye Go!"), sum))
	// Output:
	// 20 26f145fad89981009615a9af777a10ef700a3bba
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum() {
	// Calculate HMAC-RIPEMD160 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -ripemd160

	key := []byte("secret-key")
	sum := zripemd160.HMACSum([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zripemd160.HMACEqualSum([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zripemd160.HMACEqualSum([]byte("Bye Go!"), key, sum))
	// Output:
	// 20 3528a42c06620b2162e5ad8311e40a8f5bded07b
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}
