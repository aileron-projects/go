package zuid

import (
	"crypto/rand"
	"encoding/binary"
	"hash/fnv"
	"os"
	"sync/atomic"
	"time"

	"github.com/aileron-projects/go/internal/helper"
)

var (
	hostnameFNV1a64 = mustFNV1a64Hash(os.Hostname())
	counter         = atomic.Uint64{}
	timeNow         = time.Now
)

func mustFNV1a64Hash(s string, err error) []byte {
	if err != nil {
		panic(err)
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum(nil)
}

// NewTimeBase returns a new 30 bytes timestamp based ID.
// The returned ID is sortable by time, higher entropy than UUIDv4
// and is hard to guess because of the random value.
// IDs are hard to duplicate but the possibility is not zero.
// The returned 30 bytes IDs do not contain any padding string
// when encoded with [encoding/hex], [encoding/base64] and [encoding/base32].
//
// ID consists of:
//
//   - 8 bytes unix time in microsecond. Valid until January 10th, 294247.
//   - 22 bytes random value read from [crypto/rand.Reader].
//
// Bit arrangements:
//
//	0                   1                   2                   3
//	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                         timestamp_high                        |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                         timestamp_low                         |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|             random            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// Encoded ID examples: (time:2000-01-01, random:1234567890123456789012)
//
//   - [encoding/hex]: 00035d013b37e00031323334353637383930313233343536373839303132
//   - [encoding/base32.StdEncoding]: AABV2AJ3G7QAAMJSGM2DKNRXHA4TAMJSGM2DKNRXHA4TAMJS
//   - [encoding/base32.HexEncoding]: 001LQ09R6VG00C9I6CQ3ADHN70SJ0C9I6CQ3ADHN70SJ0C9I
//   - [encoding/base64.StdEncoding]: AANdATs34AAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEy
//   - [encoding/base64.URLEncoding]: AANdATs34AAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEy
func NewTimeBase() []byte {
	// x is a 30 bytes unique ID.
	id := [30]byte{}
	binary.BigEndian.PutUint64(id[0:], uint64(timeNow().UnixMicro())) // Initial 8 bytes timestamp.
	_, err := rand.Read(id[8:])                                       // Rest of 22 bytes random.
	helper.MustNil(err)
	return id[:]
}

// NewHostBase returns a new 30 bytes hostname based ID.
// The returned ID is sortable by time, nearly has the same entropy as UUIDv4
// and is hard to guess because of the random value.
// IDs are hard to duplicate but the possibility is not zero.
// The returned 30 bytes IDs do not contain any padding string
// when encoded with [encoding/hex], [encoding/base64] and [encoding/base32].
//
// ID consists of:
//
//   - 8 bytes unix time in microsecond. Valid until January 10th, 294247.
//   - 8 bytes FNV1a/64 hash of the hostname.
//   - 14 bytes random value read from [crypto/rand.Reader].
//
// Bit arrangements:
//
//	0                   1                   2                   3
//	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                         timestamp_high                        |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                         timestamp_low                         |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                    FNV1a/64(hostname)_high                    |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                    FNV1a/64(hostname)_low                     |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|             random            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// Encoded ID examples: (time:2000-01-01, random:12345678901234 hostFNV:host1234)
//
//   - [encoding/hex]: 00035d013b37e000686f7374313233343132333435363738393031323334
//   - [encoding/base32.StdEncoding]: AABV2AJ3G7QAA2DPON2DCMRTGQYTEMZUGU3DOOBZGAYTEMZU
//   - [encoding/base32.HexEncoding]: 001LQ09R6VG00Q3FEDQ32CHJ6GOJ4CPK6KR3EE1P60OJ4CPK
//   - [encoding/base64.StdEncoding]: AANdATs34ABob3N0MTIzNDEyMzQ1Njc4OTAxMjM0
//   - [encoding/base64.URLEncoding]: AANdATs34ABob3N0MTIzNDEyMzQ1Njc4OTAxMjM0
func NewHostBase() []byte {
	id := [30]byte{}
	binary.BigEndian.PutUint64(id[0:], uint64(timeNow().UnixMicro())) // Initial 8 bytes timestamp.
	copy(id[8:], hostnameFNV1a64)                                     // Next 8 bytes hostname hash.
	_, err := rand.Read(id[16:])                                      // Last 14 bytes random.
	helper.MustNil(err)
	return id[:]
}

// NewCountBase returns a new 30 bytes counter based ID.
// The returned ID is sortable by time, nearly has the same entropy as UUIDv4
// and is hard to guess because of the random value.
// IDs are hard to duplicate but the possibility is not zero.
// The returned 30 bytes IDs do not contain any padding string
// when encoded with [encoding/hex], [encoding/base64] and [encoding/base32].
//
// ID consists of:
//
//   - 8 bytes unix time in microsecond. Valid until January 10th, 294247.
//   - 14 bytes random value read from [crypto/rand.Reader].
//   - 8 bytes unsigned integer counter (reset to zero when overflow).
//
// Bit arrangements:
//
//	0                   1                   2                   3
//	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                         timestamp_high                        |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                         timestamp_low                         |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                             random                            |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|            random             |         counter_high          |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                        counter_middle                         |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|          counter_low          |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// Encoded ID examples: (time:2000-01-01, random:12345678901234 counter:1)
//
//   - [encoding/hex]: 00035d013b37e00031323334353637383930313233340000000000000001
//   - [encoding/base32.StdEncoding]: AABV2AJ3G7QAAMJSGM2DKNRXHA4TAMJSGM2AAAAAAAAAAAAB
//   - [encoding/base32.HexEncoding]: 001LQ09R6VG00C9I6CQ3ADHN70SJ0C9I6CQ0000000000001
//   - [encoding/base64.StdEncoding]: AANdATs34AAxMjM0NTY3ODkwMTIzNAAAAAAAAAAB
//   - [encoding/base64.URLEncoding]: AANdATs34AAxMjM0NTY3ODkwMTIzNAAAAAAAAAAB
func NewCountBase() []byte {
	id := [30]byte{}
	binary.BigEndian.PutUint64(id[0:], uint64(timeNow().UnixMicro())) // Initial 8 bytes timestamp.
	_, err := rand.Read(id[8:22])                                     // Next 14 bytes random.
	helper.MustNil(err)
	binary.BigEndian.PutUint64(id[22:], counter.Add(1)) // Last 8 bytes counter.
	return id[:]
}
