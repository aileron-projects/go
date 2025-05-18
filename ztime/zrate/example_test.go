package zrate_test

import (
	"fmt"

	"github.com/aileron-projects/go/ztime/zrate"
)

func ExampleConcurrentLimiter() {
	limiter := zrate.NewConcurrentLimiter(3)
	for i := range 5 {
		token := limiter.AllowNow()
		defer token.Release() // Occupied token must be released.
		fmt.Println(i, token.OK())
	}
	// Output:
	// 0 true
	// 1 true
	// 2 true
	// 3 false
	// 4 false
}

func ExampleNewFixedWindowLimiter() {
	limiter := zrate.NewFixedWindowLimiter(3)
	for i := range 5 {
		token := limiter.AllowNow()
		fmt.Println(i, token.OK())
	}
	// Output:
	// 0 true
	// 1 true
	// 2 true
	// 3 false
	// 4 false
}

func ExampleNewSlidingWindowLimiter() {
	limiter := zrate.NewSlidingWindowLimiter(3)
	for i := range 5 {
		token := limiter.AllowNow()
		fmt.Println(i, token.OK())
	}
	// Output:
	// 0 true
	// 1 true
	// 2 true
	// 3 false
	// 4 false
}

func ExampleNewTokenBucketLimiter() {
	limiter := zrate.NewTokenBucketLimiter(3, 2)
	for i := range 5 {
		token := limiter.AllowNow()
		fmt.Println(i, token.OK())
	}
	// Output:
	// 0 true
	// 1 true
	// 2 true
	// 3 false
	// 4 false
}
