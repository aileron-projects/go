package ziotest

import "io"

// ShortWriter accepts n bytes at maximum.
// The returned writer does not return any error
// even the written bytes reached to the n.
func ShortWriter(w io.Writer, n int64) io.Writer {
	return &errWriter{
		writer:   w,
		err:      nil,
		errAfter: n,
	}
}

// ErrWriter returns an [io.Writer] that returns [io.ErrClosedPipe] after n bytes written.
// The returned writer returns the error if the inner [io.Writer] returned an error
// before n bytes were written.
// It ignores w if it is nil.
// Use [ErrWriterWith] to specify returned error instead of [io.ErrClosedPipe].
func ErrWriter(w io.Writer, n int64) io.Writer {
	return &errWriter{
		writer:   w,
		err:      io.ErrClosedPipe,
		errAfter: n,
	}
}

// ErrWriterWith works almost the same as [ErrWriter]
// except for the point that the returned writer returns
// the given error after n bytes were written.
func ErrWriterWith(w io.Writer, n int64, err error) io.Writer {
	return &errWriter{
		writer:   w,
		err:      err,
		errAfter: n,
	}
}

type errWriter struct {
	writer   io.Writer
	err      error
	errAfter int64
	written  int64
}

func (w *errWriter) Write(p []byte) (n int, err error) {
	if w.written >= w.errAfter {
		return 0, w.err
	}
	left := w.errAfter - w.written
	write := min(len(p), int(left))
	if w.writer != nil {
		n, err := w.writer.Write(p[:write])
		if err != nil {
			return n, err
		}
	}
	w.written += int64(write)
	if w.written >= w.errAfter {
		return write, w.err
	}
	return write, nil
}
