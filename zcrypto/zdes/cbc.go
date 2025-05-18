package zdes

import (
	"crypto/cipher"
	"crypto/des"

	"github.com/aileron-projects/go/zcrypto/internal"
)

// EncryptCBC encrypts plaintext with DES CBC cipher.
// The key must be the DES key, 8 bytes.
// 8 bytes iv read from [crypto/rand.Reader] is used.
// See more details at [crypto/des] and [crypto/cipher].
//
// Returned slice consists of:
//   - Initial 8 bytes ([crypto/des.BlockSize]): Initial vector generated using [crypto/rand.Read].
//   - Rest of the bytes: DES CBC encrypted plaintext (PKCS#7 padding applied).
func EncryptCBC(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := newDES(key)
	if err != nil {
		return nil, err
	}
	plaintext, _ = internal.PadPKCS7(des.BlockSize, plaintext)
	ciphertext := make([]byte, len(plaintext))
	cipher.NewCBCEncrypter(c, iv).CryptBlocks(ciphertext, plaintext)
	return append(iv, ciphertext...), nil
}

// DecryptCBC decrypts ciphertext encrypted with DES CBC cipher.
// The key must be the DES key, 8 bytes.
// See more details at [crypto/des] and [crypto/cipher].
func DecryptCBC(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < 2*des.BlockSize || n%des.BlockSize != 0 {
		return nil, ErrCipherLength(n)
	}
	iv, ciphertext := ciphertext[:des.BlockSize], ciphertext[des.BlockSize:]

	c, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(c, iv).CryptBlocks(plaintext, ciphertext)
	return internal.UnpadPKCS7(des.BlockSize, plaintext)
}

// EncryptCBC3 encrypts plaintext with 3DES CBC cipher.
// The key must be the 3DES key, 24 bytes.
// 8 bytes iv read from [crypto/rand.Reader] is used.
// See more details at [crypto/des] and [crypto/cipher].
//
// Returned slice consists of:
//   - Initial 8 bytes ([crypto/des.BlockSize]): Initial vector generated using [crypto/rand.Read].
//   - Rest of the bytes: 3DES CBC encrypted plaintext (PKCS#7 padding applied).
func EncryptCBC3(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := new3DES(key)
	if err != nil {
		return nil, err
	}
	plaintext, _ = internal.PadPKCS7(des.BlockSize, plaintext)
	ciphertext := make([]byte, len(plaintext))
	cipher.NewCBCEncrypter(c, iv).CryptBlocks(ciphertext, plaintext)
	return append(iv, ciphertext...), nil
}

// DecryptCBC3 decrypts ciphertext encrypted with 3DES CBC cipher.
// The key must be the 3DES key, 24 bytes.
// See more details at [crypto/des] and [crypto/cipher].
func DecryptCBC3(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < 2*des.BlockSize || n%des.BlockSize != 0 {
		return nil, ErrCipherLength(n)
	}
	iv, ciphertext := ciphertext[:des.BlockSize], ciphertext[des.BlockSize:]

	c, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(c, iv).CryptBlocks(plaintext, ciphertext)
	return internal.UnpadPKCS7(des.BlockSize, plaintext)
}
