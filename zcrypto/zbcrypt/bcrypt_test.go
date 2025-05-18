package zbcrypt_test

import (
	"strings"
	"testing"

	"github.com/aileron-projects/go/zcrypto/zbcrypt"
	"github.com/aileron-projects/go/ztesting"
	"golang.org/x/crypto/bcrypt"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("1234567890123456789012345678901234567890"))
		defer done()
		b, err := zbcrypt.New(10)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "hash not match", "$2a$10$Lxe3KBCwKxOzLha2MR.vKePt0LpiRJVu6NMWU/u7ddWJx0sm6kBb.", string(h))
	})
	t.Run("min cost", func(t *testing.T) {
		done := ztesting.ReplaceRandReader(strings.NewReader("1234567890123456789012345678901234567890"))
		defer done()
		b, err := zbcrypt.New(4)
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
		h, _ := b.Sum([]byte("test"))
		ztesting.AssertEqual(t, "hash not match", "$2a$04$Lxe3KBCwKxOzLha2MR.vKeIyr1H1Cb4ubhc770bjAQSAVlORxnG2S", string(h))
	})
	t.Run("param invalid cost=32", func(t *testing.T) {
		b, err := zbcrypt.New(32)
		ztesting.AssertEqual(t, "non nil hasher returned", nil, b)
		ztesting.AssertEqualErr(t, "nil error returned", bcrypt.InvalidCostError(32), err)
	})
}

func TestEqual(t *testing.T) {
	done := ztesting.ReplaceRandReader(strings.NewReader("1234567890123456789012345678901234567890"))
	defer done()
	b, err := zbcrypt.New(10)
	ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	h, _ := b.Sum([]byte("test"))
	ztesting.AssertEqual(t, "incorrect equal result", true, b.Equal(h, []byte("test")))
	ztesting.AssertEqual(t, "incorrect equal result", false, b.Equal(h, []byte("wrong")))
}

func TestCompare(t *testing.T) {
	done := ztesting.ReplaceRandReader(strings.NewReader("1234567890123456789012345678901234567890"))
	defer done()
	b, err := zbcrypt.New(10)
	ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	h, _ := b.Sum([]byte("test"))
	ztesting.AssertEqualErr(t, "incorrect compare result", nil, b.Compare(h, []byte("test")))
	ztesting.AssertEqualErr(t, "incorrect compare result", bcrypt.ErrMismatchedHashAndPassword, b.Compare(h, []byte("wrong")))
}
