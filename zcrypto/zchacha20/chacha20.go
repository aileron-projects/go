package zchacha20

import (
	"crypto/cipher"
	"io"

	"github.com/aileron-projects/go/zcrypto/internal"
	"golang.org/x/crypto/chacha20"
)

// NewStreamReader returns a new instance of [crypto/cipher.StreamReader]
// with [golang.org/x/crypto/chacha20.Cipher].
// The key must be 32 bytes and the nonce must be 12 or 24 bytes.
func NewStreamReader(key, nonce []byte, r io.Reader) (*cipher.StreamReader, error) {
	s, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}
	return &cipher.StreamReader{
		S: s,
		R: r,
	}, nil
}

// NewStreamWriter returns a new instance of [crypto/cipher.StreamWriter]
// with [golang.org/x/crypto/chacha20.Cipher].
// The key must be 32 bytes and the nonce must be 12 or 24 bytes.
func NewStreamWriter(key, nonce []byte, w io.Writer) (*cipher.StreamWriter, error) {
	s, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}
	return &cipher.StreamWriter{
		S: s,
		W: w,
	}, nil
}

// Copy copies from src to dst.
// If src is encrypted, then decrypted bytes are written into the dst.
// If src is not encrypted, then encrypted bytes are written into the dst.
// The key must be 32 bytes and the nonce must be 12 or 24 bytes.
func Copy(key, nonce []byte, dst io.Writer, src io.Reader) error {
	s, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return err
	}
	return internal.Copy(s, dst, src)
}
