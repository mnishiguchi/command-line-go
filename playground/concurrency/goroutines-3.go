package main

import (
	"fmt"
	"sync"
)

// worker processes jobs from the queue and prints received values.
// Demonstrates the **fan-out** pattern: multiple goroutines consume from the same channel.
// Use Case: Task processing, load balancing, parallel job execution.
func worker(workerID int, jobQueue <-chan int, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this worker as done when finished

	// Receive and process jobs until the jobQueue is closed.
	for job := range jobQueue {
		fmt.Printf("Worker %d processed job %d\n", workerID, job)
	}
}

func main() {
	// WaitGroup ensures main() waits for all workers to finish.
	var wg sync.WaitGroup

	// Buffered channel (jobQueue) stores jobs to be processed.
	// Use Case: Prevents excessive blocking when producing jobs faster than consuming.
	jobQueue := make(chan int, 10)

	// Number of worker goroutines (concurrent job processors)
	numWorkers := 3

	// Increment WaitGroup counter for each worker
	wg.Add(numWorkers)

	// Start multiple worker goroutines
	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobQueue, &wg) // Launch worker
	}

	// Producer: Sends 100 jobs into the jobQueue
	for i := 0; i < 100; i++ {
		jobQueue <- i
	}

	close(jobQueue) // Close channel to signal workers no more jobs are coming.

	// Wait until all workers finish processing their jobs
	wg.Wait()
	fmt.Println("All workers have finished processing.")
}
