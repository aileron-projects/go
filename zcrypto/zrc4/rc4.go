package zrc4

import (
	"crypto/cipher"
	"crypto/rc4"
	"io"

	"github.com/aileron-projects/go/zcrypto/internal"
)

// NewStreamReader returns a new instance of [crypto/cipher.StreamReader]
// with [crypto/rc4.Cipher].
// As commended on the [crypto/rc4.NewCipher], the key length must be
// at least 1 byte and at most 256 bytes.
func NewStreamReader(key []byte, r io.Reader) (*cipher.StreamReader, error) {
	s, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &cipher.StreamReader{
		S: s,
		R: r,
	}, nil
}

// NewStreamWriter returns a new instance of [crypto/cipher.StreamWriter]
// with [crypto/rc4.Cipher].
// As commended on the [crypto/rc4.NewCipher], the key length must be
// at least 1 byte and at most 256 bytes.
func NewStreamWriter(key []byte, w io.Writer) (*cipher.StreamWriter, error) {
	s, err := rc4.NewCipher(key)
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
// As commended on the [crypto/rc4.NewCipher], the key length must be
// at least 1 byte and at most 256 bytes.
func Copy(key []byte, dst io.Writer, src io.Reader) error {
	s, err := rc4.NewCipher(key)
	if err != nil {
		return err
	}
	return internal.Copy(s, dst, src)
}
