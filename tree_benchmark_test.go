package bptree

import (
	"fmt"
	"testing"
)

func BenchmarkInsert(b *testing.B) {
	t := NewTree()

	// insert b.N times
	for n := 0; n < b.N; n++ {
		key := n
		value := n + 1

		err := t.Insert(key, value)
		if err != nil {
			fmt.Printf("error: %s\n\n", err)
		}
	}
}

func BenchmarkInsertFind(b *testing.B) {
	t := NewTree()

	// insert b.N times
	for n := 0; n < b.N; n++ {
		key := n
		value := n + 1

		err := t.Insert(key, value)

		if err != nil {
			fmt.Printf("error: %s\n\n", err)
		}
	}

	// find one by one
	for n := 0; n < b.N; n++ {
		_, err = t.Find(n, false)
		if err != nil {
			fmt.Printf("error: %s\n\n", err)
		}
	}
}

func BenchmarkInsertDelete(b *testing.B) {
	t := NewTree()

	// insert b.N times
	for n := 0; n < b.N; n++ {
		key := n
		value := n + 1

		err := t.Insert(key, value)
		if err != nil {
			fmt.Printf("error: %s\n\n", err)
		}
	}

	// delete them
	for n := 0; n < b.N; n++ {
		err = t.Delete(n)
		if err != nil {
			fmt.Printf("error: %s\n\n", err)
		}
	}
}
