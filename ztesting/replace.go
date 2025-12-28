package ztesting

import (
	"crypto/rand"
	"io"
	"os"
)

// ReplaceRandReader replaces [rand.Reader] with r.
// Do not run test parallel when using this.
// Call done after the test finished to set original rand.Reader.
//
// Example:
//
//	done := ztesting.ReplaceRandReader(YourReader)
//	defer done()
//	// Write your test here.
func ReplaceRandReader(r io.Reader) (done func()) {
	tmp := rand.Reader
	rand.Reader = r
	return func() {
		rand.Reader = tmp
	}
}

// ReplaceStdout replaces [os.Stdout] and return reader.
// Do not run test parallel when using this.
// Call done after the test finished to set original Stdout.
//
// Example:
//
//	r, done := ztesting.ReplaceStdout()
//	defer done()
//	// Write your test here.
func ReplaceStdout() (r *os.File, done func()) {
	tmp := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	return r, func() {
		os.Stdout = tmp
	}
}

// ReplaceStderr replaces [os.Stderr] and return reader.
// Do not run test parallel when using this.
// Call done after the test finished to set original Stderr.
//
// Example:
//
//	r, done := ztesting.ReplaceStderr()
//	defer done()
//	// Write your test here.
func ReplaceStderr() (r *os.File, done func()) {
	tmp := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	return r, func() {
		os.Stderr = tmp
	}
}

// ReplaceStdin replaces [os.Stdin] and return writer.
// Do not run test parallel when using this.
// Call done after the test finished to set original Stdin.
//
// Example:
//
//	w, done := ztesting.ReplaceStdin()
//	defer done()
//	// Write your test here.
func ReplaceStdin() (w *os.File, done func()) {
	tmp := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	return w, func() {
		os.Stdin = tmp
	}
}
