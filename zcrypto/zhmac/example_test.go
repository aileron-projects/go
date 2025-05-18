package zhmac_test

import (
	"crypto"
	_ "crypto/sha256"
	_ "crypto/sha3"
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zhmac"
	_ "golang.org/x/crypto/blake2b"
	_ "golang.org/x/crypto/blake2s"
)

func ExampleSum_sha256() {
	// The hmac value can be validated using openssl.
	// echo -n "hello world" | openssl dgst -hmac "1234567890" -sha256
	msg := []byte("hello world")
	key := []byte("1234567890")
	b := zhmac.Sum(crypto.SHA256, msg, key)
	fmt.Println(hex.EncodeToString(b))
	// Output:
	// e45f6072617cce5e06a95be481c43351023d99233599eac9fdaffc958142629a
}

func ExampleSum_sha3() {
	// The hmac value can be validated using openssl.
	// echo -n "hello world" | openssl dgst -hmac "1234567890" -sha3-256
	msg := []byte("hello world")
	key := []byte("1234567890")
	b := zhmac.Sum(crypto.SHA3_256, msg, key)
	fmt.Println(hex.EncodeToString(b))
	// Output:
	// fa66d35e0fba7a2ca62b61b0fb837ad280c837a09f250f4774c4b0253e325f7e
}

func ExampleSum_blake2b() {
	// The hmac value can be validated using python.
	// Note that it does not use blake2b native HMAC mechanism
	//   >> import hmac, hashlib
	//   >> hmac.new(b"1234567890", b"hello world", lambda:hashlib.blake2b(digest_size=32)).hexdigest()
	msg := []byte("hello world")
	key := []byte("1234567890")
	b := zhmac.Sum(crypto.BLAKE2b_256, msg, key)
	fmt.Println(hex.EncodeToString(b))
	// Output:
	// 2c7e515a2659e64117f1111d0db87194b2dfe6aa84c7946fcef09afb5e1a864a
}

func ExampleSum_blake22() {
	// The hmac value can be validated using python.
	// Note that it does not use blake2s native HMAC mechanism
	//   >> import hmac, hashlib
	//   >> hmac.new(b"1234567890", b"hello world", lambda:hashlib.blake2s(digest_size=32)).hexdigest()
	msg := []byte("hello world")
	key := []byte("1234567890")
	b := zhmac.Sum(crypto.BLAKE2s_256, msg, key)
	fmt.Println(hex.EncodeToString(b))
	// Output:
	// dfb87acb873eebfb2b8cb06465dbbea440519fb39fae2b5d1f63a165ab331a4f
}
