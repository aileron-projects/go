package ziotest

import "io"

type Charset string

var (
	ASCII = []byte{
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
		52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
		72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91,
		92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
		112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
	}
	CharsetAscii   = string(ASCII)
	CharsetDigit   = "0123456789"
	CharsetLetter1 = "abcdefghijklmnopqrstuvwxyz"
	CharsetLetter2 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetLetter  = CharsetLetter1 + CharsetLetter2
	CharsetSymbol  = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

// CharsetReader returns [io.Reader] that reads from the given charset.
// The reader reads the charset repeatedly if the loop is true.
func CharsetReader(charset string, loop bool) io.Reader {
	return &charsetReader{
		chars: charset,
		loop:  loop,
	}
}

type charsetReader struct {
	chars   string
	loop    bool
	current int
}

func (r *charsetReader) Read(p []byte) (n int, err error) {
	length := len(r.chars)
	written := 0
	for i := range len(p) {
		if r.current >= length {
			if r.loop {
				r.current = 0
			} else {
				return written, io.EOF
			}
		}
		p[i] = r.chars[r.current]
		written += 1
		r.current += 1
	}
	return written, nil
}

// ShortReader reads n bytes at maximum.
// The returned reader does not return any error
// even the read bytes reached to the n.
func ShortReader(r io.Reader, n int64) io.Reader {
	return &errReader{
		reader:   r,
		err:      nil,
		errAfter: n,
	}
}

// ErrReader returns an [io.Reader] that returns [io.ErrClosedPipe] after n bytes read.
// The returned reader returns the error if the inner [io.Reader] returned an error
// before n bytes were read.
// Use [ErrReaderWith] to specify returned error instead of [io.ErrClosedPipe].
// nil reader will cause panic.
func ErrReader(r io.Reader, n int64) io.Reader {
	return &errReader{
		reader:   r,
		err:      io.ErrClosedPipe,
		errAfter: n,
	}
}

// ErrReaderWith works almost the same as [ErrReader]
// except for the point that the returned reader returns
// the given error after n bytes were read.
func ErrReaderWith(r io.Reader, n int64, err error) io.Reader {
	return &errReader{
		reader:   r,
		err:      err,
		errAfter: n,
	}
}

type errReader struct {
	reader   io.Reader
	err      error
	errAfter int64
	read     int64
}

func (r *errReader) Read(p []byte) (n int, err error) {
	if r.read >= r.errAfter {
		return 0, r.err
	}
	left := r.errAfter - r.read
	nn := len(p)
	if len(p) > int(left) {
		nn = int(left)
	}
	n, err = r.reader.Read(p[:nn])
	r.read += int64(n)
	if err != nil {
		return n, err
	}
	if r.read >= r.errAfter {
		return n, r.err
	}
	return n, nil
}
