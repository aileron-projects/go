package zsync

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztesting"
)

type testRunner struct {
	sleepMillis int
	runErr      error // Returned error.
	panicErr    error // If non-nil, panic immediately.

	isEntered         bool
	isExited          bool
	isSleepCompleted  bool
	isContextCanceled bool
}

func (r *testRunner) Run(ctx context.Context) error {
	r.isEntered = true
	defer func() { r.isExited = true }()
	if r.panicErr != nil {
		panic(r.panicErr)
	}
	select {
	case <-ctx.Done():
		r.isContextCanceled = true
	case <-time.After(time.Duration(r.sleepMillis) * time.Millisecond):
		r.isSleepCompleted = true
	}
	return r.runErr
}

func TestRegisterFunc(t *testing.T) {
	t.Run("register nil", func(t *testing.T) {
		g := &RunGroup{}
		g.RegisterFunc(nil)
		err := g.Run(context.Background())
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("register func", func(t *testing.T) {
		tr := &testRunner{sleepMillis: 1}
		g := &RunGroup{}
		g.RegisterFunc(tr.Run)
		err := g.Run(context.Background())
		ztesting.AssertEqual(t, true, tr.isEntered)
		ztesting.AssertEqual(t, true, tr.isExited)
		ztesting.AssertEqualErr(t, nil, err)
	})
}

func TestRegister(t *testing.T) {
	t.Run("register nil", func(t *testing.T) {
		g := &RunGroup{}
		g.Register(nil)
		err := g.Run(context.Background())
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("register runner", func(t *testing.T) {
		tr := &testRunner{sleepMillis: 1}
		g := &RunGroup{}
		g.Register(tr)
		err := g.Run(context.Background())
		ztesting.AssertEqual(t, true, tr.isEntered)
		ztesting.AssertEqual(t, true, tr.isExited)
		ztesting.AssertEqualErr(t, nil, err)
	})
}

func TestRunAndFailFast(t *testing.T) {
	t.Run("no runners", func(t *testing.T) {
		g := &RunGroup{}
		err := g.RunAndFailFast(context.Background())
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("nil context", func(t *testing.T) {
		tr := &testRunner{sleepMillis: 1}
		g := &RunGroup{}
		g.Register(tr)
		err := g.RunAndFailFast(nil)
		ztesting.AssertEqual(t, true, tr.isEntered)
		ztesting.AssertEqual(t, true, tr.isExited)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("nil error", func(t *testing.T) {
		tr1 := &testRunner{sleepMillis: 100}
		tr2 := &testRunner{sleepMillis: 500}
		g := &RunGroup{}
		g.Register(tr1, tr2)
		err := g.RunAndFailFast(context.Background())
		ztesting.AssertEqual(t, true, tr1.isEntered)
		ztesting.AssertEqual(t, true, tr1.isExited)
		ztesting.AssertEqual(t, true, tr2.isEntered)
		ztesting.AssertEqual(t, true, tr2.isExited)
		ztesting.AssertEqual(t, true, tr1.isSleepCompleted)
		ztesting.AssertEqual(t, true, tr2.isSleepCompleted)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("non nil error 1", func(t *testing.T) {
		tr1 := &testRunner{sleepMillis: 100, runErr: io.EOF}
		tr2 := &testRunner{sleepMillis: 10_000}
		g := &RunGroup{}
		g.Register(tr1, tr2)
		err := g.RunAndFailFast(context.Background())
		ztesting.AssertEqual(t, true, tr1.isEntered)
		ztesting.AssertEqual(t, true, tr1.isExited)
		ztesting.AssertEqual(t, true, tr2.isEntered)
		ztesting.AssertEqual(t, true, tr2.isExited)
		ztesting.AssertEqual(t, true, tr1.isSleepCompleted)
		ztesting.AssertEqual(t, true, tr2.isContextCanceled)
		ztesting.AssertEqualErr(t, io.EOF, err)
	})
	t.Run("non nil error 2", func(t *testing.T) {
		tr1 := &testRunner{sleepMillis: 10_000}
		tr2 := &testRunner{sleepMillis: 100, runErr: io.EOF}
		g := &RunGroup{}
		g.Register(tr1, tr2)
		err := g.RunAndFailFast(context.Background())
		ztesting.AssertEqual(t, true, tr1.isEntered)
		ztesting.AssertEqual(t, true, tr1.isExited)
		ztesting.AssertEqual(t, true, tr2.isEntered)
		ztesting.AssertEqual(t, true, tr2.isExited)
		ztesting.AssertEqual(t, true, tr2.isSleepCompleted)
		ztesting.AssertEqual(t, true, tr1.isContextCanceled)
		ztesting.AssertEqualErr(t, io.EOF, err)
	})
	t.Run("runner panic", func(t *testing.T) {
		tr1 := &testRunner{panicErr: io.EOF}
		tr2 := &testRunner{sleepMillis: 10_000}
		g := &RunGroup{}
		g.Register(tr1, tr2)
		err := g.RunAndFailFast(context.Background())
		ztesting.AssertEqual(t, true, tr1.isEntered)
		ztesting.AssertEqual(t, true, tr1.isExited)
		ztesting.AssertEqual(t, true, tr2.isEntered)
		ztesting.AssertEqual(t, true, tr2.isExited)
		ztesting.AssertEqual(t, false, tr1.isSleepCompleted)
		ztesting.AssertEqual(t, false, tr1.isContextCanceled)
		ztesting.AssertEqual(t, false, tr2.isSleepCompleted)
		ztesting.AssertEqual(t, true, tr2.isContextCanceled)
		errStr := errors.New("zsync: runner exit with panic. [EOF]")
		ztesting.AssertEqualErr(t, errStr, err) // Compare in string.
	})
}

func TestRunAndWaitAll(t *testing.T) {
	t.Run("no runners", func(t *testing.T) {
		g := &RunGroup{}
		err := g.RunAndWaitAll(context.Background())
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("nil context", func(t *testing.T) {
		tr := &testRunner{sleepMillis: 1}
		g := &RunGroup{}
		g.Register(tr)
		err := g.RunAndWaitAll(nil)
		ztesting.AssertEqual(t, true, tr.isEntered)
		ztesting.AssertEqual(t, true, tr.isExited)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("nil error", func(t *testing.T) {
		tr1 := &testRunner{sleepMillis: 100}
		tr2 := &testRunner{sleepMillis: 500}
		g := &RunGroup{}
		g.Register(tr1, tr2)
		err := g.RunAndWaitAll(context.Background())
		ztesting.AssertEqual(t, true, tr1.isEntered)
		ztesting.AssertEqual(t, true, tr1.isExited)
		ztesting.AssertEqual(t, true, tr2.isEntered)
		ztesting.AssertEqual(t, true, tr2.isExited)
		ztesting.AssertEqual(t, true, tr1.isSleepCompleted)
		ztesting.AssertEqual(t, true, tr2.isSleepCompleted)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("non nil error 1", func(t *testing.T) {
		tr1 := &testRunner{sleepMillis: 100, runErr: io.EOF}
		tr2 := &testRunner{sleepMillis: 500}
		g := &RunGroup{}
		g.Register(tr1, tr2)
		err := g.RunAndWaitAll(context.Background())
		ztesting.AssertEqual(t, true, tr1.isEntered)
		ztesting.AssertEqual(t, true, tr1.isExited)
		ztesting.AssertEqual(t, true, tr2.isEntered)
		ztesting.AssertEqual(t, true, tr2.isExited)
		ztesting.AssertEqual(t, true, tr1.isSleepCompleted)
		ztesting.AssertEqual(t, true, tr2.isSleepCompleted)
		ztesting.AssertEqualErr(t, io.EOF, err)
	})
	t.Run("non nil error 2", func(t *testing.T) {
		tr1 := &testRunner{sleepMillis: 100}
		tr2 := &testRunner{sleepMillis: 10, runErr: io.EOF}
		g := &RunGroup{}
		g.Register(tr1, tr2)
		err := g.RunAndWaitAll(context.Background())
		ztesting.AssertEqual(t, true, tr1.isEntered)
		ztesting.AssertEqual(t, true, tr1.isExited)
		ztesting.AssertEqual(t, true, tr2.isEntered)
		ztesting.AssertEqual(t, true, tr2.isExited)
		ztesting.AssertEqual(t, true, tr1.isSleepCompleted)
		ztesting.AssertEqual(t, true, tr2.isSleepCompleted)
		ztesting.AssertEqualErr(t, io.EOF, err)
	})
	t.Run("1 of 2 runner panics", func(t *testing.T) {
		tr1 := &testRunner{panicErr: io.EOF}
		tr2 := &testRunner{sleepMillis: 100}
		g := &RunGroup{}
		g.Register(tr1, tr2)
		err := g.RunAndWaitAll(context.Background())
		ztesting.AssertEqual(t, true, tr1.isEntered)
		ztesting.AssertEqual(t, true, tr1.isExited)
		ztesting.AssertEqual(t, true, tr2.isEntered)
		ztesting.AssertEqual(t, true, tr2.isExited)
		ztesting.AssertEqual(t, false, tr1.isSleepCompleted)
		ztesting.AssertEqual(t, false, tr1.isContextCanceled)
		ztesting.AssertEqual(t, true, tr2.isSleepCompleted)
		ztesting.AssertEqual(t, false, tr2.isContextCanceled)
		errStr := errors.New("zsync: runner exit with panic. [EOF]")
		ztesting.AssertEqualErr(t, errStr, err) // Compare in string.
	})
}

func TestAwakeRunner(t *testing.T) {
	t.Run("OnStart", func(t *testing.T) {
		tr := &testRunner{sleepMillis: 1}
		var count atomic.Int32
		g := &RunGroup{
			OnStart: func(r Runner) {
				count.Add(1)
				rr := r.(*testRunner)
				ztesting.AssertEqual(t, false, rr.isEntered)
				ztesting.AssertEqual(t, false, rr.isExited)
			},
		}
		g.Register(tr)
		err := g.RunAndWaitAll(context.Background())
		ztesting.AssertEqual(t, 1, count.Load())
		ztesting.AssertEqual(t, true, tr.isEntered)
		ztesting.AssertEqual(t, true, tr.isExited)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("OnExit", func(t *testing.T) {
		tr := &testRunner{sleepMillis: 1}
		var count atomic.Int32
		g := &RunGroup{
			OnExit: func(r Runner, err error) {
				count.Add(1)
				rr := r.(*testRunner)
				ztesting.AssertEqual(t, true, rr.isEntered)
				ztesting.AssertEqual(t, true, rr.isExited)
			},
		}
		g.Register(tr)
		err := g.RunAndWaitAll(context.Background())
		ztesting.AssertEqual(t, 1, count.Load())
		ztesting.AssertEqual(t, true, tr.isEntered)
		ztesting.AssertEqual(t, true, tr.isExited)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("OnStart 2 runners", func(t *testing.T) {
		tr1 := &testRunner{sleepMillis: 1}
		tr2 := &testRunner{sleepMillis: 1}
		var count atomic.Int32
		g := &RunGroup{
			OnStart: func(r Runner) {
				count.Add(1)
				rr := r.(*testRunner)
				ztesting.AssertEqual(t, false, rr.isEntered)
				ztesting.AssertEqual(t, false, rr.isExited)
			},
		}
		g.Register(tr1, tr2)
		err := g.RunAndWaitAll(context.Background())
		ztesting.AssertEqual(t, 2, count.Load())
		ztesting.AssertEqual(t, true, tr1.isEntered)
		ztesting.AssertEqual(t, true, tr1.isExited)
		ztesting.AssertEqual(t, true, tr2.isEntered)
		ztesting.AssertEqual(t, true, tr2.isExited)
		ztesting.AssertEqual(t, true, tr1.isSleepCompleted)
		ztesting.AssertEqual(t, true, tr2.isSleepCompleted)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("OnExit 2 runners", func(t *testing.T) {
		tr1 := &testRunner{sleepMillis: 1}
		tr2 := &testRunner{sleepMillis: 1}
		var count atomic.Int32
		g := &RunGroup{
			OnExit: func(r Runner, err error) {
				count.Add(1)
				rr := r.(*testRunner)
				ztesting.AssertEqual(t, true, rr.isEntered)
				ztesting.AssertEqual(t, true, rr.isExited)
			},
		}
		g.Register(tr1, tr2)
		err := g.RunAndWaitAll(context.Background())
		ztesting.AssertEqual(t, 2, count.Load())
		ztesting.AssertEqual(t, true, tr1.isEntered)
		ztesting.AssertEqual(t, true, tr1.isExited)
		ztesting.AssertEqual(t, true, tr2.isEntered)
		ztesting.AssertEqual(t, true, tr2.isExited)
		ztesting.AssertEqual(t, true, tr1.isSleepCompleted)
		ztesting.AssertEqual(t, true, tr2.isSleepCompleted)
		ztesting.AssertEqualErr(t, nil, err)
	})
	t.Run("OnExit with panic", func(t *testing.T) {
		tr := &testRunner{panicErr: io.EOF} // Panics dummy error.
		var count atomic.Int32
		g := &RunGroup{
			OnExit: func(r Runner, err error) {
				count.Add(1)
				rr := r.(*testRunner)
				ztesting.AssertEqual(t, true, rr.isEntered)
				ztesting.AssertEqual(t, true, rr.isExited)
			},
		}
		g.Register(tr)
		err := g.RunAndWaitAll(context.Background())
		ztesting.AssertEqual(t, 1, count.Load())
		ztesting.AssertEqual(t, true, tr.isEntered)
		ztesting.AssertEqual(t, true, tr.isExited)
		ztesting.AssertEqualErr(t, io.EOF, err)
	})
	t.Run("panics string", func(t *testing.T) {
		g := &RunGroup{}
		g.RegisterFunc(func(ctx context.Context) error { panic("string panic") })
		err := g.RunAndWaitAll(context.Background())
		errStr := errors.New("zsync: runner exit with panic. [string panic]")
		ztesting.AssertEqualErr(t, errStr, err) // Compare in string.
	})
	t.Run("panics error", func(t *testing.T) {
		g := &RunGroup{}
		g.RegisterFunc(func(ctx context.Context) error { panic(io.EOF) })
		err := g.RunAndWaitAll(context.Background())
		errStr := errors.New("zsync: runner exit with panic. [EOF]")
		ztesting.AssertEqualErr(t, errStr, err) // Compare in string.
	})
}
