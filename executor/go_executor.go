package executor

import (
	"context"
	"sync"

	"github.com/GustavoKatel/asyncutils/event"
	"github.com/GustavoKatel/asyncutils/executor/interfaces"
	"github.com/GustavoKatel/asyncutils/queue"
)

var _ interfaces.Executor = &goExecutor{}

type goExecutor struct {
	queue      queue.Queue
	queueMutex *sync.Mutex

	hasJobsEvent event.Event

	workers int

	errorChs      []chan error
	errorChsMutex *sync.RWMutex

	ctx       context.Context
	ctxCancel context.CancelFunc
}

// NewDefaultExecutor creates a new default executor which maps workers as gorountines
func NewDefaultExecutor(workers int) (interfaces.Executor, error) {
	return NewDefaultExecutorContext(context.Background(), workers)
}

// NewDefaultExecutorContext creates a new default executor which maps workers as gorountines
func NewDefaultExecutorContext(ctx context.Context, workers int) (interfaces.Executor, error) {
	ctx, cancel := context.WithCancel(ctx)

	exec := &goExecutor{
		queue:      queue.New(),
		queueMutex: &sync.Mutex{},

		hasJobsEvent: event.NewEvent(false),

		workers: workers,

		errorChs:      []chan error{},
		errorChsMutex: &sync.RWMutex{},

		ctx:       ctx,
		ctxCancel: cancel,
	}

	return exec, nil
}

func (ge *goExecutor) Start() error {
	go ge.background()
	return nil
}

func (ge *goExecutor) Stop() error {
	ge.ctxCancel()
	ge.hasJobsEvent.Set()
	return nil
}

func (ge *goExecutor) background() {
	for i := 0; i < ge.workers; i++ {
		go ge.worker(i)
	}
}

func (ge *goExecutor) worker(id int) {
	for ge.ctx.Err() == nil {
		ge.hasJobsEvent.Wait()

		jobI := ge.queue.PopFront()

		// The queue is empty
		if jobI == nil {
			ge.hasJobsEvent.Reset()
			continue
		}

		job := jobI.(*jobImpl)

		if err := job.jobFn(ge.ctx); err != nil {
			ge.emitError(err)
		}
	}
}

func (ge *goExecutor) emitError(err error) {
	ge.errorChsMutex.RLock()
	defer ge.errorChsMutex.RUnlock()

	for _, ch := range ge.errorChs {
		ch <- err
	}
}

func (ge *goExecutor) ErrorChan(ch chan error) {
	ge.errorChsMutex.RLock()
	defer ge.errorChsMutex.RUnlock()

	ge.errorChs = append(ge.errorChs, ch)
}

func (ge *goExecutor) PostJob(job interfaces.JobFn) error {
	if ge.ctx.Err() != nil {
		return ErrExecutorStopped
	}

	jobSpec := &jobImpl{
		jobFn: job,
	}

	ge.queue.PushBack(jobSpec)

	ge.hasJobsEvent.SetOne()
	return nil
}

func (ge *goExecutor) CollectChan(jobs ...interfaces.JobWithResultFn) <-chan interface{} {
	ch := make(chan interface{})

	results := sync.Map{}
	var currentPos int
	closed := false
	publishMutex := &sync.Mutex{}

	publish := func() {
		publishMutex.Lock()
		defer publishMutex.Unlock()

		for ; currentPos < len(jobs); currentPos++ {
			r, prs := results.Load(currentPos)
			if !prs {
				return
			}
			ch <- r
		}

		if !closed && currentPos == len(jobs) {
			close(ch)
			closed = true
		}
	}

	for i, job := range jobs {
		pos := i
		jobFn := job
		jobSpec := &jobImpl{
			jobFn: func(ctx context.Context) error {
				r, err := jobFn(ctx)

				if ctx.Err() == nil {
					results.Store(pos, r)
					go publish()
				}

				return err
			},
		}

		ge.queue.PushBack(jobSpec)
	}

	ge.hasJobsEvent.Set()

	return (<-chan interface{})(ch)
}

func (ge *goExecutor) Collect(jobs ...interfaces.JobWithResultFn) ([]interface{}, error) {
	ch := ge.CollectChan(jobs...)
	results := []interface{}{}

	for r := range ch {
		results = append(results, r)
	}

	return results, nil
}

func (ge *goExecutor) Len() int {
	return ge.queue.Size()
}
