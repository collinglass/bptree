package main

import (
	"fmt"
	"github.com/collinglass/bptree"
)

func main() {
	key := 1
	value := 2

	t := bptree.NewTree()

	err := t.Insert(key, value)
	if err != nil {
		fmt.Printf("error: %s\n\n", err)
	}

	r, err := t.Find(key, true)
	if err != nil {
		fmt.Printf("error: %s\n\n", err)
	}

	fmt.Printf("%v\n\n", r.Value)

	t.FindAndPrint(key, true)
}
