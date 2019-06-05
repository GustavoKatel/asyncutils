package interfaces

import (
	"time"

	executorIfaces "github.com/GustavoKatel/asyncutils/executor/interfaces"
)

// Scheduler provides a scheduling work channel
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
