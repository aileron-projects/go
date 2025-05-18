package zlog_test

import (
	"context"
	"fmt"

	"github.com/aileron-projects/go/zlog"
)

func ExampleRange() {
	fmt.Printf("%s %d\n", zlog.Trace, zlog.Trace)
	fmt.Printf("%s %d\n", zlog.Debug, zlog.Debug)
	fmt.Printf("%s %d\n", zlog.Info, zlog.Info)
	fmt.Printf("%s %d\n", zlog.Warn, zlog.Warn)
	fmt.Printf("%s %d\n", zlog.Error, zlog.Error)
	fmt.Printf("%s %d\n", zlog.Fatal, zlog.Fatal)
	fmt.Printf("%s %d\n", zlog.Undefined, zlog.Undefined)
	fmt.Printf("%s %d\n", zlog.Range(99), zlog.Range(99))
	// Output:
	// TRACE 2
	// DEBUG 4
	// INFO 8
	// WARN 16
	// ERROR 32
	// FATAL 64
	// UNDEFINED 1
	// UNDEFINED 99
}

func ExampleLevel() {
	fmt.Printf("%s %d\n", zlog.LvTrace, zlog.LvTrace)
	fmt.Printf("%s %d\n", zlog.LvDebug, zlog.LvDebug)
	fmt.Printf("%s %d\n", zlog.LvInfo, zlog.LvInfo)
	fmt.Printf("%s %d\n", zlog.LvWarn, zlog.LvWarn)
	fmt.Printf("%s %d\n", zlog.LvError, zlog.LvError)
	fmt.Printf("%s %d\n", zlog.LvFatal, zlog.LvFatal)
	fmt.Printf("%s %d\n", zlog.LvUndef, zlog.LvUndef)
	fmt.Printf("%s %d\n", zlog.Level(99), zlog.Level(99))
	// Output:
	// TRACE 1
	// DEBUG 5
	// INFO 9
	// WARN 13
	// ERROR 17
	// FATAL 21
	// UNDEFINED 0
	// UNDEFINED 99
}

func ExampleLevel_compare() {
	fmt.Printf("INFO >= WARN : %v\n", zlog.LvInfo.HigherEqual(zlog.LvWarn))
	fmt.Printf("INFO >= INFO : %v\n", zlog.LvInfo.HigherEqual(zlog.LvInfo))
	fmt.Printf("INFO >= DEBUG : %v\n", zlog.LvInfo.HigherEqual(zlog.LvDebug))
	// Output:
	// INFO >= WARN : false
	// INFO >= INFO : true
	// INFO >= DEBUG : true
}

func ExampleContextWithAttrs() {
	// Store attributes in the context.
	ctx := context.Background()
	ctx = zlog.ContextWithAttrs(ctx, "foo", 123)
	// Extract attributes from the context.
	attrs := zlog.AttrsFromContext(ctx)

	fmt.Println(attrs)
	// Output:
	// [foo 123]
}

func ExampleContextWithLevel() {
	// Store log level in the context.
	ctx := context.Background()
	ctx = zlog.ContextWithLevel(ctx, zlog.LvError)
	// Extract attributes from the context.
	lv := zlog.LevelFromContext(ctx)

	fmt.Println(lv)
	// Output:
	// ERROR
}
