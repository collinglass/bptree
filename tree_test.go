package bptree

import (
	"fmt"
	"testing"
)

func hello() {
	fmt.Println("bptree says 'hello friend'")
}

func TestInsertNilRoot(t *testing.T) {
	tree := NewTree()
	hello()

	key := 1
	value := 5

	err := tree.Insert(key, value)

	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := tree.Find(key, false)
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
	tree := NewTree()

	key := 1
	value := 5

	err := tree.Insert(key, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := tree.Find(key, false)
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
	tree := NewTree()

	key := 1
	value := 5

	err := tree.Insert(key, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = tree.Insert(key, value+1)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	r, err := tree.Find(key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	if tree.Root.num_keys > 1 {
		t.Errorf("expected 1 key and got %d", tree.Root.num_keys)
	}
}

func TestInsertSameValueTwice(t *testing.T) {
	tree := NewTree()

	key := 1
	value := 5

	err := tree.Insert(key, value)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = tree.Insert(key+1, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := tree.Find(key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	if tree.Root.num_keys <= 1 {
		t.Errorf("expected more than 1 key and got %d", tree.Root.num_keys)
	}
}

func TestFindNilRoot(t *testing.T) {
	tree := NewTree()

	r, err := tree.Find(1, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("expected nil got %s \n", r)
	}
}

func TestFind(t *testing.T) {
	tree := NewTree()

	key := 1
	value := 5

	err := tree.Insert(key, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := tree.Find(key, false)
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
	tree := NewTree()

	key := 1

	err := tree.Delete(key)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	r, err := tree.Find(key, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("returned struct after delete \n")
	}
}

func TestDelete(t *testing.T) {
	tree := NewTree()

	key := 1
	value := 5

	err := tree.Insert(key, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := tree.Find(key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	err = tree.Delete(key)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	r, err = tree.Find(key, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("returned struct after delete \n")
	}
}

func TestDeleteNotFound(t *testing.T) {
	tree := NewTree()

	key := 1
	value := 5

	err := tree.Insert(key, value)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := tree.Find(key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	err = tree.Delete(key + 1)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	r, err = tree.Find(key+1, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}
}

func TestMultiInsertSingleDelete(t *testing.T) {
	tree := NewTree()

	key := 1
	value := 5

	err := tree.Insert(key, value)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = tree.Insert(key+1, value+1)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = tree.Insert(key+2, value+2)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = tree.Insert(key+3, value+3)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = tree.Insert(key+4, value+4)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := tree.Find(key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	err = tree.Delete(key)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	r, err = tree.Find(key, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("returned struct after delete - %v \n", r)
	}
}

func TestMultiInsertMultiDelete(t *testing.T) {
	tree := NewTree()

	key := 1
	value := 5

	err := tree.Insert(key, value)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = tree.Insert(key+1, value+1)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = tree.Insert(key+2, value+2)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = tree.Insert(key+3, value+3)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = tree.Insert(key+4, value+4)
	if err != nil {
		t.Errorf("%s", err)
	}

	r, err := tree.Find(key, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	err = tree.Delete(key)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	r, err = tree.Find(key, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("returned struct after delete - %v \n", r)
	}

	r, err = tree.Find(key+3, false)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	if r == nil {
		t.Errorf("returned nil \n")
	}

	if r.Value != value+3 {
		t.Errorf("expected %v and got %v \n", value, r.Value)
	}

	err = tree.Delete(key + 3)
	if err != nil {
		t.Errorf("%s\n", err)
	}

	r, err = tree.Find(key+3, false)
	if err == nil {
		t.Errorf("expected error and got nil")
	}

	if r != nil {
		t.Errorf("returned struct after delete - %v \n", r)
	}
}
