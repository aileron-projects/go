package zblake2s_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zblake2s"
)

func ExampleSum256() {
	// Calculate BLAKE2s-256 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -blake2s256

	sum := zblake2s.Sum256([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zblake2s.EqualSum256([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zblake2s.EqualSum256([]byte("Bye Go!"), sum))
	// Output:
	// 32 0eb35c8350638fdb21de90833229bad5b19a0d3a9e61cc307c3c9b59859ea24c
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum256() {
	// Calculate HMAC-BLAKE2s-256 hash.
	// Validation data can be generated with:
	// 	- https://emn178.github.io/online-tools/blake2s/
	// 	- https://toolkitbay.com/tkb/tool/BLAKE2s_256

	key := []byte("secret-key")
	sum := zblake2s.HMACSum256([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zblake2s.HMACEqualSum256([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zblake2s.HMACEqualSum256([]byte("Bye Go!"), key, sum))
	// Output:
	// 32 3f010ec9bff3d4f9d0f6d020dcacfe9e2a2da397bce893b3bdfdb038d628875d
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}
