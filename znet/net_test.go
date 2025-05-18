package znet_test

import (
	"testing"

	"github.com/aileron-projects/go/znet"
	"github.com/aileron-projects/go/ztesting"
)

func TestParseNetAddr(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		addr             string
		network, address string
	}{
		"addr case01": {"127.0.0.1", "", "127.0.0.1"},
		"addr case02": {"127.0.0.1:80", "", "127.0.0.1:80"},
		"addr case03": {"localhost:80", "", "localhost:80"},
		"addr case04": {"http://localhost", "", "http://localhost"},
		"ip":          {"ip://address", "ip", "address"},
		"ip4":         {"ip4://address", "ip4", "address"},
		"ip6":         {"ip6://address", "ip6", "address"},
		"tcp":         {"tcp://address", "tcp", "address"},
		"tcp4":        {"tcp4://address", "tcp4", "address"},
		"tcp6":        {"tcp6://address", "tcp6", "address"},
		"udp":         {"udp://address", "udp", "address"},
		"udp4":        {"udp4://address", "udp4", "address"},
		"udp6":        {"udp6://address", "udp6", "address"},
		"unix":        {"unix://address", "unix", "address"},
		"unixgram":    {"unixgram://address", "unixgram", "address"},
		"unixpacket":  {"unixpacket://address", "unixpacket", "address"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			n, a := znet.ParseNetAddr(tc.addr)
			ztesting.AssertEqual(t, "network not match", tc.network, n)
			ztesting.AssertEqual(t, "address not match", tc.address, a)
		})
	}
}
