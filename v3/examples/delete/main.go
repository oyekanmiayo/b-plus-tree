package main

import "slices"

type NodeType int

// using the same example as search
// global constants
const (
	MAX_DEGREE = 3
)

const (
	ROOT_NODE NodeType = iota + 1
	INTERNAL_NODE
	LEAF_NODE
)

type Node struct {
	kind     NodeType
	parent   *Node
	keys     []int
	children []*Node
	data     []int
}

func main() {
	BasicDeleteExample()
}

func (n *Node) Search(key int) (*Node, int, error) {
	idx, found := slices.BinarySearch(n.keys, key)

	if found {
		if n.kind == LEAF_NODE {
			return n, idx, nil
		} else {
			return n.children[idx+1].Search(key)
		}
	}

	if len(n.children) == 0 {
		return n, 0, nil
	}

	return n.children[idx].Search(key)
}

func Splice(idx int, elems []int) []int {
	if len(elems) == 1 {
		return nil
	} else {
		return append(elems[:idx], elems[idx+1:]...)
	}
}
