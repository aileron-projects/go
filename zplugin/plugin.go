package zplugin

import (
	"errors"
	"plugin"
)

var (
	// ErrAssert represents type assertion such like value.(T) failed.
	ErrAssert = errors.New("zplugin: failed to assert type")
)

// pluginOpen opens plugin.
// This is defined so it can be replaced when testing.
var pluginOpen = func(path string) (Lookupper, error) {
	return plugin.Open(path)
}

// Lookupper lookups exported package symbols.
// [plugin.Plugin] implements the interface.
type Lookupper interface {
	Lookup(symName string) (plugin.Symbol, error)
}

// WithTestOpenFunc replaces [plugin.Open] func for testing.
// Note that this changes global state.
func WithTestOpenFunc(f func(path string) (Lookupper, error)) (done func()) {
	tmp := pluginOpen
	pluginOpen = f
	return func() {
		pluginOpen = tmp
	}
}

// Lookup reads Go [plugin] from the path and lookups symbol by the name.
// It is short for [plugin.Open] and [plugin.Plugin.Lookup].
//
// For example, a plugin defined as
//
//	package main
//	import "fmt"
//
//	func MyFunc() {
//	  fmt.Printf("Hello, number %d\n", V)
//	}
//
// may be loaded with the [Lookup] function and then the exported package
// symbols MyFunc can be accessed
//
//	sym, err := zplugin.Lookup("plugin_name.so", "MyFunc")
//	if err != nil {
//		panic(err)
//	}
//	myFunc := sym.(func()) // Now MyFunc() is available with myFunc.
func Lookup(path, symName string) (plugin.Symbol, error) {
	p, err := pluginOpen(path)
	if err != nil {
		return nil, err
	}
	s, err := p.Lookup(symName)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// LookupAll reads Go [plugin] from the given path and lookups all given symbol names.
// See the comment on [Lookup].
func LookupAll(path string, symNames ...string) (map[string]plugin.Symbol, error) {
	p, err := pluginOpen(path)
	if err != nil {
		return nil, err
	}
	symbols := make(map[string]plugin.Symbol, len(symNames))
	for _, name := range symNames {
		s, err := p.Lookup(name)
		if err != nil {
			return nil, err
		}
		symbols[name] = s
	}
	return symbols, nil
}

// Use reads Go [plugin] from the path and lookups symbol by the name.
// Then it also assert symbol into the type T.
// It is short for [Lookup] and type assertion.
//
// For example, a plugin defined as
//
//	package main
//	import "fmt"
//
//	func MyFunc() {
//	  fmt.Printf("Hello, number %d\n", V)
//	}
//
// may be loaded with the [Use] function and then the exported package
// symbols MyFunc can be accessed
//
//	// MyFunc() is available with myFunc.
//	myFunc, err := zplugin.Use[func()]("plugin_name.so", "MyFunc")
//	if err != nil {
//		panic(err)
//	}
func Use[T any](path, symName string) (plugin.Symbol, T, error) {
	sym, err := Lookup(path, symName)
	if err != nil {
		var t T
		return nil, t, err
	}
	t, ok := sym.(T)
	if !ok {
		var t T
		return sym, t, ErrAssert
	}
	return sym, t, nil
}
