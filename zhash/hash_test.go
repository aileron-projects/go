package zhash

import (
	"encoding/hex"
	"hash"
	"hash/fnv"
	"slices"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestRegisterHash(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		h := hash.Hash(fnv.New32())
		RegisterHash(FNV32, func() hash.Hash { return h })
		got := FNV32.New()
		ztesting.AssertEqual(t, "hash not match", true, slices.Equal(h.Sum(nil), got.Sum(nil)))
	})
	t.Run("failed", func(t *testing.T) {
		defer func() {
			r := recover()
			ztesting.AssertEqualErr(t, "error not match", ErrUnknown, r.(error))
		}()
		RegisterHash(maxHash, nil)
	})
}

func TestHash_String(t *testing.T) {
	t.Parallel()
	testCase := map[string]struct {
		h Hash
		s string
	}{
		"FNV32":           {FNV32, "FNV1/32"},
		"FNV32a":          {FNV32a, "FNV1a/32"},
		"FNV64":           {FNV64, "FNV1/64"},
		"FNV64a":          {FNV64a, "FNV1a/64"},
		"FNV128":          {FNV128, "FNV1/128"},
		"FNV128a":         {FNV128a, "FNV1a/128"},
		"CRC32IEEE":       {CRC32IEEE, "CRC32-IEEE"},
		"CRC32Castagnoli": {CRC32Castagnoli, "CRC32-Castagnoli"},
		"CRC32Koopman":    {CRC32Koopman, "CRC32-Koopman"},
		"CRC64ISO":        {CRC64ISO, "CRC64-ISO"},
		"CRC64ECMA":       {CRC64ECMA, "CRC64-ECMA"},
		"Unknown":         {maxHash, "unknown hash value 31"},
	}
	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			got := tc.h.String()
			ztesting.AssertEqual(t, "name does not match", tc.s, got)
		})
	}
}

func TestHash_Size(t *testing.T) {
	t.Parallel()
	t.Run("panic unknown", func(t *testing.T) {
		defer func() {
			r := recover()
			ztesting.AssertEqualErr(t, "error not match", ErrUnknown, r.(error))
		}()
		maxHash.Size()
	})

	testCase := map[string]struct {
		h Hash
		s int
	}{
		"FNV32":           {FNV32, 4},
		"FNV32a":          {FNV32a, 4},
		"FNV64":           {FNV64, 8},
		"FNV64a":          {FNV64a, 8},
		"FNV128":          {FNV128, 16},
		"FNV128a":         {FNV128a, 16},
		"CRC32IEEE":       {CRC32IEEE, 4},
		"CRC32Castagnoli": {CRC32Castagnoli, 4},
		"CRC32Koopman":    {CRC32Koopman, 4},
		"CRC64ISO":        {CRC64ISO, 8},
		"CRC64ECMA":       {CRC64ECMA, 8},
	}
	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			got := tc.h.Size()
			ztesting.AssertEqual(t, "size does not match", tc.s, got)
		})
	}
}

func TestHash_Available(t *testing.T) {
	t.Parallel()
	t.Run("max hash", func(t *testing.T) {
		got := maxHash.Available()
		ztesting.AssertEqual(t, "unknown alg is available", false, got)
	})
	t.Run("available", func(t *testing.T) {
		RegisterHash(FNV32, func() hash.Hash { return fnv.New32() })
		got := FNV32.Available()
		ztesting.AssertEqual(t, "unknown alg is available", true, got)
	})
}

func TestHash_New(t *testing.T) {
	t.Parallel()
	t.Run("not available", func(t *testing.T) {
		defer func() {
			r := recover()
			ztesting.AssertEqualErr(t, "error not match", ErrNotAvailable, r.(error))
		}()
		maxHash.New()
	})
	t.Run("hash available", func(t *testing.T) {
		RegisterHash(FNV32, func() hash.Hash { return fnv.New32() })
		h := FNV32.New()
		h.Write([]byte("test"))
		ztesting.AssertEqual(t, "invalid hash", "bc2c0be9", hex.EncodeToString(h.Sum(nil)))
	})
}

func TestHash_NewFunc(t *testing.T) {
	t.Parallel()
	t.Run("not available", func(t *testing.T) {
		defer func() {
			r := recover()
			ztesting.AssertEqualErr(t, "error not match", ErrNotAvailable, r.(error))
		}()
		maxHash.NewFunc()
	})
	t.Run("hash available", func(t *testing.T) {
		RegisterHash(FNV32, func() hash.Hash { return fnv.New32() })
		f := FNV32.NewFunc()
		h := f()
		h.Write([]byte("test"))
		ztesting.AssertEqual(t, "invalid hash", "bc2c0be9", hex.EncodeToString(h.Sum(nil)))
	})
}

func TestHash_Sum(t *testing.T) {
	t.Parallel()
	t.Run("not available", func(t *testing.T) {
		defer func() {
			r := recover()
			ztesting.AssertEqualErr(t, "error not match", ErrNotAvailable, r.(error))
		}()
		maxHash.Sum([]byte("test"))
	})
	t.Run("available", func(t *testing.T) {
		RegisterHash(FNV32, func() hash.Hash { return fnv.New32() })
		got := FNV32.Sum([]byte("test"))
		ztesting.AssertEqual(t, "invalid returned value", "bc2c0be9", hex.EncodeToString(got))
	})
}

func TestHash_Equal(t *testing.T) {
	t.Parallel()
	t.Run("equal", func(t *testing.T) {
		RegisterHash(FNV32, func() hash.Hash { return fnv.New32() })
		b, _ := hex.DecodeString("bc2c0be9")
		got := FNV32.Equal(b, []byte("test"))
		ztesting.AssertEqual(t, "invalid compare result", true, got)
	})
	t.Run("not equal", func(t *testing.T) {
		RegisterHash(FNV32, func() hash.Hash { return fnv.New32() })
		b, _ := hex.DecodeString("bc2c0be9")
		got := FNV32.Equal(b, []byte("wrong"))
		ztesting.AssertEqual(t, "invalid compare result", false, got)
	})
}

func TestHash_Compare(t *testing.T) {
	t.Parallel()
	t.Run("equal", func(t *testing.T) {
		RegisterHash(FNV32, func() hash.Hash { return fnv.New32() })
		b, _ := hex.DecodeString("bc2c0be9")
		got := FNV32.Compare(b, []byte("test"))
		ztesting.AssertEqualErr(t, "invalid compare result", nil, got)
	})
	t.Run("not equal", func(t *testing.T) {
		RegisterHash(FNV32, func() hash.Hash { return fnv.New32() })
		b, _ := hex.DecodeString("bc2c0be9")
		got := FNV32.Compare(b, []byte("wrong"))
		ztesting.AssertEqualErr(t, "invalid compare result", ErrNotMatch, got)
	})
	t.Run("alg unknown", func(t *testing.T) {
		err := maxHash.Compare([]byte("xxx"), []byte("yyy"))
		ztesting.AssertEqualErr(t, "error not match", ErrNotAvailable, err)
	})
}
