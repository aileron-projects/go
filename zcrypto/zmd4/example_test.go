package zmd4_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zmd4"
)

func ExampleSum() {
	// Calculate MD4 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -md4

	sum := zmd4.Sum([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zmd4.EqualSum([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zmd4.EqualSum([]byte("Bye Go!"), sum))
	// Output:
	// 16 38b77b3efd9016c334cb3b45da2bb00e
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum() {
	// Calculate HMAC-MD4 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -md4

	key := []byte("secret-key")
	sum := zmd4.HMACSum([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zmd4.HMACEqualSum([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zmd4.HMACEqualSum([]byte("Bye Go!"), key, sum))
	// Output:
	// 16 ac740e3ffa02eabe6ffd06a216bd1447
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}
