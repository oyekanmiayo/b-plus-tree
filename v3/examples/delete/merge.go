package main

import (
	"errors"
	"fmt"
)

type BPlusTree struct {
	root *Node
}

func (t *BPlusTree) Delete(key int) error {
	if t.root == nil {
		return errors.New("empty tree")
	} else {
		n, _, err := t.root.Search(key)

		if err == nil {
			return n.delete(t, key)
		}

		return errors.New("key not in tree")
	}
}

func (n *Node) delete(t *BPlusTree, key int) error {
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

	if err := n.mergeSiblings(t, key); err == nil {
		return nil
	}

	return nil
}

func (n *Node) mergeSiblings(t *BPlusTree, key int) error {
	return nil
}

func BasicMergeExample() {
	var tree BPlusTree

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
	// delete simple cascade with root
	fmt.Println(tree.Delete(3))

	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
}
