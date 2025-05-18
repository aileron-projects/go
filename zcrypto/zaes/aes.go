package zaes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"strconv"

	"github.com/aileron-projects/go/internal/helper"
	"github.com/aileron-projects/go/zcrypto/internal"
)

var (
	_ internal.EncryptFunc = EncryptECB
	_ internal.EncryptFunc = EncryptCFB
	_ internal.EncryptFunc = EncryptCTR
	_ internal.EncryptFunc = EncryptOFB
	_ internal.EncryptFunc = EncryptCBC
	_ internal.EncryptFunc = EncryptGCM
	_ internal.DecryptFunc = DecryptECB
	_ internal.DecryptFunc = DecryptCFB
	_ internal.DecryptFunc = DecryptCTR
	_ internal.DecryptFunc = DecryptOFB
	_ internal.DecryptFunc = DecryptCBC
	_ internal.DecryptFunc = DecryptGCM
)

// ErrCipherLength is an error type that tells
// the length of ciphertext is incorrect.
type ErrCipherLength int

func (e ErrCipherLength) Error() string {
	return "zaes: incorrect ciphertext length. got:" + strconv.Itoa(int(e))
}

// newAES returns a new AES cipher block and initial vector.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
func newAES(key []byte) (cipher.Block, []byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	helper.MustNil(err)
	return c, iv, nil
}
