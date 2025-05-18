package zaes_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zcrypto/zaes"
)

func ExampleEncryptCBC() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("1234567890123456")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("1234567890123456") // 16 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zaes.EncryptCBC(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext))

	// Decrypt
	decrypted, _ := zaes.DecryptCBC(key, ciphertext)
	fmt.Println(string(decrypted))
	// Output:
	// 31323334353637383930313233343536b72e4b49f65e860eb712f1fbb2b3cc1b
	// plaintext
}

func ExampleEncryptCFB() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("1234567890123456")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("1234567890123456") // 16 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zaes.EncryptCFB(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext))

	// Decrypt
	decrypted, _ := zaes.DecryptCFB(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 313233343536373839303132333435360510ac65b228f592af
	// plaintext
}

func ExampleEncryptCTR() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("1234567890123456")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("1234567890123456") // 16 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zaes.EncryptCTR(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext))

	// Decrypt
	decrypted, _ := zaes.DecryptCTR(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 313233343536373839303132333435360510ac65b228f592af
	// plaintext
}

func ExampleEncryptECB() {
	key := []byte("1234567890123456") // 16 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zaes.EncryptECB(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext))

	// Decrypt
	decrypted, _ := zaes.DecryptECB(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// d2394f972abe3523f8d8b43cd007ff0b
	// plaintext
}

func ExampleEncryptGCM() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("1234567890123456")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("1234567890123456") // 16 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zaes.EncryptGCM(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext))

	// Decrypt
	decrypted, _ := zaes.DecryptGCM(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 313233343536373839303132c11d0fc1b6e53170d378161ce9b917f4df1827853a378a6686
	// plaintext
}

func ExampleEncryptOFB() {
	tmp := rand.Reader
	rand.Reader = bytes.NewReader([]byte("1234567890123456")) // Fix value for test.
	defer func() { rand.Reader = tmp }()

	key := []byte("1234567890123456") // 16 bytes.
	plaintext := []byte("plaintext")

	// Encrypt
	ciphertext, _ := zaes.EncryptOFB(key, plaintext)
	fmt.Println(hex.EncodeToString(ciphertext))

	// Decrypt
	decrypted, _ := zaes.DecryptOFB(key, ciphertext)
	fmt.Println(string(decrypted))

	// Output:
	// 313233343536373839303132333435360510ac65b228f592af
	// plaintext
}
