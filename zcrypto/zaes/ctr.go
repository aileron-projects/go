package zaes

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"github.com/aileron-projects/go/zcrypto/internal"
)

// EncryptCTR encrypts plaintext with AES CTR cipher.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
// See also [crypto/aes] and [crypto/cipher].
//
// Returned slice consists of:
//   - Initial 16 bytes ([crypto/aes.BlockSize]): Initial vector generated using [crypto/rand.Read].
//   - Rest of the bytes: AES CTR encrypted plaintext.
func EncryptCTR(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := newAES(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(plaintext))
	cipher.NewCTR(c, iv).XORKeyStream(ciphertext, plaintext)
	return append(iv, ciphertext...), nil
}

// DecryptCTR decrypts ciphertext encrypted with AES CTR cipher.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
// See also [crypto/aes] and [crypto/cipher].
func DecryptCTR(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < aes.BlockSize {
		return nil, ErrCipherLength(n)
	}
	iv, ciphertext := ciphertext[:aes.BlockSize], ciphertext[aes.BlockSize:]
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCTR(c, iv).XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

// CopyCTR copies from src to dst.
// If src is encrypted, then decrypted bytes are written into the dst.
// If src is not encrypted, then encrypted bytes are written into the dst.
func CopyCTR(key, iv []byte, dst io.Writer, src io.Reader) error {
	c, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	return internal.Copy(cipher.NewCTR(c, iv), dst, src)
}
