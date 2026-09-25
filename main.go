package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Seed the table with initial data.
	table = append(table, User{1, "Alice"})
	table = append(table, User{2, "Bob"})
	table = append(table, User{3, "Charlie"})
	table = append(table, User{3, "Charlie"})
	table = append(table, User{4, "David"})
	table = append(table, User{5, "Eve"})

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
		insertUser(6, "Frank")
	}()

	// Wait for both goroutines to finish (like writer.join(); builder.join()).
	wg.Wait()

	// Print the final index.
	fmt.Println("\nFinal index:")
	for name, id := range idx {
		fmt.Printf("%s -> %d\n", name, id)
	}
}