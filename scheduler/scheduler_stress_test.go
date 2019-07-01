package scheduler

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/GustavoKatel/asyncutils/event"
	"github.com/stretchr/testify/assert"
)

type config struct {
	Name    string
	Verbose bool

	Producers int
	K         int
	UnitK     int
	Delay     time.Duration
	MinExec   int
}

func _log(c *config, f string, args ...interface{}) {
	if c.Verbose {
		format := fmt.Sprintf("[%s] %s", c.Name, f)
		log.Printf(format, args...)
	}
}

func doFixedStressTest(t *testing.T, c *config) {
	event := event.NewEvent(false)

	ctx, cancel := context.WithCancel(context.Background())

	scheduler, err := NewWithContext(ctx)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	if err := scheduler.Start(); err != nil {
		t.Fatalf("Err: %v", err)
	}

	var count int64
	var postCount int64

	var wg sync.WaitGroup
	wg.Add(c.Producers)

	for i := 0; i < c.Producers; i++ {
		go func(id int) {
			event.Wait()
			d := rand.Intn(c.K)
			<-time.After(time.Duration(d) * time.Millisecond)

			jobctx, cancel := context.WithTimeout(ctx, 25*time.Second)

			go func() {
				<-jobctx.Done()
				wg.Done()
			}()

			_log(c, "Post: %d", id)
			scheduler.PostThrottledJob(
				func(ctx context.Context) error {
					<-time.After(time.Duration(c.UnitK) * time.Millisecond)
					atomic.AddInt64(&count, 1)
					_log(c, "Done: %v", id)
					cancel()
					return nil
				},
				c.Delay,
			)
			atomic.AddInt64(&postCount, 1)

		}(i)
	}

	event.Set()

	wg.Wait()

	cancel()

	_log(c, "Total posted: %v", postCount)
	_log(c, "Total executed: %v", count)
	assert := assert.New(t)
	assert.True(count >= int64(c.MinExec), "Expected to have executed at least %v", c.MinExec)
	assert.True(postCount >= int64(c.Producers), "Expected to have posted %v", c.Producers)
}

func TestRunThrottled100(t *testing.T) {
	doFixedStressTest(t, &config{
		Name: "TestRunThrottled100",

		Producers: 100,
		K:         10,
		UnitK:     1500,
		Delay:     1 * time.Second,
		MinExec:   2,
	})
}

func TestRunThrottled100KHigh(t *testing.T) {
	doFixedStressTest(t, &config{
		Name: "TestRunThrottled100KHigh",

		Producers: 100,
		K:         1500,
		UnitK:     15000,
		Delay:     1 * time.Second,
		MinExec:   1,
	})
}

func TestRunThrottled1000(t *testing.T) {
	doFixedStressTest(t, &config{
		Name: "TestRunThrottled1000",

		Producers: 1000,
		K:         5000,
		UnitK:     500,
		Delay:     1 * time.Second,
		MinExec:   1,
	})
}

func TestRunThrottled1000KLow(t *testing.T) {
	doFixedStressTest(t, &config{
		Name: "TestRunThrottled1000KLow",

		Producers: 1000,
		K:         10,
		UnitK:     500,
		Delay:     1 * time.Second,
		MinExec:   1,
	})
}
