package main

import (
	"fmt"
	"time"
)

// synchronize merges pending entries into the index(second scan)
// Called after the index builder has finished scanning the snapshot(first scan) and
// all writers have completed their current transactions.
func synchronize() {
	mtx.Lock()
	defer mtx.Unlock()

	fmt.Println("Synchronizing...")
	for _, user := range pending {
		idx[user.Name] = user.ID
	}
	// Clear the pending list (equivalent to pending.clear() in C++).
	pending = nil
}

// buildIndex builds the fake index by:
//  1. Taking a snapshot of the current table.
//  2. Scanning every row in the snapshot and inserting it into idx (with a
//     1-second delay per row to simulate a huge table scan).
//  3. Waiting until all writers have finished (activeWriters == 0).
//  4. Merging any rows that were inserted during the build (synchronize).
//
// Equivalent to the C++ buildIndex().
func buildIndex() {
	mtx.Lock()
	building = true
	mtx.Unlock()

	// snapshot the current table.
	dbSnapshot := getCurrentSnapShot()

	//index every row from the snapshot.
	for _, user := range dbSnapshot {
		fmt.Printf("Building index for %s\n", user.Name)

		mtx.Lock()
		idx[user.Name] = user.ID
		mtx.Unlock()

		// Simulate the time needed to scan & index millions of rows.
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Builder waiting...")

	// wait until every writer that started has completed.
	// cv.Wait() atomically releases the mutex and blocks; when woken
	// (by cv.Signal()), it re-acquires the mutex.
	mtx.Lock()
	for activeWriters > 0 {
		cv.Wait()
	}
	fmt.Println("Builder woke up!")
	mtx.Unlock()

	// merge pending writes into the index.
	synchronize()

	mtx.Lock()
	building = false
	mtx.Unlock()
}