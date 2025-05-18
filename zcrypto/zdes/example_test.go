package zdes_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zdes"
)

func ExampleEncryptECB() {
	key := []byte("12345678") // 8 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zdes.EncryptECB(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext))

	// Decrypt
	decrypted, _ := zdes.DecryptECB(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 5c74163a9d085a2bf035d51988b42379
	// plaintext
}

func ExampleEncryptECB3() {
	key := []byte("123456789012345678901234") // 24 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zdes.EncryptECB3(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext))

	// Decrypt
	decrypted, _ := zdes.DecryptECB3(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 5f05e14052c798b8c6bf86b2c0295a1f
	// plaintext
}

func ExampleEncryptCBC() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("12345678")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("12345678") // 8 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zdes.EncryptCBC(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext)) // iv is left side of the output.

	// Decrypt
	decrypted, _ := zdes.DecryptCBC(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 31323334353637385194bfcb8072ea4b1e3d83016ef65d8a
	// plaintext
}

func ExampleEncryptCBC3() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("12345678")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("123456789012345678901234") // 24 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zdes.EncryptCBC3(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext)) // iv is left side of the output.

	// Decrypt
	decrypted, _ := zdes.DecryptCBC3(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 3132333435363738f68b09f9c8e1893b8682059b7db02962
	// plaintext
}

func ExampleEncryptCFB() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("12345678")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("12345678") // 8 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zdes.EncryptCFB(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext)) // iv is left side of the output.

	// Decrypt
	decrypted, _ := zdes.DecryptCFB(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 3132333435363738e6bc63e116a1e9f15c
	// plaintext
}

func ExampleEncryptCFB3() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("12345678")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("123456789012345678901234") // 24 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zdes.EncryptCFB3(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext)) // iv is left side of the output.

	// Decrypt
	decrypted, _ := zdes.DecryptCFB3(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 31323334353637386db9685cd9766d33a2
	// plaintext
}

func ExampleEncryptCTR() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("12345678")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("12345678") // 8 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zdes.EncryptCTR(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext)) // iv is left side of the output.

	// Decrypt
	decrypted, _ := zdes.DecryptCTR(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 3132333435363738e6bc63e116a1e9f135
	// plaintext
}

func ExampleEncryptCTR3() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("12345678")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("123456789012345678901234") // 24 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zdes.EncryptCTR3(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext)) // iv is left side of the output.

	// Decrypt
	decrypted, _ := zdes.DecryptCTR3(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 31323334353637386db9685cd9766d33d3
	// plaintext
}

func ExampleEncryptOFB() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("12345678")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("12345678") // 8 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zdes.EncryptOFB(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext)) // iv is left side of the output.

	// Decrypt
	decrypted, _ := zdes.DecryptOFB(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 3132333435363738e6bc63e116a1e9f19d
	// plaintext
}

func ExampleEncryptOFB3() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("12345678")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("123456789012345678901234") // 24 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zdes.EncryptOFB3(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext)) // iv is left side of the output.

	// Decrypt
	decrypted, _ := zdes.DecryptOFB3(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 31323334353637386db9685cd9766d337e
	// plaintext
}
