package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
========================================================
💡 Go Concurrency Demo: Fan-In Pattern
========================================================
- This program demonstrates the **Fan-In** concurrency pattern in Go.
- Multiple goroutines (producers) send data into a shared results channel.
- A single goroutine (consumer) aggregates results from the producers.

✅ Use Cases:
   - Aggregating API responses from multiple microservices
   - Merging logs from multiple sources
   - Parallel computations with results collection
*/

const numProducers = 3 // Number of concurrent producer goroutines
const numMessages = 5  // Number of messages each producer sends

// producer simulates generating messages and sending them to the results channel.
// Each producer waits for a random delay before sending each message to simulate real-world delays.
func producer(id int, results chan<- string, wg *sync.WaitGroup, rng *rand.Rand) {
	defer wg.Done() // Mark this producer as finished when function exits

	for i := 1; i <= numMessages; i++ {
		// Simulate processing time (random delay between 500-1000ms)
		processTime := time.Duration(rng.Intn(500)+500) * time.Millisecond
		time.Sleep(processTime) // Simulate a time-consuming operation

		// Send message to the shared results channel
		message := fmt.Sprintf("📩 Producer %d sent message %d after %v", id, i, processTime)
		results <- message
	}
	fmt.Printf("[Producer] ✅ Producer %d has finished sending messages\n", id)
}

func main() {
	// ✅ Create a new random generator to avoid deprecated rand.Seed()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// ✅ Create a buffered channel for results to prevent blocking
	results := make(chan string, numProducers*numMessages) // Stores messages from all producers
	var wg sync.WaitGroup                                  // Used to wait for all producers to finish

	startTime := time.Now() // Start tracking execution time

	// ✅ Step 1: Launch multiple producer goroutines (Fan-Out)
	wg.Add(numProducers) // Ensure all producers are counted in the WaitGroup
	for i := 1; i <= numProducers; i++ {
		go producer(i, results, &wg, rng) // Each producer runs concurrently
	}

	// ✅ Step 2: Launch a goroutine to close the results channel when all producers finish
	go func() {
		wg.Wait()      // Wait for all producers to complete
		close(results) // Close the results channel to signal no more messages will be sent
		fmt.Println("[Main] All producers finished. Closing results channel.")
	}()

	// ✅ Step 3: Read results from the results channel (Fan-In)
	for message := range results {
		fmt.Println("[Consumer] Received:", message) // Aggregates results
	}

	// ✅ Step 4: Measure execution time
	elapsedTime := time.Since(startTime)
	fmt.Printf("\n🏁 All messages received in %v\n", elapsedTime)
}
