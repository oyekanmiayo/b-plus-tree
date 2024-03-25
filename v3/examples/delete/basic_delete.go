package main

import (
	"errors"
	"fmt"
)

type BasicDelBtree struct {
	root *Node
}

func (t *BasicDelBtree) Delete(key int) error {
	if t.root == nil {
		return errors.New("empty tree")
	} else {
		// find leaf Node to delete
		n, _, err := t.root.Search(key)

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

func BasicDeleteExample() {
	var tree BasicDelBtree

	root := &Node{
		kind: ROOT_NODE,
		keys: []int{2, 3}}

	tree.root = root

	// scaffolding..
	root.children = []*Node{
		{
			kind:   LEAF_NODE,
			data:   []int{1},
			parent: root,
		},
		{
			kind:   LEAF_NODE,
			data:   []int{2},
			parent: root,
		},
		{
			kind:   LEAF_NODE,
			data:   []int{3, 4},
			parent: root,
		},
	}

	// delete no cascade
	fmt.Println(tree.Delete(4))

	// delete simple cascade root
	fmt.Println(tree.Delete(3))

	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
}
