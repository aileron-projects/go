package zcrc64

import (
	"bytes"
	"hash"
	"hash/crc64"

	"github.com/aileron-projects/go/internal/ihash"
	"github.com/aileron-projects/go/zhash"
)

var (
	_ ihash.SumFunc = SumISO
	_ ihash.SumFunc = SumECMA

	_ ihash.EqualSumFunc = EqualSumISO
	_ ihash.EqualSumFunc = EqualSumECMA
)

var (
	ISOTable  = crc64.MakeTable(crc64.ISO)
	ECMATable = crc64.MakeTable(crc64.ECMA)
)

func init() {
	zhash.RegisterHash(zhash.CRC64ISO, func() hash.Hash { return crc64.New(ISOTable) })
	zhash.RegisterHash(zhash.CRC64ECMA, func() hash.Hash { return crc64.New(ECMATable) })
}

// SumISO returns CRC64 hash using [crc64.ISO] table.
// It uses [hash/crc64.New].
func SumISO(b []byte) []byte {
	h := crc64.New(ISOTable)
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, crc64.Size))
}

// SumECMA returns CRC64 hash using [crc64.ECMA] table.
// It uses [hash/crc64.New].
func SumECMA(b []byte) []byte {
	h := crc64.New(ECMATable)
	_, _ = h.Write(b)
	return h.Sum(make([]byte, 0, crc64.Size))
}

// EqualSumISO compares CRC64 hash using [crc64.ISO] table.
// It returns if the sum matches to the hash of b.
func EqualSumISO(b []byte, sum []byte) bool {
	return bytes.Equal(SumISO(b), sum)
}

// EqualSumECMA compares CRC64 hash using [crc64.ECMA] table.
// It returns if the sum matches to the hash of b.
func EqualSumECMA(b []byte, sum []byte) bool {
	return bytes.Equal(SumECMA(b), sum)
}
