package znet

import (
	"fmt"
)

func Example() {
	// https://en.wikipedia.org/wiki/IPv4
	// https://en.wikipedia.org/wiki/IPv6
	wl := NewWhiteList()
	err := wl.Allow("127.0.0.0/8", "192.168.0.0/16", "fd00:0:0::/48")
	if err != nil {
		panic(err)
	}

	targetIPs := []string{
		"127.0.0.1",     // OK
		"192.168.1.1",   // OK
		"126.0.0.1",     // NG
		"192.169.1.1",   // NG
		"fd00:0:0::1",   // OK
		"fd00:0:0:1::1", // OK
		"fd00:0:1::1",   // NG
		"fc00:0:0::1",   // NG

	}
	for _, ip := range targetIPs {
		fmt.Printf("%s is allowed? %v\n", ip, wl.Allowed(ip))
	}
	// Output:
	// 127.0.0.1 is allowed? true
	// 192.168.1.1 is allowed? true
	// 126.0.0.1 is allowed? false
	// 192.169.1.1 is allowed? false
	// fd00:0:0::1 is allowed? true
	// fd00:0:0:1::1 is allowed? true
	// fd00:0:1::1 is allowed? false
	// fc00:0:0::1 is allowed? false
}

func ExampleWhiteList_ipv4() {
	// https://en.wikipedia.org/wiki/IPv4
	prefixes := []string{
		"10.0.0.0/8",     // 10.0.0.0–10.255.255.255 Private network
		"100.64.0.0/10",  // 100.64.0.0–100.127.255.255 Private network
		"127.0.0.0/8",    // 127.0.0.0–127.255.255.255 Host
		"172.16.0.0/12",  // 172.16.0.0–172.31.255.255 Private network
		"192.0.0.0/24",   // 192.0.0.0–192.0.0.255 Private network
		"192.168.0.0/16", // 192.168.0.0–192.168.255.255 Private network
		"198.18.0.0/15",  // 198.18.0.0–198.19.255.255 Private network
	}

	wl := NewWhiteList()
	_ = wl.Allow(prefixes...)      // Ignore error.
	_ = wl.Disallow("10.1.2.3/32") // Ignore error.

	targetIPs := []string{
		"10.255.255.1",  // OK
		"127.0.0.1",     // OK
		"192.168.1.2",   // OK
		"192.88.10.20",  // NG
		"224.10.20.30",  // NG
		"255.255.10.20", // NG
		"10.1.2.3",      // NG (Disallow list is prior to allow list)
	}
	for _, ip := range targetIPs {
		fmt.Printf("%s is allowed? %v\n", ip, wl.Allowed(ip))
	}
	// Output:
	// 10.255.255.1 is allowed? true
	// 127.0.0.1 is allowed? true
	// 192.168.1.2 is allowed? true
	// 192.88.10.20 is allowed? false
	// 224.10.20.30 is allowed? false
	// 255.255.10.20 is allowed? false
	// 10.1.2.3 is allowed? false
}

func ExampleBlackList_ipv4() {
	// https://en.wikipedia.org/wiki/IPv4
	prefixes := []string{
		"10.0.0.0/8",     // 10.0.0.0–10.255.255.255 Private network
		"100.64.0.0/10",  // 100.64.0.0–100.127.255.255 Private network
		"127.0.0.0/8",    // 127.0.0.0–127.255.255.255 Host
		"172.16.0.0/12",  // 172.16.0.0–172.31.255.255 Private network
		"192.0.0.0/24",   // 192.0.0.0–192.0.0.255 Private network
		"192.168.0.0/16", // 192.168.0.0–192.168.255.255 Private network
		"198.18.0.0/15",  // 198.18.0.0–198.19.255.255 Private network
	}

	bl := NewBlackList()
	_ = bl.Disallow(prefixes...) // Ignore error.
	_ = bl.Allow("10.1.2.3/32")  // Ignore error.

	targetIPs := []string{
		"10.255.255.1",  // NG
		"127.0.0.1",     // NG
		"192.168.1.2",   // NG
		"192.88.10.20",  // OK
		"224.10.20.30",  // OK
		"255.255.10.20", // OK
		"10.1.2.3",      // OK (Allow list is prior to disallow list)
	}
	for _, ip := range targetIPs {
		fmt.Printf("%s is allowed? %v\n", ip, bl.Allowed(ip))
	}
	// Output:
	// 10.255.255.1 is allowed? false
	// 127.0.0.1 is allowed? false
	// 192.168.1.2 is allowed? false
	// 192.88.10.20 is allowed? true
	// 224.10.20.30 is allowed? true
	// 255.255.10.20 is allowed? true
	// 10.1.2.3 is allowed? true
}
