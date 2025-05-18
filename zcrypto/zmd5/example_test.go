package zmd5_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zmd5"
)

func ExampleSum() {
	// Calculate MD5 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -md5

	sum := zmd5.Sum([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zmd5.EqualSum([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zmd5.EqualSum([]byte("Bye Go!"), sum))
	// Output:
	// 16 73a6333befc4397176b63dd4540ad5fe
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum() {
	// Calculate HMAC-MD5 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -md5

	key := []byte("secret-key")
	sum := zmd5.HMACSum([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zmd5.HMACEqualSum([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zmd5.HMACEqualSum([]byte("Bye Go!"), key, sum))
	// Output:
	// 16 951010d3df8a7f9a2e1ea2a8db158dc6
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}
