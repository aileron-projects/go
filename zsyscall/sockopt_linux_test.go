//go:build linux

package zsyscall

import (
	"io/fs"
	"syscall"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztesting"
)

func TestSockSOOption_Controllers(t *testing.T) {
	t.Parallel()
	opt := &SockSOOption{
		BindToIFindex:      10,
		BindToDevice:       "eth0",
		Debug:              true,
		KeepAlive:          true,
		Linger:             11,
		Mark:               12,
		ReceiveBuffer:      13,
		ReceiveBufferForce: 14,
		ReceiveTimeout:     15,
		SendTimeout:        16,
		ReuseAddr:          true,
		ReusePort:          true,
		SendBuffer:         17,
		SendBufferForce:    18,
	}
	cs := opt.Controllers()
	ztesting.AssertEqual(t, "number of controllers not match", 14, len(cs))
}

func TestSockIPOption_Controllers(t *testing.T) {
	t.Parallel()
	opt := &SockIPOption{
		BindAddressNoPort:   true,
		FreeBind:            true,
		LocalPortRangeUpper: 10,
		LocalPortRangeLower: 11,
		Transparent:         true,
		TTL:                 12,
	}
	cs := opt.Controllers()
	ztesting.AssertEqual(t, "number of controllers not match", 5, len(cs))
}

func TestSockIPV6Option_Controllers(t *testing.T) {
	t.Parallel()
	opt := &SockIPV6Option{
		V6Only: true,
	}
	cs := opt.Controllers()
	ztesting.AssertEqual(t, "number of controllers not match", 1, len(cs))
}

func TestSockTCPOption_Controllers(t *testing.T) {
	t.Parallel()
	opt := &SockTCPOption{
		CORK:            true,
		DeferAccept:     10,
		KeepCount:       11,
		KeepIdle:        12,
		KeepInterval:    13,
		Linger2:         14,
		MaxSegment:      15,
		NoDelay:         true,
		QuickAck:        true,
		SynCount:        16,
		UserTimeout:     17,
		WindowClamp:     18,
		FastOpen:        true,
		FastOpenConnect: true,
	}
	cs := opt.Controllers()
	ztesting.AssertEqual(t, "number of controllers not match", 14, len(cs))
}

func TestSockUDPOption_Controllers(t *testing.T) {
	t.Parallel()
	opt := &SockUDPOption{
		CORK:    true,
		Segment: 10,
		GRO:     true,
	}
	cs := opt.Controllers()
	ztesting.AssertEqual(t, "number of controllers not match", 3, len(cs))
}

func TestSoBindToIFindex(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soBindToIFindex(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", SO_BINDTOIFINDEX, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soBindToIFindex(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soBindToIFindex(1)
		want := &SocketError{Opts: "SOL_SOCKET.SO_BINDTOIFINDEX"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoBindToDevice(t *testing.T) {
	defer func() {
		setsockoptString = syscall.SetsockoptString // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soBindToDevice("")
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptString = func(fd int, level, opt int, s string) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_BINDTODEVICE, opt)
			ztesting.AssertEqual(t, "value not match", "eth0", s)
			return nil
		}
		c := soBindToDevice("eth0")
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptString = func(fd int, level, opt int, s string) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soBindToDevice("eth0")
		want := &SocketError{Opts: "SOL_SOCKET.SO_BINDTODEVICE"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoDebug(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soDebug(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_DEBUG, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soDebug(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soDebug(true)
		want := &SocketError{Opts: "SOL_SOCKET.SO_DEBUG"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoKeepAlive(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soKeepAlive(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_KEEPALIVE, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soKeepAlive(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soKeepAlive(true)
		want := &SocketError{Opts: "SOL_SOCKET.SO_KEEPALIVE"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoLinger(t *testing.T) {
	defer func() {
		setsockoptLinger = syscall.SetsockoptLinger // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soLinger(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("on", func(t *testing.T) {
		setsockoptLinger = func(fd int, level, opt int, l *syscall.Linger) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_LINGER, opt)
			ztesting.AssertEqual(t, "onoff not match", 1, l.Onoff)
			ztesting.AssertEqual(t, "linger not match", 9, l.Linger)
			return nil
		}
		c := soLinger(9)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("off", func(t *testing.T) {
		setsockoptLinger = func(fd int, level, opt int, l *syscall.Linger) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_LINGER, opt)
			ztesting.AssertEqual(t, "onoff not match", 0, l.Onoff)
			ztesting.AssertEqual(t, "linger not match", 0, l.Linger)
			return nil
		}
		c := soLinger(-9)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptLinger = func(fd int, level, opt int, l *syscall.Linger) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soLinger(9)
		want := &SocketError{Opts: "SOL_SOCKET.SO_LINGER"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoMark(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soMark(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_MARK, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soMark(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soMark(1)
		want := &SocketError{Opts: "SOL_SOCKET.SO_MARK"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoRcvbuf(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soRcvbuf(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_RCVBUF, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soRcvbuf(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soRcvbuf(1)
		want := &SocketError{Opts: "SOL_SOCKET.SO_RCVBUF"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoRcvbufForce(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soRcvbufForce(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_RCVBUFFORCE, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soRcvbufForce(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soRcvbufForce(1)
		want := &SocketError{Opts: "SOL_SOCKET.SO_RCVBUFFORCE"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoSndbuf(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soSndbuf(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_SNDBUF, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soSndbuf(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soSndbuf(1)
		want := &SocketError{Opts: "SOL_SOCKET.SO_SNDBUF"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoSndbufForce(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soSndbufForce(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_SNDBUFFORCE, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soSndbufForce(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soSndbufForce(1)
		want := &SocketError{Opts: "SOL_SOCKET.SO_SNDBUFFORCE"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoSndtimeo(t *testing.T) {
	defer func() {
		setsockoptTimeval = syscall.SetsockoptTimeval // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soSndtimeo(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptTimeval = func(fd int, level, opt int, tv *syscall.Timeval) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_SNDTIMEO, opt)
			ztesting.AssertEqual(t, "value not match", syscall.NsecToTimeval(int64(time.Second)), *tv)
			return nil
		}
		c := soSndtimeo(time.Second)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptTimeval = func(fd int, level, opt int, tv *syscall.Timeval) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soSndtimeo(time.Second)
		want := &SocketError{Opts: "SOL_SOCKET.SO_SNDTIMEO"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoRcvtimeo(t *testing.T) {
	defer func() {
		setsockoptTimeval = syscall.SetsockoptTimeval // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soRcvtimeo(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptTimeval = func(fd int, level, opt int, tv *syscall.Timeval) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_RCVTIMEO, opt)
			ztesting.AssertEqual(t, "value not match", syscall.NsecToTimeval(int64(time.Second)), *tv)
			return nil
		}
		c := soRcvtimeo(time.Second)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptTimeval = func(fd int, level, opt int, tv *syscall.Timeval) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soRcvtimeo(time.Second)
		want := &SocketError{Opts: "SOL_SOCKET.SO_RCVTIMEO"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoReuseaddr(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soReuseaddr(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_REUSEADDR, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soReuseaddr(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soReuseaddr(true)
		want := &SocketError{Opts: "SOL_SOCKET.SO_REUSEADDR"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestSoReuseport(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := soReuseport(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", SO_REUSEPORT, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soReuseport(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soReuseport(true)
		want := &SocketError{Opts: "SOL_SOCKET.SO_REUSEPORT"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestIPBindAddressNoPort(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := ipBindAddressNoPort(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_IP, level)
			ztesting.AssertEqual(t, "option not match", IP_BIND_ADDRESS_NO_PORT, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := ipBindAddressNoPort(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := ipBindAddressNoPort(true)
		want := &SocketError{Opts: "IPPROTO_IP.IP_BIND_ADDRESS_NO_PORT"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestIPFreeBind(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := ipFreeBind(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_IP, level)
			ztesting.AssertEqual(t, "option not match", syscall.IP_FREEBIND, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := ipFreeBind(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := ipFreeBind(true)
		want := &SocketError{Opts: "IPPROTO_IP.IP_FREEBIND"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestIPLocalPortRange(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := ipLocalPortRange(0, 0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_IP, level)
			ztesting.AssertEqual(t, "option not match", IP_LOCAL_PORT_RANGE, opt)
			ztesting.AssertEqual(t, "value not match", int(uint32(2048)<<16|uint32(2000)), value)
			return nil
		}
		c := ipLocalPortRange(2048, 2000)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := ipLocalPortRange(2048, 2000)
		want := &SocketError{Opts: "IPPROTO_IP.IP_LOCAL_PORT_RANGE"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestIPTransparent(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := ipTransparent(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_IP, level)
			ztesting.AssertEqual(t, "option not match", syscall.IP_TRANSPARENT, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := ipTransparent(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := ipTransparent(true)
		want := &SocketError{Opts: "IPPROTO_IP.IP_TRANSPARENT"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestIPTTL(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := ipTTL(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_IP, level)
			ztesting.AssertEqual(t, "option not match", syscall.IP_TTL, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := ipTTL(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := ipTTL(1)
		want := &SocketError{Opts: "IPPROTO_IP.IP_TTL"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestIPV6V6Only(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := ipv6V6Only(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_IPV6, level)
			ztesting.AssertEqual(t, "option not match", syscall.IPV6_V6ONLY, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := ipv6V6Only(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := ipv6V6Only(true)
		want := &SocketError{Opts: "IPPROTO_IPV6.IPV6_V6ONLY"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPCORK(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpCORK(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_CORK, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpCORK(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpCORK(true)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_CORK"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPDeferAccept(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpDeferAccept(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_DEFER_ACCEPT, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpDeferAccept(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpDeferAccept(1)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_DEFER_ACCEPT"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPKeepCount(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpKeepCount(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_KEEPCNT, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpKeepCount(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpKeepCount(1)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_KEEPCNT"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPKeepIdle(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpKeepIdle(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_KEEPIDLE, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpKeepIdle(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpKeepIdle(1)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_KEEPIDLE"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPKeepInterval(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpKeepInterval(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_KEEPINTVL, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpKeepInterval(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpKeepInterval(1)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_KEEPINTVL"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPLinger2(t *testing.T) {
	defer func() {
		setsockoptLinger = syscall.SetsockoptLinger // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpLinger2(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("on", func(t *testing.T) {
		setsockoptLinger = func(fd int, level, opt int, l *syscall.Linger) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_LINGER2, opt)
			ztesting.AssertEqual(t, "onoff not match", 1, l.Onoff)
			ztesting.AssertEqual(t, "linger not match", 9, l.Linger)
			return nil
		}
		c := tcpLinger2(9)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("off", func(t *testing.T) {
		setsockoptLinger = func(fd int, level, opt int, l *syscall.Linger) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_LINGER2, opt)
			ztesting.AssertEqual(t, "onoff not match", 0, l.Onoff)
			ztesting.AssertEqual(t, "linger not match", 0, l.Linger)
			return nil
		}
		c := tcpLinger2(-9)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptLinger = func(fd int, level, opt int, l *syscall.Linger) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpLinger2(9)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_LINGER2"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPMaxSegment(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpMaxSegment(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_MAXSEG, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpMaxSegment(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpMaxSegment(1)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_MAXSEG"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPNoDelay(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpNoDelay(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_NODELAY, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpNoDelay(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpNoDelay(true)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_NODELAY"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPQuickAck(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpQuickAck(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_QUICKACK, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpQuickAck(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpQuickAck(true)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_QUICKACK"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPSynCount(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpSynCount(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_SYNCNT, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpSynCount(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpSynCount(1)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_SYNCNT"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPUserTimeout(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpUserTimeout(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", TCP_USER_TIMEOUT, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpUserTimeout(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpUserTimeout(1)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_USER_TIMEOUT"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPWindowClamp(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpWindowClamp(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_WINDOW_CLAMP, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpWindowClamp(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpWindowClamp(1)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_WINDOW_CLAMP"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPFastOpen(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpFastOpen(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", TCP_FASTOPEN, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpFastOpen(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpFastOpen(true)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_FASTOPEN"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestTCPFastOpenConnect(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := tcpFastOpenConnect(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", TCP_FASTOPEN_CONNECT, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpFastOpenConnect(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpFastOpenConnect(true)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_FASTOPEN_CONNECT"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestUDPCORK(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := udpCORK(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_UDP, level)
			ztesting.AssertEqual(t, "option not match", UDP_CORK, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := udpCORK(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := udpCORK(true)
		want := &SocketError{Opts: "IPPROTO_UDP.UDP_CORK"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestUDPSegment(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := udpSegment(0)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_UDP, level)
			ztesting.AssertEqual(t, "option not match", UDP_SEGMENT, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := udpSegment(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := udpSegment(1)
		want := &SocketError{Opts: "IPPROTO_UDP.UDP_SEGMENT"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}

func TestUDPGRO(t *testing.T) {
	defer func() {
		setsockoptInt = syscall.SetsockoptInt // Reset
	}()
	t.Run("disabled", func(t *testing.T) {
		c := udpGRO(false)
		ztesting.AssertEqual(t, "controller should be nil", true, c == nil)
	})
	t.Run("enabled", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_UDP, level)
			ztesting.AssertEqual(t, "option not match", UDP_GRO, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := udpGRO(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd int, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := udpGRO(true)
		want := &SocketError{Opts: "IPPROTO_UDP.UDP_GRO"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}
