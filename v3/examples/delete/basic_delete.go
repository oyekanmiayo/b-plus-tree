package main

import (
	"errors"
	"slices"
)

type BPlusTree struct {
	root *Node
}

type Node struct {
	kind     NodeType
	parent   *Node
	keys     []int
	children []*Node
	data     []int
}

func (t *BPlusTree) Delete(key int) error {
	if t.root == nil {
		return errors.New("empty tree")
	} else {
		// find leaf node to delete
		n, _, err := t.root.search(key)

		if err == nil {
			return n.delete(t, key)
		}

		return errors.New("key not in tree")
	}
}

func (n *Node) delete(t *BPlusTree, key int) error {
	idx := 0

	for i, v := range n.data {
		if v == key {
			idx = i
		}
	}

	n.data = append(n.data[:idx], n.data[idx+1:]...)

	for i, k := range n.parent.keys {
		if k == key {
			n.parent.keys = append(n.data[:i], n.data[i+1:]...)
		}
	}

	if err := n.mergeSiblings(t, key); err == nil {
		return nil
	}
	return errors.New("see: rebalancing.go")
}

func BasicDeleteExample() {
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
			data:   []int{1},
			parent: root,
		},
		{
			kind:   LEAF_NODE,
			data:   []int{3, 4},
			parent: root,
		},
	}

	// delete no cascade
	tree.Delete(4)
	/*
		fmt.Println(tree.root)
		fmt.Println(tree.root.children[0])
		fmt.Println(tree.root.children[1])
		fmt.Println(tree.root.children[2])
	*/

	/*
		// delete causes cascade/merge
		tree.Delete(4)
		fmt.Println(tree.root)
		fmt.Println(tree.root.children[0])
		fmt.Println(tree.root.children[1])
		fmt.Println(tree.root.children[2])
	*/
}

func (n *Node) mergeSiblings(t *BPlusTree, key int) error {
	return nil
}

func (n *Node) search(key int) (*Node, int, error) {
	idx, found := slices.BinarySearch(n.keys, key)

	if found {
		if n.kind == LEAF_NODE {
			return n, idx, nil
		} else {
			return nil, 0, errors.ErrUnsupported
		}

	}

	if len(n.children) == 0 {
		return n, 0, nil
	}

	return n.children[idx].search(key)
}
