package zaes

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/aileron-projects/go/internal/helper"
)

// EncryptGCM encrypts plaintext with AES GCM cipher.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
// It uses 12 bytes nonce and 16 bytes tag set in the [crypto/cipher.NewGCM].
// See also [crypto/aes] and [crypto/cipher].
//
// Returned slice consists of:
//   - Initial 12 bytes: Nonce generated using [crypto/rand.Read].
//   - Rest of the bytes: AES GCM encrypted plaintext.
func EncryptGCM(key []byte, plaintext []byte) ([]byte, error) {
	c, nonce, err := newAES(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(c) // cipher with 12 bytes nonce, 16 bytes tag.
	helper.MustNil(err)
	// nonce should be an unique value. No need to be a random value.
	// NewGCM uses standard nonce size of 12 bytes.
	nonce = nonce[:aead.NonceSize()] // Adjust size.
	// Append encrypted text to the nonce slice.
	// That means ciphertext = append(nonce, encrypt(plaintext)).
	ciphertext := aead.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptGCM decrypts ciphertext encrypted with AES GCM cipher.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
// See also [crypto/aes] and [crypto/cipher].
func DecryptGCM(key []byte, ciphertext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(c) // cipher with 12 bytes nonce, 16 bytes tag.
	helper.MustNil(err)
	if len(ciphertext) < aead.NonceSize() {
		return nil, ErrCipherLength(len(ciphertext))
	}
	nonce, ciphertext := ciphertext[:aead.NonceSize()], ciphertext[aead.NonceSize():]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}
