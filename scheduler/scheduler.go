package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/GustavoKatel/asyncutils/executor"
	executorIfaces "github.com/GustavoKatel/asyncutils/executor/interfaces"
	"github.com/GustavoKatel/asyncutils/scheduler/interfaces"
)

// Scheduler provides a scheduling work channel
type Scheduler interfaces.Scheduler

var _ interfaces.Scheduler = &schedulerImpl{}

type schedulerImpl struct {
	worker executor.Executor

	ctx       context.Context
	ctxCancel context.CancelFunc

	lastExecution      *time.Time
	lastExecutionMutex *sync.RWMutex

	// Maps tag to the last recorded call
	throttleLast      executorIfaces.JobFn
	throttleLastMutex *sync.RWMutex
}

// New creates a new working channel
func New() (interfaces.Scheduler, error) {
	return NewWithContext(context.Background())
}

// NewWithContext creates a new working channel
func NewWithContext(ctx context.Context) (interfaces.Scheduler, error) {
	worker, err := executor.NewDefaultExecutorContext(ctx, 1)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	aw := &schedulerImpl{
		worker: worker,

		ctx:       ctx,
		ctxCancel: cancel,

		lastExecution:      nil,
		lastExecutionMutex: &sync.RWMutex{},

		throttleLast:      nil,
		throttleLastMutex: &sync.RWMutex{},
	}

	return aw, nil
}

func (aw *schedulerImpl) Start() error {
	return aw.worker.Start()
}

func (aw *schedulerImpl) Stop() error {
	aw.ctxCancel()
	return aw.worker.Stop()
}

func (aw *schedulerImpl) ErrorChan(ch chan error) {
	aw.worker.ErrorChan(ch)
}

func (aw *schedulerImpl) waitAndSchedule(job executorIfaces.JobFn, diffDelay time.Duration, delay time.Duration) {
	aw.throttleLastMutex.Lock()

	shouldWait := aw.throttleLast == nil
	aw.throttleLast = job

	aw.throttleLastMutex.Unlock()

	if !shouldWait {
		return
	}

	go func() {
		select {
		case <-time.After(diffDelay):
			aw.throttleLastMutex.Lock()

			if aw.throttleLast == nil {
				aw.throttleLastMutex.Unlock()
				return
			}

			job := aw.throttleLast
			aw.throttleLast = nil

			aw.throttleLastMutex.Unlock()

			aw.PostThrottledJob(job, delay)
		case <-aw.ctx.Done():
			return
		}
	}()
}

func (aw *schedulerImpl) PostJob(job executorIfaces.JobFn) error {
	return aw.PostThrottledJob(job, 0)
}

func (aw *schedulerImpl) PostThrottledJob(job executorIfaces.JobFn, delay time.Duration) error {
	aw.lastExecutionMutex.Lock()
	defer aw.lastExecutionMutex.Unlock()

	now := time.Now()

	if aw.lastExecution != nil && now.Sub(*aw.lastExecution) < delay {
		aw.waitAndSchedule(job, delay-now.Sub(*aw.lastExecution), delay)
		return nil
	}

	aw.lastExecution = &now

	return aw.worker.PostJob(job)
}

func (aw *schedulerImpl) Len() int {
	return aw.worker.Len()
}
