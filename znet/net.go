package znet

import "strings"

const (
	NetIP         = "ip"
	NetIP4        = "ip4"
	NetIP6        = "ip6"
	NetTCP        = "tcp"
	NetTCP4       = "tcp4"
	NetTCP6       = "tcp6"
	NetUDP        = "udp"
	NetUDP4       = "udp4"
	NetUDP6       = "udp6"
	NetUnix       = "unix"
	NetUnixgram   = "unixgram"
	NetUnixpacket = "unixpacket"
)

// ParseNetAddr parses network and address from addr.
// Given addr should be in "<ADDRESS>" or "<NETWORK>://<ADDRESS>" format.
// Returned network will be empty when the network is unknown of not found.
// For unix domain sockets, use "unix:///var/run/example.sock" for path name socket
// and use "unix://@example" for abstract socket.
// Note that the address format of unix sockets depends on applications.
// For examples, curl has following interface for unix sockets.
//
//	curl --unix-socket "/var/run/example.sock" http://foo.bar.com
//	curl --abstract-unix-socket "example" http://foo.bar.com
//
// Following networks are supported and compatible with [net.Dial] and [net.Listen].
//
//   - <ADDRESS>
//   - ip://<ADDRESS>
//   - ip4://<ADDRESS>
//   - ip6://<ADDRESS>
//   - tcp://<ADDRESS>
//   - tcp4://<ADDRESS>
//   - tcp6://<ADDRESS>
//   - udp://<ADDRESS>
//   - udp4://<ADDRESS>
//   - udp6://<ADDRESS>
//   - unix://<ADDRESS>
//   - unixgram://<ADDRESS>
//   - unixpacket://<ADDRESS>
func ParseNetAddr(addr string) (network, address string) {
	before, after, found := strings.Cut(addr, "://")
	if !found {
		return "", addr
	}
	switch before {
	case "ip", "ip4", "ip6",
		"tcp", "tcp4", "tcp6",
		"udp", "udp4", "udp6",
		"unix", "unixgram", "unixpacket":
		return before, after
	}
	return "", addr
}
