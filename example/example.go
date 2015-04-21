package main

import (
	"fmt"
	"github.com/collinglass/bptree"
)

func main() {
	key := 1
	value := []byte("hello friend")

	t := bptree.NewTree()

	err := t.Insert(key, value)
	if err != nil {
		fmt.Printf("error: %s\n\n", err)
	}

	r, err := t.Find(key, true)
	if err != nil {
		fmt.Printf("error: %s\n\n", err)
	}

	fmt.Printf("%s\n\n", r.Value)

	t.FindAndPrint(key, true)
}
