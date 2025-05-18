package zsha256_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zsha256"
)

func ExampleSum224() {
	// Calculate SHA256/224 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha224

	sum := zsha256.Sum224([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha256.EqualSum224([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha256.EqualSum224([]byte("Bye Go!"), sum))
	// Output:
	// 28 58a808579d6f7c610cd462d89564c502e45db9c45ed4b376303af7a8
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSum256() {
	// Calculate SHA256 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha256

	sum := zsha256.Sum256([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha256.EqualSum256([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha256.EqualSum256([]byte("Bye Go!"), sum))
	// Output:
	// 32 040bf11d7b007ee960e2b0bed48db8e729c2ba8f1efd649505332c40541673cd
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum224() {
	// Calculate HMAC-SHA256/224 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha224

	key := []byte("secret-key")
	sum := zsha256.HMACSum224([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha256.HMACEqualSum224([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha256.HMACEqualSum224([]byte("Bye Go!"), key, sum))
	// Output:
	// 28 1e8a9de1ec6ba6c1223a95c2ce34679de9f93ac0cc764b4cbf170f0a
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum256() {
	// Calculate HMAC-SHA256 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha256

	key := []byte("secret-key")
	sum := zsha256.HMACSum256([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha256.HMACEqualSum256([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha256.HMACEqualSum256([]byte("Bye Go!"), key, sum))
	// Output:
	// 32 95f09ba10c62fe594ed3439c707c28f44f611547857c0720823a7636dae851c2
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}
