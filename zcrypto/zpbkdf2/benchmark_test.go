package zpbkdf2_test

import (
	"crypto"
	"encoding/hex"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zpbkdf2"
)

func BenchmarkSum(b *testing.B) {
	crypt, _ := zpbkdf2.New(10, 4096, 32, crypto.SHA256)
	b.ResetTimer()
	for b.Loop() {
		crypt.Sum([]byte("password"))
	}
}

func BenchmarkCompare(b *testing.B) {
	crypt, _ := zpbkdf2.New(10, 4096, 32, crypto.SHA256)
	hashedPW, _ := hex.DecodeString("31323334353637383930ec2ca18ec7a42b048b171021a6acda3f953615a38c11d08c503cfa5dfd10ebfa")
	b.ResetTimer()
	for b.Loop() {
		crypt.Compare(hashedPW, []byte("password"))
	}
}
