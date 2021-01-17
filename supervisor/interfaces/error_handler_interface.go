package interfaces

// ErrorHandler is responsible for handling supervisor errors
type ErrorHandler interface {
	OnServiceError(token ServiceToken, err error)
}
