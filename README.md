# AsyncUtils

[![GoDoc](https://godoc.org/github.com/GustavoKatel/asyncutils?status.svg)](https://godoc.org/github.com/GustavoKatel/asyncutils)
![Main](https://github.com/GustavoKatel/asyncutils/workflows/Main/badge.svg)

Synchornization and asynchronous operations utilities in golang

## Queue

Thread-safe queue implementation

```go
type Queue interface {
	type Queue interface {
	PushBack(el interface{})
	PushFront(el interface{})

	// PopBack removes an element from the back of the queue. Returns nil if queue is empty
	PopBack() interface{}

	// PopFront removes an element from the head of the queue. Returns nil if queue is empty
	PopFront() interface{}

	// Get returns an element in position "pos" or nil if "pos" is out of bounds
	Get(pos int) interface{}

	Size() int
}
}
```

## Event

Event synchronizes goroutines with a set-reset flag style

```go
type Event interface {
	// IsSet returns true if set has been called
	IsSet() bool

	// Set sets the flag to true and awake pending goroutines
	Set()

	// SetOne sets the flag to true and awake only one pending goroutines
	SetOne()

	// Wait waits this flag to be set
	Wait()

	// WaitTimeout waits this flag to be set or timeout
	WaitTimeout(d time.Duration)

	// Reset resets this flag
	Reset()
}
```

## Executor

Asynchronous function execution

```go
type Executor interface {
	// Start starts the executor
	Start() error

	// Stop stops the executor and all the pending jobs
	Stop() error

	// PostJob enqueue a job
	PostJob(job JobFn) error

	// Collect executes all jobs posted and return the results in order
	// if an error happens, the resulting slice will contain less elements than jobs
	// please check ErrorChan
	Collect(jobs ...JobWithResultFn) ([]interface{}, error)

	// CollectChan same as Collect but return a channel with the results
	CollectChan(jobs ...JobWithResultFn) <-chan interface{}

	// ErrorChan registers an error emitting channel
	ErrorChan(ch chan error)

	// Len size of the pending queue
	Len() int
}
```

### Example:

#### Collect results

`Collect` and `CollectChan` keeps the order of the results (Similar to `Promise.all` in js)

```go
// One worker (goroutine)
exc, err := NewDefaultExecutor(1)
assert.Nil(err)

assert.Nil(exc.Start())
defer exc.Stop()

job1 := func(ctx context.Context) (interface{}, error) {
    <-time.After(500 * time.Millisecond)
    return 1, nil
}
job2 := func(ctx context.Context) (interface{}, error) {
    return 2, nil
}

results, err := exc.Collect(job1, job2)
assert.Nil(err)
assert.Equal(2, len(results))
assert.Equal(1, results[0])
assert.Equal(2, results[1])
```

#### Enqueue

```go
// Two workers (goroutines)
exc, err := NewDefaultExecutor(2)
assert.Nil(err)

assert.Nil(exc.Start())
defer exc.Stop()

results := make(chan int, 2)
exc.PostJob(func(ctx context.Context) error {
    <-time.After(500 * time.Millisecond)
    results <- 2
    return nil
})

exc.PostJob(func(ctx context.Context) error {
    results <- 1
    return nil
})

re := <-results
assert.Equal(1, re)

re = <-results
assert.Equal(2, re)
```

## Scheduler

Scheduler and throttler. See [Scheduler](https://github.com/GustavoKatel/asyncutils/blob/master/scheduler/README.md)

```go
type Scheduler interface {
	// Start starts the worker channel
	Start() error

	// Stop stops all running and scheduled jobs
	Stop() error

	// PostJob schedules a job execution
	PostJob(job executorIfaces.JobFn) error

	// PostThrottledJob posts a job only and only if the time span of its last execution was greater than "duration"
	PostThrottledJob(job executorIfaces.JobFn, delay time.Duration) error

	// Len returns the number of jobs scheduled
	Len() int

	// ErrorChan registers an error emitting channel
	ErrorChan(ch chan error)
}
```

