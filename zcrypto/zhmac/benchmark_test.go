package zhmac_test

import (
	"crypto"
	"testing"

	_ "github.com/aileron-projects/go/zcrypto/zblake2b"
	_ "github.com/aileron-projects/go/zcrypto/zblake2s"
	"github.com/aileron-projects/go/zcrypto/zhmac"
	_ "github.com/aileron-projects/go/zcrypto/zmd4"
	_ "github.com/aileron-projects/go/zcrypto/zmd5"
	_ "github.com/aileron-projects/go/zcrypto/zripemd160"
	_ "github.com/aileron-projects/go/zcrypto/zsha1"
	_ "github.com/aileron-projects/go/zcrypto/zsha256"
	_ "github.com/aileron-projects/go/zcrypto/zsha3"
	_ "github.com/aileron-projects/go/zcrypto/zsha512"
)

var (
	msg = []byte("Go is an open source programming language that makes it simple to build secure, scalable systems.")
	key = []byte("12345678901234567890123456789012")
)

func BenchmarkSumMD4(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.MD4, msg, key)
	}
}

func BenchmarkSumMD5(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.MD5, msg, key)
	}
}

func BenchmarkSumSHA1(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA1, msg, key)
	}
}

func BenchmarkSumSHA224(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA1, msg, key)
	}
}

func BenchmarkSumSHA256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA256, msg, key)
	}
}

func BenchmarkSumSHA384(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA384, msg, key)
	}
}

func BenchmarkSumSHA512(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA512, msg, key)
	}
}

func BenchmarkSumRipemd160(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.RIPEMD160, msg, key)
	}
}

func BenchmarkSumSHA3_224(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA3_224, msg, key)
	}
}

func BenchmarkSumSHA3_256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA3_256, msg, key)
	}
}

func BenchmarkSumSHA3_384(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA3_384, msg, key)
	}
}

func BenchmarkSumSHA3_512(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA3_512, msg, key)
	}
}

func BenchmarkSumSHA512_224(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA512_224, msg, key)
	}
}

func BenchmarkSumSHA512_256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.SHA512_256, msg, key)
	}
}

func BenchmarkSumBlake2s_256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.BLAKE2s_256, msg, key)
	}
}

func BenchmarkSumBlake2b_256(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.BLAKE2b_256, msg, key)
	}
}

func BenchmarkSumBlake2b_384(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.BLAKE2b_384, msg, key)
	}
}

func BenchmarkSumBlake2b_512(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zhmac.Sum(crypto.BLAKE2b_512, msg, key)
	}
}
