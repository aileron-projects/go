package internal

// EncryptFunc is the type of function that encrypts the given plaintext.
type EncryptFunc func(key []byte, plaintext []byte) (ciphertext []byte, err error)

// DecryptFunc is the type of function that decrypts the given ciphertext.
type DecryptFunc func(key []byte, ciphertext []byte) (plaintext []byte, err error)
