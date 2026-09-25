package main

import (
	"fmt"
)

func beginTransaction() *Transaction {
	mtx.Lock()
	defer mtx.Unlock()

	tx := &Transaction{
		ID:    nextTxId,
		State: TxActive,
	}

	nextTxId++

	transactions[tx.ID] = tx

	fmt.Printf("Transaction %d started\n", tx.ID)

	return tx

}

func commitTransaction(tx *Transaction) {
	mtx.Lock()
	defer mtx.Unlock()

	tx.State = TxCommitted

	fmt.Printf("Transaction %d committed\n", tx.ID)

	cv.Broadcast()
}
