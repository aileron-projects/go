package zdes

import (
	"crypto/des"

	"github.com/aileron-projects/go/zcrypto/internal"
)

// EncryptECB encrypts plaintext with DES ECB cipher.
// The key must be the DES key, 8 bytes.
// ECB mode does not use iv.
// See more details at [crypto/des] and [crypto/cipher].
//
// Returned slice consists of:
//   - All parts: DES ECB encrypted plaintext (PKCS#7 padding applied).
//
// Deprecated: ECB mode is not recommended because it does not uses iv.
// Same plaintext is alway be encrypted into the same ciphertext.
// If block cipher is necessary, use CBC mode instead.
func EncryptECB(key []byte, plaintext []byte) ([]byte, error) {
	c, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext, _ = internal.PadPKCS7(des.BlockSize, plaintext)
	ciphertext := make([]byte, len(plaintext))
	for i := range len(plaintext) / des.BlockSize {
		start, end := i*des.BlockSize, (i+1)*des.BlockSize
		c.Encrypt(ciphertext[start:end], plaintext[start:end])
	}
	return ciphertext, nil
}

// DecryptECB decrypts ciphertext encrypted with DES ECB cipher.
// The key must be the DES key, 8 bytes.
// See more details at [crypto/des] and [crypto/cipher].
//
// Deprecated: ECB mode is not recommended because it does not uses iv.
// Same plaintext is alway be encrypted into the same ciphertext.
// If block cipher is necessary, use CBC mode instead.
func DecryptECB(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < des.BlockSize || n%des.BlockSize != 0 {
		return nil, ErrCipherLength(n)
	}
	c, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	for i := range len(ciphertext) / des.BlockSize {
		start, end := i*des.BlockSize, (i+1)*des.BlockSize
		c.Decrypt(plaintext[start:end], ciphertext[start:end])
	}
	return internal.UnpadPKCS7(des.BlockSize, plaintext)
}

// EncryptECB3 encrypts plaintext with 3DES ECB cipher.
// The key must be the 3DES key, 24 bytes.
// ECB mode does not use iv.
// See more details at [crypto/des] and [crypto/cipher].
//
// Returned slice consists of:
//   - All parts: 3DES ECB encrypted plaintext (PKCS#7 padding applied).
//
// Deprecated: ECB mode is not recommended because it does not uses iv.
// Same plaintext is alway be encrypted into the same ciphertext.
// If block cipher is necessary, use CBC mode instead.
func EncryptECB3(key []byte, plaintext []byte) ([]byte, error) {
	c, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext, _ = internal.PadPKCS7(des.BlockSize, plaintext)
	ciphertext := make([]byte, len(plaintext))
	for i := range len(plaintext) / des.BlockSize {
		start, end := i*des.BlockSize, (i+1)*des.BlockSize
		c.Encrypt(ciphertext[start:end], plaintext[start:end])
	}
	return ciphertext, nil
}

// DecryptECB3 decrypts ciphertext encrypted with 3DES ECB cipher.
// The key must be the 3DES key, 24 bytes.
// See more details at [crypto/des] and [crypto/cipher].
//
// Deprecated: ECB mode is not recommended because it does not uses iv.
// Same plaintext is alway be encrypted into the same ciphertext.
// If block cipher is necessary, use CBC mode instead.
func DecryptECB3(key []byte, ciphertext []byte) ([]byte, error) {
	n := len(ciphertext)
	if n < des.BlockSize || n%des.BlockSize != 0 {
		return nil, ErrCipherLength(n)
	}
	c, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	for i := range len(ciphertext) / des.BlockSize {
		start, end := i*des.BlockSize, (i+1)*des.BlockSize
		c.Decrypt(plaintext[start:end], ciphertext[start:end])
	}
	return internal.UnpadPKCS7(des.BlockSize, plaintext)
}
