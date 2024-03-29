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
			return n.delete(key)
		}

		return errors.New("key not in tree")
	}
}

func (n *Node) delete(key int) error {
	for i, v := range n.data {
		if v == key {
			n.data = splice(i, n.data)
		}
	}

	if n.kind == ROOT_NODE {
		return nil
	}

	if len(n.data) == 0 {
		if sibling, err := n.preMerge(); err == nil {
			return n.mergeSibling(sibling, key)
		} else {
			return errors.New("see rebalancing.go")
		}
	}

	// nodes other than root must be at least half full.
	if (n.kind != ROOT_NODE) && (len(n.parent.keys) < (MAX_DEGREE-1)/2) {
		if sibling, err := n.preMerge(); err == nil {
			return n.mergeSibling(sibling, key)
		} else {
			return errors.New("see rebalancing.go")
		}
	}

	// check tree invariants and recurse upward
	for i, k := range n.parent.keys {
		if k == key {
			n.parent.keys = splice(i, n.parent.keys)
			newSplit := len(n.data) / 2
			n.parent.keys = append(n.parent.keys, n.data[newSplit])

			if len(n.parent.keys) < ((MAX_DEGREE - 1) / 2) {
				if sibling, err := n.parent.preMerge(); err == nil {
					return n.parent.mergeSibling(sibling, key)
				} else {
					return errors.New("see rebalancing.go")
				}
			}
		}
	}

	return nil
}

/*
contents should be merged. if their contents do not fit into a single node
eys are redistributed - rebalancing.go.
*/
// merging can be... very interesting.
// you can slap on an iter api like(rust):
// https://github.com/rust-lang/rust/blob/1c19595575968ea77c7f85e97c67d44d8c0f9a68/library/alloc/src/collections/btree/merge_iter.rs#L41
// and maybe... just maybe, stream/lift that iter out to a scheduler/async runtime -- complex, magical, do not do this, but neat to know.

// go/pebble
// iterator/cursor: https://github.com/cockroachdb/pebble/blob/c4daad9128e053e496fa7916fda8b6df57256823/internal/manifest/btree.go#L973 &&
// https://github.com/cockroachdb/pebble/blob/c4daad9128e053e496fa7916fda8b6df57256823/internal/manifest/btree.go#L891

// the actual merge operation
// https://github.com/cockroachdb/pebble/blob/c4daad9128e053e496fa7916fda8b6df57256823/internal/manifest/btree.go#L620
func (n *Node) mergeSibling(sibling *Node, key int) error {
	if n.parent != sibling.parent {
		panic("sibling invariant not satisfied")
	}

	switch n.kind {
	case LEAF_NODE:
		sibling.keys = append(sibling.keys, n.keys...)

		for i, node := range sibling.parent.children {
			if node == n {
				n.parent.children = append(n.parent.children[:i], n.parent.children[i+1:]...)
			}
		}

		for i, k := range n.parent.keys {
			if k == key {
				n.parent.keys = splice(i, n.parent.keys)
				newSplit := len(n.data) / 2

				if len(n.data) != 0 {
					n.parent.keys = append(n.parent.keys, n.data[newSplit])
				}

				if len(n.parent.keys) < ((MAX_DEGREE - 1) / 2) {
					if sibling, err := n.parent.preMerge(); err == nil {
						return n.parent.mergeSibling(sibling, key)
					} else {
						return errors.New("see rebalancing.go")
					}
				}
			}
		}

	case INTERNAL_NODE:
		if len(n.parent.keys) < ((MAX_DEGREE - 1) / 2) {
			if sibling, err := n.parent.preMerge(); err == nil {
				return n.parent.mergeSibling(sibling, key)
			} else {
				return errors.New("see rebalancing.go")
			}
		}
	}

	return nil
}

// preMerge if two adjacent leaf nodes have a common parent and their contents fit into a single node
func (n *Node) preMerge() (*Node, error) {
	switch n.kind {
	case INTERNAL_NODE:
		// no sibling pointers so we have to go up to parent
		// we need to find the previous sibling or next sibling.
		// prev sibling:
		for i, node := range n.parent.children {
			if n == node {
				leftSibling, foundLeft := boundCheck(i-1 >= 0, n.parent.children[i-1])
				// rightSibling, foundRight := boundCheck(i+1 < len(n.parent.children), n.parent.children[i+1])

				if foundLeft && len(leftSibling.keys)+1 < MAX_DEGREE {
					mergePoint := len(n.parent.keys) / 2
					leftSibling.keys = append(leftSibling.keys, n.parent.keys[mergePoint])
					n.parent.keys = append(n.parent.keys[:mergePoint], n.parent.keys[mergePoint+1:]...)

					n.parent.children = append(n.parent.children[:i], n.parent.children[i+1:]...)
					leftSibling.keys = append(leftSibling.keys, n.keys...)
					leftSibling.children = append(leftSibling.children, n.children...)

					return leftSibling, nil

				} else {
					return nil, errors.New("cannot merge internal node, must redistribute")
				}
			}
		}

	case LEAF_NODE:
		if n.previous != nil {
			if len(n.previous.data)+1 < MAX_DEGREE {
				n.previous.next = n.next
				return n.previous, nil
			}
		}

		if n.next != nil {
			if len(n.next.data)+1 < MAX_DEGREE {
				n.next.previous = n.previous
				return n.next, nil
			}
		}
	}

	return nil, errors.New("cannot merge with sibling")
}

func boundCheck(cond bool, node *Node) (*Node, bool) {
	if cond {
		return node, true
	} else {
		return nil, false
	}
}

func MergeDeleteExample(tree *BTree) {
	// delete no cascade

	// delete causes cascade/merge
	tree.Delete(4)
	tree.Delete(5)
	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
}
