package ztesting

import (
	"crypto/rand"
	"io"
	"os"
)

// ReplaceRandReader replaces [rand.Reader] to the r.
// Do not run test parallel when using this.
// Call done function returned from it when tests are completed.
//
//	done := ztesting.ReplaceRandReader(YourReader)
//	defer done()
func ReplaceRandReader(r io.Reader) (done func()) {
	tmp := rand.Reader
	rand.Reader = r
	return func() {
		rand.Reader = tmp
	}
}

// ReplaceStdout replaces [os.Stdout] and return reader.
// Do not run test parallel when using this.
// Call done function returned from it when tests are completed.
//
//	r, done := ztesting.ReplaceStdout()
//	defer done()
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
// Call done function returned from it when tests are completed.
//
//	r, done := ztesting.ReplaceStderr()
//	defer done()
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
// Call done function returned from it when tests are completed.
//
//	w, done := ztesting.ReplaceStdin()
//	defer done()
func ReplaceStdin() (w *os.File, done func()) {
	tmp := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	return w, func() {
		os.Stdin = tmp
	}
}
