package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2) // Set counter for 2 goroutines

	go func() {
		defer wg.Done() // Decrement counter when done
		fmt.Println("Hello")
	}()

	go func() {
		defer wg.Done() // Decrement counter when done
		fmt.Println("World")
	}()

	wg.Wait() // Wait for all goroutines to finish
}
