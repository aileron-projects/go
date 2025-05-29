package zbcrypt_test

import (
	"testing"

	"github.com/aileron-projects/go/zcrypto/zbcrypt"
)

func BenchmarkSum(b *testing.B) {
	crypt, _ := zbcrypt.New(10)
	b.ResetTimer()
	for range b.N {
		crypt.Sum([]byte("password"))
	}
}

func BenchmarkCompare(b *testing.B) {
	crypt, _ := zbcrypt.New(10)
	hashedPW := []byte("$2a$10$Lxe3KBCwKxOzLha2MR.vKeGtY6qvS27BZYaR11ITbrDDr2OYCgirC")
	b.ResetTimer()
	for range b.N {
		crypt.Compare(hashedPW, []byte("password"))
	}
}
