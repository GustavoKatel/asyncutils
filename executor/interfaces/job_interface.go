package interfaces

import "context"

// JobFn job interface
type JobFn func(ctx context.Context) error

// JobWithResultFn job with have result
type JobWithResultFn func(ctx context.Context) (interface{}, error)

// JobResultIndexed matches a job result to its index in the job list
type JobResultIndexed struct {
	Index  int
	Result interface{}
}
