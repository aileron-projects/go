package zrc4_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aileron-projects/go/zcrypto/zrc4"
)

func ExampleNewStreamReader_encrypt() {
	msg := "Hello Go!"
	key := []byte("secret-key")
	r, _ := zrc4.NewStreamReader(key, strings.NewReader(msg))

	buf := make([]byte, 100)
	n, _ := r.Read(buf)
	fmt.Println("Encrypted:", hex.EncodeToString(buf[:n]))
	// Output:
	// Encrypted: f86d22f75b8d4bb719
}

func ExampleNewStreamReader_decrypt() {
	msg, _ := hex.DecodeString("f86d22f75b8d4bb719")
	key := []byte("secret-key")
	r, _ := zrc4.NewStreamReader(key, bytes.NewReader(msg))

	buf := make([]byte, 100)
	n, _ := r.Read(buf)
	fmt.Println("Decrypted:", string(buf[:n]))
	// Output:
	// Decrypted: Hello Go!
}

func ExampleNewStreamWriter_encrypt() {
	msg := []byte("Hello Go!")
	key := []byte("secret-key")

	var buf bytes.Buffer
	w, _ := zrc4.NewStreamWriter(key, &buf)
	w.Write(msg)
	fmt.Println("Encrypted:", hex.EncodeToString(buf.Bytes()))
	// Output:
	// Encrypted: f86d22f75b8d4bb719
}

func ExampleNewStreamWriter_decrypt() {
	msg, _ := hex.DecodeString("f86d22f75b8d4bb719")
	key := []byte("secret-key")

	var buf bytes.Buffer
	w, _ := zrc4.NewStreamWriter(key, &buf)
	w.Write(msg)
	fmt.Println("Decrypted:", buf.String())
	// Output:
	// Decrypted: Hello Go!
}

func ExampleCopy_encrypt() {
	msg := []byte("Hello Go!")
	key := []byte("secret-key")

	var buf bytes.Buffer
	zrc4.Copy(key, &buf, bytes.NewReader(msg))
	fmt.Println("Encrypted:", hex.EncodeToString(buf.Bytes()))
	// Output:
	// Encrypted: f86d22f75b8d4bb719
}

func ExampleCopy_decrypt() {
	msg, _ := hex.DecodeString("f86d22f75b8d4bb719")
	key := []byte("secret-key")

	var buf bytes.Buffer
	zrc4.Copy(key, &buf, bytes.NewReader(msg))
	fmt.Println("Decrypted:", buf.String())
	// Output:
	// Decrypted: Hello Go!
}
