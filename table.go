package main

import (
	"fmt"
	"time"
)

// getCurrentSnapShot returns a copy of the table under a lock.
func getCurrentSnapShot() []User {
	mtx.Lock()
	defer mtx.Unlock()

	// slices the table and copies it into snap
	snap := make([]User, len(table))
	copy(snap, table)
	return snap
}

// insertUser simulates a write transaction: inserts a new user into the table.
// If an index build is in progress, the user is saved in pending so it can be
// merged into the index later (synchronize).
func insertUser(id int, name string) {
	// initialize a new transaction
	tx := beginTransaction()

	fmt.Printf("Writer started: %s\n", name)

	// Simulate a long-running write transaction (3 seconds).
	time.Sleep(3 * time.Second)

	// Commit the new row.
	mtx.Lock()

	user := User{
		ID:   id,
		Name: name,
	}

	table = append(table, user)
	fmt.Printf("Inserted: %s\n", name)

	if building {
		pending = append(pending, user) // remember for later index sync
	}

	mtx.Unlock()

	commitTransaction(tx)
}
