package event

import "time"

type EventWaiter interface {
	// Wait waits this flag to be set
	Wait()

	// WaitTimeout waits this flag to be set or timeout
	WaitTimeout(d time.Duration)
}

// Event synchronizes goroutines with a set-reset flag style
type Event interface {
	EventWaiter

	// IsSet returns true if set has been called
	IsSet() bool

	// Set sets the flag to true and awake pending goroutines
	Set()

	// SetOne sets the flag to true and awake only one pending goroutines
	SetOne()

	// Reset resets this flag
	Reset()
}
