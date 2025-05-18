package zsha512_test

import (
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zsha512"
)

func ExampleSum224() {
	// Calculate SHA512/224 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha512-224

	sum := zsha512.Sum224([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha512.EqualSum224([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha512.EqualSum224([]byte("Bye Go!"), sum))
	// Output:
	// 28 56724525f8dd95140fbe710684017a97fb4f9419d1b8bc6e62a41bf3
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSum256() {
	// Calculate SHA512/256 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha512-256

	sum := zsha512.Sum256([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha512.EqualSum256([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha512.EqualSum256([]byte("Bye Go!"), sum))
	// Output:
	// 32 ac553f354e4db90319932627c16f2bea019920ff8b9023cf51aff240bd2f4289
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSum384() {
	// Calculate SHA512/384 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha384

	sum := zsha512.Sum384([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha512.EqualSum384([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha512.EqualSum384([]byte("Bye Go!"), sum))
	// Output:
	// 48 def7f878cc707817ba94a613b6475d0f53fc5844c78d7e26f9b9b13a3b3437a92b01fe4446259170163bfd3f7b33315f
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleSum512() {
	// Calculate SHA512 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -sha512

	sum := zsha512.Sum512([]byte("Hello Go!"))
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha512.EqualSum512([]byte("Hello Go!"), sum))
	fmt.Println("`Bye Go!` match?", zsha512.EqualSum512([]byte("Bye Go!"), sum))
	// Output:
	// 64 77ada5c6a5d44b4fe6a6f6509310be924a882317da25b9002a69f2a5bae2588f728d573472d5bba7317f330f48fa15a1b5f4f31494051f29479a71ed310f44da
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum224() {
	// Calculate HMAC-SHA512/224 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha512-224

	key := []byte("secret-key")
	sum := zsha512.HMACSum224([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha512.HMACEqualSum224([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha512.HMACEqualSum224([]byte("Bye Go!"), key, sum))
	// Output:
	// 28 4b1e5cb41a254bd4e1290c15c0a1861be0259711ca074dd1c77723fb
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum256() {
	// Calculate HMAC-SHA512/256 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha512-256

	key := []byte("secret-key")
	sum := zsha512.HMACSum256([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha512.HMACEqualSum256([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha512.HMACEqualSum256([]byte("Bye Go!"), key, sum))
	// Output:
	// 32 a5f9097d0b51c0f867edb0689f1dcf84da7e0a427543d1c39fef60d23d9d1707
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum384() {
	// Calculate HMAC-SHA512/384 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha384

	key := []byte("secret-key")
	sum := zsha512.HMACSum384([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha512.HMACEqualSum384([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha512.HMACEqualSum384([]byte("Bye Go!"), key, sum))
	// Output:
	// 48 d9ecbd065ed77bdd2db653e1dd18601ead2c1b7537509c537f990d70cce66fceed42e04249494fc6f1809324d097224f
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}

func ExampleHMACSum512() {
	// Calculate HMAC-SHA512 hash.
	// Validation data can be generated with:
	// 	- echo -n "Hello Go!" | openssl dgst -hmac "secret-key" -sha512

	key := []byte("secret-key")
	sum := zsha512.HMACSum512([]byte("Hello Go!"), key)
	encoded := hex.EncodeToString(sum)
	fmt.Println(len(sum), encoded)
	fmt.Println("`Hello Go!` match?", zsha512.HMACEqualSum512([]byte("Hello Go!"), key, sum))
	fmt.Println("`Bye Go!` match?", zsha512.HMACEqualSum512([]byte("Bye Go!"), key, sum))
	// Output:
	// 64 bd0c9d973e78f334f98074ee6353a902399176f2a12beb85ae4b7e3a8cd262388980f13150ea8fcee1c76975e6db5facf3d4ec99ddeff5cad9d954cdd72a0abd
	// `Hello Go!` match? true
	// `Bye Go!` match? false
}
