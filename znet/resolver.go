package znet

import (
	"cmp"
	"context"
	"fmt"
	"net"
	"slices"
	"sync"
)

// IPResolver resolves host to IP addresses.
type IPResolver struct {
	// Network is the network name.
	// Network must be one of "ip", "ip4" or "ip6".
	// If not set, "ip" is used.
	Network string
	// Host is the host name resolved to ip addresses.
	Host string
	// LookupIP is the optional function to lookup IP addresses from host name.
	// If not set, net.DefaultResolver.LookupIP is used.
	LookupIP func(ctx context.Context, network, host string) ([]net.IP, error)

	mu      sync.RWMutex
	ips     []string
	current int
}

// IPs returns the list of ip addresses.
func (r *IPResolver) IPs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ips
}

// Next returns the next ip address selected with round-robin algorithm.
// Next is safe for concurrent call.
// An empty string "" is returned when no addresses available.
func (r *IPResolver) Next() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.ips) == 0 {
		return ""
	}
	r.current += 1
	if r.current >= len(r.ips) {
		r.current = 0
	}
	return r.ips[r.current]
}

// Resolve resolves host to ip addresses.
// Resolve is safe for concurrent call.
func (r *IPResolver) Resolve(ctx context.Context) error {
	lookup := r.LookupIP
	if lookup == nil {
		lookup = net.DefaultResolver.LookupIP
	}
	ips, err := lookup(ctx, cmp.Or(r.Network, "ip"), r.Host)
	if err != nil {
		return err
	}
	addrs := make([]string, 0, len(ips))
	for _, ip := range ips {
		addrs = append(addrs, ip.String())
	}
	slices.Sort(addrs) // Sort to keep the order.
	fmt.Println(addrs)
	r.mu.Lock()
	r.ips = addrs
	r.mu.Unlock()
	return nil
}
