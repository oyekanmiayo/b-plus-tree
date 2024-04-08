package main

import (
	"errors"
	"fmt"
	"slices"
)

/*
type Node struct {
	kind   NodeType
	// we introduce the first convenience pointer (parent)
	// consider the _cost_ of all the pointer bookeeping we pay
	// for maintaining what could/should have been put on the callstack
	parent *Node
	keys     []int
	children []*Node
	data     []int
}
*/

func (t *BTree) Insert(key int) error {
	if t.root == nil {
		t.root = &Node{kind: ROOT_NODE}
		t.root.insert(t, key)

		return nil
	} else {
		// find leaf node to insert into or root at first
		n, _, err := t.root.Search(key)

		if err == nil {
			return errors.New("duplicate key/value")
		}

		return n.insert(t, key)
	}
}

func (n *Node) insert(t *BTree, key int) error {
	if n.kind == ROOT_NODE && len(n.children) == 0 {
		n.data = append(n.data, key)
		n.keys = append(n.keys, key)
	}

	if n.kind == LEAF_NODE {
		n.data = append(n.data, key)
	}

	if len(n.data) < MAX_DEGREE {
		return nil
	} else {

		/*
			uncomment to see the splitting per node
			fmt.Printf("node overfull now splitting leaf %v", n.data)
			fmt.Println()
		*/
		n.split(t, len(n.data)/2)
	}

	return nil
}

/*
see what a 'production' split looks like, the difference is night and day :)
https://github.com/cockroachdb/pebble/blob/c4daad9128e053e496fa7916fda8b6df57256823/internal/manifest/btree.go#L403
*/
func (n *Node) split(t *BTree, midIdx int) error {
	switch n.kind {
	case LEAF_NODE:
		splitPoint := n.data[midIdx]
		left, right := n.data[:midIdx], n.data[midIdx:]
		n.data = left

		newNode := &Node{kind: LEAF_NODE, parent: n.parent, data: right}

		n.parent.children = append(n.parent.children, newNode)
		n.parent.keys = append(n.parent.keys, splitPoint)

		// sibling pointers - only on leaf nodes
		n.next = newNode
		newNode.previous = n

	case INTERNAL_NODE:
		splitPoint := n.keys[midIdx]

		// NB: note it's index/key + 1 for internal
		left, right := n.keys[:midIdx], n.keys[midIdx+1:]
		n.keys = left

		newNode := &Node{kind: INTERNAL_NODE, keys: right, parent: n.parent}
		n.parent.children = append(n.parent.children, newNode)
		n.parent.keys = append(n.parent.keys, splitPoint)

		/*
			Notice that the splitting operation modifies three nodes.
			 This is why it is important that the (internal) nodes of a B-tree DO NOT maintain parent pointers.
		*/
		// pointer relocation/bookkeeping
		mid := len(n.children) / 2
		leftPointers, rightPointers := n.children[:mid], n.children[mid:]

		for _, child := range rightPointers {
			child.parent = newNode
		}

		n.children, newNode.children = leftPointers, rightPointers

	case ROOT_NODE:
		if len(n.data) == 0 {
			splitPoint := n.keys[midIdx]
			left, right := n.keys[:midIdx], n.keys[midIdx+1:]

			// demote current root
			newRoot := &Node{kind: ROOT_NODE, parent: nil}
			newRoot.keys = append(newRoot.keys, splitPoint)
			t.root = newRoot

			// pointer relocation/bookkeeping
			mid := len(n.children) / 2
			leftPointers, rightPointers := n.children[:mid], n.children[mid:]
			sibling := &Node{kind: INTERNAL_NODE, keys: left, children: leftPointers, parent: newRoot}
			n.kind, n.keys, n.children, n.parent = INTERNAL_NODE, right, rightPointers, newRoot
			newRoot.children = append(newRoot.children, sibling, n)

			for _, child := range leftPointers {
				child.parent = sibling
			}

		} else {
			// demote current root to a leaf
			n.keys = []int{}
			n.kind = LEAF_NODE
			newRoot := &Node{kind: ROOT_NODE, parent: nil}
			n.parent = newRoot
			t.root = newRoot

			newRoot.children = append(newRoot.children, n)

			n.split(t, len(n.data)/2)
		}

	}

	if len(n.parent.keys) > MAX_DEGREE-1 {
		n.parent.split(t, len(n.parent.keys)/2)
	}

	return nil
}

func BasicInsertLeafExample() {
	var tree BTree

	// NB: keys are values and vice versa
	for i := 1; i <= 8; i++ {
		tree.Insert(i)
	}

	// with the correct storage of elements in our leaf node
	// let's add elements to resemble the example image
	fmt.Println()
	fmt.Println("----built the example 2,3-TREE----")
	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
	fmt.Println(tree.root.children[2])
}

// public api for fuzzer
func KeyExists(t *BTree, key int) bool {
	n, _, err := t.root.Search(key)

	if err != nil {
		return false
	}

	_, found := slices.BinarySearch(n.data, key)

	return found
}
