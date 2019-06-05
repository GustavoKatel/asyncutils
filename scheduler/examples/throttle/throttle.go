package main

import (
	"context"
	"log"
	"sync"
	"time"

	scheduler "github.com/GustavoKatel/asyncutils/scheduler"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	worker, err := scheduler.New()
	if err != nil {
		panic(err)
	}
	worker.Start()
	defer worker.Stop()

	worker.PostJob(func(ctx context.Context) error {
		// Long operation 1
		log.Printf("Operation1")
		wg.Done()
		return nil
	})

	worker.PostThrottledJob(func(ctx context.Context) error {
		// Long operation 2 is not executed due to throttle
		log.Printf("Operation2")
		return nil
	}, 500*time.Millisecond)

	worker.PostThrottledJob(func(ctx context.Context) error {
		// Long operation 3
		log.Printf("Operation3")
		return nil
	}, 500*time.Millisecond)

	<-time.After(600 * time.Millisecond)

	worker.PostThrottledJob(func(ctx context.Context) error {
		// Long operation 4
		log.Printf("Operation4")
		wg.Done()
		return nil
	}, 500*time.Millisecond)

	wg.Wait()
	log.Printf("Pending: %v", worker.Len())
}
