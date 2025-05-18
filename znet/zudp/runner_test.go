package zudp_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/aileron-projects/go/znet/zudp"
	"github.com/aileron-projects/go/ztesting"
)

func TestServerRunner(t *testing.T) {
	t.Parallel()
	t.Run("no error", func(t *testing.T) {
		r := &zudp.ServerRunner{
			Serve:    func() error { return nil },
			Shutdown: func(context.Context) error { return nil },
		}
		err := r.Run(context.Background())
		ztesting.AssertEqualErr(t, "error not match", nil, err)
	})
	t.Run("serve error", func(t *testing.T) {
		t.Run("already closed", func(t *testing.T) {
			r := &zudp.ServerRunner{
				Serve:    func() error { return net.ErrClosed },
				Shutdown: func(context.Context) error { return nil },
			}
			err := r.Run(context.Background())
			ztesting.AssertEqualErr(t, "error not match", net.ErrClosed, err)
		})
		t.Run("non-nil error", func(t *testing.T) {
			testErr := errors.New("serve error")
			r := &zudp.ServerRunner{
				Serve:    func() error { return testErr },
				Shutdown: func(context.Context) error { return nil },
			}
			err := r.Run(context.Background())
			ztesting.AssertEqualErr(t, "error not match", testErr, err)
		})
	})
	t.Run("shutdown error", func(t *testing.T) {
		t.Run("timeout", func(t *testing.T) {
			closeCalled := false
			r := &zudp.ServerRunner{
				Serve:           func() error { return net.ErrClosed },
				Shutdown:        func(ctx context.Context) error { <-ctx.Done(); return ctx.Err() },
				Close:           func() error { closeCalled = true; return nil },
				ShutdownTimeout: 10 * time.Millisecond,
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err := r.Run(ctx)
			ztesting.AssertEqualErr(t, "error not match", context.DeadlineExceeded, err)
			ztesting.AssertEqual(t, "close is not called", true, closeCalled)
		})
		t.Run("non-nil error", func(t *testing.T) {
			testErr := errors.New("shutdown error")
			closeCalled := false
			r := &zudp.ServerRunner{
				Serve:    func() error { return net.ErrClosed },
				Shutdown: func(context.Context) error { return testErr },
				Close:    func() error { closeCalled = true; return nil },
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err := r.Run(ctx)
			ztesting.AssertEqualErr(t, "error not match", testErr, err)
			ztesting.AssertEqual(t, "close is called", false, closeCalled)
		})
	})
}
