//go:build linux

package zsyscall

import (
	"cmp"
	"syscall"
	"time"
)

const (
	// Socket options.
	// Values are copied from golang.org/x/sys/unix

	SO_BINDTOIFINDEX = 0x3e
	SO_REUSEPORT     = 0xf
	SO_ZEROCOPY      = 0x3c

	IP_BIND_ADDRESS_NO_PORT = 0x18
	IP_LOCAL_PORT_RANGE     = 0x33

	TCP_FASTOPEN         = 0x17
	TCP_FASTOPEN_CONNECT = 0x1e
	TCP_USER_TIMEOUT     = 0x12

	UDP_CORK    = 0x1
	UDP_GRO     = 0x68
	UDP_SEGMENT = 0x67
)

// setsockoptString is the re-mapped function
// of [syscall.SetsockoptString] for testing.
var setsockoptString = syscall.SetsockoptString

func (c *SockSOOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendNonNil(controllers, soBindToIFindex(c.BindToIFindex))
	controllers = appendNonNil(controllers, soBindToDevice(c.BindToDevice))
	controllers = appendNonNil(controllers, soDebug(c.Debug))
	controllers = appendNonNil(controllers, soKeepAlive(c.KeepAlive))
	controllers = appendNonNil(controllers, soLinger(c.Linger))
	controllers = appendNonNil(controllers, soMark(c.Mark))
	controllers = appendNonNil(controllers, soRcvbuf(c.ReceiveBuffer))
	controllers = appendNonNil(controllers, soRcvbufForce(c.ReceiveBufferForce))
	controllers = appendNonNil(controllers, soSndtimeo(c.SendTimeout))
	controllers = appendNonNil(controllers, soRcvtimeo(c.ReceiveTimeout))
	controllers = appendNonNil(controllers, soReuseaddr(c.ReuseAddr))
	controllers = appendNonNil(controllers, soReuseport(c.ReusePort))
	controllers = appendNonNil(controllers, soSndbuf(c.SendBuffer))
	controllers = appendNonNil(controllers, soSndbufForce(c.SendBufferForce))
	return controllers
}

func soBindToIFindex(index int) Controller {
	if index <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, SO_BINDTOIFINDEX, index); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_BINDTOIFINDEX"}
		}
		return nil
	}
}

func soBindToDevice(value string) Controller {
	if value == "" {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptString(int(fd), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, value); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_BINDTODEVICE"}
		}
		return nil
	}
}

func soDebug(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_DEBUG, 1); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_DEBUG"}
		}
		return nil
	}
}

func soKeepAlive(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_KEEPALIVE"}
		}
		return nil
	}
}

func soLinger(value int32) Controller {
	if value == 0 {
		return nil
	}
	onoff := int32(1)
	if value < 0 {
		onoff = 0
	}
	l := &syscall.Linger{
		Onoff:  onoff,
		Linger: max(0, value),
	}
	return func(fd uintptr) error {
		if err := setsockoptLinger(int(fd), syscall.SOL_SOCKET, syscall.SO_LINGER, l); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_LINGER"}
		}
		return nil
	}
}

func soMark(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, value); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_MARK"}
		}
		return nil
	}
}

func soRcvbuf(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, value); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_RCVBUF"}
		}
		return nil
	}
}

func soRcvbufForce(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUFFORCE, value); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_RCVBUFFORCE"}
		}
		return nil
	}
}

func soSndbuf(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF, value); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_SNDBUF"}
		}
		return nil
	}
}

func soSndbufForce(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUFFORCE, value); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_SNDBUFFORCE"}
		}
		return nil
	}
}

func soSndtimeo(value time.Duration) Controller {
	if value <= 0 {
		return nil
	}
	tv := syscall.NsecToTimeval(int64(value))
	return func(fd uintptr) error {
		if err := setsockoptTimeval(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDTIMEO, &tv); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_SNDTIMEO"}
		}
		return nil
	}
}

func soRcvtimeo(value time.Duration) Controller {
	if value <= 0 {
		return nil
	}
	tv := syscall.NsecToTimeval(int64(value))
	return func(fd uintptr) error {
		if err := setsockoptTimeval(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_RCVTIMEO"}
		}
		return nil
	}
}

func soReuseaddr(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_REUSEADDR"}
		}
		return nil
	}
}

func soReuseport(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.SOL_SOCKET, SO_REUSEPORT, 1); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_REUSEPORT"}
		}
		return nil
	}
}

func (c *SockIPOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendNonNil(controllers, ipBindAddressNoPort(c.BindAddressNoPort))
	controllers = appendNonNil(controllers, ipFreeBind(c.FreeBind))
	controllers = appendNonNil(controllers, ipLocalPortRange(c.LocalPortRangeUpper, c.LocalPortRangeLower))
	controllers = appendNonNil(controllers, ipTransparent(c.Transparent))
	controllers = appendNonNil(controllers, ipTTL(c.TTL))
	return controllers
}

func ipBindAddressNoPort(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IP, IP_BIND_ADDRESS_NO_PORT, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_IP.IP_BIND_ADDRESS_NO_PORT"}
		}
		return nil
	}
}

func ipFreeBind(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_FREEBIND, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_IP.IP_FREEBIND"}
		}
		return nil
	}
}

func ipLocalPortRange(upper, lower uint16) Controller {
	upper = cmp.Or(upper, 0)
	lower = cmp.Or(lower, 0)
	if upper <= 0 && lower <= 0 {
		return nil
	}
	v := uint32(upper)<<16 | uint32(lower)
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IP, IP_LOCAL_PORT_RANGE, int(v)); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_IP.IP_LOCAL_PORT_RANGE"}
		}
		return nil
	}
}

func ipTransparent(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TRANSPARENT, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_IP.IP_TRANSPARENT"}
		}
		return nil
	}
}

func ipTTL(ttl int) Controller {
	if ttl <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TTL, ttl); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_IP.IP_TTL"}
		}
		return nil
	}
}

func (c *SockIPV6Option) Controllers() []Controller {
	var controllers []Controller
	controllers = appendNonNil(controllers, ipv6V6Only(c.V6Only))
	return controllers
}

func ipv6V6Only(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_IPV6.IPV6_V6ONLY"}
		}
		return nil
	}
}

func (c *SockTCPOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendNonNil(controllers, tcpCORK(c.CORK))
	controllers = appendNonNil(controllers, tcpDeferAccept(c.DeferAccept))
	controllers = appendNonNil(controllers, tcpKeepCount(c.KeepCount))
	controllers = appendNonNil(controllers, tcpKeepIdle(c.KeepIdle))
	controllers = appendNonNil(controllers, tcpKeepInterval(c.KeepInterval))
	controllers = appendNonNil(controllers, tcpLinger2(c.Linger2))
	controllers = appendNonNil(controllers, tcpMaxSegment(c.MaxSegment))
	controllers = appendNonNil(controllers, tcpNoDelay(c.NoDelay))
	controllers = appendNonNil(controllers, tcpQuickAck(c.QuickAck))
	controllers = appendNonNil(controllers, tcpSynCount(c.SynCount))
	controllers = appendNonNil(controllers, tcpUserTimeout(c.UserTimeout))
	controllers = appendNonNil(controllers, tcpWindowClamp(c.WindowClamp))
	controllers = appendNonNil(controllers, tcpFastOpen(c.FastOpen))
	controllers = appendNonNil(controllers, tcpFastOpenConnect(c.FastOpenConnect))
	return controllers
}

func tcpCORK(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_CORK, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_CORK"}
		}
		return nil
	}
}

func tcpDeferAccept(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_DEFER_ACCEPT, value); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_DEFER_ACCEPT"}
		}
		return nil
	}
}

func tcpKeepCount(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_KEEPCNT, value); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_KEEPCNT"}
		}
		return nil
	}
}

func tcpKeepIdle(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_KEEPIDLE, value); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_KEEPIDLE"}
		}
		return nil
	}
}

func tcpKeepInterval(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL, value); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_KEEPINTVL"}
		}
		return nil
	}
}

func tcpLinger2(value int32) Controller {
	if value == 0 {
		return nil
	}
	onoff := int32(1)
	if value < 0 {
		onoff = 0
	}
	l := &syscall.Linger{
		Onoff:  onoff,
		Linger: max(0, value),
	}
	return func(fd uintptr) error {
		if err := setsockoptLinger(int(fd), syscall.IPPROTO_TCP, syscall.TCP_LINGER2, l); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_LINGER2"}
		}
		return nil
	}
}

func tcpMaxSegment(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_MAXSEG, value); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_MAXSEG"}
		}
		return nil
	}
}

func tcpNoDelay(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_NODELAY"}
		}
		return nil
	}
}

func tcpQuickAck(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_QUICKACK, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_QUICKACK"}
		}
		return nil
	}
}

func tcpSynCount(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_SYNCNT, value); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_SYNCNT"}
		}
		return nil
	}
}

func tcpUserTimeout(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, TCP_USER_TIMEOUT, value); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_USER_TIMEOUT"}
		}
		return nil
	}
}

func tcpWindowClamp(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_WINDOW_CLAMP, value); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_WINDOW_CLAMP"}
		}
		return nil
	}
}

func tcpFastOpen(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, TCP_FASTOPEN, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_FASTOPEN"}
		}
		return nil
	}
}

func tcpFastOpenConnect(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_TCP, TCP_FASTOPEN_CONNECT, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_FASTOPEN_CONNECT"}
		}
		return nil
	}
}

func (c *SockUDPOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendNonNil(controllers, udpCORK(c.CORK))
	controllers = appendNonNil(controllers, udpSegment(c.Segment))
	controllers = appendNonNil(controllers, udpGRO(c.GRO))
	return controllers
}

func udpCORK(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_UDP, UDP_CORK, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_UDP.UDP_CORK"}
		}
		return nil
	}
}

func udpSegment(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_UDP, UDP_SEGMENT, value); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_UDP.UDP_SEGMENT"}
		}
		return nil
	}
}

func udpGRO(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(int(fd), syscall.IPPROTO_UDP, UDP_GRO, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_UDP.UDP_GRO"}
		}
		return nil
	}
}
