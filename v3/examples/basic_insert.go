package main

import (
	"fmt"
)

type BNode struct {
	kind     NodeType
	keys     []int
	children []*BNode
}

// NB: this example is "wrong"
// it is here to clearly illustrate the 'split' alogrithm using only keys.
// the actual insertion is in insert.go, but it includes too many spurious details.
func (n *BNode) insert(key int) error {
	n.keys = append(n.keys, key)

	if len(n.keys) <= MAX_DEGREE-1 {
		// does the target node have room? if yes do nothing
		return nil
	} else {
		// is it a leaf node? the bound: len(n.pointers) < n + 1
		// is it an internal node? the bound is len(n.keys) < n
		// where n >= 2, n == degree == occupancy == fanout
		n.split(len(n.keys) / 2)
	}

	return nil
}

func (n *BNode) split(midIdx int) error {
	// handle two cases:
	// one for an internal node split
	// edge case, how to handle the root node?
	// basic aglorithm.
	splitPoint := n.keys[midIdx]

	fmt.Print("split: ")
	fmt.Println(splitPoint)

	leftKeys := n.keys[:midIdx]
	rightKeys := n.keys[midIdx:]

	n.keys = []int{splitPoint}

	leftNode := &BNode{kind: LEAF_NODE, keys: leftKeys}
	rightNode := &BNode{kind: LEAF_NODE, keys: rightKeys}
	n.children = []*BNode{leftNode, rightNode}

	return nil
}

func BasicInsertExample() {
	// every node except the root node must respect the inquality:
	// branching factor - 1 <= num keys < (2 * branching factor) - 1
	// if this doesn't make sense ignore it. The take away:
	// every node except root has a minimum no. of keys or it's invalid.

	// step two: it's promotion time, split keys into two halves, if there's only halve,
	// ie the new node has < MIN_KEYS we push up one-way:
	// lastly point to data/allocate/move root node's data to a node.

	// -- LEAF
	//  (internal node)  (internal or nil)
	//   \               /
	//   (root)

	// step three, check all children:
	// recurse UP from root to new internal node(s), check that we're not full
	// if full, we split again on internal node, allocate a new node(s)

	root := &BNode{kind: ROOT_NODE}

	for key := 1; key <= 4; key++ {
		root.insert(key)

	}

	fmt.Print("root")
	fmt.Println(root.keys)

	fmt.Print("child one")
	fmt.Println(root.children[0].keys)

	fmt.Print("child two")
	fmt.Println(root.children[1].keys)

}
