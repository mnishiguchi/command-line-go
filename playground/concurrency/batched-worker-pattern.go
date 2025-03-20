package main

import (
	"fmt"
	"time"
)

/*
========================================================
💡 Go Concurrency Demo: Batched Worker Pattern
========================================================
- Workers process tasks in **batches** instead of individually.
- Ensures efficient task handling, reducing overhead.

✅ Use Cases:
   - Bulk database inserts
   - Efficient log/message processing
   - Reducing network/API call frequency
*/

const (
	totalTasks  = 16              // Number of tasks
	batchSize   = 5               // Max tasks per batch
	workerCount = 3               // Number of workers
	batchDelay  = 2 * time.Second // Simulated batch processing time
)

// worker processes tasks in batches and sends results.
func worker(workerID int, taskChannel <-chan int, resultChannel chan<- string, doneChannel chan<- bool) {
	var batchBuffer []int

	for task := range taskChannel {
		batchBuffer = append(batchBuffer, task)

		// Process batch when full
		if len(batchBuffer) == batchSize {
			fmt.Printf("🔄 Worker %d processing batch: %v\n", workerID, batchBuffer)
			time.Sleep(batchDelay)
			resultChannel <- fmt.Sprintf("✅ Worker %d completed batch: %v", workerID, batchBuffer)
			batchBuffer = nil // Reset batch buffer
		}
	}

	// Process remaining tasks
	if len(batchBuffer) > 0 {
		fmt.Printf("🔄 Worker %d processing final batch: %v\n", workerID, batchBuffer)
		time.Sleep(batchDelay)
		resultChannel <- fmt.Sprintf("✅ Worker %d completed final batch: %v", workerID, batchBuffer)
	}

	doneChannel <- true // Notify main that worker is done
}

func main() {
	taskChannel := make(chan int, totalTasks)
	resultChannel := make(chan string, totalTasks)
	doneChannel := make(chan bool, workerCount)

	// Start workers
	for i := 1; i <= workerCount; i++ {
		go worker(i, taskChannel, resultChannel, doneChannel)
	}

	// Send tasks to workers
	for i := 1; i <= totalTasks; i++ {
		taskChannel <- i
	}
	close(taskChannel)

	// Wait for all workers to finish before closing resultChannel
	for i := 1; i <= workerCount; i++ {
		<-doneChannel
	}
	close(resultChannel)

	// Read results
	for result := range resultChannel {
		fmt.Println(result)
	}

	fmt.Println("\n🏁 All tasks processed in batches.")
}

