package main

import (
	"errors"
	"fmt"
)

func (t *BTree) Delete(key int) error {
	if t.root == nil {
		return errors.New("empty tree")
	} else {
		n, _, err := t.root.SearchDelete(key)

		if err == nil {
			return n.delete(t, key)
		}

		return errors.New("key not in tree")
	}
}

func (n *Node) delete(t *BTree, key int) error {
	for i, v := range n.data {
		if v == key {
			n.data = splice(i, n.data)
		}
	}

	for i, k := range n.parent.keys {
		if k == key {
			n.parent.keys = splice(i, n.parent.keys)
			n.parent.children = append(n.parent.children[:i], n.parent.children[i+1])
		}
	}

	if err := n.mergeSiblings(t, key); err == nil {
		return nil
	}

	return nil
}

func (n *Node) mergeSiblings(t *BTree, key int) error {
	return nil
}

func MergeDeleteExample(tree *BTree) {
	// delete no cascade
	tree.Delete(1)
	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
	fmt.Println(tree.root.children[2])

	// delete causes cascade/merge
	tree.Delete(4)
	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])
	fmt.Println(tree.root.children[2])

}
