package znet

import (
	"net/netip"
)

// trieNode is the node for trie trees.
// Here, trie tree is for ip whitelist and blacklist.
// See https://en.wikipedia.org/wiki/Trie.
type trieNode struct {
	children [2]*trieNode
	isSubnet bool
}

// rootNote is the root node for a trie tree.
type rootNode trieNode

func (n *rootNode) add(bits int, ip []byte) {
	node := (*trieNode)(n)
	count := 0
Loop:
	for _, b := range ip {
		for i := range 8 { // Loop over 8 bits.
			bit := (b >> (7 - i)) & 1
			if node.children[bit] == nil {
				node.children[bit] = &trieNode{}
			}
			node = node.children[bit]
			if count++; count >= bits {
				break Loop
			}
		}
	}
	node.isSubnet = true
}

func (n *rootNode) contains(ip []byte) bool {
	node := (*trieNode)(n)
	for _, b := range ip {
		for i := range 8 { // Loop over 8 bits.
			if node == nil {
				return false
			}
			if node.isSubnet {
				return true
			}
			bit := (b >> (7 - i)) & 1
			node = node.children[bit]
		}
	}
	return node != nil && node.isSubnet
}

// allowList is the IPv4 and IPv6 allow and disallow list.
// All lists have their own trie tree.
type allowList struct {
	allowV4    *rootNode
	disallowV4 *rootNode
	allowV6    *rootNode
	disallowV6 *rootNode
}

// Allow registers network addresses to the allow list.
// Given ip addresses are parsed with [net/netip.ParsePrefix].
// When one of the ips is invalid, it returns an error
// without adding any given ips.
// CIDR "/0" matches to all IPs and "<IPv4>/32" or
// "<IPv6>/128" matches to the only specified IP.
// Invalid addresses such as "/16" (CIDR only) or "127.0.0.1" (address only)
// will results in an error because they are not accepted by [net/netip.ParsePrefix].
func (l *allowList) Allow(ps ...string) error {
	pfs := make([]netip.Prefix, 0, len(ps))
	for _, p := range ps {
		pf, err := netip.ParsePrefix(p)
		if err != nil {
			return err
		}
		pfs = append(pfs, pf)
	}
	l.AllowPrefix(pfs...)
	return nil
}

// AllowPrefix registers network addresses to the allow list.
// Invalid prefixes are returned with and error and only valid
// prefixes are registered to the allow list.
// CIDR "/0" matches to all IPs and "<IPv4>/32" or
// "<IPv6>/128" matches to the only specified IP.
// Invalid addresses such as "/16" (CIDR only) or "127.0.0.1" (address only)
// will results in an error because they are not accepted by [net/netip.ParsePrefix].
func (l *allowList) AllowPrefix(ps ...netip.Prefix) {
	for _, p := range ps {
		addr := p.Addr()
		switch {
		case addr.Is4():
			v4 := addr.As4()
			l.allowV4.add(p.Bits(), v4[:])
		case addr.Is6():
			v6 := addr.As16()
			l.allowV6.add(p.Bits(), v6[:])
		}
	}
}

// Disallow registers network addresses to the disallow list.
// Given ip addresses are parsed with [net/netip.ParsePrefix].
// When one of the ips is invalid, it returns an error
// without adding any given ips.
// Note that CIDR "/0" matches to all IPs and "<IPv4>/32" or
// "<IPv6>/128" matches to the only specified IP.
func (l *allowList) Disallow(ps ...string) error {
	pfs := make([]netip.Prefix, 0, len(ps))
	for _, p := range ps {
		pf, err := netip.ParsePrefix(p)
		if err != nil {
			return err
		}
		pfs = append(pfs, pf)
	}
	l.DisallowPrefix(pfs...)
	return nil
}

// DisallowPrefix registers network addresses to the disallow list.
// Invalid prefixes are returned with and error and only valid
// prefixes are registered to the disallow list.
// Note that CIDR "/0" matches to all IPs and "<IPv4>/32" or
// "<IPv6>/128" matches to the only specified IP.
func (l *allowList) DisallowPrefix(ps ...netip.Prefix) {
	for _, p := range ps {
		addr := p.Addr()
		switch {
		case addr.Is4():
			addr := addr.As4()
			l.disallowV4.add(p.Bits(), addr[:])
		case addr.Is6():
			addr := addr.As16()
			l.disallowV6.add(p.Bits(), addr[:])
		}
	}
}

// NewWhiteList returns a new instance of [WhiteList].
// [WhiteList] is for filtering IPv4 and IPv6 addresses
// using network addresses in whitelist.
// See the comments on [WhiteList].
func NewWhiteList() *WhiteList {
	return &WhiteList{
		allowList: &allowList{
			allowV4:    &rootNode{},
			disallowV4: &rootNode{},
			allowV6:    &rootNode{},
			disallowV6: &rootNode{},
		},
	}
}

// NewBlackList returns a new instance of [BlackList].
// [BlackList] is for filtering IPv4 and IPv6 addresses
// using network addresses in blacklist.
// See the comments on [BlackList].
func NewBlackList() *BlackList {
	return &BlackList{
		allowList: &allowList{
			allowV4:    &rootNode{},
			disallowV4: &rootNode{},
			allowV6:    &rootNode{},
			disallowV6: &rootNode{},
		},
	}
}

// WhiteList is the IP whitelist.
// Disallow list is always prior to the allow list.
// Use [NewWhiteList] to create a new instance of WhiteList.
type WhiteList struct {
	*allowList
}

// Allowed returns if the ip is allowed by the whitelist.
// Both IPv4 and IPv6 are accepted.
// It returns false when the ip is not valid ip address.
// Given ip is parsed with [net/netip.ParseAddr].
// For whitelist, disallow list is always prior to the allow list.
// It returns false for addresses that are not contained in
// both allow list and disallow list.
func (l *WhiteList) Allowed(ip string) bool {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	return l.AllowedAddr(addr)
}

// AllowedAddr returns if the addr is allowed by the whitelist.
// Both IPv4 and IPv6 are accepted.
// It returns false when the addr is not valid ip address.
// For whitelist, disallow list is always prior to the allow list.
// It returns false for addresses that are not contained in
// both allow list and disallow list.
func (l *WhiteList) AllowedAddr(addr netip.Addr) bool {
	switch {
	case addr.Is4():
		v4 := addr.As4()
		if l.disallowV4.contains(v4[:]) {
			return false
		}
		return l.allowV4.contains(v4[:])
	case addr.Is6():
		v6 := addr.As16()
		if l.disallowV6.contains(v6[:]) {
			return false
		}
		return l.allowV6.contains(v6[:])
	}
	return false
}

// BlackList is the IP blacklist.
// Allow list is always prior to the disallow list.
// Use [NewBlackList] to create a new instance of BlackList.
type BlackList struct {
	*allowList
}

// Allowed returns if the ip is allowed by the blacklist.
// Both IPv4 and IPv6 are accepted.
// It returns false when the ip is not valid ip address.
// Given ip is parsed with [net/netip.ParseAddr].
// For blacklist, allow list is always prior to the disallow list.
// It returns true for addresses that are not contained in
// both allow list and disallow list.
func (l *BlackList) Allowed(ip string) bool {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	return l.AllowedAddr(addr)
}

// AllowedAddr returns if the addr is allowed by the blacklist.
// Both IPv4 and IPv6 are accepted.
// It returns false when the addr is not valid ip address.
// For blacklist, allow list is always prior to the disallow list.
// It returns true for addresses that are not contained in
// both allow list and disallow list.
func (l *BlackList) AllowedAddr(addr netip.Addr) bool {
	switch {
	case addr.Is4():
		v4 := addr.As4()
		if l.allowV4.contains(v4[:]) {
			return true
		}
		return !l.disallowV4.contains(v4[:])
	case addr.Is6():
		v6 := addr.As16()
		if l.allowV6.contains(v6[:]) {
			return true
		}
		return !l.disallowV6.contains(v6[:])
	}
	return false
}
