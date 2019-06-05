package interfaces

// Queue thread-safe queue implementation
type Queue interface {
	PushBack(el interface{})
	PushFront(el interface{})

	PopBack() interface{}
	PopFront() interface{}

	Get(pos int) interface{}

	Size() int
}
