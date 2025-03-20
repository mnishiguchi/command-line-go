package main

import (
	"fmt"
	"time"
)

/*
========================================================
💡 Go Concurrency Demo: Timeout Pattern
========================================================
- This program demonstrates **Timeout Handling** in Go.
- If a task takes too long, it gets **automatically canceled**.

✅ Use Cases:
   - Prevent blocking forever (e.g., slow API responses)
   - Limit execution time of computations
   - Ensure the system remains responsive
*/

// slowTask simulates a long-running task.
func slowTask(id int, duration time.Duration, done chan<- string) {
	time.Sleep(duration) // Simulate task duration
	done <- fmt.Sprintf("✅ Task %d completed in %v", id, duration)
}

func main() {
	timeout := 2 * time.Second // Set timeout duration
	done := make(chan string)  // Channel to receive task completion messages

	fmt.Println("🚀 Starting task with timeout...")

	// ✅ Launch the slow task (this task will take 3 seconds)
	go slowTask(1, 3*time.Second, done)

	// ✅ Use `select` to enforce the timeout
	select {
	case result := <-done:
		fmt.Println(result) // Task completed successfully
	case <-time.After(timeout):
		fmt.Println("⏳ Timeout! Task took too long and was canceled.")
	}

	fmt.Println("\n🏁 Program finished.")
}
