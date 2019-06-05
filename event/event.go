package event

import (
	"context"
	"sync"
	"time"
)

var _ Event = &eventImpl{}

// NewEvent creates a new sync event flag
func NewEvent(initValue bool) Event {
	mutex := &sync.RWMutex{}

	return &eventImpl{
		flagMutex: mutex,
		flag:      initValue,
		cond:      sync.NewCond(mutex.RLocker()),
	}
}

type eventImpl struct {
	flagMutex *sync.RWMutex
	flag      bool
	cond      *sync.Cond
}

func (s *eventImpl) IsSet() bool {
	s.flagMutex.RLock()
	defer s.flagMutex.RUnlock()
	return s.flag
}

func (s *eventImpl) setFlag() {
	s.flagMutex.Lock()
	defer s.flagMutex.Unlock()
	s.flag = true
}

func (s *eventImpl) Set() {
	s.setFlag()
	s.cond.Broadcast()
}

func (s *eventImpl) SetOne() {
	s.setFlag()
	s.cond.Signal()
}

func (s *eventImpl) Wait() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for !s.flag {
		s.cond.Wait()
	}
}

func (s *eventImpl) WaitTimeout(d time.Duration) {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		defer cancel()
		for !s.flag && ctx.Err() == nil {
			s.cond.Wait()
		}
	}()

	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}

func (s *eventImpl) Reset() {
	s.flagMutex.Lock()
	defer s.flagMutex.Unlock()
	s.flag = false
}
