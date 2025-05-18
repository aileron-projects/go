package zpbkdf2_test

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aileron-projects/go/zcrypto/zpbkdf2"
)

func ExamplePBKDF2_Sum() {
	// Replace rand reader temporarily for testing.
	tmp := rand.Reader
	rand.Reader = strings.NewReader("12345678901234567890")
	defer func() { rand.Reader = tmp }()

	saltLen := 10
	crypt, err := zpbkdf2.New(saltLen, 4096, 32, crypto.SHA256)
	if err != nil {
		panic(err)
	}
	hashedPW, err := crypt.Sum([]byte("password"))
	if err != nil {
		panic(err)
	}
	fmt.Println("salt    :", hex.EncodeToString(hashedPW[:saltLen]))
	fmt.Println("pw hash :", hex.EncodeToString(hashedPW[saltLen:]))
	fmt.Println("overall :", hex.EncodeToString(hashedPW))
	// Output:
	// salt    : 31323334353637383930
	// pw hash : ec2ca18ec7a42b048b171021a6acda3f953615a38c11d08c503cfa5dfd10ebfa
	// overall : 31323334353637383930ec2ca18ec7a42b048b171021a6acda3f953615a38c11d08c503cfa5dfd10ebfa
}

func ExamplePBKDF2_Compare() {
	// Replace rand reader temporarily for testing.
	tmp := rand.Reader
	rand.Reader = strings.NewReader("12345678901234567890")
	defer func() { rand.Reader = tmp }()

	hashedPW, _ := hex.DecodeString("31323334353637383930ec2ca18ec7a42b048b171021a6acda3f953615a38c11d08c503cfa5dfd10ebfa")

	saltLen := 10
	crypt, err := zpbkdf2.New(saltLen, 4096, 32, crypto.SHA256)
	if err != nil {
		panic(err)
	}
	err = crypt.Compare(hashedPW, []byte("password"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Is password correct? :", err == nil)
	// Output:
	// Is password correct? : true
}
