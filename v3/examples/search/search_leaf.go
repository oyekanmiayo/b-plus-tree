package main

import (
	"fmt"
)

// Define a struct for B+ tree node
type Node struct {
	keys     []int
	children []*Node
	leaf     bool  // is a leaf node
	data     []int // Data stored
}

/*
more clearly without the generic `Node` and implict `bool`:
a leaf node:
type LeafNode struct {
	leaf    bool
	data    []int
	parent  *InternalNode // Pointer to parent node for easier navigation
	next    *LeafNode     // Pointer to the next leaf node for range queries
}

an internal node:
type InternalNode struct {
	keys     []int
	children []*InternalNode // For internal nodes, children are other internal nodes
	parent   *InternalNode   // Pointer to parent node for easier navigation
}
*/

func (n *Node) search(key int) (*Node, int) {
	if n.leaf {
		// If it's a leaf node, return the leaf node and the index where the key would be or is found
		return n, binarySearch(n.keys, key)
	}

	// If it's not a leaf node, recursively search for the appropriate child node
	i := binarySearch(n.keys, key)
	return n.children[i].search(key)
}

func binarySearch(arr []int, key int) int {
	low, high := 0, len(arr)-1
	for low <= high {
		mid := low + (high-low)/2
		if arr[mid] == key {
			return mid
		} else if arr[mid] < key {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	// see for practice: https://leetcode.com/problems/search-insert-position/
	return low // where it would be or should be inserted
}

func BasicSearchLeaf(key int) {
	root := &Node{
		keys: []int{3, 5},
		children: []*Node{
			{
				keys:     []int{2},
				leaf:     true,
				data:     []int{1},
				children: []*Node{},
			},
			{
				keys:     []int{4},
				leaf:     true,
				data:     []int{3},
				children: []*Node{},
			},
			{
				keys:     []int{6, 7},
				leaf:     true,
				data:     []int{5, 7},
				children: []*Node{},
			},
		},
	}

	result, index := root.search(key)

	fmt.Println(result)
	fmt.Println(index)
}
