package zio

import (
	"errors"
	"io"
)

var (
	// ErrReadLimit is the error that indicates read limit is reached.
	// This error is returns from [LimitReader].
	ErrReadLimit = errors.New("zio: read limit reached")
	// ErrWriteLimit is the error that indicates write limit is reached.
	// This error is returns from [LimitWriter].
	ErrWriteLimit = errors.New("zio: write limit reached")
)

// LimitReader returns a [io.Reader] that reads from r but stops
// when the limit is reached.
// Unlike [io.LimitedReader], the returned reader returns
// [ErrReadLimit] when the limit is reached.
// It returns nil of given r is nil.
func LimitReader(r io.Reader, limit int64) io.Reader {
	if r == nil {
		return nil
	}
	return &limitReader{
		Reader: r,
		limit:  limit,
	}
}

// LimitWriter returns a [io.Writer] that writes to w but stops
// when the limit is reached.
// Returned writer returns [ErrWriteLimit] when the limit is reached.
// It returns nil of given w is nil.
func LimitWriter(w io.Writer, limit int64) io.Writer {
	if w == nil {
		return nil
	}
	return &limitWriter{
		Writer: w,
		limit:  limit,
	}
}

type limitReader struct {
	io.Reader
	limit int64
}

func (l *limitReader) Read(p []byte) (n int, err error) {
	if l.limit <= 0 {
		return 0, ErrReadLimit
	}
	limited := int64(len(p)) > l.limit
	if limited {
		p = p[:l.limit]
	}
	n, err = l.Reader.Read(p)
	l.limit -= int64(n)
	if err == nil && limited {
		err = ErrReadLimit
	}
	return n, err
}

type limitWriter struct {
	io.Writer
	limit int64
}

func (l *limitWriter) Write(p []byte) (n int, err error) {
	if l.limit <= 0 {
		return 0, ErrWriteLimit
	}
	limited := int64(len(p)) > l.limit
	if limited {
		p = p[:l.limit]
	}
	n, err = l.Writer.Write(p)
	l.limit -= int64(n)
	if err == nil && limited {
		err = ErrWriteLimit
	}
	return n, err
}
