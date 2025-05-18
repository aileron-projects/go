package zdes

import (
	"crypto/cipher"
	"crypto/des"
)

// EncryptCFB encrypts plaintext with DES CFB cipher.
// The key must be the DES key, 8 bytes.
// 8 bytes iv read from [crypto/rand.Reader] is used.
// See more details at [crypto/des] and [crypto/cipher].
//
// Returned slice consists of:
//   - Initial 8 bytes ([crypto/des.BlockSize]): Initial vector generated using [crypto/rand.Read].
//   - Rest of the bytes: DES CFB encrypted plaintext.
//
// Deprecated: See the comment on [crypto/cipher.NewCFBEncrypter].
func EncryptCFB(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := newDES(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(plaintext))
	cipher.NewCFBEncrypter(c, iv).XORKeyStream(ciphertext, plaintext)
	return append(iv, ciphertext...), nil
}

// DecryptCFB decrypts ciphertext encrypted with DES CFC cipher.
// The key must be the DES key, 8 bytes.
// See more details at [crypto/des] and [crypto/cipher].
//
// Deprecated: See the comment on [crypto/cipher.NewCFBDecrypter].
func DecryptCFB(key []byte, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < des.BlockSize {
		return nil, ErrCipherLength(len(ciphertext))
	}
	iv, ciphertext := ciphertext[:des.BlockSize], ciphertext[des.BlockSize:]
	c, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCFBDecrypter(c, iv).XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

// EncryptCFB3 encrypts plaintext with 3DES CFB cipher.
// The key must be the 3DES key, 24 bytes.
// 8 bytes iv read from [crypto/rand.Reader] is used.
// See more details at [crypto/des] and [crypto/cipher].
//
// Returned slice consists of:
//   - Initial 8 bytes ([crypto/des.BlockSize]): Initial vector generated using [crypto/rand.Read].
//   - Rest of the bytes: 3DES CFB encrypted plaintext.
//
// Deprecated: See the comment on [crypto/cipher.NewCFBEncrypter].
func EncryptCFB3(key []byte, plaintext []byte) ([]byte, error) {
	c, iv, err := new3DES(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(plaintext))
	cipher.NewCFBEncrypter(c, iv).XORKeyStream(ciphertext, plaintext)
	return append(iv, ciphertext...), nil
}

// DecryptCFB3 decrypts ciphertext encrypted with 3DES CFC cipher.
// The key must be the 3DES key, 24 bytes.
// See more details at [crypto/des] and [crypto/cipher].
//
// Deprecated: See the comment on [crypto/cipher.NewCFBDecrypter].
func DecryptCFB3(key []byte, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < des.BlockSize {
		return nil, ErrCipherLength(len(ciphertext))
	}
	iv, ciphertext := ciphertext[:des.BlockSize], ciphertext[des.BlockSize:]
	c, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCFBDecrypter(c, iv).XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
