package main

import (
	"errors"
	"fmt"
	"slices"
)

func (t *BTree) BasicDelete(key int) error {
	if t.root == nil {
		return errors.New("empty tree")
	} else {
		// find leaf Node to delete
		n, _, err := t.root.SearchDelete(key)
		fmt.Println(n)
		if err == nil {
			return n.basicDelete(key)
		}

		return errors.New("key not in tree")
	}
}

func (n *Node) basicDelete(key int) error {
	for i, v := range n.data {
		if v == key {
			n.data = splice(i, n.data)
		}
	}

	for i, k := range n.parent.keys {
		if k == key {
			n.parent.keys = splice(i, n.parent.keys)
			n.parent.children = append(n.parent.children[:i], n.parent.children[i+1])
		}
	}

	return errors.New("see: merge.go for merges")
}

func splice(idx int, elems []int) []int {
	if len(elems) == 1 {
		return nil
	} else {
		return append(elems[:idx], elems[idx+1:]...)
	}
}

// HELPERS
func (n *Node) SearchDelete(key int) (*Node, int, error) {
	idx, found := slices.BinarySearch(n.keys, key)

	if found {
		if n.kind == LEAF_NODE {
			return n, idx, nil
		} else {
			return n.children[idx+1].SearchDelete(key)
		}
	}

	if len(n.children) == 0 {
		return n, 0, nil
	}

	return n.children[idx].SearchDelete(key)
}

func BasicDeleteExample(tree *BTree) {

	fmt.Println(tree)
	// delete no cascade
	fmt.Println(tree.BasicDelete(4))

	// delete simple cascade root
	fmt.Println(tree.BasicDelete(3))

	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
}
