bptree
====

bptree is a pure Go B+ tree implementation inspired by [Amittai Aviram's][bpt] C implementation. The goal of the project is to maintain a tiny B+ tree with a concise API to use in the composition of more complex systems.

[bpt]: http://www.amittai.com/prose/bplustree.html


## Project Roadmap

bptree currently only supports integer values, support for any value type is in the roadmap ahead. Please submit issues/PRs for anything else you'd like to see in this package.


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
```


## API

#### func Insert

```func Insert(root *node, key, value int) (*node, error)```

Insert adds a key value pair to the b+ tree.

#### func Find

```func Find(root *node, key int, verbose bool) (*record, error)```

Find returns the value associated with the key. If verbose is true, it gives some insight on its location in the tree.

#### func FindAndPrint

```func FindAndPrint(root *node, key int, verbose bool)```

FindAndPrint outputs the value associated with the key. If verbose is true, it gives added insight.

#### func FindAndPrintRange

```func FindAndPrintRange(root *node, key_start, key_end int, verbose bool)```

FindAndPrintRange outputs the values associated with the keys inside of the range inclusive.

#### func PrintTree

```func PrintTree(root *node)```

Prints the whole tree.

#### func PrintLeaves

```func PrintLeaves(root *node)```

Prints the leaves of the tree.

#### func Delete

```func Delete(root *node, key int) (*node, error)```

Deletes the key value pair associated with the given key.

#### func DestroyTree

```func DestroyTree(root *node) *node```

Destroy tree.



