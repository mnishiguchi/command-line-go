package main

import (
	"context"
	"fmt"
	"time"
)

/*
========================================================
💡 Go Concurrency Demo: Cancellation Pattern
========================================================
- This program demonstrates **Goroutine Cancellation** using `context.Context`.
- If a cancellation request is received, the task **stops immediately**.

✅ Use Cases:
   - Gracefully shutting down workers when a service stops
   - Canceling API requests when the client disconnects
   - Stopping background tasks when they are no longer needed
*/

// simulates a cancellable task.
func longRunningTask(ctx context.Context, id int) {
	for i := 1; i <= 5; i++ {
		select {
		case <-ctx.Done(): // 🔴 Task is canceled
			fmt.Printf("🛑 Task %d was canceled!\n", id)
			return
		default:
			fmt.Printf("🔄 Task %d processing step %d\n", id, i)
			time.Sleep(1 * time.Second) // Simulate work
		}
	}
	fmt.Printf("✅ Task %d completed successfully\n", id)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // ✅ Create cancelable context
	defer cancel()                                          // Ensure cancel() is called to release resources

	fmt.Println("🚀 Starting cancellable task...")

	// ✅ Launch the task
	go longRunningTask(ctx, 1)

	// ✅ Simulate user action: cancel task after 3 seconds
	time.Sleep(3 * time.Second)
	fmt.Println("🛑 Canceling the task...")
	cancel() // 🔴 Cancel the task

	// ✅ Give some time for cancellation to take effect
	time.Sleep(2 * time.Second)

	fmt.Println("\n🏁 Program finished.")
}
