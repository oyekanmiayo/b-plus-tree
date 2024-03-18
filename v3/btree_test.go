package main

import (
	"testing"
)

func TestInsert(t *testing.T) {
	tree := newBTree(3)

	tree.insert(1, []byte("hello"))
	tree.insert(2, []byte("hello"))
	tree.insert(3, []byte("hello"))
	tree.insert(4, []byte("hello"))
	tree.insert(5, []byte("hello"))
}
