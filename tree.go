package bptree

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	defaultOrder = 4
	minOrder     = 3
	maxOrder     = 20

	order          = defaultOrder
	queue          *node
	verbose_output = false
	version        = 0.1
)

type record struct {
	Value int
}

type node struct {
	pointers []interface{}
	keys     []int
	parent   *node
	is_leaf  bool
	num_keys int
	next     *node
}

func Insert(root *node, key, value int) (*node, error) {
	var pointer *record
	var leaf *node

	if _, err := Find(root, key, false); err == nil {
		return root, errors.New("key already exists")
	}

	pointer, err := makeRecord(value)
	if err != nil {
		return root, err
	}

	if root == nil {
		return startNewTree(key, pointer)
	}

	leaf = findLeaf(root, key, false)

	if leaf.num_keys < order-1 {
		insertIntoLeaf(leaf, key, pointer)
		return root, nil
	}

	return insertIntoLeafAfterSplitting(root, leaf, key, pointer)
}

func Find(root *node, key int, verbose bool) (*record, error) {
	i := 0
	c := findLeaf(root, key, verbose)
	if c == nil {
		return nil, errors.New("key not found")
	}
	for i = 0; i < c.num_keys; i++ {
		if c.keys[i] == key {
			break
		}
	}
	if i == c.num_keys {
		return nil, errors.New("key not found")
	}

	r, _ := c.pointers[i].(*record)

	return r, nil
}

func FindAndPrint(root *node, key int, verbose bool) {
	r, err := Find(root, key, verbose)

	if err != nil || r == nil {
		fmt.Printf("Record not found under key %d.\n", key)
	} else {
		fmt.Printf("Record at %d -- key %d, value %d.\n", r, key, r.Value)
	}
}

func FindAndPrintRange(root *node, key_start, key_end int, verbose bool) {
	var i int
	array_size := key_end - key_start + 1
	returned_keys := make([]int, array_size)
	returned_pointers := make([]interface{}, array_size)
	num_found := findRange(root, key_start, key_end, verbose, returned_keys, returned_pointers)
	if num_found == 0 {
		fmt.Println("None found,\n")
	} else {
		for i = 0; i < num_found; i++ {
			c, _ := returned_pointers[i].(*record)
			fmt.Printf("Key: %d  Location: %d  Value: %d\n",
				returned_keys[i],
				returned_pointers[i],
				c.Value)
		}
	}
}

func PrintTree(root *node) {
	var n *node
	i := 0
	rank := 0
	new_rank := 0

	if root == nil {
		fmt.Printf("Empty tree.\n")
		return
	}
	queue = nil
	enqueue(root)
	for queue != nil {
		n = dequeue()
		if n != nil {
			if n.parent != nil && n == n.parent.pointers[0] {
				new_rank = pathToRoot(root, n)
				if new_rank != rank {
					fmt.Printf("\n")
				}
			}
			if verbose_output {
				fmt.Printf("(%d)", n)
			}
			for i = 0; i < n.num_keys; i++ {
				if verbose_output {
					fmt.Printf("%d ", n.pointers[i])
				}
				fmt.Printf("%d ", n.keys[i])
			}
			if !n.is_leaf {
				for i = 0; i <= n.num_keys; i++ {
					c, _ := n.pointers[i].(*node)
					enqueue(c)
				}
			}
			if verbose_output {
				if n.is_leaf {
					fmt.Printf("%d ", n.pointers[order-1])
				} else {
					fmt.Printf("%d ", n.pointers[n.num_keys])
				}
			}
			fmt.Printf(" | ")
		}
	}
	fmt.Printf("\n")
}

func PrintLeaves(root *node) {
	var i int
	c := root
	if root == nil {
		fmt.Printf("Empty tree.\n")
		return
	}
	for !c.is_leaf {
		c, _ = c.pointers[0].(*node)
	}

	for true {
		for i = 0; i < c.num_keys; i++ {
			if verbose_output {
				fmt.Printf("%d ", c.pointers[i])
			}
			fmt.Printf("%d ", c.keys[i])
		}
		if verbose_output {
			fmt.Printf("%d ", c.pointers[order-1])
		}
		if c.pointers[order-1] != nil {
			fmt.Printf(" | ")
			c, _ = c.pointers[order-1].(*node)
		} else {
			break
		}
	}
	fmt.Printf("\n")
}

func Delete(root *node, key int) (*node, error) {
	key_record, err := Find(root, key, false)
	if err != nil {
		return root, err
	}
	key_leaf := findLeaf(root, key, false)
	if key_record != nil && key_leaf != nil {
		root = deleteEntry(root, key_leaf, key, key_record)
	}
	return root, nil
}

func DestroyTree(root *node) *node {
	root = nil
	return nil
}

//
//
//
//
//
//
//
// Private Functions
//
//
//
//
//
//
//
func enqueue(new_node *node) {
	var c *node
	if queue == nil {
		queue = new_node
		queue.next = nil
	} else {
		c = queue
		for c.next != nil {
			c = c.next
		}
		c.next = new_node
		new_node.next = nil
	}
}

func dequeue() *node {
	n := queue
	queue = queue.next
	n.next = nil
	return n
}

func height(root *node) int {
	h := 0
	c := root
	for !c.is_leaf {
		c, _ = c.pointers[0].(*node)
		h++
	}
	return h
}

func pathToRoot(root, child *node) int {
	length := 0
	c := child
	for c != root {
		c = c.parent
		length += 1
	}
	return length
}

func findRange(root *node, key_start, key_end int, verbose bool, returned_keys []int, returned_pointers []interface{}) int {
	var i int
	num_found := 0

	n := findLeaf(root, key_start, verbose)
	if n == nil {
		return 0
	}
	for i = 0; i < n.num_keys && n.keys[i] < key_start; i++ {
		if i == n.num_keys { // could be wrong
			return 0
		}
	}
	for n != nil {
		for i = i; i < n.num_keys && n.keys[i] <= key_end; i++ {
			returned_keys[num_found] = n.keys[i]
			returned_pointers[num_found] = n.pointers[i]
			num_found += 1
		}
		n, _ = n.pointers[order-1].(*node)
		i = 0
	}
	return num_found
}

func findLeaf(root *node, key int, verbose bool) *node {
	i := 0
	c := root
	if c == nil {
		if verbose {
			fmt.Printf("Empty tree.\n")
		}
		return c
	}
	for !c.is_leaf {
		if verbose {
			fmt.Printf("[")
			for i = 0; i < c.num_keys-1; i++ {
				fmt.Printf("%d ", c.keys[i])
			}
			fmt.Printf("%d]", c.keys[i])
		}
		i = 0
		for i < c.num_keys {
			if key >= c.keys[i] {
				i += 1
			} else {
				break
			}
		}
		if verbose {
			fmt.Printf("%d ->\n", i)
		}
		c, _ = c.pointers[i].(*node)
	}
	if verbose {
		fmt.Printf("Leaf [")
		for i = 0; i < c.num_keys-1; i++ {
			fmt.Printf("%d ", c.keys[i])
		}
		fmt.Printf("%d] ->\n", c.keys[i])
	}
	return c
}

func cut(length int) int {
	if length%2 == 0 {
		return length / 2
	}

	return length/2 + 1
}

//
//	INSERTION
//
func makeRecord(value int) (*record, error) {
	new_record := new(record)
	if new_record == nil {
		return nil, errors.New("Error: Record creation.")
	} else {
		new_record.Value = value
	}
	return new_record, nil
}

func makeNode() (*node, error) {
	new_node := new(node)
	if new_node == nil {
		return nil, errors.New("Error: Node creation.")
	}
	new_node.keys = make([]int, order-1)
	if new_node.keys == nil {
		return nil, errors.New("Error: New node keys array.")
	}
	new_node.pointers = make([]interface{}, order)
	if new_node.keys == nil {
		return nil, errors.New("Error: New node pointers array.")
	}
	new_node.is_leaf = false
	new_node.num_keys = 0
	new_node.parent = nil
	new_node.next = nil
	return new_node, nil
}

func makeLeaf() (*node, error) {
	leaf, err := makeNode()
	if err != nil {
		return nil, err
	}
	leaf.is_leaf = true
	return leaf, nil
}

func getLeftIndex(parent, left *node) int {
	left_index := 0
	for left_index <= parent.num_keys && parent.pointers[left_index] != left {
		left_index += 1
	}
	return left_index
}

func insertIntoLeaf(leaf *node, key int, pointer *record) {
	var i, insertion_point int

	for insertion_point < leaf.num_keys && leaf.keys[insertion_point] < key {
		insertion_point += 1
	}

	for i = leaf.num_keys; i > insertion_point; i-- {
		leaf.keys[i] = leaf.keys[i-1]
		leaf.pointers[i] = leaf.pointers[i-1]
	}
	leaf.keys[insertion_point] = key
	leaf.pointers[insertion_point] = pointer
	leaf.num_keys += 1
	return
}

func insertIntoLeafAfterSplitting(root *node, leaf *node, key int, pointer *record) (*node, error) {
	var new_leaf *node
	var insertion_index, split, new_key, i, j int
	var err error

	new_leaf, err = makeLeaf()
	if err != nil {
		return root, nil
	}

	temp_keys := make([]int, order)
	if temp_keys == nil {
		return root, errors.New("Error: Temporary keys array.")
	}

	temp_pointers := make([]interface{}, order)
	if temp_pointers == nil {
		return root, errors.New("Error: Temporary pointers array.")
	}

	for insertion_index < order-1 && leaf.keys[insertion_index] < key {
		insertion_index += 1
	}

	for i = 0; i < leaf.num_keys; i++ {
		if j == insertion_index {
			j += 1
		}
		temp_keys[j] = leaf.keys[i]
		temp_pointers[j] = leaf.pointers[i]
		j += 1
	}

	temp_keys[insertion_index] = key
	temp_pointers[insertion_index] = pointer

	leaf.num_keys = 0

	split = cut(order - 1)

	for i = 0; i < split; i++ {
		leaf.pointers[i] = temp_pointers[i]
		leaf.keys[i] = temp_keys[i]
		leaf.num_keys += 1
	}

	j = 0
	for i = split; i < order; i++ {
		new_leaf.pointers[j] = temp_pointers[i]
		new_leaf.keys[j] = temp_keys[i]
		new_leaf.num_keys += 1
		j += 1
	}

	new_leaf.pointers[order-1] = leaf.pointers[order-1]
	leaf.pointers[order-1] = new_leaf

	for i = leaf.num_keys; i < order-1; i++ {
		leaf.pointers[i] = nil
	}
	for i = new_leaf.num_keys; i < order-1; i++ {
		new_leaf.pointers[i] = nil
	}

	new_leaf.parent = leaf.parent
	new_key = new_leaf.keys[0]

	return insertIntoParent(root, leaf, new_key, new_leaf)
}

func insertIntoNode(root, n *node, left_index, key int, right *node) *node {
	var i int
	for i = n.num_keys; i > left_index; i-- {
		n.pointers[i+1] = n.pointers[i]
		n.keys[i] = n.keys[i-1]
	}
	n.pointers[left_index+1] = right
	n.keys[left_index] = key
	n.num_keys += 1
	return root
}

func insertIntoNodeAfterSplitting(root, old_node *node, left_index, key int, right *node) (*node, error) {
	var i, j, split, k_prime int
	var new_node, child *node
	var temp_keys []int
	var temp_pointers []interface{}
	var err error

	temp_pointers = make([]interface{}, order+1)
	if temp_pointers == nil {
		return root, errors.New("Error: Temporary pointers array for splitting nodes.")
	}

	temp_keys = make([]int, order+1)
	if temp_keys == nil {
		return root, errors.New("Error: Temporary keys array for splitting nodes.")
	}

	for i = 0; i < old_node.num_keys+1; i++ {
		if j == left_index+1 {
			j += 1
		}
		temp_pointers[j] = old_node.pointers[i]
		j += 1
	}

	for i = 0; i < old_node.num_keys+1; i++ {
		if j == left_index+1 {
			j += 1
		}
		temp_keys[j] = old_node.keys[i]
		j += 1
	}

	temp_pointers[left_index+1] = right
	temp_keys[left_index] = key

	split = cut(order)
	new_node, err = makeNode()
	if err != nil {
		return root, err
	}
	old_node.num_keys = 0
	for i = 0; i < split-1; i++ {
		old_node.pointers[i] = temp_pointers[i]
		old_node.keys[i] = temp_keys[i]
		old_node.num_keys += 1
	}
	old_node.pointers[i] = temp_pointers[i]
	k_prime = temp_keys[split-1]
	j += 1
	for i += 1; i < order; i++ {
		new_node.pointers[j] = temp_pointers[i]
		new_node.keys[j] = temp_keys[i]
		new_node.num_keys += 1
		j += 1
	}
	new_node.pointers[j] = temp_pointers[i]
	new_node.parent = old_node.parent
	for i = 0; i <= new_node.num_keys; i++ {
		child, _ = new_node.pointers[i].(*node)
		child.parent = new_node
	}

	return insertIntoParent(root, old_node, k_prime, new_node)
}

func insertIntoParent(root *node, left *node, key int, right *node) (*node, error) {
	var left_index int
	parent := left.parent

	if parent == nil {
		return insertIntoNewRoot(left, key, right)
	}
	left_index = getLeftIndex(parent, left)

	if parent.num_keys < order-1 {
		return insertIntoNode(root, parent, left_index, key, right), nil
	}

	return insertIntoNodeAfterSplitting(root, parent, left_index, key, right)
}

func insertIntoNewRoot(left *node, key int, right *node) (*node, error) {
	root, err := makeNode()
	if err != nil {
		return nil, err
	}
	root.keys[0] = key
	root.pointers[0] = left
	root.pointers[1] = right
	root.num_keys += 1
	root.parent = nil
	left.parent = root
	right.parent = root
	return root, nil
}

func startNewTree(key int, pointer *record) (*node, error) {
	root, err := makeLeaf()
	if err != nil {
		return root, err
	}
	root.keys[0] = key
	root.pointers[0] = pointer
	root.pointers[order-1] = nil
	root.parent = nil
	root.num_keys += 1
	return root, nil
}

func getNeighbourIndex(n *node) int {
	var i int

	for i = 0; i <= n.parent.num_keys; i++ {
		if reflect.DeepEqual(n.parent.pointers[i], n) {
			return i - 1
		}
	}

	return i
}

func removeEntryFromNode(n *node, key int, pointer interface{}) *node {
	var i, num_pointers int

	for n.keys[i] != key {
		i += 1
	}

	for i += 1; i < n.num_keys; i++ {
		n.keys[i-1] = n.keys[i]
	}

	if n.is_leaf {
		num_pointers = n.num_keys
	} else {
		num_pointers = n.num_keys + 1
	}

	i = 0
	for n.pointers[i] != pointer {
		i += 1
	}
	for i += 1; i < num_pointers; i++ {
		n.pointers[i-1] = n.pointers[i]
	}
	n.num_keys -= 1

	if n.is_leaf {
		for i = n.num_keys; i < order-1; i++ {
			n.pointers[i] = nil
		}
	} else {
		for i = n.num_keys + 1; i < order; i++ {
			n.pointers[i] = nil
		}
	}

	return n
}

func adjustRoot(root *node) *node {
	var new_root *node

	if root.num_keys > 0 {
		return root
	}

	if !root.is_leaf {
		new_root, _ = root.pointers[0].(*node)
		new_root.parent = nil
	} else {
		new_root = nil
	}

	return new_root
}

func coalesceNodes(root, n, neighbour *node, neighbour_index, k_prime int) *node {
	var i, j, neighbour_insertion_index, n_end int
	var tmp *node

	if neighbour_index == -1 {
		tmp = n
		n = neighbour
		neighbour = tmp
	}

	neighbour_insertion_index = neighbour.num_keys

	if !n.is_leaf {
		neighbour.keys[neighbour_insertion_index] = k_prime
		neighbour.num_keys += 1

		n_end = n.num_keys
		i = neighbour_insertion_index + 1
		for j = 0; j < n_end; j++ {
			neighbour.keys[i] = n.keys[j]
			neighbour.pointers[i] = n.pointers[j]
			neighbour.num_keys += 1
			n.num_keys -= 1
			i += 1
		}
		neighbour.pointers[i] = n.pointers[j]

		for i = 0; i < neighbour.num_keys+1; i++ {
			tmp, _ = neighbour.pointers[i].(*node)
			tmp.parent = neighbour
		}
	} else {
		i = neighbour_insertion_index
		for j = 0; j < n.num_keys; j++ {
			neighbour.keys[i] = n.keys[j]
			n.pointers[i] = n.pointers[j]
			neighbour.num_keys += 1
		}
		neighbour.pointers[order-1] = n.pointers[order-1]
	}
	root = deleteEntry(root, n.parent, k_prime, n)

	return root
}

func redistributeNodes(root, n, neighbour *node, neighbour_index, k_prime_index, k_prime int) *node {
	var i int
	var tmp *node

	if neighbour_index != -1 {
		if !n.is_leaf {
			n.pointers[n.num_keys+1] = n.pointers[n.num_keys]
		}
		for i = n.num_keys; i > 0; i-- {
			n.keys[i] = n.keys[i-1]
			n.pointers[i] = n.pointers[i-1]
		}
		if !n.is_leaf { // why the second if !n.is_leaf
			n.pointers[0] = neighbour.pointers[neighbour.num_keys]
			tmp, _ = n.pointers[0].(*node)
			tmp.parent = n
			neighbour.pointers[neighbour.num_keys] = nil
			n.keys[0] = k_prime
			n.parent.keys[k_prime_index] = neighbour.keys[neighbour.num_keys-1]
		} else {
			n.pointers[0] = neighbour.pointers[neighbour.num_keys-1]
			neighbour.pointers[neighbour.num_keys-1] = nil
			n.keys[0] = neighbour.keys[neighbour.num_keys-1]
			n.parent.keys[k_prime_index] = n.keys[0]
		}
	} else {
		if n.is_leaf {
			n.keys[n.num_keys] = neighbour.keys[0]
			n.pointers[n.num_keys] = neighbour.pointers[0]
			n.parent.keys[k_prime_index] = neighbour.keys[1]
		} else {
			n.keys[n.num_keys] = k_prime
			n.pointers[n.num_keys+1] = neighbour.pointers[0]
			tmp, _ = n.pointers[n.num_keys+1].(*node)
			tmp.parent = n
			n.parent.keys[k_prime_index] = neighbour.keys[0]
		}
		for i = 0; i < neighbour.num_keys-1; i++ {
			neighbour.keys[i] = neighbour.keys[i+1]
			neighbour.pointers[i] = neighbour.pointers[i+1]
		}
		if !n.is_leaf {
			neighbour.pointers[i] = neighbour.pointers[i+1]
		}
	}
	n.num_keys += 1
	neighbour.num_keys -= 1

	return root
}

func deleteEntry(root, n *node, key int, pointer interface{}) *node {
	var min_keys, neighbour_index, k_prime_index, k_prime, capacity int
	var neighbour *node

	n = removeEntryFromNode(n, key, pointer)

	if n == root {
		return adjustRoot(root)
	}

	if n.is_leaf {
		min_keys = cut(order - 1)
	} else {
		min_keys = cut(order) - 1
	}

	if n.num_keys >= min_keys {
		return root
	}

	neighbour_index = getNeighbourIndex(n)

	if neighbour_index == -1 {
		k_prime_index = 0
	} else {
		k_prime_index = neighbour_index
	}

	k_prime = n.parent.keys[k_prime_index]

	if neighbour_index == -1 {
		neighbour, _ = n.parent.pointers[1].(*node)
	} else {
		neighbour, _ = n.parent.pointers[neighbour_index].(*node)
	}

	if n.is_leaf {
		capacity = order
	} else {
		capacity = order - 1
	}

	if neighbour.num_keys+n.num_keys < capacity {
		return coalesceNodes(root, n, neighbour, neighbour_index, k_prime)
	} else {
		return redistributeNodes(root, n, neighbour, neighbour_index, k_prime_index, k_prime)
	}

}
