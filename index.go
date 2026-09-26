package main

import (
	"fmt"
	"time"
)

func hasActiveTransactionsLocked(excludeTxID int) bool {
	/*
	   Checks for any pending transaction befor completing the scan.
	*/
	for _, tx := range transactions {
		if tx.ID != excludeTxID && tx.State == TxActive {
			return true
		}
	}

	return false
}

// func hasPendingActiveTransactions(excludeTxID int) bool {
// 	/* helper function as we cannot call hasActiveTransactionsLocked() on previous transactions as
// 	it has a lock on mtx and we cannot call it on the same lock as it will cause a deadlock. */

// 	mtx.Lock()
// 	defer mtx.Unlock()

// 	return hasActiveTransactionsLocked(excludeTxID)
// }

// buildIndex builds the fake index by:
// 1. Starting an indexing transaction.
// 2. Taking a snapshot and performing the first scan.
// 3. Waiting for concurrent transactions to finish.
// 4. Taking a new snapshot and performing a second scan.
// 5. Completing the index build.
func buildIndex() {

	tx := beginTransaction() // as indexing is also a transaction - as per postgres documentation

	mtx.Lock()
	building = true
	mtx.Unlock()

	// snapshot the current table.
	dbSnapshot := getCurrentSnapShot()

	mtx.Lock()
	tableCopy := make([]User, len(table))
	copy(tableCopy, table)
	mtx.Unlock()

	//index every row from the snapshot.
	for _, user := range tableCopy {

		if !dbSnapshot.IsVisible(user) {
			continue
		}

		fmt.Printf(
			"Building index for %s (Tx %d)\n",
			user.Name,
			tx.ID,
		)

		mtx.Lock()
		idx[user.Name] = append(idx[user.Name], user.ID)
		indexedRows[user.ID] = true
		mtx.Unlock()

		// Simulate the time needed to scan & index millions of rows.
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Builder waiting...")

	// wait until every writer that started has completed.
	// cv.Wait() atomically releases the mutex and blocks; when woken
	// (by cv.Signal()), it re-acquires the mutex.
	mtx.Lock()

	for hasActiveTransactionsLocked(tx.ID) {
		fmt.Println("Builder waiting for active transactions...")
		cv.Wait()
	}

	fmt.Println("Builder woke up!")

	mtx.Unlock()

	// second scan to merge any pending writes into the index. This is done after all writers have completed their concurrent transactions.
	secondSnapshot := getCurrentSnapShot()

	mtx.Lock()
	secondTableCopy := make([]User, len(table))
	copy(secondTableCopy, table)
	mtx.Unlock()

	for _, user := range secondTableCopy {
		if !secondSnapshot.IsVisible(user) {
			continue
		}

		mtx.Lock()

		if !indexedRows[user.ID] {
			fmt.Printf(
				"Second scan: indexing %s (ID %d)\n",
				user.Name,
				user.ID,
			)

			idx[user.Name] = append(idx[user.Name], user.ID)
			indexedRows[user.ID] = true
		}

		mtx.Unlock()
	}

	mtx.Lock()
	building = false
	mtx.Unlock()

	commitTransaction(tx)
}
