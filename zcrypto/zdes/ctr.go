package zdes

import (
	"crypto/cipher"
	"crypto/des"
	"io"

	"github.com/aileron-projects/go/zcrypto/internal"
)

// EncryptCTR encrypts plaintext with DES CTR cipher.
// The key must be the DES key, 8 bytes.
// 8 bytes iv read from [crypto/rand.Reader] is used.
// See more details at [crypto/des] and [crypto/cipher].
//
// Returned slice consists of:
//   - Initial 8 bytes ([crypto/des.BlockSize]): Initial vector generated using [crypto/rand.Read].
//   - Rest of the bytes: DES CTR encrypted plaintext.
func EncryptCTR(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := newDES(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(plaintext))
	cipher.NewCTR(c, iv).XORKeyStream(ciphertext, plaintext)
	return append(iv, ciphertext...), nil
}

// DecryptCTR decrypts ciphertext encrypted with DES CTR cipher.
// The key must be the DES key, 8 bytes.
// See more details at [crypto/des] and [crypto/cipher].
func DecryptCTR(key []byte, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < des.BlockSize {
		return nil, ErrCipherLength(len(ciphertext))
	}
	iv, ciphertext := ciphertext[:des.BlockSize], ciphertext[des.BlockSize:]
	c, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCTR(c, iv).XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

// EncryptCTR3 encrypts plaintext with 3DES CTR cipher.
// The key must be the 3DES key, 24 bytes.
// 8 bytes iv read from [crypto/rand.Reader] is used.
// See more details at [crypto/des] and [crypto/cipher].
//
// Returned slice consists of:
//   - Initial 8 bytes ([crypto/des.BlockSize]): Initial vector generated using [crypto/rand.Read].
//   - Rest of the bytes: 3DES CTR encrypted plaintext.
func EncryptCTR3(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := new3DES(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(plaintext))
	cipher.NewCTR(c, iv).XORKeyStream(ciphertext, plaintext)
	return append(iv, ciphertext...), nil
}

// DecryptCTR3 decrypts ciphertext encrypted with 3DES CTR cipher.
// The key must be the 3DES key, 24 bytes.
// See more details at [crypto/des] and [crypto/cipher].
func DecryptCTR3(key []byte, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < des.BlockSize {
		return nil, ErrCipherLength(len(ciphertext))
	}
	iv, ciphertext := ciphertext[:des.BlockSize], ciphertext[des.BlockSize:]
	c, err := des.NewTripleDESCipher(key)
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
// Note that the OFB mode ([crypto/cipher.NewOFB]) is deprecated.
func CopyCTR(key, iv []byte, dst io.Writer, src io.Reader) error {
	c, err := des.NewCipher(key)
	if err != nil {
		return err
	}
	return internal.Copy(cipher.NewCTR(c, iv), dst, src)
}

// CopyCTR3 copies from src to dst.
// If src is encrypted, then decrypted bytes are written into the dst.
// If src is not encrypted, then encrypted bytes are written into the dst.
// Note that the OFB mode ([crypto/cipher.NewOFB]) is deprecated.
func CopyCTR3(key, iv []byte, dst io.Writer, src io.Reader) error {
	c, err := des.NewTripleDESCipher(key)
	if err != nil {
		return err
	}
	return internal.Copy(cipher.NewCTR(c, iv), dst, src)
}
