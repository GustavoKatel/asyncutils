package interfaces

// Supervisor provides an interface to supersisor
// Supervisor is responsible for running high-available services
type Supervisor interface {
	// AddErrorHandler adds a new error handler
	AddErrorHandler(handler ErrorHandler) error

	AddService(service Service) (ServiceToken, error)

	Start() error
	Stop() error
}
