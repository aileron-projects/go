package ztesting_test

import (
	"crypto/rand"
	"os"
	"strings"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestReplaceRandReader(t *testing.T) {
	done := ztesting.ReplaceRandReader(strings.NewReader("12345"))
	defer done()
	b := make([]byte, 5)
	n, err := rand.Read(b)
	ztesting.AssertEqual(t, "read content not match", "12345", string(b))
	ztesting.AssertEqual(t, "read bytes not match", 5, n)
	ztesting.AssertEqual(t, "non nil error returned", nil, err)
}

func TestReplaceStdout(t *testing.T) {
	r, done := ztesting.ReplaceStdout()
	defer done()
	os.Stdout.Write([]byte("12345"))
	b := make([]byte, 5)
	n, err := r.Read(b)
	ztesting.AssertEqual(t, "written content not match", "12345", string(b))
	ztesting.AssertEqual(t, "written bytes not match", 5, n)
	ztesting.AssertEqual(t, "non nil error returned", nil, err)
}

func TestReplaceStderr(t *testing.T) {
	r, done := ztesting.ReplaceStderr()
	defer done()
	os.Stderr.Write([]byte("12345"))
	b := make([]byte, 5)
	n, err := r.Read(b)
	ztesting.AssertEqual(t, "written content not match", "12345", string(b))
	ztesting.AssertEqual(t, "written bytes not match", 5, n)
	ztesting.AssertEqual(t, "non nil error returned", nil, err)
}

func TestReplaceStdin(t *testing.T) {
	w, done := ztesting.ReplaceStdin()
	defer done()
	w.Write([]byte("12345"))
	b := make([]byte, 5)
	n, err := os.Stdin.Read(b)
	ztesting.AssertEqual(t, "read content not match", "12345", string(b))
	ztesting.AssertEqual(t, "read bytes not match", 5, n)
	ztesting.AssertEqual(t, "non nil error returned", nil, err)
}
