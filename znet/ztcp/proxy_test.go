package ztcp

import (
	"context"
	"net"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestRoundRobinDialer(t *testing.T) {
	t.Parallel()
	t.Run("test next", func(t *testing.T) {
		rrd := &roundRobinDialer{
			index: -1,
			addrs: []string{"addr1", "addr2", "addr3"},
		}
		got := []string{}
		for range 6 {
			got = append(got, rrd.next())
		}
		want := []string{"addr1", "addr2", "addr3", "addr1", "addr2", "addr3"}
		ztesting.AssertEqualSlice(t, "address not match", want, got)
	})
	t.Run("invalid tcp address", func(t *testing.T) {
		rrd := &roundRobinDialer{addrs: []string{"tcp://12345"}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "conn should be nil", nil, conn)
		_, ok := err.(*net.AddrError)
		ztesting.AssertEqual(t, "error should be addr error", true, ok)
	})
	t.Run("dial tcp", func(t *testing.T) {
		ln, _ := net.Listen("tcp4", ":0")
		defer ln.Close()
		rrd := &roundRobinDialer{addrs: []string{"tcp4://" + ln.Addr().String()}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		conn.Close()
	})
	t.Run("dial unix", func(t *testing.T) {
		s := t.TempDir() + "/test.sock"
		ln, _ := net.Listen("unix", s)
		defer ln.Close()
		rrd := &roundRobinDialer{addrs: []string{"unix://" + s}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		conn.Close()
	})
	t.Run("dial fallback", func(t *testing.T) {
		ln, _ := net.Listen("tcp", ":0")
		defer ln.Close()
		rrd := &roundRobinDialer{addrs: []string{ln.Addr().String()}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		conn.Close()
	})
}
