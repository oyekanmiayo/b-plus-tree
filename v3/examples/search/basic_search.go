package main

import "fmt"

type node struct {
	keys     []int
	children []*node
	// leaf node data --snipped
}

func (n *node) search(key int) *node {
	if len(n.children) == 0 {
		// you are at a leaf node and can now access stuff
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

func BasicSearchExample(key int) {
	root := &node{
		keys: []int{3, 5},
		children: []*node{
			{
				keys:     []int{2},
				children: []*node{},
			},
			{
				keys:     []int{4},
				children: []*node{},
			},
			{
				keys:     []int{6, 7},
				children: []*node{},
			},
		},
	}

	result := root.search(key)
	fmt.Println(result)
}
