package zchacha20_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aileron-projects/go/zcrypto/zchacha20"
)

func ExampleNewStreamReader_encrypt() {
	msg := "Hello Go!"
	key := []byte("12345678901234567890123456789012")
	nonce := []byte("123456789012")
	r, _ := zchacha20.NewStreamReader(key, nonce, strings.NewReader(msg))

	buf := make([]byte, 100)
	n, _ := r.Read(buf)
	fmt.Println("Encrypted:", hex.EncodeToString(buf[:n]))
	// Output:
	// Encrypted: 2d8cc1577c7648274f
}

func ExampleNewStreamReader_decrypt() {
	msg, _ := hex.DecodeString("2d8cc1577c7648274f")
	key := []byte("12345678901234567890123456789012")
	nonce := []byte("123456789012")
	r, _ := zchacha20.NewStreamReader(key, nonce, bytes.NewReader(msg))

	buf := make([]byte, 100)
	n, _ := r.Read(buf)
	fmt.Println("Decrypted:", string(buf[:n]))
	// Output:
	// Decrypted: Hello Go!
}

func ExampleNewStreamWriter_encrypt() {
	msg := []byte("Hello Go!")
	key := []byte("12345678901234567890123456789012")
	nonce := []byte("123456789012")

	var buf bytes.Buffer
	w, _ := zchacha20.NewStreamWriter(key, nonce, &buf)
	w.Write(msg)
	fmt.Println("Encrypted:", hex.EncodeToString(buf.Bytes()))
	// Output:
	// Encrypted: 2d8cc1577c7648274f
}

func ExampleNewStreamWriter_decrypt() {
	msg, _ := hex.DecodeString("2d8cc1577c7648274f")
	key := []byte("12345678901234567890123456789012")
	nonce := []byte("123456789012")

	var buf bytes.Buffer
	w, _ := zchacha20.NewStreamWriter(key, nonce, &buf)
	w.Write(msg)
	fmt.Println("Decrypted:", buf.String())
	// Output:
	// Decrypted: Hello Go!
}

func ExampleCopy_encrypt() {
	msg := []byte("Hello Go!")
	key := []byte("12345678901234567890123456789012")
	nonce := []byte("123456789012")

	var buf bytes.Buffer
	zchacha20.Copy(key, nonce, &buf, bytes.NewReader(msg))
	fmt.Println("Encrypted:", hex.EncodeToString(buf.Bytes()))
	// Output:
	// Encrypted: 2d8cc1577c7648274f
}

func ExampleCopy_decrypt() {
	msg, _ := hex.DecodeString("2d8cc1577c7648274f")
	key := []byte("12345678901234567890123456789012")
	nonce := []byte("123456789012")

	var buf bytes.Buffer
	zchacha20.Copy(key, nonce, &buf, bytes.NewReader(msg))
	fmt.Println("Decrypted:", buf.String())
	// Output:
	// Decrypted: Hello Go!
}
