package zudp

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
	t.Run("invalid udp address", func(t *testing.T) {
		rrd := &roundRobinDialer{addrs: []string{"udp://12345"}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "conn should be nil", nil, conn)
		_, ok := err.(*net.AddrError)
		ztesting.AssertEqual(t, "error should be addr error", true, ok)
	})
	t.Run("dial udp", func(t *testing.T) {
		pc, _ := net.ListenPacket("udp4", ":0")
		defer pc.Close()
		rrd := &roundRobinDialer{addrs: []string{"udp4://" + pc.LocalAddr().String()}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		conn.Close()
	})
	t.Run("dial unix", func(t *testing.T) {
		// In this case, make dial fail because the windows does not support unixgram.
		s := t.TempDir() + "/test.sock"
		rrd := &roundRobinDialer{addrs: []string{"unixgram://" + s}}
		_, err := rrd.dial(context.Background(), nil)
		_, ok := err.(*net.OpError)
		ztesting.AssertEqual(t, "error should be net op error", true, ok)
	})
	t.Run("dial fallback", func(t *testing.T) {
		pc, _ := net.ListenPacket("udp", ":0")
		defer pc.Close()
		rrd := &roundRobinDialer{addrs: []string{pc.LocalAddr().String()}}
		conn, err := rrd.dial(context.Background(), nil)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		conn.Close()
	})
}
