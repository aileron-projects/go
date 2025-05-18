package zdes_test

import (
	"testing"

	"github.com/aileron-projects/go/zcrypto/zdes"
)

var (
	key8      = []byte("12345678")
	key24     = []byte("123456789012345678901234")
	plaintext = []byte("Go is an open source programming language that makes it simple to build secure, scalable systems.")
)

func BenchmarkECB(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zdes.EncryptECB(key8, plaintext)
		decrypted, _ := zdes.DecryptECB(key8, ciphertext)
		_ = decrypted
	}
}

func BenchmarkCBC(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zdes.EncryptCBC(key8, plaintext)
		decrypted, _ := zdes.DecryptCBC(key8, ciphertext)
		_ = decrypted
	}
}

func BenchmarkCFB(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zdes.EncryptCFB(key8, plaintext)
		decrypted, _ := zdes.DecryptCFB(key8, ciphertext)
		_ = decrypted
	}
}

func BenchmarkCTR(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zdes.EncryptCTR(key8, plaintext)
		decrypted, _ := zdes.DecryptCTR(key8, ciphertext)
		_ = decrypted
	}
}

func BenchmarkOFB(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zdes.EncryptOFB(key8, plaintext)
		decrypted, _ := zdes.DecryptOFB(key8, ciphertext)
		_ = decrypted
	}
}

func BenchmarkECB3(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zdes.EncryptECB3(key24, plaintext)
		decrypted, _ := zdes.DecryptECB3(key24, ciphertext)
		_ = decrypted
	}
}

func BenchmarkCBC3(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zdes.EncryptCBC3(key24, plaintext)
		decrypted, _ := zdes.DecryptCBC3(key24, ciphertext)
		_ = decrypted
	}
}

func BenchmarkCFB3(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zdes.EncryptCFB3(key24, plaintext)
		decrypted, _ := zdes.DecryptCFB3(key24, ciphertext)
		_ = decrypted
	}
}

func BenchmarkCTR3(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zdes.EncryptCTR3(key24, plaintext)
		decrypted, _ := zdes.DecryptCTR3(key24, ciphertext)
		_ = decrypted
	}
}

func BenchmarkOFB3(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		ciphertext, _ := zdes.EncryptOFB3(key24, plaintext)
		decrypted, _ := zdes.DecryptOFB3(key24, ciphertext)
		_ = decrypted
	}
}
