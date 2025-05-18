package zencoding

import (
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// NewDecodeReader returns a [io.Reader] that reads from r and decode it with d.
// See also [golang.org/x/text/transform.NewReader].
func NewDecodeReader(d *encoding.Decoder, r io.Reader) io.Reader {
	return transform.NewReader(r, d)
}

// NewDecodeWriter returns a [io.Writer] that decode written bytes before write to w.
// See also [golang.org/x/text/transform.NewWriter].
func NewDecodeWriter(d *encoding.Decoder, w io.Writer) io.Writer {
	return transform.NewWriter(w, d)
}

// NewEncodeReader returns a [io.Reader] that reads from r and encode it with e.
// See also [golang.org/x/text/transform.NewReader].
func NewEncodeReader(e *encoding.Encoder, r io.Reader) io.Reader {
	return transform.NewReader(r, e)
}

// NewEncodeWriter returns a [io.Writer] that encode written bytes before write to w.
// See also [golang.org/x/text/transform.NewWriter].
func NewEncodeWriter(e *encoding.Encoder, w io.Writer) io.Writer {
	return transform.NewWriter(w, e)
}
