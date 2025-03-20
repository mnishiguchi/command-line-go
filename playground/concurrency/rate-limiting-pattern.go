package main

import (
	"fmt"
	"time"
)

/*
========================================================
💡 Go Concurrency Demo: Rate Limiting Pattern
========================================================
- This program demonstrates **Rate Limiting** using a ticker.
- Tasks execute at **a controlled pace** to avoid overloading the system.

✅ Use Cases:
   - API request throttling (prevent exceeding request limits)
   - Controlled job execution (batch processing, background tasks)
   - Regulating data ingestion to prevent server overload
*/

// simulates processing a task with controlled execution speed.
func processTask(id int) {
	fmt.Printf("✅ Task %d executed at %s\n", id, time.Now().Format("15:04:05.000"))
}

func main() {
	const totalTasks = 10    // Number of tasks to execute
	const tasksPerSecond = 2 // Number of tasks allowed per second

	ticker := time.NewTicker(time.Second / tasksPerSecond) // Creates a pacing mechanism
	defer ticker.Stop()                                    // Stop the ticker when done

	fmt.Println("🚀 Starting rate-limited task execution...")

	for taskID := 1; taskID <= totalTasks; taskID++ {
		<-ticker.C // Wait for the next tick (controls execution rate)
		processTask(taskID)
	}

	fmt.Println("\n🏁 All tasks executed at a controlled rate.")
}
