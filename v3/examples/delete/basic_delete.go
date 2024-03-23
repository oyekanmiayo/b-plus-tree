package main

import (
	"errors"
	"fmt"
	"slices"
)

type BPlusTree struct {
	root *Node
}

// Define a struct for B+ tree node
type Node struct {
	keys     []int
	children []*Node
	leaf     bool  // is a leaf node
	data     []int // Data stored
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
	if err := n.mergeSiblings(t, key); err == nil {
		return nil
	}
	return errors.New("see: rebalancing.go")
}

func BasicDelete(key int) {
	var tree BPlusTree

	// scaffolding..
	root := &Node{
		leaf: false,
		keys: []int{3, 5},
		children: []*Node{
			{
				keys: []int{2},
				leaf: false,
				data: nil,
				children: []*Node{
					{
						leaf: true,
						data: []int{1},
					},
					{
						leaf: true,
						data: []int{2},
					},
				},
			},
			{
				keys: []int{4},
				leaf: false,
				data: nil,
				children: []*Node{
					{
						leaf: true,
						data: []int{3},
					},
					{
						leaf: true,
						data: []int{4},
					},
				},
			},
			{
				keys: []int{6, 7},
				leaf: false,
				data: nil,
				children: []*Node{
					{
						leaf: true,
						data: []int{5},
					},
					{
						leaf: true,
						data: []int{6},
					},
					{
						leaf: true,
						data: []int{7, 8},
					},
				},
			},
		},
	}

	tree.root = root

	// delete no cascade
	tree.Delete(1)
	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
	fmt.Println(tree.root.children[2])

	// delete causes cascade/merge
	tree.Delete(4)
	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
	fmt.Println(tree.root.children[2])
}

func (n *Node) mergeSiblings(t *BPlusTree, key int) error {
	return nil
}

// see: basic search of how/why this works
func (n *Node) search(key int) (*Node, int, error) {
	idx, found := slices.BinarySearch(n.keys, key)

	if found {
		if len(n.children) == 0 {
			return n, idx, nil
		} else {
			return nil, 0, errors.ErrUnsupported
		}

	}

	if len(n.children) == 0 {
		return n, 0, errors.New("key not found, at leaf containing key")
	}

	return n.children[idx].search(key)
}
