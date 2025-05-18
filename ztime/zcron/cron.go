package zcron

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// Event is the event definition on cron.
// Following types of events are defined.
//
//   - [OnJobRun]
//   - [OnJobAccepted]
//   - [OnJobDeclined]
//   - [OnJobStarted]
//   - [OnJobExited]
//   - [OnJobFailed]
//   - [OnJobPanicked]
type Event int

const (
	Undefined     Event = iota
	OnJobRun            // OnJobRun triggered when job is scheduled and run.
	OnJobAccepted       // OnJobAccepted triggered when job is accepted.
	OnJobDeclined       // OnJobAccepted triggered when job is declined (Max concurrency reached).
	OnJobStarted        // OnJobStarted triggered just before [Config.JobFunc] is called.
	OnJobExited         // OnJobExited triggered just after [Config.JobFunc] is returned.
	OnJobFailed         // OnJobFailed triggered when [Config.JobFunc] returned an error. This trigger is before [OnJobExited].
	OnJobPanicked       // OnJobPanicked triggered when [Config.JobFunc] panicked. This trigger is before [OnJobExited].
)

var (
	// ErrNilConfig indicates nil was given as config.
	ErrNilConfig = errors.New("ztime/zcron: nil config")
	// ErrNilJob indicates nil job function was given.
	ErrNilJob = errors.New("ztime/zcron: nil job function")
)

// Config is the configuration for the [Cron].
type Config struct {
	// Crontab is the cron expression.
	// See [Parse] for the syntax.
	Crontab string
	// MaxConcurrency is the maximum concurrency
	// of the currently running JobFunc.
	// Values should be 1, 2, ...
	// 1 means exactly 1 job at a time.
	// If less than 1, 1 is used.
	MaxConcurrency int
	// MaxRetry is the maximum number to retry
	// when the JobFunc returned non-nil error.
	// Values should be 0, 1, 2, ...
	// 0 means no retry. If less than 0, 0 is used.
	MaxRetry int
	// JobFunc is the job to run.
	JobFunc func(context.Context) error
	// WithContext provides a context which passed to the [Job.Run].
	// This is called once for a run and not called for retry.
	// [context.Background] is used when WithContext is nil.
	WithContext func() context.Context
	// EventHook is the function that hooks [Event]s.
	// The [Event] is notified through the first argument.
	// Additional information such as error is passed by a
	// depending on the event type.
	// 	- For OnJobFailed: error is given by a[0].
	// 	- For OnJobPanicked: recovered value given by a[0].
	EventHook func(e Event, a ...any)
}

// NewCron returns a new instance of [Cron].
func NewCron(c *Config) (*Cron, error) {
	if c == nil {
		return nil, ErrNilConfig
	}
	if c.JobFunc == nil {
		return nil, ErrNilJob
	}
	ct, err := Parse(c.Crontab)
	if err != nil {
		return nil, err
	}
	return &Cron{
		cron:      ct,
		timeAfter: time.After,
		runner: &runner{
			maxRetry:    max(0, c.MaxRetry),
			queue:       make(chan struct{}, max(1, c.MaxConcurrency)),
			jobFunc:     c.JobFunc,
			withContext: c.WithContext,
			eventFunc:   c.EventHook,
		},
	}, nil
}

// Cron schedules and runs jobs based on [Crontab].
type Cron struct {
	// started represents if the cron is started or not.
	// It cannot start
	started atomic.Bool
	cron    *Crontab
	runner  *runner
	// stop stops the cron working.
	stop      chan struct{}
	timeAfter func(time.Duration) <-chan time.Time
}

// WithTimeAfterFunc replaces internal wait functions.
func (c *Cron) WithTimeAfterFunc(timeAfter func(time.Duration) <-chan time.Time) {
	c.timeAfter = timeAfter
}

// WithTimeFunc replaces internal clock.
func (c *Cron) WithTimeFunc(timeNow func() time.Time) {
	c.cron.WithTimeFunc(timeNow)
}

// Start starts the cron.
// The cron scheduling can only run 1 process at a time.
// Calling Start multiple times does not run the multiple crons.
// If the next scheduling is after 10 minutes or more,
// the internal timer calibrates the time after 95% of the times.
// For examples, now a job is scheduled after 10 minutes,
// the con calibrates timer to fire the job after 9min30sec.
// It blocks the process. Run in a new goroutine if blocking is
// not necessary like below.
//
//	cron, _ := zcron.NewCron(.....)
//	go cron.Start()
func (c *Cron) Start() {
	if c.started.Swap(true) {
		return // already running.
	}
	defer c.started.Store(false)
	c.stop = make(chan struct{}, 1)

	for {
		now := c.cron.Now()
		next := c.cron.NextAfter(now)
		wait := min(next.Sub(now), time.Hour)
		calibrate := false
		if wait >= 10*time.Minute {
			calibrate = true // Calibrate time after 95% of it.
			wait = 95 * wait / 100
		}
		select {
		case <-c.timeAfter(wait):
			if calibrate {
				continue
			}
		case <-c.stop:
			return
		}
		c.runner.Run() // Run must not block.
	}
}

// Stop stops new scheduling of jobs.
// Stopping cron does not stop the already running jobs.
func (c *Cron) Stop() {
	if c.started.Load() {
		c.stop <- struct{}{}
	}
}

// runner runs jobFunc and controls
// concurrency and retries.
type runner struct {
	queue       chan struct{}
	maxRetry    int
	jobFunc     func(context.Context) error
	withContext func() context.Context
	eventFunc   func(Event, ...any)
}

func (r *runner) onEvent(e Event, a ...any) {
	if r.eventFunc == nil {
		return
	}
	r.eventFunc(e, a...)
}

func (r *runner) Run() {
	r.onEvent(OnJobStarted)

	select {
	case r.queue <- struct{}{}:
		r.onEvent(OnJobAccepted)
	default:
		r.onEvent(OnJobDeclined)
		return // Max concurrency exceeded.
	}

	go func() { // Run the job.
		r.onEvent(OnJobStarted)
		defer func() {
			<-r.queue // Remove lock.
		}()
		defer func() {
			if rec := recover(); rec != nil {
				r.onEvent(OnJobPanicked, rec)
			}
			r.onEvent(OnJobExited)
		}()
		ctx := context.Background()
		if r.withContext != nil {
			ctx = r.withContext()
		}
		for range r.maxRetry + 1 {
			err := r.jobFunc(ctx)
			if err != nil {
				r.onEvent(OnJobFailed, err)
				return
			}
		}
	}()
}
