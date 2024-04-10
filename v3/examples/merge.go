package main

import (
	"errors"
	"fmt"
)

// Deletion is the most complicated operation for a B-Tree.
// this covers part one, "merging"
// step one: find leaf node delete data
// see: https://opendatastructures.org/ods-python/14_2_B_Trees.html#SECTION001723000000000000000
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
			n.data = cut(i, n.data)
		}
	}

	// is the leaf empty? we delete key from parent & merge
	if n.kind == LEAF_NODE && len(n.data) == 0 {
		if sibling, idx, err := n.preMerge(); err == nil {
			return n.mergeSibling(sibling, idx, key)
		} else {
			return errors.New("see rebalancing.go")
		}
	} else {
		// should we update the parent's seperator?
		if n.parent.keys[0] < n.data[0] {
			// delete the key from the parent
			for i, k := range n.parent.keys {
				if k == key {
					n.parent.keys = cut(i, n.parent.keys)
					newSeperator := len(n.data) / 2
					n.parent.keys = append(n.parent.keys, n.data[newSeperator])

					// overflow or underflow triggers a merge cascade recurse to parent
					if len(n.parent.keys) < ((MAX_DEGREE-1)/2) || len(n.parent.keys) > MAX_DEGREE {
						if sibling, idx, err := n.parent.preMerge(); err == nil {
							return n.parent.mergeSibling(sibling, idx, key)
						} else {
							return errors.New("see rebalancing.go")
						}
					}
				}
			}
		}
	}

	// recurse UPWARD and check invariants
	if len(n.parent.keys) < ((MAX_DEGREE - 1) / 2) {
		if sibling, idx, err := n.parent.preMerge(); err == nil {
			return n.parent.mergeSibling(sibling, idx, key)
		} else {
			return errors.New("see rebalancing.go")
		}
	}
	return nil
}

// merging can be... very interesting.
// you can slap on an iter api like(rust):
// https://github.com/rust-lang/rust/blob/1c19595575968ea77c7f85e97c67d44d8c0f9a68/library/alloc/src/collections/btree/merge_iter.rs#L41
// and maybe... just maybe, stream/lift that iter out to a scheduler/async runtime -- complex, magical, do not do this, but neat to know.
// NB/Warning if you want to do it anyway: You need to be careful when providing a cursor/iter api that it is re-entrant & thread safe.
// Doing this may likely lead to needing to do _unsafe memory things_ carefully consider the invariants.

// go/pebble
// iterator/cursor: https://github.com/cockroachdb/pebble/blob/c4daad9128e053e496fa7916fda8b6df57256823/internal/manifest/btree.go#L973 &&
// https://github.com/cockroachdb/pebble/blob/c4daad9128e053e496fa7916fda8b6df57256823/internal/manifest/btree.go#L891

// the actual merge operation
// https://github.com/cockroachdb/pebble/blob/c4daad9128e053e496fa7916fda8b6df57256823/internal/manifest/btree.go#L620

/*
contents should be merged. if their contents do not fit into a single node
else are redistributed - rebalancing.go.
*/
func (n *Node) mergeSibling(sibling *Node, idx, key int) error {
	if n.parent != sibling.parent {
		panic("sibling invariant not satisfied")
	}

	switch n.kind {
	case LEAF_NODE:
		sibling.data = append(sibling.data, n.data...)

		// deallocate
		for i, node := range sibling.parent.children {
			if node == n {
				n.parent.children = append(n.parent.children[:i], n.parent.children[i+1:]...)
			}
		}

		for i, k := range n.parent.keys {
			if k == key {
				n.parent.keys = cut(i, n.parent.keys)
				newSplit := len(n.data) / 2

				if len(n.data) != 0 {
					n.parent.keys = append(n.parent.keys, n.data[newSplit])
				}

				if len(n.parent.keys) < ((MAX_DEGREE - 1) / 2) {
					if sibling, idx, err := n.parent.preMerge(); err == nil {
						return n.parent.mergeSibling(sibling, idx, key)
					} else {
						return errors.New("see rebalancing.go")
					}
				}
			}
		}

	case INTERNAL_NODE:
		sibling.keys = append(sibling.keys, n.parent.keys...)
		n.parent.children = append(n.parent.children[:idx], n.parent.children[idx+1:]...)
		sibling.keys = append(sibling.keys, n.keys...)
		sibling.children = append(sibling.children, n.children...)

		// recursive case
		if len(n.parent.keys) < ((MAX_DEGREE - 1) / 2) {
			if sibling, idx, err := n.parent.preMerge(); err == nil {
				return n.parent.mergeSibling(sibling, key, idx)
			} else {
				return errors.New("see rebalancing.go")
			}
		}
	case ROOT_NODE:
		panic("todo: root cascade")
	}

	return nil
}

// preMerge if two adjacent leaf nodes have a common parent and their contents fit into a single node
func (n *Node) preMerge() (*Node, int, error) {
	switch n.kind {
	case INTERNAL_NODE:
		// no sibling pointers so we have to go up to parent
		// we check all our siblings if we can re-distribute
		for i, sibling := range n.parent.children {
			if n == sibling {
				// cannot merge with self
				continue
			} else {
				// can merge with sibling?
				if len(sibling.keys)+len(n.keys) < MAX_DEGREE {
					return sibling, i, nil

				}
			}
		}

	case LEAF_NODE:
		if n.previous != nil {
			if len(n.previous.data)+len(n.data) < MAX_DEGREE {
				n.previous.next = n.next
				return n.previous, 0, nil
			}
		}

		if n.next != nil {
			if len(n.next.data)+len(n.data) < MAX_DEGREE {
				n.next.previous = n.previous
				return n.next, 0, nil
			}
		}
	default:
		panic("operation not supported on root")
	}

	return nil, 0, errors.New("cannot merge with sibling")
}

func MergeDeleteExample(tree *BTree) {
	// delete no cascade, just updates
	tree.Delete(4)

	// delete's causes consequetive cascade/merge to root
	// tree.Delete(5)

	fmt.Println(tree.root)
	fmt.Println(tree.root.children[0])
	fmt.Println(tree.root.children[1])

}
