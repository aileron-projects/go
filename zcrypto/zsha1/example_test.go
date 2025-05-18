package zsha1_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zsha1"
)

func ExampleSum() {
	// Calculate SHA1 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha1

	sum := zsha1.Sum([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha1.EqualSum([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha1.EqualSum([]byte("Bye Go!"), sum))
	// Output:
	// 20 cd4ce271604bd88f9ba37f5584814f0560dd7b7b
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum() {
	// Calculate HMAC-SHA1 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha1

	key := []byte("secret-key")
	sum := zsha1.HMACSum([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha1.HMACEqualSum([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha1.HMACEqualSum([]byte("Bye Go!"), key, sum))
	// Output:
	// 20 f897f4708e3a29962f2f65616533ce3a94ec4a26
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}
