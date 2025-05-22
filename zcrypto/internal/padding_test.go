package internal_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/internal"
	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/ztesting/ziotest"
)

func TestPadPKCS7(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		blockSize int
		data      []byte
		want      []byte
		err       error
	}{
		"size=0": {
			blockSize: 0,
			data:      []byte{},
			err:       internal.ErrBlockSize,
		},
		"size=256": {
			blockSize: 256,
			data:      []byte{},
			err:       internal.ErrBlockSize,
		},
		"size=1": {
			blockSize: 1,
			data:      []byte{},
			want:      bytes.Repeat([]byte{'\x01'}, 1),
		},
		"size=1 data=1": {
			blockSize: 1,
			data:      bytes.Repeat([]byte{'\x00'}, 1),
			want:      append(bytes.Repeat([]byte{'\x00'}, 1), bytes.Repeat([]byte{'\x01'}, 1)...),
		},
		"size=1 data=100": {
			blockSize: 1,
			data:      bytes.Repeat([]byte{'\x00'}, 100),
			want:      append(bytes.Repeat([]byte{'\x00'}, 100), bytes.Repeat([]byte{'\x01'}, 1)...),
		},
		"size=100": {
			blockSize: 100,
			data:      []byte{},
			want:      bytes.Repeat([]byte{'\x64'}, 100),
		},
		"size=100 data=1": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x00'}, 1),
			want:      append(bytes.Repeat([]byte{'\x00'}, 1), bytes.Repeat([]byte{'\x63'}, 99)...),
		},
		"size=100 data=50": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x00'}, 50),
			want:      append(bytes.Repeat([]byte{'\x00'}, 50), bytes.Repeat([]byte{'\x32'}, 50)...),
		},
		"size=100 data=100": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x00'}, 100),
			want:      append(bytes.Repeat([]byte{'\x00'}, 100), bytes.Repeat([]byte{'\x64'}, 100)...),
		},
		"size=100 data=150": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x00'}, 150),
			want:      append(bytes.Repeat([]byte{'\x00'}, 150), bytes.Repeat([]byte{'\x32'}, 50)...),
		},
		"size=255": {
			blockSize: 255,
			data:      []byte{},
			want:      bytes.Repeat([]byte{'\xff'}, 255),
		},
		"size=255 data=1": {
			blockSize: 255,
			data:      bytes.Repeat([]byte{'\x00'}, 1),
			want:      append(bytes.Repeat([]byte{'\x00'}, 1), bytes.Repeat([]byte{'\xfe'}, 254)...),
		},
		"size=255 data=100": {
			blockSize: 255,
			data:      bytes.Repeat([]byte{'\x00'}, 100),
			want:      append(bytes.Repeat([]byte{'\x00'}, 100), bytes.Repeat([]byte{'\x9b'}, 155)...),
		},
		"size=255 data=260": {
			blockSize: 255,
			data:      bytes.Repeat([]byte{'\x00'}, 260),
			want:      append(bytes.Repeat([]byte{'\x00'}, 260), bytes.Repeat([]byte{'\xfa'}, 250)...),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			b, err := internal.PadPKCS7(tc.blockSize, tc.data)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
			ztesting.AssertEqual(t, "data not match", tc.want, b)
			if len(b) > 0 {
				bb, err := internal.UnpadPKCS7(tc.blockSize, b)
				ztesting.AssertEqualErr(t, "non nil error", nil, err)
				ztesting.AssertEqual(t, "data not match", tc.data, bb)
			}
		})
	}
}

func TestUnpadPKCS7(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		blockSize int
		data      []byte
		want      []byte
		err       error
	}{
		"size=0":                {blockSize: 0, err: internal.ErrBlockSize},
		"size=256":              {blockSize: 256, err: internal.ErrBlockSize},
		"size=50 invalid data":  {blockSize: 100, data: bytes.Repeat([]byte{'\xff'}, 100), err: internal.ErrPaddingSize},
		"size=100 invalid data": {blockSize: 100, data: bytes.Repeat([]byte{'\x00'}, 99), err: internal.ErrDataLength},
		"size=255 invalid data": {blockSize: 100, data: bytes.Repeat([]byte{'\x00'}, 256), err: internal.ErrDataLength},
		"size=1": {
			blockSize: 1,
			data:      bytes.Repeat([]byte{'\x01'}, 1),
			want:      []byte{},
		},
		"size=1 data=1": {
			blockSize: 1,
			data:      append(bytes.Repeat([]byte{'\x00'}, 1), bytes.Repeat([]byte{'\x01'}, 1)...),
			want:      bytes.Repeat([]byte{'\x00'}, 1),
		},
		"size=1 data=100": {
			blockSize: 1,
			data:      append(bytes.Repeat([]byte{'\x00'}, 100), bytes.Repeat([]byte{'\x01'}, 1)...),
			want:      bytes.Repeat([]byte{'\x00'}, 100),
		},
		"size=100": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x64'}, 100),
			want:      []byte{},
		},
		"size=100 data=1": {
			blockSize: 100,
			data:      append(bytes.Repeat([]byte{'\x00'}, 1), bytes.Repeat([]byte{'\x63'}, 99)...),
			want:      bytes.Repeat([]byte{'\x00'}, 1),
		},
		"size=100 data=50": {
			blockSize: 100,
			data:      append(bytes.Repeat([]byte{'\x00'}, 50), bytes.Repeat([]byte{'\x32'}, 50)...),
			want:      bytes.Repeat([]byte{'\x00'}, 50),
		},
		"size=100 data=100": {
			blockSize: 100,
			data:      append(bytes.Repeat([]byte{'\x00'}, 100), bytes.Repeat([]byte{'\x64'}, 100)...),
			want:      bytes.Repeat([]byte{'\x00'}, 100),
		},
		"size=100 data=150": {
			blockSize: 100,
			data:      append(bytes.Repeat([]byte{'\x00'}, 150), bytes.Repeat([]byte{'\x32'}, 50)...),
			want:      bytes.Repeat([]byte{'\x00'}, 150),
		},
		"size=255": {
			blockSize: 255,
			data:      bytes.Repeat([]byte{'\xff'}, 255),
			want:      []byte{},
		},
		"size=255 data=1": {
			blockSize: 255,
			data:      append(bytes.Repeat([]byte{'\x00'}, 1), bytes.Repeat([]byte{'\xfe'}, 254)...),
			want:      bytes.Repeat([]byte{'\x00'}, 1),
		},
		"size=255 data=100": {
			blockSize: 255,
			data:      append(bytes.Repeat([]byte{'\x00'}, 100), bytes.Repeat([]byte{'\x9b'}, 155)...),
			want:      bytes.Repeat([]byte{'\x00'}, 100),
		},
		"size=255 data=260": {
			blockSize: 255,
			data:      append(bytes.Repeat([]byte{'\x00'}, 260), bytes.Repeat([]byte{'\xfa'}, 250)...),
			want:      bytes.Repeat([]byte{'\x00'}, 260),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			b, err := internal.UnpadPKCS7(tc.blockSize, tc.data)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
			ztesting.AssertEqual(t, "data not match", tc.want, b)
			if tc.err == nil {
				bb, err := internal.PadPKCS7(tc.blockSize, b)
				ztesting.AssertEqualErr(t, "non nil error", nil, err)
				ztesting.AssertEqual(t, "data not match", tc.data, bb)
			}
		})
	}
}

func TestPadISO7816(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		blockSize int
		data      []byte
		want      []byte
		err       error
	}{
		"size=0":   {blockSize: 0, err: internal.ErrBlockSize},
		"size=256": {blockSize: 256, err: internal.ErrBlockSize},
		"size=1": {
			blockSize: 1,
			data:      []byte{},
			want:      bytes.Repeat([]byte{'\x80'}, 1),
		},
		"size=1 data=1": {
			blockSize: 1,
			data:      bytes.Repeat([]byte{'\x01'}, 1),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 1), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 0)...),
		},
		"size=1 data=100": {
			blockSize: 1,
			data:      bytes.Repeat([]byte{'\x01'}, 100),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 100), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 0)...),
		},
		"size=100": {
			blockSize: 100,
			data:      []byte{},
			want:      append([]byte{'\x80'}, bytes.Repeat([]byte{'\x00'}, 99)...),
		},
		"size=100 data=1": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x01'}, 1),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 1), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 98)...),
		},
		"size=100 data=50": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x01'}, 50),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 50), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 49)...),
		},
		"size=100 data=100": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x01'}, 100),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 100), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 99)...),
		},
		"size=100 data=150": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x01'}, 150),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 150), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 49)...),
		},
		"size=255": {
			blockSize: 255,
			data:      []byte{},
			want:      append([]byte{'\x80'}, bytes.Repeat([]byte{'\x00'}, 254)...),
		},
		"size=255 data=1": {
			blockSize: 255,
			data:      bytes.Repeat([]byte{'\x01'}, 1),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 1), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 253)...),
		},
		"size=255 data=100": {
			blockSize: 255,
			data:      bytes.Repeat([]byte{'\x01'}, 100),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 100), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 154)...),
		},
		"size=255 data=260": {
			blockSize: 255,
			data:      bytes.Repeat([]byte{'\x01'}, 260),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 260), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 249)...),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			b, err := internal.PadISO7816(tc.blockSize, tc.data)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
			ztesting.AssertEqual(t, "data not match", tc.want, b)
			if len(b) > 0 {
				bb, err := internal.UnpadISO7816(tc.blockSize, b)
				ztesting.AssertEqualErr(t, "non nil error", nil, err)
				ztesting.AssertEqual(t, "data not match", tc.data, bb)
			}
		})
	}
}

func TestUnpadISO7816(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		blockSize int
		data      []byte
		want      []byte
		err       error
	}{
		"size=0":                {blockSize: 0, err: internal.ErrBlockSize},
		"size=256":              {blockSize: 256, err: internal.ErrBlockSize},
		"size=50 invalid data":  {blockSize: 100, data: bytes.Repeat([]byte{'\xff'}, 100), err: internal.ErrPaddingSize},
		"size=100 invalid data": {blockSize: 100, data: bytes.Repeat([]byte{'\x00'}, 99), err: internal.ErrDataLength},
		"size=255 invalid data": {blockSize: 100, data: bytes.Repeat([]byte{'\x00'}, 256), err: internal.ErrDataLength},
		"size=1": {
			blockSize: 1,
			data:      bytes.Repeat([]byte{'\x80'}, 1),
			want:      []byte{},
		},
		"size=1 data=1": {
			blockSize: 1,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 1), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 0)...),
			want:      bytes.Repeat([]byte{'\x01'}, 1),
		},
		"size=1 data=100": {
			blockSize: 1,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 100), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 0)...),
			want:      bytes.Repeat([]byte{'\x01'}, 100),
		},
		"size=100": {
			blockSize: 100,
			data:      append([]byte{'\x80'}, bytes.Repeat([]byte{'\x00'}, 99)...),
			want:      []byte{},
		},
		"size=100 data=1": {
			blockSize: 100,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 1), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 98)...),
			want:      bytes.Repeat([]byte{'\x01'}, 1),
		},
		"size=100 data=50": {
			blockSize: 100,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 50), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 49)...),
			want:      bytes.Repeat([]byte{'\x01'}, 50),
		},
		"size=100 data=100": {
			blockSize: 100,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 100), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 99)...),
			want:      bytes.Repeat([]byte{'\x01'}, 100),
		},
		"size=100 data=150": {
			blockSize: 100,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 150), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 49)...),
			want:      bytes.Repeat([]byte{'\x01'}, 150),
		},
		"size=255": {
			blockSize: 255,
			data:      append([]byte{'\x80'}, bytes.Repeat([]byte{'\x00'}, 254)...),
			want:      []byte{},
		},
		"size=255 data=1": {
			blockSize: 255,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 1), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 253)...),
			want:      bytes.Repeat([]byte{'\x01'}, 1),
		},
		"size=255 data=100": {
			blockSize: 255,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 100), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 154)...),
			want:      bytes.Repeat([]byte{'\x01'}, 100),
		},
		"size=255 data=260": {
			blockSize: 255,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 260), []byte{'\x80'}...), bytes.Repeat([]byte{'\x00'}, 249)...),
			want:      bytes.Repeat([]byte{'\x01'}, 260),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			b, err := internal.UnpadISO7816(tc.blockSize, tc.data)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
			ztesting.AssertEqual(t, "data not match", tc.want, b)
			if tc.err == nil {
				bb, err := internal.PadISO7816(tc.blockSize, b)
				ztesting.AssertEqualErr(t, "non nil error", nil, err)
				ztesting.AssertEqual(t, "data not match", tc.data, bb)
			}
		})
	}
}

func TestPadISO10126(t *testing.T) {
	done := ztesting.ReplaceRandReader(ziotest.CharsetReader(string([]byte{'\x00'}), true))
	defer done()
	testCases := map[string]struct {
		blockSize int
		data      []byte
		want      []byte
		err       error
	}{
		"size=0":   {blockSize: 0, err: internal.ErrBlockSize},
		"size=256": {blockSize: 256, err: internal.ErrBlockSize},
		"size=1": {
			blockSize: 1,
			data:      []byte{},
			want:      bytes.Repeat([]byte{'\x01'}, 1),
		},
		"size=1 data=1": {
			blockSize: 1,
			data:      bytes.Repeat([]byte{'\x01'}, 1),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 1), bytes.Repeat([]byte{'\x00'}, 0)...), []byte{'\x01'}...),
		},
		"size=1 data=100": {
			blockSize: 1,
			data:      bytes.Repeat([]byte{'\x01'}, 100),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 100), bytes.Repeat([]byte{'\x00'}, 0)...), []byte{'\x01'}...),
		},
		"size=100": {
			blockSize: 100,
			data:      []byte{},
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 0), bytes.Repeat([]byte{'\x00'}, 99)...), []byte{'\x64'}...),
		},
		"size=100 data=1": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x01'}, 1),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 1), bytes.Repeat([]byte{'\x00'}, 98)...), []byte{'\x63'}...),
		},
		"size=100 data=50": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x01'}, 50),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 50), bytes.Repeat([]byte{'\x00'}, 49)...), []byte{'\x32'}...),
		},
		"size=100 data=100": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x01'}, 100),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 100), bytes.Repeat([]byte{'\x00'}, 99)...), []byte{'\x64'}...),
		},
		"size=100 data=150": {
			blockSize: 100,
			data:      bytes.Repeat([]byte{'\x01'}, 150),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 150), bytes.Repeat([]byte{'\x00'}, 49)...), []byte{'\x32'}...),
		},
		"size=255": {
			blockSize: 255,
			data:      []byte{},
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 0), bytes.Repeat([]byte{'\x00'}, 254)...), []byte{'\xff'}...),
		},
		"size=255 data=1": {
			blockSize: 255,
			data:      bytes.Repeat([]byte{'\x01'}, 1),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 1), bytes.Repeat([]byte{'\x00'}, 253)...), []byte{'\xfe'}...),
		},
		"size=255 data=100": {
			blockSize: 255,
			data:      bytes.Repeat([]byte{'\x01'}, 100),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 100), bytes.Repeat([]byte{'\x00'}, 154)...), []byte{'\x9b'}...),
		},
		"size=255 data=260": {
			blockSize: 255,
			data:      bytes.Repeat([]byte{'\x01'}, 260),
			want:      append(append(bytes.Repeat([]byte{'\x01'}, 260), bytes.Repeat([]byte{'\x00'}, 249)...), []byte{'\xfa'}...),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			b, err := internal.PadISO10126(tc.blockSize, tc.data)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
			ztesting.AssertEqual(t, "data not match", tc.want, b)
			if len(b) > 0 {
				bb, err := internal.UnpadISO10126(tc.blockSize, b)
				ztesting.AssertEqualErr(t, "non nil error", nil, err)
				ztesting.AssertEqual(t, "data not match", tc.data, bb)
			}
		})
	}
}

func TestPadISO10126_ReadError(t *testing.T) {
	done := ztesting.ReplaceRandReader(strings.NewReader("1"))
	defer done()
	b, err := internal.PadISO10126(10, []byte("12345"))
	ztesting.AssertEqualErr(t, "error not match", io.ErrUnexpectedEOF, err)
	ztesting.AssertEqual(t, "data not match", nil, b)
}

func TestUnpadISO10126(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		blockSize int
		data      []byte
		want      []byte
		err       error
	}{
		"size=0":                {blockSize: 0, err: internal.ErrBlockSize},
		"size=256":              {blockSize: 256, err: internal.ErrBlockSize},
		"size=50 invalid data":  {blockSize: 100, data: bytes.Repeat([]byte{'\xff'}, 100), err: internal.ErrPaddingSize},
		"size=100 invalid data": {blockSize: 100, data: bytes.Repeat([]byte{'\x00'}, 99), err: internal.ErrDataLength},
		"size=255 invalid data": {blockSize: 100, data: bytes.Repeat([]byte{'\x00'}, 256), err: internal.ErrDataLength},
		"size=1": {
			blockSize: 1,
			data:      bytes.Repeat([]byte{'\x01'}, 1),
			want:      []byte{},
		},
		"size=1 data=1": {
			blockSize: 1,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 1), bytes.Repeat([]byte{'\x00'}, 0)...), []byte{'\x01'}...),
			want:      bytes.Repeat([]byte{'\x01'}, 1),
		},
		"size=1 data=100": {
			blockSize: 1,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 100), bytes.Repeat([]byte{'\x00'}, 0)...), []byte{'\x01'}...),
			want:      bytes.Repeat([]byte{'\x01'}, 100),
		},
		"size=100": {
			blockSize: 100,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 0), bytes.Repeat([]byte{'\x00'}, 99)...), []byte{'\x64'}...),
			want:      []byte{},
		},
		"size=100 data=1": {
			blockSize: 100,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 1), bytes.Repeat([]byte{'\x00'}, 98)...), []byte{'\x63'}...),
			want:      bytes.Repeat([]byte{'\x01'}, 1),
		},
		"size=100 data=50": {
			blockSize: 100,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 50), bytes.Repeat([]byte{'\x00'}, 49)...), []byte{'\x32'}...),
			want:      bytes.Repeat([]byte{'\x01'}, 50),
		},
		"size=100 data=100": {
			blockSize: 100,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 100), bytes.Repeat([]byte{'\x00'}, 99)...), []byte{'\x64'}...),
			want:      bytes.Repeat([]byte{'\x01'}, 100),
		},
		"size=100 data=150": {
			blockSize: 100,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 150), bytes.Repeat([]byte{'\x00'}, 49)...), []byte{'\x32'}...),
			want:      bytes.Repeat([]byte{'\x01'}, 150),
		},
		"size=255": {
			blockSize: 255,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 0), bytes.Repeat([]byte{'\x00'}, 254)...), []byte{'\xff'}...),
			want:      []byte{},
		},
		"size=255 data=1": {
			blockSize: 255,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 1), bytes.Repeat([]byte{'\x00'}, 253)...), []byte{'\xfe'}...),
			want:      bytes.Repeat([]byte{'\x01'}, 1),
		},
		"size=255 data=100": {
			blockSize: 255,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 100), bytes.Repeat([]byte{'\x00'}, 154)...), []byte{'\x9b'}...),
			want:      bytes.Repeat([]byte{'\x01'}, 100),
		},
		"size=255 data=260": {
			blockSize: 255,
			data:      append(append(bytes.Repeat([]byte{'\x01'}, 260), bytes.Repeat([]byte{'\x00'}, 249)...), []byte{'\xfa'}...),
			want:      bytes.Repeat([]byte{'\x01'}, 260),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			b, err := internal.UnpadISO10126(tc.blockSize, tc.data)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
			ztesting.AssertEqual(t, "data not match", tc.want, b)
			if tc.err == nil {
				bb, err := internal.PadISO10126(tc.blockSize, b)
				ztesting.AssertEqualErr(t, "non nil error", nil, err)
				ztesting.AssertEqual(t, "data not match", tc.data, bb)
			}
		})
	}
}
