package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
========================================================
💡 Go Concurrency Demo: Fan-Out / Worker Pool Pattern
========================================================
- This program demonstrates the **Fan-Out / Worker Pool** concurrency pattern in Go.
- Multiple worker goroutines consume tasks from a shared job queue (channel).
- This pattern increases **throughput** by parallelizing work execution.

✅ Use Cases:
   - Efficient task processing (e.g., log processing, API request handling)
   - Distributing computational workload across multiple CPU cores
   - Background job execution (e.g., batch database updates)
*/

const (
	numWorkers = 3  // Number of worker goroutines (parallel processors)
	numJobs    = 10 // Total number of jobs to process
)

// worker processes jobs from the jobs channel and sends results to the results channel.
// It simulates processing time by sleeping for a random duration.
func worker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this worker as done when it exits

	for job := range jobs { // Receive jobs from the queue
		processTime := time.Duration(rand.Intn(500)+500) * time.Millisecond // Random delay (500-1000ms)
		time.Sleep(processTime)                                             // Simulate job execution time

		// Generate and send the job completion result
		result := fmt.Sprintf("✅ Worker %d completed job %d in %v", id, job, processTime)
		results <- result
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Initialize random seed for varied processing times

	// ✅ Step 1: Create job & result channels
	jobs := make(chan int, numJobs)       // Job queue (task input)
	results := make(chan string, numJobs) // Results queue (task output)
	var wg sync.WaitGroup                 // WaitGroup to synchronize worker completion

	startTime := time.Now() // Record start time for performance measurement

	// ✅ Step 2: Launch worker goroutines (Fan-Out)
	wg.Add(numWorkers) // Ensure all workers are counted before launching them
	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobs, results, &wg) // Start worker goroutine
	}

	// ✅ Step 3: Send jobs to the job queue (Producer)
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	close(jobs) // ✅ No more jobs will be sent; workers should complete their tasks.

	// ✅ Step 4: Wait for all workers to finish processing
	wg.Wait()
	close(results) // ✅ Close the results channel after all workers have completed.

	// ✅ Step 5: Collect and display results
	for result := range results {
		fmt.Println(result)
	}

	// ✅ Step 6: Measure and display total execution time
	elapsedTime := time.Since(startTime)
	fmt.Printf("\n🏁 All jobs processed in %v\n", elapsedTime)
}
