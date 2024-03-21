package main

func main() {
	// the keys are: 1, 2, 3, 4, 5, 6, 7, 8
	// this is pre-populated
	// this is a simple 2-way b-tree
	// factoid: this is a 2-3 tree, B-Trees generalise 2-3 Trees
	// this example will break with a bigger degree, it's just to show the operations

	BasicSearchExample(6) // returns the node with key 6

	//these return values in the leaf, notice the n, n+1 relationship between keys and child pointers
	BasicSearchLeaf(3)
	BasicSearchLeaf(6)
	BasicSearchLeaf(8)
	// returns the value of key 6, this is kind of a badly contrived example
	// tldr; key 6 points to the value of 5
	// because to get to the data 6, you'd need a key between 6 and 7
	// to follow that pointer and get to 6
}
