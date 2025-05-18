package znet_test

import (
	"net/netip"
	"testing"

	"github.com/aileron-projects/go/znet"
	"github.com/aileron-projects/go/ztesting"
)

func TestAllowList(t *testing.T) {
	t.Parallel()
	list := znet.NewWhiteList()
	t.Run("allow fails", func(t *testing.T) {
		err := list.Allow("127.0.0.1") // No CIDR
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("disallow fails", func(t *testing.T) {
		err := list.Disallow("127.0.0.1") // No CIDR
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
}

func TestWhiteList(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		allow, disallow     []string
		allowed, disallowed []string
	}{
		"case01": {
			[]string{}, []string{}, []string{},
			[]string{"0.0.0.0", "127.0.0.1", "196.168.1.1", "::1", "fd00::1", "invalid"},
		},
		"v4 case01": {
			[]string{"127.0.0.0/8"}, []string{"127.0.0.2/32"},
			[]string{"127.0.0.1"}, []string{"0.0.0.0", "126.0.0.1", "128.0.0.1", "196.168.1.1"},
		},
		"v4 case02": {
			[]string{"127.0.0.0/8", "192.168.0.0/16"}, []string{"127.0.0.2/32", "192.168.0.2/32"},
			[]string{"127.0.0.1", "192.168.0.1"},
			[]string{"0.0.0.0", "126.0.0.1", "196.168.0.2", "128.0.0.1", "192.167.0.1", "192.169.0.1"},
		},
		"v6 case01": {
			[]string{"fd00:0:0::/48"}, []string{"fd00:0:0::2/128"},
			[]string{"fd00:0:0::1", "fd00:0:0:1::1"}, []string{"::1", "fd00:0:0::2", "fd00:0:1::1", "fd00:1:0::1", "fc00:0:0::1"},
		},
		"v6 case02": {
			[]string{"fd00:0:0::/48", "fc00:0:0::/48"}, []string{"fd00:0:0::2/128", "fc00:0:0::2/128"},
			[]string{"fd00:0:0::1", "fd00:0:0:1::1", "fc00:0:0::1", "fc00:0:0:1::1"},
			[]string{"::1", "fd00:0:0::2", "fc00:0:0::2", "fd00:0:1::1", "fd00:1:0::1", "fc00:1:0::1"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			l := znet.NewWhiteList()
			l.Allow(tc.allow...)
			l.Disallow(tc.disallow...)
			for _, ip := range tc.allowed {
				allowed := l.Allowed(ip)
				t.Log(ip, allowed)
				ztesting.AssertEqual(t, "ip should be allowed", true, allowed)
			}
			for _, ip := range tc.disallowed {
				allowed := l.Allowed(ip)
				t.Log(ip, allowed)
				ztesting.AssertEqual(t, "ip should not be allowed", false, allowed)
			}
		})
	}
}

func TestBlackList(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		disallow, allow     []string
		disallowed, allowed []string
	}{
		"case01": {
			[]string{}, []string{}, []string{"invalid"},
			[]string{"0.0.0.0", "127.0.0.1", "196.168.1.1", "::1", "fd00::1"},
		},
		"v4 case01": {
			[]string{"127.0.0.0/8"}, []string{"127.0.0.2/32"},
			[]string{"127.0.0.1"}, []string{"0.0.0.0", "126.0.0.1", "128.0.0.1", "196.168.1.1"},
		},
		"v4 case02": {
			[]string{"127.0.0.0/8", "192.168.0.0/16"}, []string{"127.0.0.2/32", "192.168.0.2/32"},
			[]string{"127.0.0.1", "192.168.0.1"},
			[]string{"0.0.0.0", "126.0.0.1", "196.168.0.2", "128.0.0.1", "192.167.0.1", "192.169.0.1"},
		},
		"v6 case01": {
			[]string{"fd00:0:0::/48"}, []string{"fd00:0:0::2/128"},
			[]string{"fd00:0:0::1", "fd00:0:0:1::1"}, []string{"::1", "fd00:0:0::2", "fd00:0:1::1", "fd00:1:0::1", "fc00:0:0::1"},
		},
		"v6 case02": {
			[]string{"fd00:0:0::/48", "fc00:0:0::/48"}, []string{"fd00:0:0::2/128", "fc00:0:0::2/128"},
			[]string{"fd00:0:0::1", "fd00:0:0:1::1", "fc00:0:0::1", "fc00:0:0:1::1"},
			[]string{"::1", "fd00:0:0::2", "fc00:0:0::2", "fd00:0:1::1", "fd00:1:0::1", "fc00:1:0::1"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			l := znet.NewBlackList()
			l.Allow(tc.allow...)
			l.Disallow(tc.disallow...)
			for _, ip := range tc.allowed {
				allowed := l.Allowed(ip)
				t.Log(ip, allowed)
				ztesting.AssertEqual(t, "ip should be allowed", true, allowed)
			}
			for _, ip := range tc.disallowed {
				allowed := l.Allowed(ip)
				t.Log(ip, allowed)
				ztesting.AssertEqual(t, "ip should not be allowed", false, allowed)
			}
		})
	}
}

func TestAllowAddr(t *testing.T) {
	t.Parallel()
	t.Run("blacklist zero addr", func(t *testing.T) {
		l := znet.NewBlackList()
		allowed := l.AllowedAddr(netip.Addr{})
		ztesting.AssertEqual(t, "ip should not be allowed", false, allowed)
	})
	t.Run("whitelist zero addr", func(t *testing.T) {
		l := znet.NewWhiteList()
		allowed := l.AllowedAddr(netip.Addr{})
		ztesting.AssertEqual(t, "ip should not be allowed", false, allowed)
	})
}
