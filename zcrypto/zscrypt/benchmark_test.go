package zscrypt_test

import (
	"encoding/hex"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zscrypt"
)

func BenchmarkSum(b *testing.B) {
	crypt, _ := zscrypt.New(32, 32768, 8, 1, 32)
	b.ResetTimer()
	for b.Loop() {
		crypt.Sum([]byte("password"))
	}
}

func BenchmarkCompare(b *testing.B) {
	crypt, _ := zscrypt.New(32, 32768, 8, 1, 32)
	hashedPW, _ := hex.DecodeString("313233343536373839302c717880aa40e2d11cc48c684730cd6f6d87fd2bc7f2aea2ec1246d7019461d2")
	b.ResetTimer()
	for b.Loop() {
		crypt.Compare(hashedPW, []byte("password"))
	}
}
