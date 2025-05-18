package zio

import (
	"io"
	"sync"
)

// pool is the 4kiB buffer pool.
//
// Usage:
//
//	buf := *pool.Get().(*[]byte)
//	defer pool.Put(&buf)
var pool = sync.Pool{
	New: func() any {
		buf := make([]byte, 1<<12) // 4kiB
		return &buf
	},
}

// Copy copies from src to dst until either EOF is reached
// on src or an error occurs. It returns the number of bytes
// copied and the first error encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == EOF.
// Because Copy is defined to read from src until EOF, it does
// not treat an EOF from Read as an error to be reported.
//
// It panics if nil writer or nil reader is given.
//
// See also [io.Copy] and [io.CopyBuffer].
// [Copy] uses buffer pool internally for performance improvement.
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := *pool.Get().(*[]byte)
	defer pool.Put(&buf)
	for {
		nRead, readErr := src.Read(buf)
		if nRead > 0 {
			nWrite, writeErr := dst.Write(buf[:nRead])
			written += int64(nWrite)
			if writeErr != nil {
				return written, writeErr
			}
			if nRead != nWrite {
				return written, io.ErrShortWrite
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				return written, nil
			}
			return written, readErr
		}
	}
}
