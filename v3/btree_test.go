package main

import (
	"testing"
)

func TestInsertRoot(t *testing.T) {
	tree := newBTree(2)

	tree.insert(1, []byte("hello"))
}

func TestInsertAfterSplit(t *testing.T) {
	tree := newBTree(3)

	for i := 1; i <= 10; i++ {
		tree.insert(i, []byte("hello"))
	}
}

/*
// maybe bench the bin search vs a linear
func TestFind(t *testing.T) {
	tree := newBTree(3)

	for i := 1; i <= 10; i++ {
		tree.insert(i, []byte("hello"))
	}

	fmt.Println("? wtf")

	for i := 1; i <= 10; i++ {
		res := tree.find(i)

		if !bytes.Equal(res, []byte("hello")) {
			t.Error("could not find node")
		}
	}
}

func TestDelete(t *testing.T) {
	tree := newBTree(3)

	for i := 1; i <= 10; i++ {
		tree.insert(i, []byte("hello"))
	}

	for i := 1; i <= 10; i++ {
		tree.delete(i)
	}

	for i := 1; i <= 10; i++ {
		res := tree.find(i)

		if res != nil {
			return t.Error("did not remove node correctly")
		}
	}

}

*/
