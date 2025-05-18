package zuid

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"math"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztesting"
)

func TestMustFNV1a64Hash(t *testing.T) {
	t.Parallel()
	t.Run("nil error", func(t *testing.T) {
		b := mustFNV1a64Hash("test", nil)
		ztesting.AssertEqual(t, "hash not match", "f9e6e6ef197c2b25", hex.EncodeToString(b))
	})
	t.Run("non-nil error", func(t *testing.T) {
		defer func() {
			err := recover().(error)
			ztesting.AssertEqualErr(t, "panicked error not match", io.EOF, err)
		}()
		b := mustFNV1a64Hash("test", io.EOF)
		ztesting.AssertEqual(t, "hash not match", "", hex.EncodeToString(b))
	})
}

func TestNewTimeBase(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) // Unix 946684800 second.
	}
	defer func() { timeNow = time.Now }()

	rr := strings.NewReader("123456789012345678901234567890")
	done := ztesting.ReplaceRandReader(rr)
	defer done()
	id := NewTimeBase()
	ztesting.AssertEqual(t, "id length not match", 30, len(id))
	ztesting.AssertEqual(t, "time not match", 946684800_000_000, binary.BigEndian.Uint64(id[0:8]))
}

func TestNewHostBase(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) // Unix 946684800 second.
	}
	hostnameFNV1a64 = []byte("host1234")
	defer func() { timeNow = time.Now }()

	rr := strings.NewReader("123456789012345678901234567890")
	done := ztesting.ReplaceRandReader(rr)
	defer done()
	id := NewHostBase()
	ztesting.AssertEqual(t, "id length not match", 30, len(id))
	ztesting.AssertEqual(t, "time not match", 946684800_000_000, binary.BigEndian.Uint64(id[0:8]))
	ztesting.AssertEqual(t, "hostname not match", "host1234", string(id[8:16]))
}

func TestNewCountBase(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) // Unix 946684800 second.
	}
	defer func() { timeNow = time.Now }()

	rr := strings.NewReader("123456789012345678901234567890")
	done := ztesting.ReplaceRandReader(rr)
	defer done()
	counter = atomic.Uint64{}
	counter.Store(math.MaxUint64 - 1)
	id := NewCountBase()
	ztesting.AssertEqual(t, "id length not match", 30, len(id))
	ztesting.AssertEqual(t, "time not match", 946684800_000_000, binary.BigEndian.Uint64(id[0:8]))
	ztesting.AssertEqual(t, "counter not match", uint64(math.MaxUint64), binary.BigEndian.Uint64(id[22:]))
}
