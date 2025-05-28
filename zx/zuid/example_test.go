package zuid_test

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/aileron-projects/go/zx/zuid"
)

func ExampleContextWithID() {
	ctx := zuid.ContextWithID(context.Background(), "key", "example-unique-id")
	uid := zuid.FromContext(ctx, "key")
	fmt.Println(uid)
	// Output:
	// example-unique-id
}

func ExampleNewTime() {
	id := zuid.NewTime()
	fmt.Println("Len :", len(id))
	fmt.Println("hex        :", hex.EncodeToString(id))
	fmt.Println("base32.Std :", base32.StdEncoding.EncodeToString(id))
	fmt.Println("base32.Hex :", base32.HexEncoding.EncodeToString(id))
	fmt.Println("base64.Std :", base64.StdEncoding.EncodeToString(id))
	fmt.Println("base64.URL :", base64.URLEncoding.EncodeToString(id))
}

func ExampleNewHost() {
	id := zuid.NewHost()
	fmt.Println("Len :", len(id))
	fmt.Println("hex        :", hex.EncodeToString(id))
	fmt.Println("base32.Std :", base32.StdEncoding.EncodeToString(id))
	fmt.Println("base32.Hex :", base32.HexEncoding.EncodeToString(id))
	fmt.Println("base64.Std :", base64.StdEncoding.EncodeToString(id))
	fmt.Println("base64.URL :", base64.URLEncoding.EncodeToString(id))
}

func ExampleNewCount() {
	id := zuid.NewCount()
	fmt.Println("Len :", len(id))
	fmt.Println("hex        :", hex.EncodeToString(id))
	fmt.Println("base32.Std :", base32.StdEncoding.EncodeToString(id))
	fmt.Println("base32.Hex :", base32.HexEncoding.EncodeToString(id))
	fmt.Println("base64.Std :", base64.StdEncoding.EncodeToString(id))
	fmt.Println("base64.URL :", base64.URLEncoding.EncodeToString(id))
}
