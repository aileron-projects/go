package zsyscall

import (
	"io"
	"syscall"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestSockOption_ControlFunc(t *testing.T) {
	t.Parallel()
	t.Run("nil option", func(t *testing.T) {
		var o *SockOption
		f := o.ControlFunc(0)
		ztesting.AssertEqual(t, "control func is not nil", true, f == nil)
	})
	t.Run("non-nil option", func(t *testing.T) {
		o := &SockOption{}
		f := o.ControlFunc(0)
		ztesting.AssertEqual(t, "control func is not nil", true, f == nil)
	})
	t.Run("non-nil option", func(t *testing.T) {
		o := &SockOption{}
		f := o.ControlFunc(0)
		ztesting.AssertEqual(t, "control func is not nil", true, f == nil)
	})
	t.Run("so option", func(t *testing.T) {
		o := &SockOption{SO: &SockSOOption{}}
		f := o.ControlFunc(SockOptSO)
		ztesting.AssertEqual(t, "control func is not nil", true, f == nil)
	})
	t.Run("ip option", func(t *testing.T) {
		o := &SockOption{IP: &SockIPOption{}}
		f := o.ControlFunc(SockOptIP)
		ztesting.AssertEqual(t, "control func is not nil", true, f == nil)
	})
	t.Run("ipv6 option", func(t *testing.T) {
		o := &SockOption{IPV6: &SockIPV6Option{}}
		f := o.ControlFunc(SockOptIPV6)
		ztesting.AssertEqual(t, "control func is not nil", true, f == nil)
	})
	t.Run("tcp option", func(t *testing.T) {
		o := &SockOption{TCP: &SockTCPOption{}}
		f := o.ControlFunc(SockOptTCP)
		ztesting.AssertEqual(t, "control func is not nil", true, f == nil)
	})
	t.Run("udp option", func(t *testing.T) {
		o := &SockOption{UDP: &SockUDPOption{}}
		f := o.ControlFunc(SockOptUDP)
		ztesting.AssertEqual(t, "control func is not nil", true, f == nil)
	})
}

type testRawConn struct {
	syscall.RawConn
	fd  uintptr
	err error
}

func (c *testRawConn) Control(f func(fd uintptr)) error {
	f(c.fd)
	return c.err
}

func TestControllers(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		cs := controllers([]Controller{})
		err := cs.control("network", "address", &testRawConn{})
		ztesting.AssertEqualErr(t, "error not match", nil, err)
	})
	t.Run("1 control func", func(t *testing.T) {
		count := 0
		cs := controllers([]Controller{
			func(fd uintptr) error { count++; return nil },
		})
		err := cs.control("network", "address", &testRawConn{})
		ztesting.AssertEqual(t, "count not match", 1, count)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
	})
	t.Run("2 control func", func(t *testing.T) {
		count := 0
		cs := controllers([]Controller{
			func(fd uintptr) error { count++; return nil },
			func(fd uintptr) error { count++; return nil },
		})
		err := cs.control("network", "address", &testRawConn{})
		ztesting.AssertEqual(t, "count not match", 2, count)
		ztesting.AssertEqualErr(t, "error not match", nil, err)
	})
	t.Run("control func error", func(t *testing.T) {
		count := 0
		cs := controllers([]Controller{
			func(fd uintptr) error { count++; return io.EOF }, // Return dummy error.
			func(fd uintptr) error { count++; return nil },
		})
		err := cs.control("network", "address", &testRawConn{})
		ztesting.AssertEqual(t, "count not match", 1, count)
		ztesting.AssertEqualErr(t, "error not match", io.EOF, err)
	})
	t.Run("conn error", func(t *testing.T) {
		count := 0
		cs := controllers([]Controller{
			func(fd uintptr) error { count++; return nil },
			func(fd uintptr) error { count++; return nil },
		})
		err := cs.control("network", "address", &testRawConn{err: io.EOF}) // Return dummy error.
		ztesting.AssertEqual(t, "count not match", 2, count)
		ztesting.AssertEqualErr(t, "error not match", io.EOF, err)
	})
}

func TestAppendNonNil(t *testing.T) {
	t.Parallel()
	arr := []Controller{}
	arr = appendNonNil(arr, nil)
	ztesting.AssertEqual(t, "nil should not be appended", 0, len(arr))
	arr = appendNonNil(arr, func(fd uintptr) error { return nil })
	ztesting.AssertEqual(t, "non-nil should be appended", 1, len(arr))
}

func TestSocketError(t *testing.T) {
	t.Parallel()
	t.Run("error msg", func(t *testing.T) {
		err := &SocketError{Err: io.EOF, Opts: "FOO.BAR"}
		msg := err.Error()
		want := "zsyscall: fail to apply socket option FOO.BAR [EOF]"
		ztesting.AssertEqual(t, "error message not match", want, msg)
	})
	t.Run("same error", func(t *testing.T) {
		err1 := &SocketError{Opts: "FOO.BAR"}
		err2 := &SocketError{Opts: "FOO.BAR"}
		is := err1.Is(err2)
		ztesting.AssertEqual(t, "error not match", true, is)
	})
	t.Run("different error", func(t *testing.T) {
		err1 := &SocketError{Opts: "FOO.BAR"}
		err2 := &SocketError{Opts: "FOO.BAZ"}
		is := err1.Is(err2)
		ztesting.AssertEqual(t, "error not match", false, is)
	})
	t.Run("non socket error", func(t *testing.T) {
		err1 := &SocketError{Opts: "FOO.BAR"}
		is := err1.Is(io.EOF)
		ztesting.AssertEqual(t, "error not match", false, is)
	})
	t.Run("nil error", func(t *testing.T) {
		err1 := &SocketError{Opts: "FOO.BAR"}
		is := err1.Is(nil)
		ztesting.AssertEqual(t, "error not match", false, is)
	})
}

// func TestSockOptionFromSpec(t *testing.T) {
// 	type condition struct {
// 		opt *k.SockOption
// 	}

// 	type action struct {
// 		opt *SockOption
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	cndNil := tb.Condition("nil option", "input nil as an option")
// 	cndZero := tb.Condition("zero option", "input zero value as an option")
// 	actCheckNil := tb.Action("check nil", "check that nil was returned")
// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"nil",
// 			[]string{cndNil},
// 			[]string{actCheckNil},
// 			&condition{
// 				opt: nil,
// 			},
// 			&action{
// 				opt: nil,
// 			},
// 		),
// 		gen(
// 			"zero",
// 			[]string{cndZero},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockOption{},
// 			},
// 			&action{
// 				opt: &SockOption{},
// 			},
// 		),
// 		gen(
// 			"with so",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockOption{
// 					SOOption: &k.SockSOOption{},
// 				},
// 			},
// 			&action{
// 				opt: &SockOption{
// 					SO: &SockSOOption{},
// 				},
// 			},
// 		),
// 		gen(
// 			"with ip",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockOption{
// 					IPOption: &k.SockIPOption{},
// 				},
// 			},
// 			&action{
// 				opt: &SockOption{
// 					IP: &SockIPOption{},
// 				},
// 			},
// 		),
// 		gen(
// 			"with ipv6",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockOption{
// 					IPV6Option: &k.SockIPV6Option{},
// 				},
// 			},
// 			&action{
// 				opt: &SockOption{
// 					IPV6: &SockIPV6Option{},
// 				},
// 			},
// 		),
// 		gen(
// 			"with tcp",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockOption{
// 					TCPOption: &k.SockTCPOption{},
// 				},
// 			},
// 			&action{
// 				opt: &SockOption{
// 					TCP: &SockTCPOption{},
// 				},
// 			},
// 		),
// 		gen(
// 			"with udp",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockOption{
// 					UDPOption: &k.SockUDPOption{},
// 				},
// 			},
// 			&action{
// 				opt: &SockOption{
// 					UDP: &SockUDPOption{},
// 				},
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			opt := SockOptionFromSpec(tt.C().opt)
// 			testutil.Diff(t, tt.A().opt, opt)
// 		})
// 	}
// }

// func TestSockSOOptionFromSpec(t *testing.T) {
// 	type condition struct {
// 		opt *k.SockSOOption
// 	}

// 	type action struct {
// 		opt *SockSOOption
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	cndNil := tb.Condition("nil option", "input nil as an option")
// 	cndZero := tb.Condition("zero option", "input zero value as an option")
// 	actCheckNil := tb.Action("check nil", "check that nil was returned")
// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"nil",
// 			[]string{cndNil},
// 			[]string{actCheckNil},
// 			&condition{
// 				opt: nil,
// 			},
// 			&action{
// 				opt: nil,
// 			},
// 		),
// 		gen(
// 			"zero",
// 			[]string{cndZero},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockSOOption{},
// 			},
// 			&action{
// 				opt: &SockSOOption{},
// 			},
// 		),
// 		gen(
// 			"all",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockSOOption{
// 					BindToDevice:       "eth0",
// 					Debug:              true,
// 					IncomingCPU:        true,
// 					KeepAlive:          true,
// 					Linger:             10,
// 					Mark:               11,
// 					ReceiveBuffer:      12,
// 					ReceiveBufferForce: 13,
// 					ReceiveTimeout:     14,
// 					SendTimeout:        15,
// 					ReuseAddr:          true,
// 					ReusePort:          true,
// 					SendBuffer:         16,
// 					SendBufferForce:    17,
// 				},
// 			},
// 			&action{
// 				opt: &SockSOOption{
// 					BindToDevice:       "eth0",
// 					Debug:              true,
// 					IncomingCPU:        true,
// 					KeepAlive:          true,
// 					Linger:             10,
// 					Mark:               11,
// 					ReceiveBuffer:      12,
// 					ReceiveBufferForce: 13,
// 					ReceiveTimeout:     14,
// 					SendTimeout:        15,
// 					ReuseAddr:          true,
// 					ReusePort:          true,
// 					SendBuffer:         16,
// 					SendBufferForce:    17,
// 				},
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			opt := SockSOOptionFromSpec(tt.C().opt)
// 			testutil.Diff(t, tt.A().opt, opt)
// 		})
// 	}
// }

// func TestSockIPOptionFromSpec(t *testing.T) {
// 	type condition struct {
// 		opt *k.SockIPOption
// 	}

// 	type action struct {
// 		opt *SockIPOption
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	cndNil := tb.Condition("nil option", "input nil as an option")
// 	cndZero := tb.Condition("zero option", "input zero value as an option")
// 	actCheckNil := tb.Action("check nil", "check that nil was returned")
// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"nil",
// 			[]string{cndNil},
// 			[]string{actCheckNil},
// 			&condition{
// 				opt: nil,
// 			},
// 			&action{
// 				opt: nil,
// 			},
// 		),
// 		gen(
// 			"zero",
// 			[]string{cndZero},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockIPOption{},
// 			},
// 			&action{
// 				opt: &SockIPOption{},
// 			},
// 		),
// 		gen(
// 			"all",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockIPOption{
// 					BindAddressNoPort:   true,
// 					FreeBind:            true,
// 					LocalPortRangeUpper: 10,
// 					LocalPortRangeLower: 11,
// 					Transparent:         true,
// 					TTL:                 12,
// 				},
// 			},
// 			&action{
// 				opt: &SockIPOption{
// 					BindAddressNoPort:   true,
// 					FreeBind:            true,
// 					LocalPortRangeUpper: 10,
// 					LocalPortRangeLower: 11,
// 					Transparent:         true,
// 					TTL:                 12,
// 				},
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			opt := SockIPOptionFromSpec(tt.C().opt)
// 			testutil.Diff(t, tt.A().opt, opt)
// 		})
// 	}
// }

// func TestSockIPV6OptionFromSpec(t *testing.T) {
// 	type condition struct {
// 		opt *k.SockIPV6Option
// 	}

// 	type action struct {
// 		opt *SockIPV6Option
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	cndNil := tb.Condition("nil option", "input nil as an option")
// 	cndZero := tb.Condition("zero option", "input zero value as an option")
// 	actCheckNil := tb.Action("check nil", "check that nil was returned")
// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"nil",
// 			[]string{cndNil},
// 			[]string{actCheckNil},
// 			&condition{
// 				opt: nil,
// 			},
// 			&action{
// 				opt: nil,
// 			},
// 		),
// 		gen(
// 			"zero",
// 			[]string{cndZero},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockIPV6Option{},
// 			},
// 			&action{
// 				opt: &SockIPV6Option{},
// 			},
// 		),
// 		gen(
// 			"all",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockIPV6Option{},
// 			},
// 			&action{
// 				opt: &SockIPV6Option{},
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			opt := SockIPV6OptionFromSpec(tt.C().opt)
// 			testutil.Diff(t, tt.A().opt, opt)
// 		})
// 	}
// }

// func TestSockTCPOptionFromSpec(t *testing.T) {
// 	type condition struct {
// 		opt *k.SockTCPOption
// 	}

// 	type action struct {
// 		opt *SockTCPOption
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	cndNil := tb.Condition("nil option", "input nil as an option")
// 	cndZero := tb.Condition("zero option", "input zero value as an option")
// 	actCheckNil := tb.Action("check nil", "check that nil was returned")
// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"nil",
// 			[]string{cndNil},
// 			[]string{actCheckNil},
// 			&condition{
// 				opt: nil,
// 			},
// 			&action{
// 				opt: nil,
// 			},
// 		),
// 		gen(
// 			"zero",
// 			[]string{cndZero},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockTCPOption{},
// 			},
// 			&action{
// 				opt: &SockTCPOption{},
// 			},
// 		),
// 		gen(
// 			"all",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockTCPOption{
// 					CORK:            true,
// 					DeferAccept:     10,
// 					KeepCount:       11,
// 					KeepIdle:        12,
// 					KeepInterval:    13,
// 					Linger2:         14,
// 					MaxSegment:      15,
// 					NoDelay:         true,
// 					QuickAck:        true,
// 					SynCount:        16,
// 					UserTimeout:     17,
// 					WindowClamp:     18,
// 					FastOpen:        true,
// 					FastOpenConnect: true,
// 				},
// 			},
// 			&action{
// 				opt: &SockTCPOption{
// 					CORK:            true,
// 					DeferAccept:     10,
// 					KeepCount:       11,
// 					KeepIdle:        12,
// 					KeepInterval:    13,
// 					Linger2:         14,
// 					MaxSegment:      15,
// 					NoDelay:         true,
// 					QuickAck:        true,
// 					SynCount:        16,
// 					UserTimeout:     17,
// 					WindowClamp:     18,
// 					FastOpen:        true,
// 					FastOpenConnect: true,
// 				},
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			opt := SockTCPOptionFromSpec(tt.C().opt)
// 			testutil.Diff(t, tt.A().opt, opt)
// 		})
// 	}
// }

// func TestSockUDPOptionFromSpec(t *testing.T) {
// 	type condition struct {
// 		opt *k.SockUDPOption
// 	}

// 	type action struct {
// 		opt *SockUDPOption
// 	}

// 	tb := testutil.NewTableBuilder[*condition, *action]()
// 	tb.Name(t.Name())
// 	cndNil := tb.Condition("nil option", "input nil as an option")
// 	cndZero := tb.Condition("zero option", "input zero value as an option")
// 	actCheckNil := tb.Action("check nil", "check that nil was returned")
// 	table := tb.Build()

// 	gen := testutil.NewCase[*condition, *action]
// 	testCases := []*testutil.Case[*condition, *action]{
// 		gen(
// 			"nil",
// 			[]string{cndNil},
// 			[]string{actCheckNil},
// 			&condition{
// 				opt: nil,
// 			},
// 			&action{
// 				opt: nil,
// 			},
// 		),
// 		gen(
// 			"zero",
// 			[]string{cndZero},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockUDPOption{},
// 			},
// 			&action{
// 				opt: &SockUDPOption{},
// 			},
// 		),
// 		gen(
// 			"all",
// 			[]string{},
// 			[]string{},
// 			&condition{
// 				opt: &k.SockUDPOption{
// 					CORK:    true,
// 					Segment: 10,
// 					GRO:     true,
// 				},
// 			},
// 			&action{
// 				opt: &SockUDPOption{
// 					CORK:    true,
// 					Segment: 10,
// 					GRO:     true,
// 				},
// 			},
// 		),
// 	}

// 	testutil.Register(table, testCases...)

// 	for _, tt := range table.Entries() {
// 		tt := tt
// 		t.Run(tt.Name(), func(t *testing.T) {
// 			opt := SockUDPOptionFromSpec(tt.C().opt)
// 			testutil.Diff(t, tt.A().opt, opt)
// 		})
// 	}
// }
