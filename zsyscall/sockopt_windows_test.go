//go:build windows

package zsyscall

import (
	"io/fs"
	"syscall"
	"testing"

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
	ztesting.AssertEqual(t, "number of controllers not match", 6, len(cs))
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
	ztesting.AssertEqual(t, "number of controllers not match", 1, len(cs))
}

func TestSockIPV6Option_Controllers(t *testing.T) {
	t.Parallel()
	opt := &SockIPV6Option{
		V6Only: true,
	}
	cs := opt.Controllers()
	ztesting.AssertEqual(t, "number of controllers not match", 0, len(cs))
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
	ztesting.AssertEqual(t, "number of controllers not match", 1, len(cs))
}

func TestSockUDPOption_Controllers(t *testing.T) {
	t.Parallel()
	opt := &SockUDPOption{
		CORK:    true,
		Segment: 10,
		GRO:     true,
	}
	cs := opt.Controllers()
	ztesting.AssertEqual(t, "number of controllers not match", 0, len(cs))
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
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", SO_DEBUG, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soDebug(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
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
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_KEEPALIVE, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soKeepAlive(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
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
		setsockoptLinger = func(fd syscall.Handle, level, opt int, l *syscall.Linger) (err error) {
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
		setsockoptLinger = func(fd syscall.Handle, level, opt int, l *syscall.Linger) (err error) {
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
		setsockoptLinger = func(fd syscall.Handle, level, opt int, l *syscall.Linger) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soLinger(9)
		want := &SocketError{Opts: "SOL_SOCKET.SO_LINGER"}
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
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_RCVBUF, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soRcvbuf(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soRcvbuf(1)
		want := &SocketError{Opts: "SOL_SOCKET.SO_RCVBUF"}
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
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_SNDBUF, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soSndbuf(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soSndbuf(1)
		want := &SocketError{Opts: "SOL_SOCKET.SO_SNDBUF"}
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
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.SOL_SOCKET, level)
			ztesting.AssertEqual(t, "option not match", syscall.SO_REUSEADDR, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := soReuseaddr(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := soReuseaddr(true)
		want := &SocketError{Opts: "SOL_SOCKET.SO_REUSEADDR"}
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
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_IP, level)
			ztesting.AssertEqual(t, "option not match", syscall.IP_TTL, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := ipTTL(1)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := ipTTL(1)
		want := &SocketError{Opts: "IPPROTO_IP.IP_TTL"}
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
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			ztesting.AssertEqual(t, "level not match", syscall.IPPROTO_TCP, level)
			ztesting.AssertEqual(t, "option not match", syscall.TCP_NODELAY, opt)
			ztesting.AssertEqual(t, "value not match", 1, value)
			return nil
		}
		c := tcpNoDelay(true)
		ztesting.AssertEqualErr(t, "error not match", nil, c(0))
	})
	t.Run("error", func(t *testing.T) {
		setsockoptInt = func(fd syscall.Handle, level, opt, value int) (err error) {
			return fs.ErrClosed // Dummy error.
		}
		c := tcpNoDelay(true)
		want := &SocketError{Opts: "IPPROTO_TCP.TCP_NODELAY"}
		ztesting.AssertEqualErr(t, "error not match", want, c(0))
	})
}
