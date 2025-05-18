package zio

import "io"

// NopReadCloser wraps given Reader and returns ReadCloser
// that do nothing when [io.ReadCloser].Close is called.
func NopReadCloser(r io.Reader) io.ReadCloser {
	return &nopReadCloser{Reader: r}
}

// NopWriteCloser wraps given Writer and returns WriteCloser
// that do nothing when [io.WriteCloser].Close is called.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{Writer: w}
}

type nopReadCloser struct {
	io.Reader
}

func (c *nopReadCloser) Close() error {
	return nil
}

type nopWriteCloser struct {
	io.Writer
}

func (c *nopWriteCloser) Close() error {
	return nil
}
