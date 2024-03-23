package main

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

func main() {
	BasicDeleteExample()
}
