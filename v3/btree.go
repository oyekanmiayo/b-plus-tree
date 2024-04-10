package main

import (
	"errors"
)

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

func (t *BTree) Search(key int) ([]int, int, error) {
	if t.root == nil {
		return nil, 0, errors.New("empty tree")
	} else {
		node, idx, err := t.root.Search(key)

		if err == nil {
			return nil, 0, errors.New("not found")
		} else {
			return node.data, idx, nil
		}
	}
}

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

func (t *BTree) Delete(key int) error {
	if t.root == nil {
		return errors.New("empty tree")
	} else {
		// find leaf node to delete from or root
		n, _, err := t.root.SearchDelete(key)

		if err == nil {
			return n.delete(t, key)
		}

		return errors.New("key not in tree")
	}
}
