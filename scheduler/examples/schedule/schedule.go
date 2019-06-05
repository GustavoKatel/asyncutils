package main

import (
	"context"
	"log"
	"sync"

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

	worker.PostJob(func(ctx context.Context) error {
		// Long operation 2
		log.Printf("Operation2")
		wg.Done()
		return nil
	})

	log.Printf("Pending: %v", worker.Len())
	wg.Wait()
	log.Printf("Pending: %v", worker.Len())
}
