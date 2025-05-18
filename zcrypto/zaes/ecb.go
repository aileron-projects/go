package zaes

import (
	"crypto/aes"

	"github.com/aileron-projects/go/zcrypto/internal"
)

// EncryptECB encrypts plaintext with AES ECB cipher.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
// See also [crypto/aes] and [crypto/cipher].
//
// Returned slice consists of:
//   - All parts: AES ECB encrypted plaintext (PKCS#7 padding applied).
func EncryptECB(key []byte, plaintext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext, _ = internal.PadPKCS7(aes.BlockSize, plaintext)
	ciphertext := make([]byte, len(plaintext))
	for i := range len(plaintext) / aes.BlockSize {
		start, end := i*aes.BlockSize, (i+1)*aes.BlockSize
		c.Encrypt(ciphertext[start:end], plaintext[start:end])
	}
	return ciphertext, nil
}

// DecryptECB decrypts ciphertext encrypted with AES ECB cipher.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
// See also [crypto/aes] and [crypto/cipher].
func DecryptECB(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < aes.BlockSize || n%aes.BlockSize != 0 {
		return nil, ErrCipherLength(n)
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	for i := range len(ciphertext) / aes.BlockSize {
		start, end := i*aes.BlockSize, (i+1)*aes.BlockSize
		c.Decrypt(plaintext[start:end], ciphertext[start:end])
	}
	return internal.UnpadPKCS7(aes.BlockSize, plaintext)
}
