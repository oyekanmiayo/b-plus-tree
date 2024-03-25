package main

import (
	"errors"
	"fmt"
	"slices"
)

/*
// rebalance attempts to combine the node with sibling nodes if the node fill
// size is below a threshold or if there are not enough keys.
func (n *node) rebalance() {
	if !n.unbalanced {
		return
	}
	n.unbalanced = false

	// Update statistics.
	n.bucket.tx.stats.Rebalance++

	// Ignore if node is above threshold (25%) and has enough keys.
	var threshold = n.bucket.tx.db.pageSize / 4
	if n.size() > threshold && len(n.inodes) > n.minKeys() {
		return
	}

	// Root node has special handling.
	if n.parent == nil {
		// If root node is a branch and only has one node then collapse it.
		if !n.isLeaf && len(n.inodes) == 1 {
			// Move root's child up.
			child := n.bucket.node(n.inodes[0].pgid, n)
			n.isLeaf = child.isLeaf
			n.inodes = child.inodes[:]
			n.children = child.children

			// Reparent all child nodes being moved.
			for _, inode := range n.inodes {
				if child, ok := n.bucket.nodes[inode.pgid]; ok {
					child.parent = n
				}
			}

			// Remove old child.
			child.parent = nil
			delete(n.bucket.nodes, child.pgid)
			child.free()
		}

		return
	}

	// If node has no keys then just remove it.
	if n.numChildren() == 0 {
		n.parent.del(n.key)
		n.parent.removeChild(n)
		delete(n.bucket.nodes, n.pgid)
		n.free()
		n.parent.rebalance()
		return
	}

	_assert(n.parent.numChildren() > 1, "parent must have at least 2 children")

	// Destination node is right sibling if idx == 0, otherwise left sibling.
	var target *node
	var useNextSibling = (n.parent.childIndex(n) == 0)
	if useNextSibling {
		target = n.nextSibling()
	} else {
		target = n.prevSibling()
	}

	// If both this node and the target node are too small then merge them.
	if useNextSibling {
		// Reparent all child nodes being moved.
		for _, inode := range target.inodes {
			if child, ok := n.bucket.nodes[inode.pgid]; ok {
				child.parent.removeChild(child)
				child.parent = n
				child.parent.children = append(child.parent.children, child)
			}
		}

		// Copy over inodes from target and remove target.
		n.inodes = append(n.inodes, target.inodes...)
		n.parent.del(target.key)
		n.parent.removeChild(target)
		delete(n.bucket.nodes, target.pgid)
		target.free()
	} else {
		// Reparent all child nodes being moved.
		for _, inode := range n.inodes {
			if child, ok := n.bucket.nodes[inode.pgid]; ok {
				child.parent.removeChild(child)
				child.parent = target
				child.parent.children = append(child.parent.children, child)
			}
		}

		// Copy over inodes to target and remove node.
		target.inodes = append(target.inodes, n.inodes...)
		n.parent.del(n.key)
		n.parent.removeChild(n)
		delete(n.bucket.nodes, n.pgid)
		n.free()
	}

	// Either this node or the target node was deleted from the parent so rebalance it.
	n.parent.rebalance()
}
*/

func (t *BTree) RebalanceDelete(key int) error {
	if t.root == nil {
		return errors.New("empty tree")
	} else {
		// find leaf node to delete
		n, _, err := t.root.Search(key)
		if err == nil {
			return n.delete(t, key)
		}

		return errors.New("key not in tree")
	}
}

func (n *node) delete(t *BTree, key int) error {
	if err := n.mergeSiblings(t, key); err == nil {
		return nil
	}
	return n.rebalanceDel(key)
}

func (n *node) rebalanceDel(key int) error {
	return errors.New("wip")
}

func RebalancingDelete(tree *BTree) {
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

func (n *node) mergeSiblings(t *BTree, key int) error {
	return nil
}

// see: basic search of how/why this works
func (n *node) search(key int) (*node, int, error) {
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
