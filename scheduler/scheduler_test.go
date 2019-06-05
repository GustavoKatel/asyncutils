package scheduler

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/GustavoKatel/asyncutils/scheduler/interfaces"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func setUpSchedulerTest(t *testing.T) (interfaces.Scheduler, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	aw, err := New()
	assert.Nil(err)

	return aw, ctrl
}

func setDownSchedulerTest(ctrl *gomock.Controller) {
	ctrl.Finish()
}

func TestStartStop(t *testing.T) {
	aw, ctrl := setUpSchedulerTest(t)
	defer setDownSchedulerTest(ctrl)
	assert := assert.New(t)

	assert.Nil(aw.Start())
	assert.Nil(aw.Stop())
}

func TestRunOne(t *testing.T) {
	aw, ctrl := setUpSchedulerTest(t)
	defer setDownSchedulerTest(ctrl)
	assert := assert.New(t)

	assert.Nil(aw.Start())
	defer aw.Stop()

	waitCh := make(chan interface{}, 1)
	assert.Nil(aw.PostJob(func(ctx context.Context) error {
		waitCh <- nil
		return nil
	}))

	<-waitCh
}

func TestErrorOne(t *testing.T) {
	aw, ctrl := setUpSchedulerTest(t)
	defer setDownSchedulerTest(ctrl)
	assert := assert.New(t)

	assert.Nil(aw.Start())
	defer aw.Stop()

	errCh := make(chan error, 1)
	aw.ErrorChan(errCh)

	assert.Nil(aw.PostJob(func(ctx context.Context) error {
		return fmt.Errorf("Test")
	}))

	err := <-errCh
	assert.Equal(err.Error(), "Test")
}

func TestLen(t *testing.T) {
	aw, ctrl := setUpSchedulerTest(t)
	defer setDownSchedulerTest(ctrl)
	assert := assert.New(t)

	assert.Nil(aw.Start())
	defer aw.Stop()

	waitCh := make(chan interface{}, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	assert.Nil(aw.PostJob(func(ctx context.Context) error {
		<-waitCh
		wg.Done()
		return nil
	}))

	assert.Nil(aw.PostJob(func(ctx context.Context) error {
		<-waitCh
		wg.Done()
		return nil
	}))

	assert.Equal(2, aw.Len())
	waitCh <- nil
	waitCh <- nil
	wg.Wait()
	assert.Equal(0, aw.Len())
}

func TestRunThrottledOne(t *testing.T) {
	aw, ctrl := setUpSchedulerTest(t)
	defer setDownSchedulerTest(ctrl)
	assert := assert.New(t)

	assert.Nil(aw.Start())
	defer aw.Stop()

	var wg sync.WaitGroup
	wg.Add(2)

	count := int64(0)

	start := time.Now()

	assert.Nil(aw.PostThrottledJob(func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		wg.Done()
		return nil
	}, 1*time.Second))

	assert.Nil(aw.PostThrottledJob(func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		wg.Done()
		return nil
	}, 1*time.Second))

	wg.Wait()
	diff := time.Now().Sub(start)
	assert.True(diff >= 1*time.Second)

	assert.Equal(int64(2), count)
}

func TestRunThrottledTwo(t *testing.T) {
	aw, ctrl := setUpSchedulerTest(t)
	defer setDownSchedulerTest(ctrl)
	assert := assert.New(t)

	assert.Nil(aw.Start())
	defer aw.Stop()

	var wg sync.WaitGroup
	wg.Add(2)

	count := int64(0)

	assert.Nil(aw.PostThrottledJob(func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		wg.Done()
		return nil
	}, 500*time.Millisecond))

	<-time.After(501 * time.Millisecond)

	assert.Nil(aw.PostThrottledJob(func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		wg.Done()
		return nil
	}, 500*time.Millisecond))

	wg.Wait()

	assert.Equal(int64(2), count)
}

func TestRunThrottledLastExec(t *testing.T) {
	aw, ctrl := setUpSchedulerTest(t)
	defer setDownSchedulerTest(ctrl)
	assert := assert.New(t)

	assert.Nil(aw.Start())
	defer aw.Stop()

	var wg sync.WaitGroup
	wg.Add(2)

	count := int64(0)

	assert.Nil(aw.PostThrottledJob(func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		<-time.After(300 * time.Millisecond)
		wg.Done()
		return nil
	}, 500*time.Millisecond))

	assert.Nil(aw.PostThrottledJob(func(ctx context.Context) error {
		atomic.AddInt64(&count, 2)
		wg.Done()
		return nil
	}, 500*time.Millisecond))

	assert.Nil(aw.PostThrottledJob(func(ctx context.Context) error {
		atomic.AddInt64(&count, 3)
		wg.Done()
		return nil
	}, 500*time.Millisecond))

	assert.Equal(1, aw.Len())
	wg.Wait()
	assert.Equal(int64(4), count)
	wg.Add(1)

	assert.Nil(aw.PostThrottledJob(func(ctx context.Context) error {
		atomic.AddInt64(&count, 10)
		wg.Done()
		return nil
	}, 500*time.Millisecond))

	wg.Wait()

	assert.Equal(int64(14), count)
}
