package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
========================================================
💡 Go Concurrency Demo: Dynamic Worker Pool Pattern
========================================================
- The number of workers dynamically scales based on workload.
- Workers pick up tasks from a shared queue and exit when no more work is available.

✅ Use Cases:
   - Processing workloads that change over time
   - Efficiently handling variable loads
   - Preventing excessive idle workers or worker starvation
*/

const (
	minWorkers = 2  // Minimum active workers
	maxWorkers = 5  // Maximum allowed workers
	totalTasks = 20 // Total number of tasks
)

// worker processes tasks from the queue
func worker(id int, taskChannel <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChannel {
		fmt.Printf("🔄 Worker %d processing task %d\n", id, task)
		time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond) // Simulated variable work time
		fmt.Printf("✅ Worker %d completed task %d\n", id, task)
	}
}

func main() {
	taskChannel := make(chan int, totalTasks)
	var wg sync.WaitGroup

	// ✅ Determine number of workers dynamically based on workload
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	numWorkers := minWorkers + rng.Intn(maxWorkers-minWorkers+1) // Random between min and max

	fmt.Printf("🚀 Starting with %d dynamic workers\n", numWorkers)

	// ✅ Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, taskChannel, &wg)
	}

	// ✅ Send tasks to the queue
	for i := 1; i <= totalTasks; i++ {
		taskChannel <- i
	}
	close(taskChannel) // No more tasks will be sent

	// ✅ Wait for all workers to complete
	wg.Wait()

	fmt.Println("\n🏁 All tasks processed with dynamic workers.")
}

