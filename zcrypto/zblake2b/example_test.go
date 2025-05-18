package zblake2b_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zblake2b"
)

func ExampleSum256() {
	// Calculate BLAKE2b-256 hash.
	// Validation data can be generated with:
	// 	- https://emn178.github.io/online-tools/blake2b/
	// 	- https://toolkitbay.com/tkb/tool/BLAKE2b_256

	sum := zblake2b.Sum256([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zblake2b.EqualSum256([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zblake2b.EqualSum256([]byte("Bye Go!"), sum))
	// Output:
	// 32 c871a43cf07e8af76cd5b2ba4cb2da991fb40641bd34b0e906b8cce28febba87
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSum384() {
	// Calculate BLAKE2b-384 hash.
	// Validation data can be generated with:
	// 	- https://emn178.github.io/online-tools/blake2b/
	// 	- https://toolkitbay.com/tkb/tool/BLAKE2b_384

	sum := zblake2b.Sum384([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zblake2b.EqualSum384([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zblake2b.EqualSum384([]byte("Bye Go!"), sum))
	// Output:
	// 48 779700376f9b32c6b7843811a1116570d3ffdd599a7b587fe081aa920829d3ac65e8f6e0bf4d9bcd73eeb5f3f9906577
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSum512() {
	// Calculate BLAKE2b-512 hash.
	// Validation data can be generated with:
	// 	- https://emn178.github.io/online-tools/blake2b/
	// 	- https://toolkitbay.com/tkb/tool/BLAKE2b_512

	sum := zblake2b.Sum512([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zblake2b.EqualSum512([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zblake2b.EqualSum512([]byte("Bye Go!"), sum))
	// Output:
	// 64 02fcc96a71f7df80c2bf8046ef93d2deac992980cef3ea8f7bedd9d4c1e2f87596fb23d527508e8a1c7617510686801a5f714444070f5625572cad1b9fb51bab
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum256() {
	// Calculate HMAC-BLAKE2b-256 hash.
	// Validation data can be generated with:
	// 	- https://emn178.github.io/online-tools/blake2b/
	// 	- https://toolkitbay.com/tkb/tool/BLAKE2b_256

	key := []byte("secret-key")
	sum := zblake2b.HMACSum256([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zblake2b.HMACEqualSum256([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zblake2b.HMACEqualSum256([]byte("Bye Go!"), key, sum))
	// Output:
	// 32 58c69dd63e35195bbda689a35be453959e06e1b9a7de3bb9c41b7f392f80830a
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum384() {
	// Calculate HMAC-BLAKE2b-384 hash.
	// Validation data can be generated with:
	// 	- https://emn178.github.io/online-tools/blake2b/
	// 	- https://toolkitbay.com/tkb/tool/BLAKE2b_384

	key := []byte("secret-key")
	sum := zblake2b.HMACSum384([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zblake2b.HMACEqualSum384([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zblake2b.HMACEqualSum384([]byte("Bye Go!"), key, sum))
	// Output:
	// 48 ccecb9fac6ef0a2955ccf3530c36888cd656c281afd11bff65b7e443465d68ef7c92c0818764e733c5e3a1bea75b4876
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum512() {
	// Calculate HMAC-BLAKE2b-512 hash.
	// Validation data can be generated with:
	// 	- https://emn178.github.io/online-tools/blake2b/
	// 	- https://toolkitbay.com/tkb/tool/BLAKE2b_512

	key := []byte("secret-key")
	sum := zblake2b.HMACSum512([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zblake2b.HMACEqualSum512([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zblake2b.HMACEqualSum512([]byte("Bye Go!"), key, sum))
	// Output:
	// 64 c77b41dd2ae78c64f7655d34da36de59948498056cfdd706fde5edb94fc621aeefcebfa31e8294bbf6eddc1c6398e820700518e8fd7ee5157290a3a57e33525b
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}
