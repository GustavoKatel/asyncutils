package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushBack(t *testing.T) {
	assert := assert.New(t)
	q := New()

	assert.Equal(0, q.Size())

	q.PushBack("abc")

	assert.Equal(1, q.Size())

	q.PushBack("back")

	assert.Equal(2, q.Size())

	assert.Equal("abc", q.Get(0))
	assert.Equal("back", q.Get(1))
}

func TestPushFront(t *testing.T) {
	assert := assert.New(t)
	q := New()

	assert.Equal(0, q.Size())

	q.PushFront("abc")

	assert.Equal(1, q.Size())

	q.PushFront("front")

	assert.Equal(2, q.Size())

	assert.Equal("front", q.Get(0))
	assert.Equal("abc", q.Get(1))
}

func TestPopBack(t *testing.T) {
	assert := assert.New(t)
	q := New()

	assert.Equal(0, q.Size())

	q.PushBack("abc")
	q.PushBack("back")

	assert.Equal(2, q.Size())

	assert.Equal("back", q.PopBack())
	assert.Equal(1, q.Size())

	assert.Equal("abc", q.PopBack())
	assert.Equal(0, q.Size())

	assert.Nil(q.PopBack())
	assert.Equal(0, q.Size())
}

func TestPopFront(t *testing.T) {
	assert := assert.New(t)
	q := New()

	assert.Equal(0, q.Size())

	q.PushBack("abc")
	q.PushBack("back")

	assert.Equal(2, q.Size())

	assert.Equal("abc", q.PopFront())
	assert.Equal(1, q.Size())

	assert.Equal("back", q.PopFront())
	assert.Equal(0, q.Size())

	assert.Nil(q.PopFront())
	assert.Equal(0, q.Size())
}

func TestGetNil(t *testing.T) {
	assert := assert.New(t)
	q := New()

	assert.Equal(0, q.Size())

	assert.Nil(q.Get(0))

	q.PushBack("abc")
	assert.Equal(1, q.Size())

	assert.Equal("abc", q.Get(0))

	assert.Equal(1, q.Size())

	assert.Nil(q.Get(1))
	assert.Nil(q.Get(2))
}
