package zaes

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/aileron-projects/go/zcrypto/internal"
)

// EncryptCBC encrypts plaintext with AES CBC cipher.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
// See also [crypto/aes] and [crypto/cipher].
//
// Returned slice consists of:
//   - Initial 16 bytes ([crypto/aes.BlockSize]): Initial vector generated using [crypto/rand.Read].
//   - Rest of the bytes: AES CBC encrypted plaintext (PKCS#7 padding applied).
func EncryptCBC(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := newAES(key)
	if err != nil {
		return nil, err
	}
	plaintext, _ = internal.PadPKCS7(aes.BlockSize, plaintext)
	ciphertext := make([]byte, len(plaintext))
	cipher.NewCBCEncrypter(c, iv).CryptBlocks(ciphertext, plaintext)
	return append(iv, ciphertext...), nil
}

// DecryptCBC decrypts ciphertext encrypted with AES CBC cipher.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
// See also [crypto/aes] and [crypto/cipher].
func DecryptCBC(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < aes.BlockSize || n%aes.BlockSize != 0 {
		return nil, ErrCipherLength(n)
	}
	iv, ciphertext := ciphertext[:aes.BlockSize], ciphertext[aes.BlockSize:]
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(c, iv).CryptBlocks(plaintext, ciphertext)
	return internal.UnpadPKCS7(aes.BlockSize, plaintext)
}
