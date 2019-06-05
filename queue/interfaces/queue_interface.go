package interfaces

// Queue thread-safe queue implementation
type Queue interface {
	PushBack(el interface{})
	PushFront(el interface{})

	// PopBack removes an element from the back of the queue. Returns nil if queue is empty
	PopBack() interface{}
	// PopFront removes an element from the head of the queue. Returns nil if queue is empty
	PopFront() interface{}

	// Get returns an element in position "pos" or nil if "pos" is out of bounds
	Get(pos int) interface{}

	Size() int
}
