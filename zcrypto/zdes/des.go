package zdes

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"strconv"

	"github.com/aileron-projects/go/internal/helper"
	"github.com/aileron-projects/go/zcrypto/internal"
)

var (
	// DES
	_ internal.EncryptFunc = EncryptECB
	_ internal.EncryptFunc = EncryptCFB
	_ internal.EncryptFunc = EncryptCTR
	_ internal.EncryptFunc = EncryptOFB
	_ internal.EncryptFunc = EncryptCBC
	_ internal.DecryptFunc = DecryptECB
	_ internal.DecryptFunc = DecryptCFB
	_ internal.DecryptFunc = DecryptCTR
	_ internal.DecryptFunc = DecryptOFB
	_ internal.DecryptFunc = DecryptCBC
	// 3DES
	_ internal.EncryptFunc = EncryptECB3
	_ internal.EncryptFunc = EncryptCFB3
	_ internal.EncryptFunc = EncryptCTR3
	_ internal.EncryptFunc = EncryptOFB3
	_ internal.EncryptFunc = EncryptCBC3
	_ internal.DecryptFunc = DecryptECB3
	_ internal.DecryptFunc = DecryptCFB3
	_ internal.DecryptFunc = DecryptCTR3
	_ internal.DecryptFunc = DecryptOFB3
	_ internal.DecryptFunc = DecryptCBC3
)

// ErrCipherLength is an error type that tells
// the length of ciphertext is incorrect.
type ErrCipherLength int

func (e ErrCipherLength) Error() string {
	return "zdes: incorrect ciphertext length. got:" + strconv.Itoa(int(e))
}

// newDES returns a new instance of cipher block for DES.
// Returned iv, or initial vector, has the same length
// with the cipher's block size [crypto/des.BlockSize].
// Entropy source for iv is [crypto/rand.Reader].
// Use [new3DES] for 3DES.
func newDES(key []byte) (cipher.Block, []byte, error) {
	c, err := des.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	iv := make([]byte, des.BlockSize)
	_, err = rand.Read(iv)
	helper.MustNil(err)
	return c, iv, nil
}

// new3DES returns a new instance of cipher block for 3DES.
// Returned iv, or initial vector, has the same length
// with the cipher's block size [crypto/des.BlockSize].
// Entropy source for iv is [crypto/rand.Reader].
// Use [newDES] for DES.
func new3DES(key []byte) (cipher.Block, []byte, error) {
	c, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, nil, err
	}
	iv := make([]byte, des.BlockSize)
	_, err = rand.Read(iv)
	helper.MustNil(err)
	return c, iv, nil
}
