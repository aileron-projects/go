// Package zaes provides functionality for AES encryption/decryption.
//
// AES
//   - has 16 bytes block size (See [crypto/aes.BlockSize])
//   - uses 16, 24, or 32 bytes for AES-128, AES-192, or AES-256 each
//
// See the following table for modes.
// Currently padding method is not changeable.
//
//	| Mode | Block/Stream | Padding | Use iv/nonce | Symmetry? | Status     |
//	| ---- | ------------ | ------- | ------------ | --------- | ---------- |
//	| ECB  | Block        | PKCS#7  | No           | Symmetry  | Deprecated |
//	| CBC  | Block        | PKCS#7  | Yes          | Asymmetry |            |
//	| CFB  | Stream       |         | Yes          | Asymmetry | Deprecated |
//	| CTR  | Stream       |         | Yes          | Symmetry  |            |
//	| OFB  | Stream       |         | Yes          | Symmetry  | Deprecated |
//	| GCM  | Stream       |         | Yes          | Symmetry  |            |
//
// If you need to validate AES, you can use openssl:
//
//	[Encryption]
//
//	echo -n "${PLAINTEXT}" | openssl enc -aes-128-cbc -K "${KEY}" -iv "${INITIAL_VEC}" | xxd -p
//	-K  : Hex encoded 16, 24 or 32 bytes key.
//	-iv : Hex encoded 8 bytes initial vector.
//	modes:
//	  -aes-128-cbc
//	  -aes-128-cfb
//	  -aes-128-ctr
//	  -aes-128-ecb
//	  -aes-128-ofb
//	  -aes-128-gcm
//	  -aes-196-cbc
//	  -aes-196-cfb
//	  -aes-196-ctr
//	  -aes-196-ecb
//	  -aes-196-ofb
//	  -aes-196-gcm
//	  -aes-256-cbc
//	  -aes-256-cfb
//	  -aes-256-ctr
//	  -aes-256-ecb
//	  -aes-256-ofb
//	  -aes-256-gcm
//
//	Example:
//	  export PLAINTEXT="plaintext"
//	  export KEY=$(echo -n "16bytessecretkey" | xxd -p)
//	  export INITIAL_VEC=$(echo -n "1234567890123456" | xxd -p)
//	  echo -n "${PLAINTEXT}" | openssl enc -aes-128-cbc -K "${KEY}" -iv "${INITIAL_VEC}" | xxd -p
//	  >> 053da3b8836a0a64105bf9815db5f86b
//
//	[Decryption]
//
//	echo -n "${CIPHERTEXT}" | openssl enc -d -aes-128-cbc -K "${KEY}" -iv "${INITIAL_VEC}"
//	-K  : Hex encoded 16, 24 or 32 bytes key.
//	-iv : Hex encoded 8 bytes initial vector.
//
//	Example:
//	  export CIPHERTEXT=$(echo -n "053da3b8836a0a64105bf9815db5f86b" | xxd -p -r)
//	  export KEY=$(echo -n "16bytessecretkey" | xxd -p)
//	  export INITIAL_VEC=$(echo -n "1234567890123456" | xxd -p)
//	  echo -n "${CIPHERTEXT}" | openssl enc -d -aes-128-cbc -K "${KEY}" -iv "${INITIAL_VEC}"
//	  >> plaintext
//
// Also, some online tools available.
//
//   - https://emn178.github.io/online-tools/aes/encrypt/
//   - https://www.toolhelper.cn/en/SymmetricEncryption/AES
//   - https://tophix.com/development-tools/encrypt-text
//   - https://www.javainuse.com/aesgenerator
//   - https://anycript.com/crypto
package zaes
