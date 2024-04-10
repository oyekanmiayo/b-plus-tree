package main

import (
	"errors"
	"fmt"
	"slices"
)

func (t *BTree) BasicDelete(key int) error {
	if t.root == nil {
		return errors.New("empty tree")
	} else {
		// find leaf Node to delete
		n, _, err := t.root.SearchDelete(key)

		if err == nil {
			return n.basicDelete(key)
		}

		return errors.New("key not in tree")
	}
}

func (n *Node) basicDelete(key int) error {
	for i, v := range n.data {
		if v == key {
			n.data = cut(i, n.data)
		}
	}

	for i, k := range n.parent.keys {
		if k == key {
			n.parent.keys = cut(i, n.parent.keys)
		}
	}

	return errors.New("see: merge.go for merges")
}

func cut(idx int, elems []int) []int {
	if len(elems) == 1 {
		return nil
	} else {
		return append(elems[:idx], elems[idx+1:]...)
	}
}

func cutPointer(idx int, elems []*Node) []*Node {
	if len(elems) == 1 {
		return nil
	} else {
		return append(elems[:idx], elems[idx+1:]...)
	}
}

// HELPERS
func (n *Node) SearchDelete(key int) (*Node, int, error) {
	idx, found := slices.BinarySearch(n.keys, key)

	if found {
		if n.kind == LEAF_NODE {
			return n, idx, nil
		} else {
			if len(n.children) == 0 {
				return n, 0, nil
			}

			if idx+1 > len(n.children) {
				return n.children[idx].SearchDelete(key)
			}

			return n.children[idx+1].SearchDelete(key)
		}
	}

	if len(n.children) == 0 {
		return n, 0, nil
	}

	return n.children[idx].SearchDelete(key)
}

func BasicDeleteExample(tree *BTree) {
	fmt.Println("---initial tree shape---")
	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
	fmt.Println(tree.root.children[2])

	fmt.Println("---delete:4 no cascade---")
	tree.BasicDelete(4)

	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
	fmt.Println(tree.root.children[2])
	fmt.Println()

	tree.BasicDelete(3)
	fmt.Println("---delete:3 simple cascade root---")
	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
	fmt.Println("oh no! we now have an orphaned node")

}
