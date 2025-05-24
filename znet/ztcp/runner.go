package ztcp

import (
	"context"
	"time"

	"github.com/aileron-projects/go/znet/internal"
)

// ServerRunner runs a server with an ability of graceful shutdown.
// See the following usage example.
//
// Example:
//
//	svr := &ztcp.Server{Addr: ":8080"} // Register a handler.
//	r := &ServerRunner{
//		Serve:           svr.ListenAndServe,
//		Shutdown:        svr.Shutdown,
//		Close:           svr.Close,
//		ShutdownTimeout: 30 * time.Second,
//	}
//
//	sigCtx, cancel := signal.NotifyContext(context.Background(),
//		os.Interrupt, syscall.SIGTERM)
//	defer cancel()
//
//	if err := r.Run(sigCtx); err != nil {
//		panic(err)
//	}
type ServerRunner struct {
	// Server starts a server.
	// Serve must not be nil.
	Serve func() error
	// Shutdown gracefully shutdowns the server.
	// Shutdown must not be nil.
	// It will be called only when the context given to
	// the ServerRunner.Run is done.
	// Otherwise, Shutdown is not called even the Server exited.
	// Typically [Server.Shutdown] should be set.
	Shutdown func(context.Context) error
	// Close immediately closes a server.
	// Unlike Shutdown, it should not block.
	// Close, if non-nil, will be called only when shutdown timeout occurred.
	// Typically [Server.Close] should be set.
	// Note that the [Server.Shutdown] does not close remaining
	// connection after shutdown timeout occurred but this Runner try to.
	Close func() error
	// ShutdownTimeout is the timeout duration applied
	// for the Shutdown function.
	// For ShutdownTimeout<=0, 30 seconds is used.
	ShutdownTimeout time.Duration
}

// Run runs a server.
// A server will be shutdown when the sigCtx is done.
// It returns non-nil error if r.Serve returns non-nil error.
// When a timeout occurred while shutting down,
// a [context.DeadlineExceeded] can be returned.
func (r *ServerRunner) Run(sigCtx context.Context) error {
	runner := &internal.ServerRunner{
		Serve:           r.Serve,
		Shutdown:        r.Shutdown,
		Close:           r.Close,
		ShutdownTimeout: r.ShutdownTimeout,
	}
	return runner.Run(sigCtx)
}
