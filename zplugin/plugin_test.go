package zplugin_test

import (
	"errors"
	"plugin"
	"testing"

	"github.com/aileron-projects/go/zplugin"
	"github.com/aileron-projects/go/ztesting"
)

type testPlugin struct {
	symbols map[string]plugin.Symbol
	err     error
}

func (p *testPlugin) Lookup(symName string) (plugin.Symbol, error) {
	return p.symbols[symName], p.err
}

func TestLookup(t *testing.T) {
	t.Run("open error", func(t *testing.T) {
		openErr := errors.New("open error")
		done := zplugin.WithTestOpenFunc(func(path string) (zplugin.Lookupper, error) {
			return nil, openErr
		})
		defer done()
		s, err := zplugin.Lookup("test.so", "MyFunc")
		ztesting.AssertEqual(t, "symbol not match", nil, s)
		ztesting.AssertEqualErr(t, "error not match", openErr, err)
	})
	t.Run("lookup error", func(t *testing.T) {
		lookupErr := errors.New("open error")
		done := zplugin.WithTestOpenFunc(func(path string) (zplugin.Lookupper, error) {
			return &testPlugin{err: lookupErr}, nil
		})
		defer done()
		s, err := zplugin.Lookup("test.so", "MyFunc")
		ztesting.AssertEqual(t, "symbol not match", nil, s)
		ztesting.AssertEqualErr(t, "error not match", lookupErr, err)
	})
	t.Run("lookup success", func(t *testing.T) {
		done := zplugin.WithTestOpenFunc(func(path string) (zplugin.Lookupper, error) {
			return &testPlugin{symbols: map[string]plugin.Symbol{"MyVar": "value"}}, nil
		})
		defer done()
		s, err := zplugin.Lookup("test.so", "MyVar")
		ztesting.AssertEqual(t, "symbol not match", plugin.Symbol("value"), s)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
	})
}

func TestLookupAll(t *testing.T) {
	t.Run("open error", func(t *testing.T) {
		openErr := errors.New("open error")
		done := zplugin.WithTestOpenFunc(func(path string) (zplugin.Lookupper, error) {
			return nil, openErr
		})
		defer done()
		ss, err := zplugin.LookupAll("test.so", "MyFunc", "MyVar")
		ztesting.AssertEqual(t, "unexpectedly symbol returned", 0, len(ss))
		ztesting.AssertEqualErr(t, "error not match", openErr, err)
	})
	t.Run("lookup error", func(t *testing.T) {
		lookupErr := errors.New("open error")
		done := zplugin.WithTestOpenFunc(func(path string) (zplugin.Lookupper, error) {
			return &testPlugin{err: lookupErr}, nil
		})
		defer done()
		ss, err := zplugin.LookupAll("test.so", "MyFunc", "MyVar")
		ztesting.AssertEqual(t, "unexpectedly symbol returned", 0, len(ss))
		ztesting.AssertEqualErr(t, "error not match", lookupErr, err)
	})
	t.Run("lookup success", func(t *testing.T) {
		done := zplugin.WithTestOpenFunc(func(path string) (zplugin.Lookupper, error) {
			return &testPlugin{symbols: map[string]plugin.Symbol{
				"MyVar1": "value1",
				"MyVar2": "value2",
			}}, nil
		})
		defer done()
		ss, err := zplugin.LookupAll("test.so", "MyVar1", "MyVar2")
		ztesting.AssertEqual(t, "symbol not match", plugin.Symbol("value1"), ss["MyVar1"])
		ztesting.AssertEqual(t, "symbol not match", plugin.Symbol("value2"), ss["MyVar2"])
		ztesting.AssertEqualErr(t, "error not match", nil, err)
	})
}

func TestUse(t *testing.T) {
	t.Run("open error", func(t *testing.T) {
		openErr := errors.New("open error")
		done := zplugin.WithTestOpenFunc(func(path string) (zplugin.Lookupper, error) {
			return nil, openErr
		})
		defer done()
		s, v, err := zplugin.Use[int]("test.so", "MyVar")
		ztesting.AssertEqual(t, "symbol not match", nil, s)
		ztesting.AssertEqual(t, "value not match", 0, v)
		ztesting.AssertEqualErr(t, "error not match", openErr, err)
	})
	t.Run("lookup error", func(t *testing.T) {
		lookupErr := errors.New("open error")
		done := zplugin.WithTestOpenFunc(func(path string) (zplugin.Lookupper, error) {
			return &testPlugin{err: lookupErr}, nil
		})
		defer done()
		s, v, err := zplugin.Use[int]("test.so", "MyVar")
		ztesting.AssertEqual(t, "symbol not match", nil, s)
		ztesting.AssertEqual(t, "value not match", 0, v)
		ztesting.AssertEqualErr(t, "error not match", lookupErr, err)
	})
	t.Run("type assertion error", func(t *testing.T) {
		done := zplugin.WithTestOpenFunc(func(path string) (zplugin.Lookupper, error) {
			return &testPlugin{symbols: map[string]plugin.Symbol{
				"MyVar": 12345,
			}}, nil
		})
		defer done()
		s, v, err := zplugin.Use[string]("test.so", "MyVar")
		ztesting.AssertEqual(t, "symbol not match", plugin.Symbol(12345), s)
		ztesting.AssertEqual(t, "value not match", "", v)
		ztesting.AssertEqualErr(t, "error not match", zplugin.ErrAssert, err)
	})
	t.Run("lookup success", func(t *testing.T) {
		done := zplugin.WithTestOpenFunc(func(path string) (zplugin.Lookupper, error) {
			return &testPlugin{symbols: map[string]plugin.Symbol{
				"MyVar": 12345,
			}}, nil
		})
		defer done()
		s, v, err := zplugin.Use[int]("test.so", "MyVar")
		ztesting.AssertEqual(t, "symbol not match", plugin.Symbol(12345), s)
		ztesting.AssertEqual(t, "value not match", 12345, v)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
	})
}
