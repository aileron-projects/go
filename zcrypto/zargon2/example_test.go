package zargon2_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aileron-projects/go/zcrypto/zargon2"
)

func ExampleArgon2i_Sum() {
	// Replace rand reader temporarily for testing.
	tmp := rand.Reader
	rand.Reader = strings.NewReader("12345678901234567890")
	defer func() { rand.Reader = tmp }()

	saltLen := 10
	crypt, err := zargon2.NewArgon2i(saltLen, 3, 32*1024, 4, 32)
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
	// pw hash : 28b9bb8a406e9dd05f754e60d9f5565f66a6d80a9930f67617ed4e97b710d824
	// overall : 3132333435363738393028b9bb8a406e9dd05f754e60d9f5565f66a6d80a9930f67617ed4e97b710d824
}

func ExampleArgon2i_Compare() {
	// Replace rand reader temporarily for testing.
	tmp := rand.Reader
	rand.Reader = strings.NewReader("12345678901234567890")
	defer func() { rand.Reader = tmp }()

	hashedPW, _ := hex.DecodeString("3132333435363738393028b9bb8a406e9dd05f754e60d9f5565f66a6d80a9930f67617ed4e97b710d824")

	crypt, err := zargon2.NewArgon2i(10, 3, 32*1024, 4, 32)
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

func ExampleArgon2id_Sum() {
	// Replace rand reader temporarily for testing.
	tmp := rand.Reader
	rand.Reader = strings.NewReader("12345678901234567890")
	defer func() { rand.Reader = tmp }()

	saltLen := 10
	crypt, err := zargon2.NewArgon2id(saltLen, 1, 64*1024, 4, 32)
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
	// pw hash : 58b9aa1ca3f06a6615309ff1da98a41296da26ae4d89d31063e2b01d540a8a01
	// overall : 3132333435363738393058b9aa1ca3f06a6615309ff1da98a41296da26ae4d89d31063e2b01d540a8a01
}

func ExampleArgon2id_Compare() {
	// Replace rand reader temporarily for testing.
	tmp := rand.Reader
	rand.Reader = strings.NewReader("12345678901234567890")
	defer func() { rand.Reader = tmp }()

	hashedPW, _ := hex.DecodeString("3132333435363738393058b9aa1ca3f06a6615309ff1da98a41296da26ae4d89d31063e2b01d540a8a01")

	crypt, err := zargon2.NewArgon2id(10, 1, 64*1024, 4, 32)
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
