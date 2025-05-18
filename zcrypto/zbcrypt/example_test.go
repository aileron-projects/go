package zbcrypt_test

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/aileron-projects/go/zcrypto/zbcrypt"
)

func ExampleBCrypt_Sum() {
	// Replace rand reader temporarily for testing.
	tmp := rand.Reader
	rand.Reader = strings.NewReader("1234567890123456789012345678901234567890")
	defer func() { rand.Reader = tmp }()

	crypt, err := zbcrypt.New(10)
	if err != nil {
		panic(err)
	}
	hashedPW, err := crypt.Sum([]byte("password"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPW))
	// Output:
	// $2a$10$Lxe3KBCwKxOzLha2MR.vKeGtY6qvS27BZYaR11ITbrDDr2OYCgirC
}

func ExampleBCrypt_Compare() {
	// Replace rand reader temporarily for testing.
	tmp := rand.Reader
	rand.Reader = strings.NewReader("1234567890123456789012345678901234567890")
	defer func() { rand.Reader = tmp }()

	hashedPW := []byte("$2a$10$Lxe3KBCwKxOzLha2MR.vKeGtY6qvS27BZYaR11ITbrDDr2OYCgirC")

	crypt, err := zbcrypt.New(10)
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
