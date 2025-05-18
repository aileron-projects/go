package internal

import (
	"crypto/cipher"
	"io"
)

// Copy copies from src to dst.
// If src is encrypted, then decrypted bytes are written into the dst.
// If src is not encrypted, then encrypted bytes are written into the dst.
// It uses 1 kiB internal buffer.
// s, dst, src must not be nil, otherwise it panics.
func Copy(s cipher.Stream, dst io.Writer, src io.Reader) error {
	buf := make([]byte, 1<<10)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			s.XORKeyStream(buf[:n], buf[:n])
			m, err := dst.Write(buf[:n])
			if err != nil {
				return err
			}
			if n != m {
				return io.ErrShortWrite
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
