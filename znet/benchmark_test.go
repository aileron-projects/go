package znet

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func ipv6Prefixes() []string {
	var ps []string
	for range 25 {
		n := 1 + rand.Uint32N(0xffff)
		ps = append(ps, "fd00:0:"+fmt.Sprintf("%04x", n)+"::/48")
	}
	for range 25 {
		n := 1 + rand.Uint32N(0xffff)
		ps = append(ps, "fd00:0:0:"+fmt.Sprintf("%04x", n)+"::/64")
	}
	for range 25 {
		n := 1 + rand.Uint32N(0xffff)
		ps = append(ps, "fd00:0:0:0:0:"+fmt.Sprintf("%04x", n)+"::/96")
	}
	for range 25 {
		n := 1 + rand.Uint32N(0xffff)
		ps = append(ps, "fd00:0:0:0:0:0:0:"+fmt.Sprintf("%04x", n)+"/128")
	}
	ps = append(ps, "fd00:0:0:0:0:0:0:0/128") // This IP is always allowed.
	return ps
}

func BenchmarkWhiteList(b *testing.B) {
	wl := NewWhiteList()
	err := wl.Allow(ipv6Prefixes()...)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for b.Loop() {
		wl.Allowed("fd00:0:0:0:0:0:0:0") // Always allowed. And the worse case.
	}
}
