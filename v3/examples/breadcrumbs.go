package main

import (
	"errors"
	"fmt"
	"slices"
)

type BreadBTree struct {
	root  *node
	stack BTstack
}

type node struct {
	kind     NodeType
	keys     []int
	children []*node
	data     []int
}

/*
Breadcrumbs contain references to the nodes followed from the root and are used
to backtrack them in reverse when propagating splits or merges.
*/
type BTstack []*node

func (t *BreadBTree) Insert(key int) error {
	if t.root == nil {
		t.stack = make(BTstack, MAX_DEGREE)

		t.root = &node{kind: ROOT_NODE}
		t.root.insert(t, key)

		return nil
	} else {
		n, _, err := t.root.breadCrumbSearch(key)
		if err == nil {
			return errors.New("duplicate key/value")
		}

		return n.insert(t, key)
	}
}

func (n *node) insert(t *BreadBTree, key int) error {
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

// TODO: refactor split OP pointers to use the BTStack
func (n *node) split(t *BreadBTree, midIdx int) error {
	return nil
}

func BreadcrumbInsertExample() {
	var tree BreadBTree

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
func (n *node) breadCrumbSearch(key int) (*node, int, error) {
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

	return n.children[idx].breadCrumbSearch(key)
}
