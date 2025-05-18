package zcrc32

import (
	"bytes"
	"hash"
	"hash/crc32"

	"github.com/aileron-projects/go/internal/ihash"
	"github.com/aileron-projects/go/zhash"
)

var (
	_ ihash.SumFunc = SumIEEE
	_ ihash.SumFunc = SumCastagnoli
	_ ihash.SumFunc = SumKoopman

	_ ihash.EqualSumFunc = EqualSumIEEE
	_ ihash.EqualSumFunc = EqualSumCastagnoli
	_ ihash.EqualSumFunc = EqualSumKoopman
)

var (
	IEEETable       = crc32.IEEETable
	CastagnoliTable = crc32.MakeTable(crc32.Castagnoli)
	KoopmanTable    = crc32.MakeTable(crc32.Koopman)
)

func init() {
	zhash.RegisterHash(zhash.CRC32IEEE, func() hash.Hash { return crc32.New(IEEETable) })
	zhash.RegisterHash(zhash.CRC32Castagnoli, func() hash.Hash { return crc32.New(CastagnoliTable) })
	zhash.RegisterHash(zhash.CRC32Koopman, func() hash.Hash { return crc32.New(KoopmanTable) })
}

// SumIEEE returns CRC32 hash using [crc32.IEEE] table.
// It uses [hash/crc32.New].
func SumIEEE(b []byte) []byte {
	h := crc32.New(IEEETable)
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, crc32.Size))
}

// SumCastagnoli returns CRC32 hash using [crc32.Castagnoli] table.
// It uses [hash/crc32.New].
func SumCastagnoli(b []byte) []byte {
	h := crc32.New(CastagnoliTable)
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, crc32.Size))
}

// SumKoopman returns CRC32 hash using [crc32.Koopman] table.
// It uses [hash/crc32.New].
func SumKoopman(b []byte) []byte {
	h := crc32.New(KoopmanTable)
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, crc32.Size))
}

// EqualSumIEEE compares CRC32 hash using [crc32.IEEE] table.
// It returns if the sum matches to the hash of b.
func EqualSumIEEE(b []byte, sum []byte) bool {
	return bytes.Equal(SumIEEE(b), sum)
}

// EqualSumECMA compares CRC32 hash using [crc32.Castagnoli] table.
// It returns if the sum matches to the hash of b.
func EqualSumCastagnoli(b []byte, sum []byte) bool {
	return bytes.Equal(SumCastagnoli(b), sum)
}

// EqualSumECMA compares CRC32 hash using [crc32.Koopman] table.
// It returns if the sum matches to the hash of b.
func EqualSumKoopman(b []byte, sum []byte) bool {
	return bytes.Equal(SumKoopman(b), sum)
}
