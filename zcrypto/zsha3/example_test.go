package zsha3_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zsha3"
)

func ExampleSum224() {
	// Calculate SHA3/224 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha3-224

	sum := zsha3.Sum224([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha3.EqualSum224([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha3.EqualSum224([]byte("Bye Go!"), sum))
	// Output:
	// 28 ff3a0dcb416052d683762cc81b8c6cdabb2946a32fc36760708daa3d
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSum256() {
	// Calculate SHA3/256 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha3-256

	sum := zsha3.Sum256([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha3.EqualSum256([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha3.EqualSum256([]byte("Bye Go!"), sum))
	// Output:
	// 32 762e68ba83959cee8da6f2255c7e6df0b1d3875b073e692d9f649dc8db4d43bc
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSum384() {
	// Calculate SHA3/384 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha3-384

	sum := zsha3.Sum384([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha3.EqualSum384([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha3.EqualSum384([]byte("Bye Go!"), sum))
	// Output:
	// 48 64cb4df91df2fe72861eb7e86361a1d2ab76618442410ab671f4b3712c0376c7cb29a0b26bed63c109a30194402e883c
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSum512() {
	// Calculate SHA3/512 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha3-512

	sum := zsha3.Sum512([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha3.EqualSum512([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha3.EqualSum512([]byte("Bye Go!"), sum))
	// Output:
	// 64 1ceb2a3b4536c1fc950273ea0c5c208f87ca444b41f2726b3c9bd7c25e5421f851b87162a7e3bc39eac0335faf8e27eb1eb42d419894dde64b3c6e877f7666d3
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSumShake128() {
	// Calculate SHAKE128 hash.
	// Validation data can be generated with(-oflen option requires openssl v3):
	// 	- echo -n "Hello Go!" | openssl dgst -shake128 -xoflen 32
	// 	- https://www.cryptool.org/en/cto/openssl/

	sum := zsha3.SumShake128([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha3.EqualSumShake128([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha3.EqualSumShake128([]byte("Bye Go!"), sum))
	// Output:
	// 32 d52544db366d156a3ed524fc3d928d489d49a85c10cd268c64c6f0055dce4a5b
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSumShake256() {
	// Calculate SHAKE256 hash.
	// Validation data can be generated with:
	// Validation data can be generated with(-oflen option requires openssl v3):
	// 	- echo -n "Hello Go!" | openssl dgst -shake256 -xoflen 64
	// 	- https://www.cryptool.org/en/cto/openssl/

	sum := zsha3.SumShake256([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha3.EqualSumShake256([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha3.EqualSumShake256([]byte("Bye Go!"), sum))
	// Output:
	// 64 755bb7002e7cbd7282cd87f4a702ba9d98573de35378fd78451bcb6a41ebf138c4035ab9b92e4d6f1041c604d0348d6d0009f1db3419560b745f3c31e6bfc463
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum224() {
	// Calculate HMAC-SHA3/224 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha3-224

	key := []byte("secret-key")
	sum := zsha3.HMACSum224([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha3.HMACEqualSum224([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha3.HMACEqualSum224([]byte("Bye Go!"), key, sum))
	// Output:
	// 28 d68000ecd1481e7f8ee9e1695a6a92cbe0d6395c046b42a61ff53b77
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum256() {
	// Calculate HMAC-SHA3/256 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha3-256

	key := []byte("secret-key")
	sum := zsha3.HMACSum256([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha3.HMACEqualSum256([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha3.HMACEqualSum256([]byte("Bye Go!"), key, sum))
	// Output:
	// 32 efdfb3b8f5075ce6d3ca8dc2dbd38cb345e2287fdeb9282d2bb6a4178f688a40
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum384() {
	// Calculate HMAC-SHA3/384 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha3-384

	key := []byte("secret-key")
	sum := zsha3.HMACSum384([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha3.HMACEqualSum384([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha3.HMACEqualSum384([]byte("Bye Go!"), key, sum))
	// Output:
	// 48 d94e646a24894579d228743fa0474f338ddfc127c689c125d31df209d7428e9427ed7a7a9d376e11f93e6c9043d3317d
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum512() {
	// Calculate HMAC-SHA3/512 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha3-512

	key := []byte("secret-key")
	sum := zsha3.HMACSum512([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha3.HMACEqualSum512([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha3.HMACEqualSum512([]byte("Bye Go!"), key, sum))
	// Output:
	// 64 aceaa099f2c7c8bc17ef485bcdab407d4f23de9b312494a99b637ca90dd176c05897417a2b066f7dd4a003b554b15bf2cb8403f43c1b39ef20f74892f45f6129
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}
