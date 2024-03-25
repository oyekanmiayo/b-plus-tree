package main

type NodeType int

/* All examples share the same basic tree structure */
const (
	MAX_DEGREE = 3
)

const (
	ROOT_NODE NodeType = iota + 1
	INTERNAL_NODE
	LEAF_NODE
)

type BTree struct {
	root *Node
}

type Node struct {
	kind     NodeType
	parent   *Node
	keys     []int
	children []*Node
	data     []int
}

func main() {
	// the keys are: 1, 2, 3, 4, 5, 6, 7, 8
	// the value of data on a node if supported by the example are the same as keys

	// this is a simple 2-way b-tree
	// factoid: this is a 2-3 tree, B-Trees generalise 2-3 Trees
	// this example will break with a bigger degree, it's just to show the operations
	var exampleSearchOne BTree
	var exampleSearchTwo BTree

	// prepopulate
	for i := 1; i <= 8; i++ {
		exampleSearchOne.Insert(i)
		exampleSearchTwo.Insert(i)
	}

	///////////////
	/// SEARCH ////
	///////////////

	//these return values in the leaf, notice the n, n+1 relationship between keys and child pointers

	BasicSearchExample(&exampleSearchOne, 6)
	BasicSearchLeaf(&exampleSearchTwo, 3)
	BasicSearchLeaf(&exampleSearchTwo, 6)
	BasicSearchLeaf(&exampleSearchTwo, 8)

	// returns the value of key 6, this is kind of a badly contrived example
	// tldr; key 6 points to the value of 5
	// because to get to the data 6, you'd need a key between 6 and 7
	// to follow that pointer and get to 6

	///////////////
	/// INSERT ////
	///////////////
	BasicInsertExample()
	BasicInsertLeafExample()
	//BreadcrumbInsertExample()
	//RebalancingExample()

	///////////////
	/// DELETE ////
	//////////////
	var exampleTreeOne BTree
	var exampleTreeTwo BTree

	for i := 1; i <= 4; i++ {
		exampleTreeOne.Insert(i)
	}

	for i := 1; i <= 4; i++ {
		exampleTreeTwo.Insert(i)
	}

	BasicDeleteExample(&exampleTreeOne)
	// MergeDeleteExample(&exampleTreeTwo)
}
