package main

import (
	"fmt"
	"github.com/collinglass/bptree"
)

func main() {
	key := 1
	value := 2

	root, err := bptree.Insert(nil, key, value)

	if err != nil {
		fmt.Printf("error: %s\n\n", err)
	}

	r, err := bptree.Find(root, key, true)
	if err != nil {
		fmt.Printf("error: %s\n\n", err)
	}

	fmt.Printf("%v\n\n", r.Value)

	bptree.FindAndPrint(root, key, true)
}
