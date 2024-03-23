package main

import (
	"errors"
	"fmt"
	"slices"
)

type BTree struct {
	root *Node
}

func (t *BTree) Insert(key int) error {
	if t.root == nil {
		t.root = &Node{kind: ROOT_NODE}
		t.root.insert(t, key)

		return nil
	} else {
		// find leaf node to insert into or root at first
		n, _, err := t.root.search(key)
		if err == nil {
			return errors.New("duplicate key/value")
		}

		return n.insert(t, key)
	}
}

type node struct {
	kind     NodeType
	keys     []int
	children []*Node
	data     []int
}

/*
Breadcrumbs contain references to the nodes followed from the root and are used
to backtrack them in reverse when propagating splits or merges.
*/
type BTstack []*node

func (n *Node) insert(t *BTree, key int) error {
	if n.kind == ROOT_NODE && len(n.children) == 0 {
		n.data = append(n.data, key)
		n.keys = append(n.keys, key)

		slices.Sort(n.keys)
	}

	if n.kind == LEAF_NODE {
		n.data = append(n.data, key)
	}

	if len(n.data) < MAX_DEGREE {
		return nil
	} else {
		n.split(t, len(n.data)/2)
	}

	return nil
}

func (n *Node) split(t *BTree, midIdx int) error {
	switch n.kind {
	case LEAF_NODE:
		splitPoint := n.data[midIdx]
		left, right := n.data[:midIdx], n.data[midIdx:]
		n.data = left

		newNode := &Node{kind: LEAF_NODE, parent: n.parent, data: right}

		n.parent.children = append(n.parent.children, newNode)
		n.parent.keys = append(n.parent.keys, splitPoint)

	case INTERNAL_NODE:
		splitPoint := n.keys[midIdx]

		// NB: note it's index/key + 1 for internal
		left, right := n.keys[:midIdx], n.keys[midIdx+1:]
		n.keys = left

		newNode := &Node{kind: INTERNAL_NODE, keys: right}
		n.parent.children = append(n.parent.children, newNode)
		n.parent.keys = append(n.parent.keys, splitPoint)

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

func BreadcrumbInsertExample() {
	var tree BTree

	// NB: keys are values and vice versa
	for i := 1; i <= 8; i++ {
		tree.Insert(i)
	}

	fmt.Println()
	fmt.Println("----built the example 2,3-TREE----")
	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
	fmt.Println(tree.root.children[2])
}

// Collect breadcrumbs on a stack
// this a global variable, for there are lots of ways to implement
// a stack, notably since you're using the callstack anyway, can stuff it in there too.
func (n *Node) search(key int) (*Node, int, error) {
	idx, found := slices.BinarySearch(n.keys, key)

	if found {
		if len(n.children) == 0 {
			return n, idx, nil
		} else {
			return nil, 0, errors.ErrUnsupported
		}

	}

	if len(n.children) == 0 {
		return n, 0, errors.New("key not found, at leaf containing key")
	}

	return n.children[idx].search(key)
}
