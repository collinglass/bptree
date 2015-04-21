bptree
====

bptree is a pure Go B+ tree implementation. The project started out as a Go copy of [Amittai Aviram's][bpt] C implementation. The goal of the project is to maintain a tiny B+ tree with a concise API to use in the composition of more complex systems.

[bpt]: http://www.amittai.com/prose/bplustree.html


## Project Roadmap

Please submit issues/PRs for anything else you'd like to see in this package.


## Getting Started

### Installing

To get started using bptree, install Go and run ```go get```

```$ go get github.com/collinglass/bptree```

This will retrieve the library, your project is now ready to use bptree.


### Example

```
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

```


## API

#### func NewTree()

```func NewTree() *Tree```

Returns a pointer to a new B+ tree

#### func (t *Tree) Insert

```func (t *Tree) Insert(key, value int) (*node, error)```

Insert adds a key value pair to the b+ tree.

#### func (t *Tree) Find

```func (t *Tree) Find(key int, verbose bool) (*record, error)```

Find returns the value associated with the key. If verbose is true, it gives some insight on its location in the tree.

#### func (t *Tree) FindAndPrint

```func (t *Tree) FindAndPrint(key int, verbose bool)```

FindAndPrint outputs the value associated with the key. If verbose is true, it gives added insight.

#### func (t *Tree) FindAndPrintRange

```func (t *Tree) FindAndPrintRange(key_start, key_end int, verbose bool)```

FindAndPrintRange outputs the values associated with the keys inside of the range inclusive.

#### func (t *Tree) PrintTree

```func (t *Tree) PrintTree(root *node)```

Prints the whole tree.

#### func (t *Tree) PrintLeaves

```func (t *Tree) PrintLeaves()```

Prints the leaves of the tree.
