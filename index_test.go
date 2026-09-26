package main

import (
	"reflect"
	"testing"
)

func TestDuplicateValuesAreIndexedByID(t *testing.T) {
	idx = make(map[string][]int)
	indexedRows = make(map[int]bool)

	table = []User{
		{ID: 3, Name: "Charlie", CreatedBy: 0},
		{ID: 4, Name: "Charlie", CreatedBy: 0},
	}

	for _, user := range table {
		idx[user.Name] = append(idx[user.Name], user.ID)
		indexedRows[user.ID] = true
	}

	if got, want := idx["Charlie"], []int{3, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("idx[\"Charlie\"] = %v, want %v", got, want)
	}

	if _, ok := indexedRows[3]; !ok {
		t.Fatal("indexedRows should be keyed by User.ID")
	}
}
