package interfaces

import "context"

// ServiceFunc provides an interface to a service run funtion
type ServiceFunc func(ctx context.Context) error

// Service interface
type Service interface {
	Init(ctx context.Context) error
	Clean(ctx context.Context) error

	Run(ctx context.Context) error
}
