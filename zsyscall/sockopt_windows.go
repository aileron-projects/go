//go:build windows

package zsyscall

import (
	"syscall"
)

const (
	SO_DEBUG = 0x1
)

func (c *SockSOOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendNonNil(controllers, soDebug(c.Debug))
	controllers = appendNonNil(controllers, soKeepAlive(c.KeepAlive))
	controllers = appendNonNil(controllers, soLinger(c.Linger))
	controllers = appendNonNil(controllers, soRcvbuf(c.ReceiveBuffer))
	controllers = appendNonNil(controllers, soReuseaddr(c.ReuseAddr))
	controllers = appendNonNil(controllers, soSndbuf(c.SendBuffer))
	return controllers
}

func soDebug(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, SO_DEBUG, 1); err != nil {
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
		if err := setsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_KEEPALIVE"}
		}
		return nil
	}
}

// soLinger returns a controller to set socket option.
//   - value=0: return nil
//   - value>0: enable linger with the value.
//   - value<0: disable linger.
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
		if err := setsockoptLinger(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_LINGER, l); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_LINGER"}
		}
		return nil
	}
}

func soRcvbuf(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, value); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_RCVBUF"}
		}
		return nil
	}
}

func soSndbuf(value int) Controller {
	if value <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF, value); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_SNDBUF"}
		}
		return nil
	}
}

func soReuseaddr(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
			return &SocketError{Err: err, Opts: "SOL_SOCKET.SO_REUSEADDR"}
		}
		return nil
	}
}

func (c *SockIPOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendNonNil(controllers, ipTTL(c.TTL))
	return controllers
}

func ipTTL(ttl int) Controller {
	if ttl <= 0 {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.IPPROTO_IP, syscall.IP_TTL, ttl); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_IP.IP_TTL"}
		}
		return nil
	}
}

func (c *SockIPV6Option) Controllers() []Controller {
	return nil
}

func (c *SockTCPOption) Controllers() []Controller {
	var controllers []Controller
	controllers = appendNonNil(controllers, tcpNoDelay(c.NoDelay))
	return controllers
}

func tcpNoDelay(enabled bool) Controller {
	if !enabled {
		return nil
	}
	return func(fd uintptr) error {
		if err := setsockoptInt(syscall.Handle(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1); err != nil {
			return &SocketError{Err: err, Opts: "IPPROTO_TCP.TCP_NODELAY"}
		}
		return nil
	}
}

func (c *SockUDPOption) Controllers() []Controller {
	return nil
}
