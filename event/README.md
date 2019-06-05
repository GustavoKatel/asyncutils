# Event

[![GoDoc](https://godoc.org/github.com/GustavoKatel/asyncutils/event?status.svg)](https://godoc.org/github.com/GustavoKatel/asyncutils/event)

Event synchronizes goroutines with a set-reset flag style

```go
type Event interface {
	// IsSet returns true if set has been called
	IsSet() bool

	// Set sets the flag to true and awake pending goroutines
	Set()

	// Wait waits this flag to be set
	Wait()

	// Reset resets this flag
	Reset()
}
```