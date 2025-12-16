package zaes_test

import (
	"testing"

	"github.com/aileron-projects/go/zcrypto/zaes"
)

var (
	key       = []byte("1234567890123456")
	plaintext = []byte("Go is an open source programming language that makes it simple to build secure, scalable systems.")
)

func BenchmarkECB(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zaes.EncryptECB(key, plaintext)
		decrypted, _ := zaes.DecryptECB(key, ciphertext)
		_ = decrypted
	}
}

func BenchmarkCBC(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zaes.EncryptCBC(key, plaintext)
		decrypted, _ := zaes.DecryptCBC(key, ciphertext)
		_ = decrypted
	}
}

func BenchmarkCFB(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zaes.EncryptCFB(key, plaintext)
		decrypted, _ := zaes.DecryptCFB(key, ciphertext)
		_ = decrypted
	}
}

func BenchmarkCTR(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zaes.EncryptCTR(key, plaintext)
		decrypted, _ := zaes.DecryptCTR(key, ciphertext)
		_ = decrypted
	}
}

func BenchmarkOFB(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zaes.EncryptOFB(key, plaintext)
		decrypted, _ := zaes.DecryptOFB(key, ciphertext)
		_ = decrypted
	}
}

func BenchmarkGCM(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zaes.EncryptGCM(key, plaintext)
		decrypted, _ := zaes.DecryptGCM(key, ciphertext)
		_ = decrypted
	}
}
