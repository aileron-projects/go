// Package zdes provides functionality for DES and 3DES encryption/decryption.
//
// DES
//   - has 8 bytes block size (See [crypto/des.BlockSize])
//   - uses 8 bytes key
//   - uses 8 bytes initial vector, iv (Except for ECB mode)
//
// 3DES
//   - has 8 bytes block size (See [crypto/des.BlockSize])
//   - uses 24 bytes key
//   - uses 8 bytes initial vector, iv (Except for ECB mode)
//
// See the following table for modes.
//
//	| Mode | Block/Stream | Padding | Use iv | Symmetry? | Status     |
//	| ---- | ------------ | ------- | ------ | --------- | ---------- |
//	| ECB  | Block        | PKCS#7  | No     | Symmetry  | Deprecated |
//	| CBC  | Block        | PKCS#7  | Yes    | Asymmetry |            |
//	| CFB  | Stream       |         | Yes    | Asymmetry | Deprecated |
//	| CTR  | Stream       |         | Yes    | Symmetry  |            |
//	| OFB  | Stream       |         | Yes    | Symmetry  | Deprecated |
//
// If you need to validate DES, you can use openssl:
//
// DES:
//
//	[Encryption]
//
//	echo -n "${PLAINTEXT}" | openssl enc -des-cbc -K "${KEY}" -iv "${INITIAL_VEC}" | xxd -p
//	-K  : Hex encoded 8 bytes key.
//	-iv : Hex encoded 8 bytes initial vector.
//	modes:
//	  -des-ecb
//	  -des-cbc
//	  -des-cfb
//	  -des-ofb
//
//	Example:
//	  export PLAINTEXT="plaintext"
//	  export KEY=$(echo -n "abcd1234" | xxd -p)
//	  export INITIAL_VEC=$(echo -n "12345678" | xxd -p)
//	  echo -n "${PLAINTEXT}" | openssl enc -des-cbc -K "${KEY}" -iv "${INITIAL_VEC}" | xxd -p
//	  >> 1f7c7fc711fb57e74891c07cb0eacc2a
//
//	[Decryption]
//
//	echo -n "${CIPHERTEXT}" | openssl enc -d -des-cbc -K "${KEY}" -iv "${INITIAL_VEC}"
//	-K  : Hex encoded 8 bytes key.
//	-iv : Hex encoded 8 bytes initial vector.
//
//	Example:
//	  export CIPHERTEXT=$(echo -n "1f7c7fc711fb57e74891c07cb0eacc2a" | xxd -p -r)
//	  export KEY=$(echo -n "abcd1234" | xxd -p)
//	  export INITIAL_VEC=$(echo -n "12345678" | xxd -p)
//	  echo -n "${CIPHERTEXT}" | openssl enc -d -des-cbc -K "${KEY}" -iv "${INITIAL_VEC}"
//	  >> plaintext
//
// 3DES:
//
//	[Encryption]
//
//	echo -n "${PLAINTEXT}" | openssl enc -des-ede3-cbc -K "${KEY}" -iv "${INITIAL_VEC}" | xxd -p
//	-K  : Hex encoded 24 bytes key.
//	-iv : Hex encoded 8 bytes initial vector.
//	modes:
//	  -des-ede3-ecb
//	  -des-ede3-cbc
//	  -des-ede3-cfb
//	  -des-ede3-ofb
//
//	Example:
//	  export PLAINTEXT="plaintext"
//	  export KEY=$(echo -n "abcdefghijkl123456789012" | xxd -p)
//	  export INITIAL_VEC=$(echo -n "12345678" | xxd -p)
//	  echo -n "${PLAINTEXT}" | openssl enc -des-ede3-cbc -K "${KEY}" -iv "${INITIAL_VEC}" | xxd -p
//	  >> feea0273aadd2f444d7e4ed28b6c58db
//
//	[Decryption]
//
//	echo -n "${CIPHERTEXT}" | openssl enc -d -des-ede3-cbc -K "${KEY}" -iv "${INITIAL_VEC}"
//	-K  : Hex encoded 24 bytes key.
//	-iv : Hex encoded 8 bytes initial vector.
//
//	Example:
//	  export CIPHERTEXT=$(echo -n "feea0273aadd2f444d7e4ed28b6c58db" | xxd -p -r)
//	  export KEY=$(echo -n "abcdefghijkl123456789012" | xxd -p)
//	  export INITIAL_VEC=$(echo -n "12345678" | xxd -p)
//	  echo -n "${CIPHERTEXT}" | openssl enc -d -des-ede3-cbc -K "${KEY}" -iv "${INITIAL_VEC}"
//	  >> plaintext
//
// Some online tools available.
//
// DES online tools:
//
//   - https://anycript.com/crypto/des
//   - https://emn178.github.io/online-tools/des/encrypt/
//   - https://www.toolhelper.cn/en/SymmetricEncryption/DES
//
// 3DES online tools:
//
//   - https://anycript.com/crypto/tripledes
//   - https://emn178.github.io/online-tools/triple-des/encrypt/
//   - https://www.toolhelper.cn/en/SymmetricEncryption/TripleDES
package zdes
