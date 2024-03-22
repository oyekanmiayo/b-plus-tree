package main

import (
	"errors"
	"fmt"
	"slices"
)

type Node struct {
	kind     NodeType
	parent   *Node // we introduce the first convenience pointer (parent)
	keys     []int
	children []*Node
	data     []int
}

func (n *Node) insert(key int) error {
	n.keys = append(n.keys, key)

	if len(n.keys) <= MAX_DEGREE-1 {
		return nil
	} else {
		n.split(len(n.keys) / 2)
	}

	return nil
}

func (n *Node) split(midIdx int) error {
	splitPoint := n.keys[midIdx]
	leftKeys := n.keys[:midIdx]
	rightKeys := n.keys[midIdx+1:]

	n.keys = []int{splitPoint}

	//we must now check MIN_KEYS otherwise our tree breaks down
	leftNode := &Node{kind: LEAF_NODE, keys: leftKeys, parent: n}
	rightNode := &Node{kind: LEAF_NODE, keys: rightKeys, parent: n}
	n.children = []*Node{leftNode, rightNode}

	// TODO: leaf node split

	return nil
}

func BasicInsertLeafExample() {
	root := &Node{kind: ROOT_NODE}

	for key := 1; key <= 8; key++ {
		root.insert(key)

	}

	// with the storage of elements in our leaf node
	// our keys resemble the example image
	fmt.Println(root.search(3))
	fmt.Println(root.search(6))

	fmt.Println(root.keys)
	fmt.Println(root.children[0].keys)
	fmt.Println(root.children[1].keys)

}

// see: basic search of how/why this works
func (n *Node) search(key int) (*Node, int, error) {
	idx, found := slices.BinarySearch(n.keys, key)

	if found {
		return n, idx, nil
	}

	if len(n.children) == 0 {
		return n, 0, errors.New("key not found or at leaf node")
	}

	return n.children[idx].search(key)
}
