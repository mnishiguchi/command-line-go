package main

import (
	"fmt"
	"sync"
)

// printHello prints "Hello" along with the goroutine ID.
func printHello(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this goroutine as done when it finishes
	fmt.Printf("Hello from goroutine %d\n", id)
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2) // Set the wait counter for two goroutines

	// Start two concurrent goroutines
	go printHello(1, &wg)
	go printHello(2, &wg)

	wg.Wait() // Wait until all goroutines finish
}
