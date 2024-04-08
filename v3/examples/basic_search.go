package main

import "fmt"

/*
this basicSearch does not concern itself with the leaf data, only keys.
type Node struct {
	keys     []int
	children []*Node
	// leaf Node data --snipped
}
*/

func (n *Node) basicSearch(key int) *Node {
	if len(n.children) == 0 {
		// you are at a leaf Node and can now access stuff
		// later when you add sibling pointers you can follow the next pointer
		return n
	}

	// alternatively more idiomatically: slices.BinarySearch(m.children, key)
	// this is here for reference/clarity
	low, high := 0, len(n.keys)-1
	for low <= high {
		mid := low + (high-low)/2
		if n.keys[mid] == key {
			return n
		} else if n.keys[mid] < key {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return n.children[low]
}

func BasicSearchExample(t *BTree, key int) {
	fmt.Println(t.root.basicSearch(key))
}
