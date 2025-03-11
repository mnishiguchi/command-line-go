package main

import (
	"fmt"
	"time"
)

/*
========================================================
💡 Go Concurrency Demo: Pipeline Pattern
========================================================
- This program demonstrates the **Pipeline Pattern** in Go.
- Data flows through multiple processing stages, each transforming it.
- Each stage runs concurrently using goroutines.

✅ Use Cases:
   - Data transformation pipelines (e.g., parsing → processing → saving)
   - Streaming data processing (e.g., logs, analytics, event processing)
   - Multi-step computations (e.g., calculations, filtering, aggregation)
*/

// Stage 1: Generate numbers (Producer)
func generateNumbers(count int) <-chan int {
	out := make(chan int) // Create an output channel

	go func() {
		for i := 1; i <= count; i++ {
			out <- i // Send number to next stage
		}

		close(out) // Close the channel after sending all numbers
	}()

	return out
}

// Stage 2: Square numbers
func squareNumbers(in <-chan int) <-chan int {
	out := make(chan int) // Create an output channel

	go func() {
		for num := range in { // Receive numbers from previous stage
			out <- num * num // Square the number and send it forward
		}

		close(out) // Close the channel after processing
	}()

	return out
}

// Stage 3: Convert numbers to formatted strings
func formatResults(in <-chan int) <-chan string {
	out := make(chan string) // Create an output channel

	go func() {
		for num := range in { // Receive squared numbers
			out <- fmt.Sprintf("Result: %d", num) // Format as string
		}

		close(out) // Close the channel after processing
	}()

	return out
}

func main() {
	startTime := time.Now() // Start measuring execution time

	// ✅ Step 1: Generate numbers
	numbers := generateNumbers(5)

	// ✅ Step 2: Square numbers
	squared := squareNumbers(numbers)

	// ✅ Step 3: Format the results as strings
	results := formatResults(squared)

	// ✅ Step 4: Collect and print results
	for result := range results {
		fmt.Println(result)
	}

	// ✅ Step 5: Measure execution time
	elapsedTime := time.Since(startTime)
	fmt.Printf("\n🏁 Pipeline completed in %v\n", elapsedTime)
}
