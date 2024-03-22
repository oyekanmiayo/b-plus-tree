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
	// locate target leaf
	// append value

	// case 1
	// does the target node have room? if yes do nothing
	// case 2
	// is it a leaf node? the bound: len(n.pointers) < n + 1
	// is it an internal node? the bound is len(n.keys) < n
	// where n >= 2, n == max degree == occupancy == fanout

	return nil
}

func BasicInsertExample() {
	// our root node is special, we store keys and store pointers to
	// data records inside it at first, when we split our root for the first time
	// we 'recurse' 'UP' pushing a new internal node upward, rearranging our pointers
	// then it 'loses' its borrow leaf capabilities.

	// B-Trees are built in-reverse of a classic binary search tree.
	// (new internal) ....can add more nodes until full, split, recurse.
	//   \               /
	//   (root)

	// splitting:
	// step one: oh crap I'm a full root node, I need to give away my keys!
	//  (nil)            (nil)
	//   \               /
	//   (root )

	// every node except the root node must respect the inquality:
	// branching factor - 1 <= num keys < (2 * branching factor) - 1
	// if this doesn't make sense ignore it. The take away:
	// every node except root has a minimum no. of keys or it's invalid.

	// step two: it's promotion time, split keys into two halves, if there's only halve,
	// ie the new node has < MIN_KEYS we push up one-way:
	// lastly point to data/allocate/move root node's data to a leaf.

	// -- LEAF
	//  (internal node)  (internal or nil)
	//   \               /
	//   (root)

	// step three, check all children:
	// recurse DOWN from root to new internal node(s), check that we're not full
	// if full, we split again on internal node, allocate a new node(s)

	root := &Node{}

	for key := 1; key <= 8; key++ {
		root.insert(key)
	}
}
