package interfaces

import "context"

// JobFn job interface
type JobFn func(ctx context.Context) error

// JobWithResultFn job with have result
type JobWithResultFn func(ctx context.Context) (interface{}, error)
