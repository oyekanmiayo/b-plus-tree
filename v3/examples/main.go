package main

type NodeType int

/* All examples share the same basic tree structure */
const (
	MAX_DEGREE int = 3
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

	// sibling pointers these help with deletions + range queries
	next     *Node
	previous *Node
}

// DISCLAIMER/NOTE:
// this implementation is, an illustrative example, probably buggy -- I don't know how.
// there are some unit tests, but not enough, there are probably hidden bugs. Bugs are candy for the mind anyway.
// it is also radically inefficient; heavy/unnecessary use of pointers.
// poorly generic structs, append everywhere rather than considering copy
// the garbage collection strategy is non-existent. There are no hints to the compiler/runtime.
// is not thread-safe. lacks tests.

func main() {
	/*
		// the keys are: 1, 2, 3, 4, 5, 6, 7, 8
		// the value of data on a node if supported by the example are the same as keys

		// this is a simple 2-way b-tree
		// factoid: this is a 2-3 tree, B-Trees generalise BSTs, 2-3 Trees
		// this example (may) break with a bigger degree

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


		// NB: this examples assumes all keys and data are sorted or added monotonically
		BasicInsertLeafExample()
		//BreadcrumbInsertExample()
		//RebalancingExample()

	*/

	///////////////
	/// DELETE ////
	//////////////
	var exampleTreeOne BTree
	var exampleTreeTwo BTree
	var exampleTreeThree BTree

	for i := 1; i <= 4; i++ {
		exampleTreeOne.Insert(i)
	}

	for i := 1; i <= 5; i++ {
		exampleTreeTwo.Insert(i)
	}

	for i := 1; i <= 8; i++ {
		exampleTreeThree.Insert(i)
	}

	// BasicDeleteExample(&exampleTreeOne)
	//MergeDeleteExample(&exampleTreeTwo)
	MergeDeleteExample(&exampleTreeThree)
}
