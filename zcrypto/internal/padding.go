package internal

import (
	"crypto/rand"
	"errors"
	"io"
)

var (
	ErrBlockSize   = errors.New("error padding. invalid block size")
	ErrDataLength  = errors.New("error padding. invalid data length")
	ErrPaddingSize = errors.New("error padding. invalid padding size")
)

// PadPKCS7 appends padding bytes to the given data and returns the resulting slice.
// This function can also be used for PKCS#5 padding, as PKCS#5 is a subset of PKCS#7.
// PKCS#7 padding values range from 0x01 to 0x10 (1 to 16 in decimal),
// while PKCS#5 is defined specifically for a block size of 8 bytes (0x08).
// The given blockSize must be in the range 1–255; otherwise, [ErrBlockSize] is returned.
// For example, when blockSize is 6 and the input data is {0x61, 0x61, 0x63},
// the returned data will be {0x61, 0x61, 0x63, 0x03, 0x03, 0x03}.
func PadPKCS7(blockSize int, data []byte) ([]byte, error) {
	if blockSize < 1 || blockSize > 255 { // Restrict block size to 1-255.
		return nil, ErrBlockSize
	}
	n := blockSize - len(data)%blockSize
	pad := make([]byte, n)
	for i := range n {
		pad[i] = byte(n)
	}
	data = append(data, pad...)
	return data, nil
}

// UnpadPKCS7 removes PKCS7 padding from the given data.
// It may return the following errors:
//   - [ErrBlockSize]: the given blockSize is not in the range 1–255.
//   - [ErrDataLength]: the length of the input data is invalid; it must be a multiple of blockSize.
//   - [ErrPaddingSize]: the padding size in the data is invalid.
func UnpadPKCS7(blockSize int, data []byte) ([]byte, error) {
	if blockSize < 1 || blockSize > 255 { // Restrict block size to 1-255.
		return nil, ErrBlockSize
	}
	n := len(data)
	if n < blockSize || n%blockSize != 0 {
		return nil, ErrDataLength
	}
	pad := int(data[n-1])
	if n < pad {
		return nil, ErrPaddingSize
	}
	return data[:n-pad], nil
}

// PadISO7816 appends padding bytes to the given data and returns the resulting slice.
// The given blockSize must be in the range 1–255; otherwise, [ErrBlockSize] is returned.
// For example, when blockSize is 6 and the input data is {0x61, 0x61, 0x63},
// the returned data will be {0x61, 0x61, 0x63, 0x80, 0x00, 0x00}.
// The byte 0x80 marks the end of the actual data, and the remaining bytes (0x00) are padding.
func PadISO7816(blockSize int, data []byte) ([]byte, error) {
	if blockSize < 1 || blockSize > 255 { // Restrict block size to 1-255.
		return nil, ErrBlockSize
	}
	n := blockSize - len(data)%blockSize
	pad := make([]byte, n)
	pad[0] = 0x80
	for i := 1; i < n; i++ {
		pad[i] = 0x00
	}
	data = append(data, pad...)
	return data, nil
}

// UnpadISO7816 removes ISO7816 padding from the given data.
// It may return the following errors:
//   - [ErrBlockSize]: the given blockSize is not in the range 1–255.
//   - [ErrDataLength]: the length of the input data is invalid; it must be a multiple of blockSize.
//   - [ErrPaddingSize]: the padding size in the data is invalid.
func UnpadISO7816(blockSize int, data []byte) ([]byte, error) {
	if blockSize < 1 || blockSize > 255 { // Restrict block size to 1-255.
		return nil, ErrBlockSize
	}
	n := len(data)
	if n < blockSize || n%blockSize != 0 {
		return nil, ErrDataLength
	}
	pad := 0
	for i := 1; i <= blockSize; i++ {
		if data[n-i] == 0x80 {
			pad = i
			break
		} else if data[n-i] != 0x00 {
			break
		}
	}
	if pad == 0 {
		return nil, ErrPaddingSize
	}
	return data[:n-pad], nil
}

// PadISO10126 appends padding bytes to given data and returns resulting slice.
// The given blockSize must be in a range of 1-255, otherwise [ErrBlockSize] is returned.
// For example, when blockSize is 6 and the input data is {0x61, 0x61, 0x63},
// then returned data will be {0x61, 0x61, 0x63, 0xAA, 0xBB, 0x03}.
// The last byte (0x03) represents the length of the padding,
// and the preceding padding bytes (0xAA, 0xBB) are random.
func PadISO10126(blockSize int, data []byte) ([]byte, error) {
	if blockSize < 1 || blockSize > 255 { // Restrict block size to 1-255.
		return nil, ErrBlockSize
	}
	n := blockSize - len(data)%blockSize
	pad := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, pad); err != nil {
		return nil, err
	}
	pad[n-1] = byte(n)
	data = append(data, pad...)
	return data, nil
}

// UnpadISO10126 removes ISO10126 padding from the given data.
// It may return the following errors:
//   - [ErrBlockSize]: the given blockSize is not in the range 1–255.
//   - [ErrDataLength]: the length of the input data is invalid; it must be a multiple of blockSize.
//   - [ErrPaddingSize]: the padding size embedded in the data is invalid.
func UnpadISO10126(blockSize int, data []byte) ([]byte, error) {
	if blockSize < 1 || blockSize > 255 { // Restrict block size to 1-255.
		return nil, ErrBlockSize
	}
	n := len(data)
	if n < blockSize || n%blockSize != 0 {
		return nil, ErrDataLength
	}
	pad := int(data[n-1])
	if n < pad {
		return nil, ErrPaddingSize
	}
	return data[:n-pad], nil
}
