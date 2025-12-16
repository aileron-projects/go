package zargon2_test

import (
	"encoding/hex"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zargon2"
)

func BenchmarkArgon2iSum(b *testing.B) {
	crypt, _ := zargon2.NewArgon2i(10, 3, 32*1024, 4, 32)
	b.ResetTimer()
	for b.Loop() {
		crypt.Sum([]byte("password"))
	}
}

func BenchmarkArgon2iCompare(b *testing.B) {
	crypt, _ := zargon2.NewArgon2i(10, 3, 32*1024, 4, 32)
	hashedPW, _ := hex.DecodeString("3132333435363738393028b9bb8a406e9dd05f754e60d9f5565f66a6d80a9930f67617ed4e97b710d824")
	b.ResetTimer()
	for b.Loop() {
		crypt.Compare(hashedPW, []byte("password"))
	}
}

func BenchmarkArgon2idSum(b *testing.B) {
	crypt, _ := zargon2.NewArgon2id(10, 1, 64*1024, 4, 32)
	b.ResetTimer()
	for b.Loop() {
		crypt.Sum([]byte("password"))
	}
}

func BenchmarkArgon2idCompare(b *testing.B) {
	crypt, _ := zargon2.NewArgon2id(10, 1, 64*1024, 4, 32)
	hashedPW, _ := hex.DecodeString("3132333435363738393058b9aa1ca3f06a6615309ff1da98a41296da26ae4d89d31063e2b01d540a8a01")
	b.ResetTimer()
	for b.Loop() {
		crypt.Compare(hashedPW, []byte("password"))
	}
}
