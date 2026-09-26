package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Seed the table with initial data. pre-existing data that is visible to everyone.
	table = append(table, User{1, "Alice", 0})
	table = append(table, User{2, "Bob", 0})
	table = append(table, User{3, "Charlie", 0})
	table = append(table, User{4, "Charlie", 0})
	table = append(table, User{5, "David", 0})
	table = append(table, User{6, "Eve", 0})

	var wg sync.WaitGroup
	wg.Add(2)

	// Start the index builder.
	go func() {
		defer wg.Done()
		buildIndex()
	}()

	// Let the builder take a snapshot before the writer starts (2 seconds).
	time.Sleep(2 * time.Second)

	// Start a writer that inserts "Frank" (6, "Frank").
	go func() {
		defer wg.Done()
		insertUser(7, "Frank")
	}()

	// Wait for both goroutines to finish (like writer.join(); builder.join()).
	wg.Wait()

	// Print the final index.
	fmt.Println("\nFinal index:")
	for name, ids := range idx {
		fmt.Printf("%s -> %v\n", name, ids)
	}
}
