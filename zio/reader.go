package zio

import (
	"errors"
	"io"
)

// TeeReader returns a [io.Reader] that writes to w what it reads from r.
// All reads from r are written to w. There is no internal buffering.
// Therefor, the write must complete before the read completes.
// If given reader r is nil, it returns nil.
// If given writer w is nil, it returns r itself.
// See also [io.TeeReader].
// The reader works like below.
//   - 1. Read from the reader r.
//   - 2. Write into the writer w. Write is performed even read operation returned an error.
//   - 3. Return written bytes and write error if any. (Write error is prior to the read error.)
//   - 4. Return read bytes and read error.
func TeeReader(r io.Reader, w io.Writer) io.Reader {
	if r == nil || w == nil {
		return r
	}
	return &teeReader{r: r, w: w}
}

type teeReader struct {
	r io.Reader
	w io.Writer
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 { // err may be [io.EOF], so write first.
		if m, er := t.w.Write(p[:n]); er != nil {
			return m, er
		}
	}
	return n, err
}

// TeeReadCloser returns a [io.ReadCloser] that writes to wc what it reads from rc.
// All reads from rc are written to wc. There is no internal buffering.
// Therefor, the write must complete before the read completes.
// If given ReadCloser rc is nil, it returns nil.
// If given WriteCloser wc is nil, it returns rc itself.
// The returned reader works like below.
//   - 1. Read from the rc.
//   - 2. Write into the writer wc. Write is performed even read operation returned an error.
//   - 3. Return written bytes and write error if any. (Write error is prior to the read error.)
//   - 4. Return read bytes and read error.
func TeeReadCloser(rc io.ReadCloser, wc io.WriteCloser) io.ReadCloser {
	if rc == nil || wc == nil {
		return rc
	}
	return &teeReadCloser{rc, wc}
}

type teeReadCloser struct {
	rc io.ReadCloser
	wc io.WriteCloser
}

func (t *teeReadCloser) Close() error {
	var errR, errW error
	errW = t.wc.Close()
	errR = t.rc.Close()
	if errR != nil && errW != nil {
		return errors.Join(errR, errW)
	}
	if errR != nil {
		return errR
	}
	if errW != nil {
		return errW
	}
	return nil
}

func (t *teeReadCloser) Read(p []byte) (n int, err error) {
	n, err = t.rc.Read(p)
	if n > 0 { // err may be [io.EOF], so write first.
		if m, er := t.wc.Write(p[:n]); er != nil {
			return m, er
		}
	}
	return n, err
}
