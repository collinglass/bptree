package bptree

import (
	"fmt"
	"testing"
)

func hello() {
	fmt.Println("bptree says 'hello friend'")
}

func TestInsertNilRoot(t *testing.T) {
	var root *node
	hello()

	key := 1
	value := 5

	root, err := Insert(nil, key, value)

	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := Find(root, key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}
}

func TestInsert(t *testing.T) {
	var root *node

	key := 1
	value := 5

	root, err := Insert(root, key, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := Find(root, key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}
}

func TestInsertSameKeyTwice(t *testing.T) {
	var root *node

	key := 1
	value := 5

	root, err := Insert(root, key, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	root, err = Insert(root, key, value+1)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	r, err := Find(root, key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	if root.num_keys > 1 {
		t.Errorf("expected 1 key and got %d", root.num_keys)
	}
}

func TestInsertSameValueTwice(t *testing.T) {
	var root *node

	key := 1
	value := 5

	root, err := Insert(root, key, value)
	if err != nil {
		t.Errorf("%s", err)
	}
	root, err = Insert(root, key+1, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := Find(root, key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	if root.num_keys <= 1 {
		t.Errorf("expected more than 1 key and got %d", root.num_keys)
	}
}

func TestFindNilRoot(t *testing.T) {
	var root *node

	r, err := Find(root, 1, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("expected nil got %s \n", r)
	}
}

func TestFind(t *testing.T) {
	var root *node

	key := 1
	value := 5

	root, err := Insert(root, key, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := Find(root, key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}
}

func TestDeleteNilTree(t *testing.T) {
	var root *node

	key := 1

	root, err := Delete(root, key)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	r, err := Find(root, key, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("returned struct after delete \n")
	}
}

func TestDelete(t *testing.T) {
	var root *node

	key := 1
	value := 5

	root, err := Insert(root, key, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := Find(root, key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	root, err = Delete(root, key)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	r, err = Find(root, key, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("returned struct after delete \n")
	}
}

func TestDeleteNotFound(t *testing.T) {
	var root *node

	key := 1
	value := 5

	root, err := Insert(root, key, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := Find(root, key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	root, err = Delete(root, key+1)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	r, err = Find(root, key+1, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}
}

func TestMultiInsertSingleDelete(t *testing.T) {
	var root *node

	key := 1
	value := 5

	root, err := Insert(root, key, value)
	if err != nil {
		t.Errorf("%s", err)
	}
	root, err = Insert(root, key+1, value+1)
	if err != nil {
		t.Errorf("%s", err)
	}
	root, err = Insert(root, key+2, value+2)
	if err != nil {
		t.Errorf("%s", err)
	}
	root, err = Insert(root, key+3, value+3)
	if err != nil {
		t.Errorf("%s", err)
	}
	root, err = Insert(root, key+4, value+4)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := Find(root, key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	root, err = Delete(root, key)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	r, err = Find(root, key, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("returned struct after delete - %v \n", r)
	}
}

func TestMultiInsertMultiDelete(t *testing.T) {
	var root *node

	key := 1
	value := 5

	root, err := Insert(root, key, value)
	if err != nil {
		t.Errorf("%s", err)
	}
	root, err = Insert(root, key+1, value+1)
	if err != nil {
		t.Errorf("%s", err)
	}
	root, err = Insert(root, key+2, value+2)
	if err != nil {
		t.Errorf("%s", err)
	}
	root, err = Insert(root, key+3, value+3)
	if err != nil {
		t.Errorf("%s", err)
	}
	root, err = Insert(root, key+4, value+4)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := Find(root, key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	root, err = Delete(root, key)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	r, err = Find(root, key, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("returned struct after delete - %v \n", r)
	}

	r, err = Find(root, key+3, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value+3 {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	root, err = Delete(root, key+3)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	r, err = Find(root, key+3, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("returned struct after delete - %v \n", r)
	}
}
