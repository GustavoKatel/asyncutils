package queue

import (
	"sync"

	"github.com/GustavoKatel/asyncutils/queue/interfaces"
)

type Queue interfaces.Queue

var _ interfaces.Queue = &queueImpl{}

type queueImpl struct {
	q     []interface{}
	mutex *sync.RWMutex
}

// New creates a new queue
func New() interfaces.Queue {
	return &queueImpl{
		q:     []interface{}{},
		mutex: &sync.RWMutex{},
	}
}

func (q *queueImpl) PushBack(el interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.q = append(q.q, el)
}

func (q *queueImpl) PushFront(el interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.q = append([]interface{}{el}, q.q...)
}

func (q *queueImpl) PopBack() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.q) == 0 {
		return nil
	}

	el := q.q[len(q.q)-1]
	q.q = q.q[:len(q.q)-1]

	return el
}

func (q *queueImpl) PopFront() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.q) == 0 {
		return nil
	}

	el := q.q[0]
	q.q = q.q[1:]

	return el
}

func (q *queueImpl) Get(pos int) interface{} {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if pos >= len(q.q) {
		return nil
	}

	return q.q[pos]
}

func (q *queueImpl) Size() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return len(q.q)
}
