//go:build !linux && !windows

package zsyscall

import (
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
	ztesting.AssertEqual(t, "number of controllers not match", 0, len(cs))
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
	ztesting.AssertEqual(t, "number of controllers not match", 0, len(cs))
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
	ztesting.AssertEqual(t, "number of controllers not match", 0, len(cs))
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
