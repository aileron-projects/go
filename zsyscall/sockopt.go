package zsyscall

import (
	"cmp"
	"errors"
	"syscall"
	"time"
)

const (
	SockOptSO uint = 1 << iota
	SockOptIP
	SockOptIPV6
	SockOptTCP
	SockOptUDP
)

var (
	setsockoptInt     = syscall.SetsockoptInt
	setsockoptLinger  = syscall.SetsockoptLinger
	setsockoptTimeval = syscall.SetsockoptTimeval
)

// Controller is a function type that controls socket.
type Controller func(fd uintptr) error

// ControlFunc is the function type that handle RawConn.
type ControlFunc func(string, string, syscall.RawConn) error

// SockOption is aggregated socket options.
//   - https://man7.org/linux/man-pages/man7/socket.7.html
//   - https://man7.org/linux/man-pages/man7/ip.7.html
//   - https://man7.org/linux/man-pages/man7/ipv6.7.html
//   - https://man7.org/linux/man-pages/man7/tcp.7.html
//   - https://man7.org/linux/man-pages/man7/udp.7.html
type SockOption struct {
	SO   *SockSOOption
	IP   *SockIPOption
	IPV6 *SockIPV6Option
	TCP  *SockTCPOption
	UDP  *SockUDPOption
}

// ControlFunc returns a control function from options.
// Only specified options are enabled.
// For example, pass the [SockOptSO] and [SockOptIP] to enable
// both options as shown below.
//
//	ControlFunc(zsyscall.SockOptSO|zsyscall.SockOptIP)
func (o *SockOption) ControlFunc(opts uint) ControlFunc {
	if o == nil {
		return nil
	}
	var cs []Controller
	if o.SO != nil && opts&SockOptSO != 0 {
		cs = append(cs, o.SO.Controllers()...)
	}
	if o.IP != nil && opts&SockOptIP != 0 {
		cs = append(cs, o.IP.Controllers()...)
	}
	if o.IPV6 != nil && opts&SockOptIPV6 != 0 {
		cs = append(cs, o.IPV6.Controllers()...)
	}
	if o.TCP != nil && opts&SockOptTCP != 0 {
		cs = append(cs, o.TCP.Controllers()...)
	}
	if o.UDP != nil && opts&SockOptUDP != 0 {
		cs = append(cs, o.UDP.Controllers()...)
	}
	if len(cs) == 0 {
		return nil
	}
	return controllers(cs).control
}

type controllers []Controller

func (cs controllers) control(network, address string, conn syscall.RawConn) (err error) {
	ctrErr := conn.Control(func(fd uintptr) {
		for _, c := range cs {
			if e := c(fd); e != nil {
				err = e
				return // Fail fast.
			}
		}
	})
	return cmp.Or(err, ctrErr)
}

// SockSOOption is socket options for SOL_SOCKET level.
//   - https://man7.org/linux/man-pages/man7/socket.7.html
//   - https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/sys_socket.h.html
type SockSOOption struct {
	BindToIFindex      int           // SO_BINDTOIFINDEX
	BindToDevice       string        // SO_BINDTODEVICE
	Debug              bool          // SO_DEBUG
	KeepAlive          bool          // SO_KEEPALIVE
	Linger             int32         // SO_LINGER
	Mark               int           // SO_MARK (since Linux 2.6.25)
	ReceiveBuffer      int           // SO_RCVBUF
	ReceiveBufferForce int           // SO_RCVBUFFORCE (since Linux 2.6.14)
	ReceiveTimeout     time.Duration // SO_RCVTIMEO
	SendTimeout        time.Duration // SO_SNDTIMEO
	ReuseAddr          bool          // SO_REUSEADDR
	ReusePort          bool          // SO_REUSEPORT (since Linux 3.9)
	SendBuffer         int           // SO_SNDBUF
	SendBufferForce    int           // SO_SNDBUFFORCE (since Linux 2.6.14)
}

// SockIPOption is socket options for IPPROTO_IP level.
//   - https://man7.org/linux/man-pages/man7/ip.7.html
type SockIPOption struct {
	BindAddressNoPort   bool   // IP_BIND_ADDRESS_NO_PORT (since Linux 4.2)
	FreeBind            bool   // IP_FREEBIND (since Linux 2.4)
	LocalPortRangeUpper uint16 // IP_LOCAL_PORT_RANGE (since Linux 6.3)
	LocalPortRangeLower uint16 // IP_LOCAL_PORT_RANGE (since Linux 6.3)
	Transparent         bool   // IP_TRANSPARENT (since Linux 2.6.24)
	TTL                 int    // IP_TTL (since Linux 1.0)
}

// SockIPV6Option is socket options for IPPROTO_IPV6 level.
//   - https://man7.org/linux/man-pages/man7/ipv6.7.html
type SockIPV6Option struct {
	V6Only bool // IPV6_V6ONLY (since Linux 2.4.21 and 2.6)
}

// SockTCPOption is socket options for IPPROTO_TCP level.
//   - https://man7.org/linux/man-pages/man7/tcp.7.html
type SockTCPOption struct {
	CORK            bool  // TCP_CORK (since Linux 2.2)
	DeferAccept     int   // TCP_DEFER_ACCEPT (since Linux 2.4)
	KeepCount       int   // TCP_KEEPCNT (since Linux 2.4)
	KeepIdle        int   // TCP_KEEPIDLE (since Linux 2.4)
	KeepInterval    int   // TCP_KEEPINTVL (since Linux 2.4)
	Linger2         int32 // TCP_LINGER2 (since Linux 2.4)
	MaxSegment      int   // TCP_MAXSEG
	NoDelay         bool  // TCP_NODELAY
	QuickAck        bool  // TCP_QUICKACK (since Linux 2.4.4)
	SynCount        int   // TCP_SYNCNT (since Linux 2.4)
	UserTimeout     int   // TCP_USER_TIMEOUT (since Linux 2.6.37)
	WindowClamp     int   // TCP_WINDOW_CLAMP (since Linux 2.4)
	FastOpen        bool  // TCP_FASTOPEN (since Linux 3.6)
	FastOpenConnect bool  // TCP_FASTOPEN_CONNECT (since Linux 4.11)
}

// SockUDPOption is socket options for IPPROTO_UDP level.
//   - https://man7.org/linux/man-pages/man7/udp.7.html
type SockUDPOption struct {
	CORK    bool // UDP_CORK (since Linux 2.5.44)
	Segment int  // UDP_SEGMENT (since Linux 4.18)
	GRO     bool // UDP_GRO (since Linux 5.0)
}

func appendNonNil(arr []Controller, target Controller) []Controller {
	if target == nil {
		return arr
	}
	return append(arr, target)
}

// SocketError is an error type that tells an application
// of a socket option failed.
type SocketError struct {
	Err  error
	Opts string
}

func (e *SocketError) Error() string {
	msg := "zsyscall: fail to apply socket option " + e.Opts
	if e.Err != nil {
		msg += " [" + e.Err.Error() + "]"
	}
	return msg
}

func (e *SocketError) Is(err error) bool {
	if err == nil || e == nil {
		return e == err
	}
	for err != nil {
		ee, ok := err.(*SocketError)
		if ok {
			return e.Opts == ee.Opts
		}
		err = errors.Unwrap(err)
	}
	return false
}
