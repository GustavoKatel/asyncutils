package executor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartStop(t *testing.T) {
	assert := assert.New(t)

	exc, err := NewDefaultExecutor(2)
	assert.Nil(err)

	assert.Nil(exc.Start())
	assert.Nil(exc.Stop())
}

func TestEnqueue(t *testing.T) {
	assert := assert.New(t)

	exc, err := NewDefaultExecutor(2)
	assert.Nil(err)

	assert.Nil(exc.Start())
	defer exc.Stop()

	results := make(chan int, 2)
	assert.Nil(exc.PostJob(func(ctx context.Context) error {
		<-time.After(500 * time.Millisecond)
		results <- 2
		return nil
	}))
	assert.Nil(exc.PostJob(func(ctx context.Context) error {
		results <- 1
		return nil
	}))

	re := <-results
	assert.Equal(1, re)

	re = <-results
	assert.Equal(2, re)
}

func TestEnqueueWithErr(t *testing.T) {
	assert := assert.New(t)

	exc, err := NewDefaultExecutor(2)
	assert.Nil(err)

	errCh := make(chan error, 1)
	exc.ErrorChan(errCh)

	assert.Nil(exc.Start())
	defer exc.Stop()

	results := make(chan int, 2)
	assert.Nil(exc.PostJob(func(ctx context.Context) error {
		<-time.After(500 * time.Millisecond)
		results <- 2
		return fmt.Errorf("test")
	}))
	assert.Nil(exc.PostJob(func(ctx context.Context) error {
		results <- 1
		return nil
	}))

	re := <-results
	assert.Equal(1, re)

	re = <-results
	assert.Equal(2, re)

	err = <-errCh
	assert.Equal("test", err.Error())
}

func TestEnqueueStopped(t *testing.T) {
	assert := assert.New(t)

	exc, err := NewDefaultExecutor(2)
	assert.Nil(err)

	assert.Nil(exc.Start())
	assert.Nil(exc.Stop())

	assert.NotNil(exc.PostJob(func(ctx context.Context) error {
		<-time.After(500 * time.Millisecond)
		return nil
	}))
}

func TestCollectChan(t *testing.T) {
	assert := assert.New(t)

	exc, err := NewDefaultExecutor(1)
	assert.Nil(err)

	assert.Nil(exc.Start())
	defer exc.Stop()

	job1 := func(ctx context.Context) (interface{}, error) {
		<-time.After(500 * time.Millisecond)
		return 1, nil
	}
	job2 := func(ctx context.Context) (interface{}, error) {
		return 2, nil
	}

	results := exc.CollectChan(job1, job2)
	r := <-results
	assert.Equal(1, r)

	r = <-results
	assert.Equal(2, r)

	r, ok := <-results
	assert.False(ok)
}

func TestCollectChanFirstServe(t *testing.T) {
	assert := assert.New(t)

	exc, err := NewDefaultExecutor(1)
	assert.Nil(err)

	assert.Nil(exc.Start())
	defer exc.Stop()

	job1 := func(ctx context.Context) (interface{}, error) {
		<-time.After(500 * time.Millisecond)
		return 4, nil
	}
	job2 := func(ctx context.Context) (interface{}, error) {
		return 8, nil
	}

	results := exc.CollectChanFirstServe(job1, job2)
	r := <-results
	assert.Equal(1, r.Index)
	assert.Equal(8, r.Result)

	r = <-results
	assert.Equal(0, r.Index)
	assert.Equal(4, r.Result)

	r, ok := <-results
	assert.False(ok)
}

func TestCollect(t *testing.T) {
	assert := assert.New(t)

	exc, err := NewDefaultExecutor(1)
	assert.Nil(err)

	assert.Nil(exc.Start())
	defer exc.Stop()

	job1 := func(ctx context.Context) (interface{}, error) {
		<-time.After(500 * time.Millisecond)
		return 1, nil
	}
	job2 := func(ctx context.Context) (interface{}, error) {
		return 2, nil
	}

	results, err := exc.Collect(job1, job2)
	assert.Nil(err)
	assert.Equal(2, len(results))
	assert.Equal(1, results[0])
	assert.Equal(2, results[1])
}
