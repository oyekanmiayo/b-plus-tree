package main

func main() {
	// the keys are: 1, 2, 3, 4, 5, 6, 7, 8
	// this is pre-populated
	// this is a simple 2-way b-tree
	// factoid: this is a 2-3 tree, B-Trees generalise 2-3 Trees
	// this example will break with a bigger degree, it's just to show the operations
	BasicSearchExample(6)
	BasicSearchLeaf(6)
}
