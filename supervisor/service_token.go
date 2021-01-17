package supervisor

import (
	"context"
	"sync"

	"github.com/GustavoKatel/asyncutils/supervisor/interfaces"
)

var _ interfaces.ServiceToken = (*serviceTokenImpl)(nil)

type serviceTokenImpl struct {
	errorsCount int64
	panicsCount int64

	ctx       context.Context
	ctxCancel context.CancelFunc

	dataMutex *sync.RWMutex

	service interfaces.Service
}

// ServiceTokenNew creates a new service token
func ServiceTokenNew(parentCtx context.Context, service interfaces.Service) (interfaces.ServiceToken, error) {
	ctx, cancel := context.WithCancel(parentCtx)

	return &serviceTokenImpl{
		errorsCount: 0,
		panicsCount: 0,

		ctx:       ctx,
		ctxCancel: cancel,

		dataMutex: &sync.RWMutex{},

		service: service,
	}, nil
}

func (st *serviceTokenImpl) ErrorsCount() int64 {
	st.dataMutex.RLock()
	defer st.dataMutex.RUnlock()

	return st.errorsCount
}

func (st *serviceTokenImpl) PanicsCount() int64 {
	st.dataMutex.RLock()
	defer st.dataMutex.RUnlock()

	return st.panicsCount
}

func (st *serviceTokenImpl) Context() context.Context {
	return st.ctx
}

func (st *serviceTokenImpl) Stop() {
	st.ctxCancel()
}

func (st *serviceTokenImpl) Service() interfaces.Service {
	return st.service
}

func (st *serviceTokenImpl) IncErrors() {
	st.dataMutex.Lock()
	defer st.dataMutex.Unlock()

	st.errorsCount++
}

func (st *serviceTokenImpl) IncPanics() {
	st.dataMutex.Lock()
	defer st.dataMutex.Unlock()

	st.panicsCount++
}
