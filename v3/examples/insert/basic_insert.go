package main

import (
	"errors"
	"slices"
)

type Node struct {
	keys     []int
	children []*Node
	leaf     bool
	data     []int
}

// see: `basic search of how/why this works`
func (n *Node) search(key int) (*Node, int, error) {
	idx, found := slices.BinarySearch(n.keys, key)
	if n.leaf {
		if !found {
			return nil, 0, errors.New("key not found")
		} else {
			return n, idx, nil
		}
	}

	return n.children[idx].search(key)
}

func (n *Node) insert(key int) error {
	return nil
}

func BasicInsertExample() {
	// our root node is special, we store keys and store pointers to
	// data records inside it at first, when we split our root for the first time
	// we 'recurse' 'up' pushing a new internal node upward, rearranging our pointers
	root := &Node{}

	for key := 1; key <= 8; key++ {
		root.insert(key)
	}
}
