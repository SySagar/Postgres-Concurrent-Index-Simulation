package main

import (
	"fmt"
	"time"
)

// getCurrentSnapShot returns a copy of the table under a lock.
func getCurrentSnapShot() Snapshot {
	mtx.Lock()
	defer mtx.Unlock()

	committed := make(map[int]bool)

	for id, tx := range transactions {
		if tx.State == TxCommitted {
			committed[id] = true
		}
	}

	// Initial rows are always visible in our simulator.
	committed[0] = true

	return Snapshot{
		committedTxIDs: committed,
	}
}

func (s Snapshot) IsVisible(user User) bool {
	return s.committedTxIDs[user.CreatedBy]
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
		ID:        id,
		Name:      name,
		CreatedBy: tx.ID,
	}

	table = append(table, user)
	fmt.Printf("Inserted: %s\n", name)

	mtx.Unlock()

	commitTransaction(tx)
}
