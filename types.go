package main

import "sync"

// User represents a row in our "database" table.
type User struct {
	ID   int
	Name string
}

// Global shared state 
var (
	mtx  sync.Mutex
	cv   = sync.NewCond(&mtx)

	table []User               // the "table" (heap storage)
	idx   = make(map[string]int) // fake index on User.Name

	pending []User // rows inserted while the index was being built
	building bool  // is an index build currently running?

	activeWriters int // number of writers that have started but not yet committed
)