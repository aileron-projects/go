package zaes

import (
	"crypto/aes"
	"crypto/cipher"
)

// EncryptCFB encrypts plaintext with AES CFB cipher.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
// See also [crypto/aes] and [crypto/cipher].
//
// Returned slice consists of:
//   - Initial 16 bytes ([crypto/aes.BlockSize]): Initial vector generated using [crypto/rand.Read].
//   - Rest of the bytes: AES CFB encrypted plaintext.
//
// Deprecated: See the comment on [crypto/cipher.NewCFBEncrypter].
func EncryptCFB(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := newAES(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(plaintext))
	cipher.NewCFBEncrypter(c, iv).XORKeyStream(ciphertext, plaintext)
	return append(iv, ciphertext...), nil
}

// DecryptCFB decrypts ciphertext encrypted with AES CFB cipher.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
// See also [crypto/aes] and [crypto/cipher].
//
// Deprecated: See the comment on [crypto/cipher.NewCFBDecrypter].
func DecryptCFB(key []byte, ciphertext []byte) ([]byte, error) {
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
	cipher.NewCFBDecrypter(c, iv).XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
