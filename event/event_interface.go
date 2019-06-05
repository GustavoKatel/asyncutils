package event

import "time"

// Event synchronizes goroutines with a set-reset flag style
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
