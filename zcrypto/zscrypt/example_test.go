package zscrypt_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aileron-projects/go/zcrypto/zscrypt"
)

func ExampleSCrypt_Sum() {
	// Replace rand reader temporarily for testing.
	tmp := rand.Reader
	rand.Reader = strings.NewReader("12345678901234567890")
	defer func() { rand.Reader = tmp }()

	saltLen := 10
	crypt, err := zscrypt.New(saltLen, 32768, 8, 1, 32)
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
	// pw hash : 2c717880aa40e2d11cc48c684730cd6f6d87fd2bc7f2aea2ec1246d7019461d2
	// overall : 313233343536373839302c717880aa40e2d11cc48c684730cd6f6d87fd2bc7f2aea2ec1246d7019461d2
}

func ExampleSCrypt_Compare() {
	// Replace rand reader temporarily for testing.
	tmp := rand.Reader
	rand.Reader = strings.NewReader("12345678901234567890")
	defer func() { rand.Reader = tmp }()

	hashedPW, _ := hex.DecodeString("313233343536373839302c717880aa40e2d11cc48c684730cd6f6d87fd2bc7f2aea2ec1246d7019461d2")

	crypt, err := zscrypt.New(10, 32768, 8, 1, 32)
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
