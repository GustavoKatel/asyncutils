package interfaces

import "context"

// ServiceToken manages a service running in the supervisor
type ServiceToken interface {
	ErrorsCount() int64
	PanicsCount() int64

	Context() context.Context

	Stop()

	Service() Service
}
