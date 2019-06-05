package interfaces

// Executor interface
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
