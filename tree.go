package bptree

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	err error

	defaultOrder = 4
	minOrder     = 3
	maxOrder     = 20

	order          = defaultOrder
	queue          *Node
	verbose_output = false
	version        = 0.1
)

type Tree struct {
	Root *Node
}

type Record struct {
	Value []byte
}

type Node struct {
	Pointers []interface{}
	Keys     []int
	Parent   *Node
	IsLeaf   bool
	NumKeys  int
	Next     *Node
}

func NewTree() *Tree {
	return &Tree{}
}

func (t *Tree) Insert(key int, value []byte) error {
	var pointer *Record
	var leaf *Node

	if _, err := t.Find(key, false); err == nil {
		return errors.New("key already exists")
	}

	pointer, err := makeRecord(value)
	if err != nil {
		return err
	}

	if t.Root == nil {
		return t.startNewTree(key, pointer)
	}

	leaf = t.findLeaf(key, false)

	if leaf.NumKeys < order-1 {
		insertIntoLeaf(leaf, key, pointer)
		return nil
	}

	return t.insertIntoLeafAfterSplitting(leaf, key, pointer)
}

func (t *Tree) Find(key int, verbose bool) (*Record, error) {
	i := 0
	c := t.findLeaf(key, verbose)
	if c == nil {
		return nil, errors.New("key not found")
	}
	for i = 0; i < c.NumKeys; i++ {
		if c.Keys[i] == key {
			break
		}
	}
	if i == c.NumKeys {
		return nil, errors.New("key not found")
	}

	r, _ := c.Pointers[i].(*Record)

	return r, nil
}

func (t *Tree) FindAndPrint(key int, verbose bool) {
	r, err := t.Find(key, verbose)

	if err != nil || r == nil {
		fmt.Printf("Record not found under key %d.\n", key)
	} else {
		fmt.Printf("Record at %d -- key %d, value %s.\n", r, key, r.Value)
	}
}

func (t *Tree) FindAndPrintRange(key_start, key_end int, verbose bool) {
	var i int
	array_size := key_end - key_start + 1
	returned_keys := make([]int, array_size)
	returned_pointers := make([]interface{}, array_size)
	num_found := t.findRange(key_start, key_end, verbose, returned_keys, returned_pointers)
	if num_found == 0 {
		fmt.Println("None found,\n")
	} else {
		for i = 0; i < num_found; i++ {
			c, _ := returned_pointers[i].(*Record)
			fmt.Printf("Key: %d  Location: %d  Value: %s\n",
				returned_keys[i],
				returned_pointers[i],
				c.Value)
		}
	}
}

func (t *Tree) PrintTree() {
	var n *Node
	i := 0
	rank := 0
	new_rank := 0

	if t.Root == nil {
		fmt.Printf("Empty tree.\n")
		return
	}
	queue = nil
	enqueue(t.Root)
	for queue != nil {
		n = dequeue()
		if n != nil {
			if n.Parent != nil && n == n.Parent.Pointers[0] {
				new_rank = t.pathToRoot(n)
				if new_rank != rank {
					fmt.Printf("\n")
				}
			}
			if verbose_output {
				fmt.Printf("(%d)", n)
			}
			for i = 0; i < n.NumKeys; i++ {
				if verbose_output {
					fmt.Printf("%d ", n.Pointers[i])
				}
				fmt.Printf("%d ", n.Keys[i])
			}
			if !n.IsLeaf {
				for i = 0; i <= n.NumKeys; i++ {
					c, _ := n.Pointers[i].(*Node)
					enqueue(c)
				}
			}
			if verbose_output {
				if n.IsLeaf {
					fmt.Printf("%d ", n.Pointers[order-1])
				} else {
					fmt.Printf("%d ", n.Pointers[n.NumKeys])
				}
			}
			fmt.Printf(" | ")
		}
	}
	fmt.Printf("\n")
}

func (t *Tree) PrintLeaves() {
	if t.Root == nil {
		fmt.Printf("Empty tree.\n")
		return
	}

	var i int
	c := t.Root
	for !c.IsLeaf {
		c, _ = c.Pointers[0].(*Node)
	}

	for {
		for i = 0; i < c.NumKeys; i++ {
			if verbose_output {
				fmt.Printf("%d ", c.Pointers[i])
			}
			fmt.Printf("%d ", c.Keys[i])
		}
		if verbose_output {
			fmt.Printf("%d ", c.Pointers[order-1])
		}
		if c.Pointers[order-1] != nil {
			fmt.Printf(" | ")
			c, _ = c.Pointers[order-1].(*Node)
		} else {
			break
		}
	}
	fmt.Printf("\n")
}

func (t *Tree) Delete(key int) error {
	key_record, err := t.Find(key, false)
	if err != nil {
		return err
	}
	key_leaf := t.findLeaf(key, false)
	if key_record != nil && key_leaf != nil {
		t.deleteEntry(key_leaf, key, key_record)
	}
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
func enqueue(new_node *Node) {
	var c *Node
	if queue == nil {
		queue = new_node
		queue.Next = nil
	} else {
		c = queue
		for c.Next != nil {
			c = c.Next
		}
		c.Next = new_node
		new_node.Next = nil
	}
}

func dequeue() *Node {
	n := queue
	queue = queue.Next
	n.Next = nil
	return n
}

func (t *Tree) height() int {
	h := 0
	c := t.Root
	for !c.IsLeaf {
		c, _ = c.Pointers[0].(*Node)
		h++
	}
	return h
}

func (t *Tree) pathToRoot(child *Node) int {
	length := 0
	c := child
	for c != t.Root {
		c = c.Parent
		length += 1
	}
	return length
}

func (t *Tree) findRange(key_start, key_end int, verbose bool, returned_keys []int, returned_pointers []interface{}) int {
	var i int
	num_found := 0

	n := t.findLeaf(key_start, verbose)
	if n == nil {
		return 0
	}
	for i = 0; i < n.NumKeys && n.Keys[i] < key_start; i++ {
		if i == n.NumKeys { // could be wrong
			return 0
		}
	}
	for n != nil {
		for i = i; i < n.NumKeys && n.Keys[i] <= key_end; i++ {
			returned_keys[num_found] = n.Keys[i]
			returned_pointers[num_found] = n.Pointers[i]
			num_found += 1
		}
		n, _ = n.Pointers[order-1].(*Node)
		i = 0
	}
	return num_found
}

func (t *Tree) findLeaf(key int, verbose bool) *Node {
	i := 0
	c := t.Root
	if c == nil {
		if verbose {
			fmt.Printf("Empty tree.\n")
		}
		return c
	}
	for !c.IsLeaf {
		if verbose {
			fmt.Printf("[")
			for i = 0; i < c.NumKeys-1; i++ {
				fmt.Printf("%d ", c.Keys[i])
			}
			fmt.Printf("%d]", c.Keys[i])
		}
		i = 0
		for i < c.NumKeys {
			if key >= c.Keys[i] {
				i += 1
			} else {
				break
			}
		}
		if verbose {
			fmt.Printf("%d ->\n", i)
		}
		c, _ = c.Pointers[i].(*Node)
	}
	if verbose {
		fmt.Printf("Leaf [")
		for i = 0; i < c.NumKeys-1; i++ {
			fmt.Printf("%d ", c.Keys[i])
		}
		fmt.Printf("%d] ->\n", c.Keys[i])
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
func makeRecord(value []byte) (*Record, error) {
	new_record := new(Record)
	if new_record == nil {
		return nil, errors.New("Error: Record creation.")
	} else {
		new_record.Value = value
	}
	return new_record, nil
}

func makeNode() (*Node, error) {
	new_node := new(Node)
	if new_node == nil {
		return nil, errors.New("Error: Node creation.")
	}
	new_node.Keys = make([]int, order-1)
	if new_node.Keys == nil {
		return nil, errors.New("Error: New node keys array.")
	}
	new_node.Pointers = make([]interface{}, order)
	if new_node.Keys == nil {
		return nil, errors.New("Error: New node pointers array.")
	}
	new_node.IsLeaf = false
	new_node.NumKeys = 0
	new_node.Parent = nil
	new_node.Next = nil
	return new_node, nil
}

func makeLeaf() (*Node, error) {
	leaf, err := makeNode()
	if err != nil {
		return nil, err
	}
	leaf.IsLeaf = true
	return leaf, nil
}

func getLeftIndex(parent, left *Node) int {
	left_index := 0
	for left_index <= parent.NumKeys && parent.Pointers[left_index] != left {
		left_index += 1
	}
	return left_index
}

func insertIntoLeaf(leaf *Node, key int, pointer *Record) {
	var i, insertion_point int

	for insertion_point < leaf.NumKeys && leaf.Keys[insertion_point] < key {
		insertion_point += 1
	}

	for i = leaf.NumKeys; i > insertion_point; i-- {
		leaf.Keys[i] = leaf.Keys[i-1]
		leaf.Pointers[i] = leaf.Pointers[i-1]
	}
	leaf.Keys[insertion_point] = key
	leaf.Pointers[insertion_point] = pointer
	leaf.NumKeys += 1
	return
}

func (t *Tree) insertIntoLeafAfterSplitting(leaf *Node, key int, pointer *Record) error {
	var new_leaf *Node
	var insertion_index, split, new_key, i, j int
	var err error

	new_leaf, err = makeLeaf()
	if err != nil {
		return nil
	}

	temp_keys := make([]int, order)
	if temp_keys == nil {
		return errors.New("Error: Temporary keys array.")
	}

	temp_pointers := make([]interface{}, order)
	if temp_pointers == nil {
		return errors.New("Error: Temporary pointers array.")
	}

	for insertion_index < order-1 && leaf.Keys[insertion_index] < key {
		insertion_index += 1
	}

	for i = 0; i < leaf.NumKeys; i++ {
		if j == insertion_index {
			j += 1
		}
		temp_keys[j] = leaf.Keys[i]
		temp_pointers[j] = leaf.Pointers[i]
		j += 1
	}

	temp_keys[insertion_index] = key
	temp_pointers[insertion_index] = pointer

	leaf.NumKeys = 0

	split = cut(order - 1)

	for i = 0; i < split; i++ {
		leaf.Pointers[i] = temp_pointers[i]
		leaf.Keys[i] = temp_keys[i]
		leaf.NumKeys += 1
	}

	j = 0
	for i = split; i < order; i++ {
		new_leaf.Pointers[j] = temp_pointers[i]
		new_leaf.Keys[j] = temp_keys[i]
		new_leaf.NumKeys += 1
		j += 1
	}

	new_leaf.Pointers[order-1] = leaf.Pointers[order-1]
	leaf.Pointers[order-1] = new_leaf

	for i = leaf.NumKeys; i < order-1; i++ {
		leaf.Pointers[i] = nil
	}
	for i = new_leaf.NumKeys; i < order-1; i++ {
		new_leaf.Pointers[i] = nil
	}

	new_leaf.Parent = leaf.Parent
	new_key = new_leaf.Keys[0]

	return t.insertIntoParent(leaf, new_key, new_leaf)
}

func insertIntoNode(n *Node, left_index, key int, right *Node) {
	var i int
	for i = n.NumKeys; i > left_index; i-- {
		n.Pointers[i+1] = n.Pointers[i]
		n.Keys[i] = n.Keys[i-1]
	}
	n.Pointers[left_index+1] = right
	n.Keys[left_index] = key
	n.NumKeys += 1
}

func (t *Tree) insertIntoNodeAfterSplitting(old_node *Node, left_index, key int, right *Node) error {
	var i, j, split, k_prime int
	var new_node, child *Node
	var temp_keys []int
	var temp_pointers []interface{}
	var err error

	temp_pointers = make([]interface{}, order+1)
	if temp_pointers == nil {
		return errors.New("Error: Temporary pointers array for splitting nodes.")
	}

	temp_keys = make([]int, order)
	if temp_keys == nil {
		return errors.New("Error: Temporary keys array for splitting nodes.")
	}

	for i = 0; i < old_node.NumKeys+1; i++ {
		if j == left_index+1 {
			j += 1
		}
		temp_pointers[j] = old_node.Pointers[i]
		j += 1
	}

	j = 0
	for i = 0; i < old_node.NumKeys; i++ {
		if j == left_index {
			j += 1
		}
		temp_keys[j] = old_node.Keys[i]
		j += 1
	}

	temp_pointers[left_index+1] = right
	temp_keys[left_index] = key

	split = cut(order)
	new_node, err = makeNode()
	if err != nil {
		return err
	}
	old_node.NumKeys = 0
	for i = 0; i < split-1; i++ {
		old_node.Pointers[i] = temp_pointers[i]
		old_node.Keys[i] = temp_keys[i]
		old_node.NumKeys += 1
	}
	old_node.Pointers[i] = temp_pointers[i]
	k_prime = temp_keys[split-1]
	j = 0
	for i += 1; i < order; i++ {
		new_node.Pointers[j] = temp_pointers[i]
		new_node.Keys[j] = temp_keys[i]
		new_node.NumKeys += 1
		j += 1
	}
	new_node.Pointers[j] = temp_pointers[i]
	new_node.Parent = old_node.Parent
	for i = 0; i <= new_node.NumKeys; i++ {
		child, _ = new_node.Pointers[i].(*Node)
		child.Parent = new_node
	}

	return t.insertIntoParent(old_node, k_prime, new_node)
}

func (t *Tree) insertIntoParent(left *Node, key int, right *Node) error {
	var left_index int
	parent := left.Parent

	if parent == nil {
		return t.insertIntoNewRoot(left, key, right)
	}
	left_index = getLeftIndex(parent, left)

	if parent.NumKeys < order-1 {
		insertIntoNode(parent, left_index, key, right)
		return nil
	}

	return t.insertIntoNodeAfterSplitting(parent, left_index, key, right)
}

func (t *Tree) insertIntoNewRoot(left *Node, key int, right *Node) error {
	t.Root, err = makeNode()
	if err != nil {
		return err
	}
	t.Root.Keys[0] = key
	t.Root.Pointers[0] = left
	t.Root.Pointers[1] = right
	t.Root.NumKeys += 1
	t.Root.Parent = nil
	left.Parent = t.Root
	right.Parent = t.Root
	return nil
}

func (t *Tree) startNewTree(key int, pointer *Record) error {
	t.Root, err = makeLeaf()
	if err != nil {
		return err
	}
	t.Root.Keys[0] = key
	t.Root.Pointers[0] = pointer
	t.Root.Pointers[order-1] = nil
	t.Root.Parent = nil
	t.Root.NumKeys += 1
	return nil
}

func getNeighbourIndex(n *Node) int {
	var i int

	for i = 0; i <= n.Parent.NumKeys; i++ {
		if reflect.DeepEqual(n.Parent.Pointers[i], n) {
			return i - 1
		}
	}

	return i
}

func removeEntryFromNode(n *Node, key int, pointer interface{}) *Node {
	var i, num_pointers int

	for n.Keys[i] != key {
		i += 1
	}

	for i += 1; i < n.NumKeys; i++ {
		n.Keys[i-1] = n.Keys[i]
	}

	if n.IsLeaf {
		num_pointers = n.NumKeys
	} else {
		num_pointers = n.NumKeys + 1
	}

	i = 0
	for n.Pointers[i] != pointer {
		i += 1
	}
	for i += 1; i < num_pointers; i++ {
		n.Pointers[i-1] = n.Pointers[i]
	}
	n.NumKeys -= 1

	if n.IsLeaf {
		for i = n.NumKeys; i < order-1; i++ {
			n.Pointers[i] = nil
		}
	} else {
		for i = n.NumKeys + 1; i < order; i++ {
			n.Pointers[i] = nil
		}
	}

	return n
}

func (t *Tree) adjustRoot() {
	var new_root *Node

	if t.Root.NumKeys > 0 {
		return
	}

	if !t.Root.IsLeaf {
		new_root, _ = t.Root.Pointers[0].(*Node)
		new_root.Parent = nil
	} else {
		new_root = nil
	}
	t.Root = new_root

	return
}

func (t *Tree) coalesceNodes(n, neighbour *Node, neighbour_index, k_prime int) {
	var i, j, neighbour_insertion_index, n_end int
	var tmp *Node

	if neighbour_index == -1 {
		tmp = n
		n = neighbour
		neighbour = tmp
	}

	neighbour_insertion_index = neighbour.NumKeys

	if !n.IsLeaf {
		neighbour.Keys[neighbour_insertion_index] = k_prime
		neighbour.NumKeys += 1

		n_end = n.NumKeys
		i = neighbour_insertion_index + 1
		for j = 0; j < n_end; j++ {
			neighbour.Keys[i] = n.Keys[j]
			neighbour.Pointers[i] = n.Pointers[j]
			neighbour.NumKeys += 1
			n.NumKeys -= 1
			i += 1
		}
		neighbour.Pointers[i] = n.Pointers[j]

		for i = 0; i < neighbour.NumKeys+1; i++ {
			tmp, _ = neighbour.Pointers[i].(*Node)
			tmp.Parent = neighbour
		}
	} else {
		i = neighbour_insertion_index
		for j = 0; j < n.NumKeys; j++ {
			neighbour.Keys[i] = n.Keys[j]
			n.Pointers[i] = n.Pointers[j]
			neighbour.NumKeys += 1
		}
		neighbour.Pointers[order-1] = n.Pointers[order-1]
	}

	t.deleteEntry(n.Parent, k_prime, n)
}

func (t *Tree) redistributeNodes(n, neighbour *Node, neighbour_index, k_prime_index, k_prime int) {
	var i int
	var tmp *Node

	if neighbour_index != -1 {
		if !n.IsLeaf {
			n.Pointers[n.NumKeys+1] = n.Pointers[n.NumKeys]
		}
		for i = n.NumKeys; i > 0; i-- {
			n.Keys[i] = n.Keys[i-1]
			n.Pointers[i] = n.Pointers[i-1]
		}
		if !n.IsLeaf { // why the second if !n.IsLeaf
			n.Pointers[0] = neighbour.Pointers[neighbour.NumKeys]
			tmp, _ = n.Pointers[0].(*Node)
			tmp.Parent = n
			neighbour.Pointers[neighbour.NumKeys] = nil
			n.Keys[0] = k_prime
			n.Parent.Keys[k_prime_index] = neighbour.Keys[neighbour.NumKeys-1]
		} else {
			n.Pointers[0] = neighbour.Pointers[neighbour.NumKeys-1]
			neighbour.Pointers[neighbour.NumKeys-1] = nil
			n.Keys[0] = neighbour.Keys[neighbour.NumKeys-1]
			n.Parent.Keys[k_prime_index] = n.Keys[0]
		}
	} else {
		if n.IsLeaf {
			n.Keys[n.NumKeys] = neighbour.Keys[0]
			n.Pointers[n.NumKeys] = neighbour.Pointers[0]
			n.Parent.Keys[k_prime_index] = neighbour.Keys[1]
		} else {
			n.Keys[n.NumKeys] = k_prime
			n.Pointers[n.NumKeys+1] = neighbour.Pointers[0]
			tmp, _ = n.Pointers[n.NumKeys+1].(*Node)
			tmp.Parent = n
			n.Parent.Keys[k_prime_index] = neighbour.Keys[0]
		}
		for i = 0; i < neighbour.NumKeys-1; i++ {
			neighbour.Keys[i] = neighbour.Keys[i+1]
			neighbour.Pointers[i] = neighbour.Pointers[i+1]
		}
		if !n.IsLeaf {
			neighbour.Pointers[i] = neighbour.Pointers[i+1]
		}
	}
	n.NumKeys += 1
	neighbour.NumKeys -= 1

	return
}

func (t *Tree) deleteEntry(n *Node, key int, pointer interface{}) {
	var min_keys, neighbour_index, k_prime_index, k_prime, capacity int
	var neighbour *Node

	n = removeEntryFromNode(n, key, pointer)

	if n == t.Root {
		t.adjustRoot()
		return
	}

	if n.IsLeaf {
		min_keys = cut(order - 1)
	} else {
		min_keys = cut(order) - 1
	}

	if n.NumKeys >= min_keys {
		return
	}

	neighbour_index = getNeighbourIndex(n)

	if neighbour_index == -1 {
		k_prime_index = 0
	} else {
		k_prime_index = neighbour_index
	}

	k_prime = n.Parent.Keys[k_prime_index]

	if neighbour_index == -1 {
		neighbour, _ = n.Parent.Pointers[1].(*Node)
	} else {
		neighbour, _ = n.Parent.Pointers[neighbour_index].(*Node)
	}

	if n.IsLeaf {
		capacity = order
	} else {
		capacity = order - 1
	}

	if neighbour.NumKeys+n.NumKeys < capacity {
		t.coalesceNodes(n, neighbour, neighbour_index, k_prime)
		return
	} else {
		t.redistributeNodes(n, neighbour, neighbour_index, k_prime_index, k_prime)
		return
	}

}
