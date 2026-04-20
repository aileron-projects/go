package znet

import (
	"context"
	"net"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestIPResolver(t *testing.T) {
	t.Parallel()
	t.Run("resolve 0 addresses", func(t *testing.T) {
		r := &IPResolver{
			LookupIP: func(ctx context.Context, network, host string) ([]net.IP, error) {
				return []net.IP{}, nil
			},
		}
		err := r.Resolve(context.Background())
		ztesting.AssertEqualErr(t, nil, err)
		ztesting.AssertEqual(t, []string{}, r.IPs())
		ztesting.AssertEqual(t, "", r.Next())
	})
	t.Run("resolve 1 address", func(t *testing.T) {
		r := &IPResolver{
			LookupIP: func(ctx context.Context, network, host string) ([]net.IP, error) {
				return []net.IP{[]byte{127, 0, 0, 1}}, nil
			},
		}
		err := r.Resolve(context.Background())
		ztesting.AssertEqualErr(t, nil, err)
		ztesting.AssertEqual(t, []string{"127.0.0.1"}, r.IPs())
		ztesting.AssertEqual(t, "127.0.0.1", r.Next())
		ztesting.AssertEqual(t, "127.0.0.1", r.Next())
	})
	t.Run("resolve 3 addresses", func(t *testing.T) {
		var n, h string
		r := &IPResolver{
			Host: "test.com",
			LookupIP: func(ctx context.Context, network, host string) ([]net.IP, error) {
				n = network
				h = host
				return []net.IP{[]byte{127, 0, 0, 1}, []byte{127, 0, 0, 2}, []byte{127, 0, 0, 3}}, nil
			},
		}
		err := r.Resolve(context.Background())
		ztesting.AssertEqualErr(t, nil, err)
		ztesting.AssertEqual(t, "ip", n)
		ztesting.AssertEqual(t, "test.com", h)
		ztesting.AssertEqual(t, []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}, r.IPs())
		ztesting.AssertEqual(t, "127.0.0.2", r.Next())
		ztesting.AssertEqual(t, "127.0.0.3", r.Next())
		ztesting.AssertEqual(t, "127.0.0.1", r.Next())
	})
	t.Run("network and host", func(t *testing.T) {
		var n, h string
		r := &IPResolver{
			Network: "ip4",
			Host:    "test.com",
			LookupIP: func(ctx context.Context, network, host string) ([]net.IP, error) {
				n = network
				h = host
				return []net.IP{}, nil
			},
		}
		err := r.Resolve(context.Background())
		ztesting.AssertEqualErr(t, nil, err)
		ztesting.AssertEqual(t, "ip4", n)
		ztesting.AssertEqual(t, "test.com", h)
	})
	t.Run("resolve error", func(t *testing.T) {
		r := &IPResolver{
			LookupIP: func(ctx context.Context, network, host string) ([]net.IP, error) {
				return []net.IP{}, &net.AddrError{Err: "test", Addr: ""}
			},
		}
		err := r.Resolve(context.Background())
		ztesting.AssertEqualErr(t, &net.AddrError{Err: "test", Addr: ""}, err)
	})
	t.Run("default resolver", func(t *testing.T) {
		r := &IPResolver{
			LookupIP: nil,
		}
		err := r.Resolve(context.Background())
		ztesting.AssertEqualErr(t, &net.DNSError{Err: "no such host"}, err)
	})
}
