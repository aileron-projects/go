package zsync

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// RunnerFunc is the runnable function.
// It implements [Runner] interface.
type RunnerFunc func(context.Context) error

func (f RunnerFunc) Run(ctx context.Context) error {
	return f(ctx)
}

// Runner run the process that implements this interface.
type Runner interface {
	// Run run the process that implements this interface.
	Run(context.Context) error
}

// RunGroup is the group of [Runner].
// It can awake all registered runners using [RunGroup.]
type RunGroup struct {
	// mu protects runners.
	mu sync.Mutex
	// runners is the list of runner to be run.
	runners []Runner
	// OnStart is called for all runners before they actually run.
	// The runner is given by the first argument.
	OnStart func(r Runner)
	// OnExit is called for all runners after they exited.
	// The runner is given by the first argument.
	// The error its runner returned is passed by the second argument.
	OnExit func(r Runner, err error)
}

// RegisterFunc registers the given function as a runner to the group.
// It's safe to call RegisterFunc from different goroutine simultaneously.
// Registered function is used from the next run and does not affect the
// function group already running.
func (g *RunGroup) RegisterFunc(f func(context.Context) error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if f == nil {
		return
	}
	g.runners = append(g.runners, RunnerFunc(f))
}

// Register registers the given runners to the group.
// It's safe to call Register from different goroutine simultaneously.
// Registered runners are used from the next run and does not affect the
// function group already running.
func (g *RunGroup) Register(rs ...Runner) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, r := range rs {
		if r == nil {
			continue
		}
		g.runners = append(g.runners, r)
	}
}

// Run is the alias for [RunGroup.RunAndFailFast].
func (g *RunGroup) Run(ctx context.Context) error {
	return g.RunAndFailFast(ctx)
}

// RunAndFailFast runs the runners in a new goroutine each.
// RunAndFailFast waits runners to return. If one of the runners returned
// non nil error, it cancel the context that passed to the all runners
// and return with the obtained non nil error.
// [RunGroup.RunAndFailFast] and [RunGroup.RunAndWaitAll] works same except
// for the case that at least one runner returned non nil error.
// It's safe to call RunAndWaitAll from different goroutine simultaneously.
// If given ctx is nil, it uses a new context created with [context.Background].
func (g *RunGroup) RunAndFailFast(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g.mu.Lock()
	if len(g.runners) == 0 {
		g.mu.Unlock()
		return nil
	}

	wg := sync.WaitGroup{}
	n := len(g.runners)
	errChan := make(chan error, n)
	for _, r := range g.runners {
		g.awakeRunner(ctx, r, errChan, &wg)
	}
	g.mu.Unlock()

	var err error
	for range n {
		err = <-errChan
		if err != nil {
			cancel()
			break
		}
	}

	wg.Wait()
	close(errChan) // no one send to it anymore.
	return err
}

// RunAndWaitAll runs the runners in a new goroutine each.
// It waits all runners to return and collects the errors returned from them.
// It returns the collected errors joined by [errors.Join].
// [RunGroup.RunAndFailFast] and [RunGroup.RunAndWaitAll] works same except
// for the case that at least one runner returned non nil error.
// It's safe to call RunAndWaitAll from different goroutine simultaneously.
// If given ctx is nil, it uses a new context created with [context.Background].
func (g *RunGroup) RunAndWaitAll(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g.mu.Lock()
	if len(g.runners) == 0 {
		g.mu.Unlock()
		return nil
	}

	wg := sync.WaitGroup{}
	errChan := make(chan error, len(g.runners))
	for _, r := range g.runners {
		g.awakeRunner(ctx, r, errChan, &wg)
	}
	g.mu.Unlock()

	wg.Wait()
	close(errChan) // no one send to it anymore.

	errs := []error{}
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (g *RunGroup) awakeRunner(ctx context.Context, r Runner, errChan chan error, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
	}
	go func(inCtx context.Context) {
		if wg != nil {
			defer wg.Done()
		}

		var runErr error
		defer func() {
			errChan <- runErr
		}()

		if g.OnStart != nil {
			g.OnStart(r)
		}
		if g.OnExit != nil {
			defer g.OnExit(r, runErr)
		}

		defer func() {
			r := recover() // Handle runner's panic.
			if r == nil {
				return
			}
			if err, ok := r.(error); ok {
				runErr = fmt.Errorf("zsync: runner exit with panic. [%w]", err)
			} else {
				runErr = fmt.Errorf("zsync: runner exit with panic. [%v]", r)
			}
		}()

		runErr = r.Run(inCtx)
	}(ctx)
}
