package supervisor

import (
	"context"
	"fmt"
	"sync"

	"github.com/GustavoKatel/asyncutils/event"

	"github.com/GustavoKatel/asyncutils/supervisor/interfaces"
)

var _ interfaces.Supervisor = (*supervisorImpl)(nil)

type supervisorImpl struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	errorsHandlers      []interfaces.ErrorHandler
	errorsHandlersMutex *sync.RWMutex

	startEvent event.Event
}

// New creates a new supervisor
func New(ctx context.Context) (interfaces.Supervisor, error) {
	ctx, cancel := context.WithCancel(ctx)

	return &supervisorImpl{
		ctx:       ctx,
		ctxCancel: cancel,

		errorsHandlers:      []interfaces.ErrorHandler{},
		errorsHandlersMutex: &sync.RWMutex{},

		startEvent: event.NewEvent(false),
	}, nil
}

func (spv *supervisorImpl) AddErrorHandler(handler interfaces.ErrorHandler) error {
	if handler == nil {
		return ErrHandlerNil
	}

	spv.errorsHandlersMutex.Lock()
	defer spv.errorsHandlersMutex.Unlock()

	spv.errorsHandlers = append(spv.errorsHandlers, handler)

	return nil
}

func (spv *supervisorImpl) Start() error {
	spv.startEvent.Set()

	return nil
}

func (spv *supervisorImpl) Stop() error {
	spv.ctxCancel()

	return nil
}

func (spv *supervisorImpl) AddService(service interfaces.Service) (interfaces.ServiceToken, error) {

	token, err := ServiceTokenNew(spv.ctx, service)
	if err != nil {
		return nil, err
	}

	go spv.serviceWorker(token)

	return token, nil
}

func (spv *supervisorImpl) serviceWorker(token interfaces.ServiceToken) {
	spv.startEvent.Wait()

	defer func() {
		if err := recover(); err != nil {
			spv.publishServiceErr(token, fmt.Errorf("%v", err))

			token.(*serviceTokenImpl).IncPanics()

			if token.Context().Err() == nil {
				go spv.serviceWorker(token)
			}
		}
	}()

	for token.Context().Err() == nil {
		if err := token.Service().Init(token.Context()); err != nil {
			spv.publishServiceErr(token, err)
			// do not call clean if init has errored
			continue
		}

		if err := token.Service().Run(token.Context()); err != nil {
			spv.publishServiceErr(token, err)
		}

		err := token.Service().Clean(token.Context())
		if err != nil {
			spv.publishServiceErr(token, err)
		}
	}
}

func (spv *supervisorImpl) publishServiceErr(token interfaces.ServiceToken, err error) {
	token.(*serviceTokenImpl).IncErrors()

	spv.errorsHandlersMutex.RLock()
	defer spv.errorsHandlersMutex.RUnlock()

	for _, handler := range spv.errorsHandlers {
		handler.OnServiceError(token, err)
	}
}
