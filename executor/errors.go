package executor

import "errors"

var (
	// ErrExecutorStopped tried to enqueue a job with the executor stopped
	ErrExecutorStopped = errors.New("Executor is stopped")
)
