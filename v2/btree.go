package main

import "math"

const (
	BTREE_MIN_DEG = 2
)

type BTreeNode struct {
	keysOccupied         int
	noOfExistingChildren int
	keys                 []int
	values               [][]byte // this is for the keys, but I think they are only useful in leaf nodes?
	children             []*BTreeNode
}

type BTree struct {
	degree    int // degree of each node in the btree?
	minDegree int
	root      *BTreeNode
}

func newBTree(degree int) *BTree {
	if degree < BTREE_MIN_DEG {
		degree = BTREE_MIN_DEG
	}

	return &BTree{
		degree:    degree,
		minDegree: int(math.Ceil(float64(degree)/float64(2)) - 1),
		root:      nil,
	}
}

func (btree *BTree) newBTreeNode() *BTreeNode {
	return &BTreeNode{
		keysOccupied:         0,
		noOfExistingChildren: 0,
		keys:                 make([]int, btree.degree-1),
		values:               make([][]byte, btree.degree-1),
		children:             make([]*BTreeNode, btree.degree),
	}
}

func (btree *BTree) newBTreeNodeItem(key int, val []byte) *BTreeNode {
	btreeNode := btree.newBTreeNode()
	btreeNode.keys[0] = key
	btreeNode.values[0] = val
	btreeNode.keysOccupied = 1
	return btreeNode
}

func (btree *BTree) find(key int) []byte {

	idx, node := nodeSearch(btree.root, key)
	if node == nil {
		// there was no match
	}

	return node.values[idx]
}

func nodeSearch(node *BTreeNode, key int) (int, *BTreeNode) {

	// key >= node.keys[idx] is true, then we don't need to keep checking future keys
	idx := 0
	for ; idx < node.keysOccupied && key >= node.keys[idx]; idx++ {
		if key == node.keys[idx] {
			return idx, node
		}
	}

	// Base case
	if len(node.children) == 0 {
		return -1, nil
	}

	return nodeSearch(node.children[idx], key)
}

func (btree *BTree) insert(key int, value []byte) {
	if btree.root == nil {
		btree.root = btree.newBTreeNodeItem(key, value)
		return
	}

	resultNode := btree.nodeInsert(btree.root, key, value)
	if resultNode != nil {
		btree.root = resultNode
	}
}

// Only return a BTreeNode is there's a value to be propagated up towards the root (i.e. when
// there's a split)
func (btree *BTree) nodeInsert(node *BTreeNode, key int, value []byte) *BTreeNode {
	// find position
	insertPos := 0
	for ; insertPos < node.keysOccupied && key > node.keys[insertPos]; insertPos++ {
	}

	// if key exists, update the value
	if insertPos < node.keysOccupied && insertPos > 0 && key == node.keys[insertPos] {
		node.values[insertPos] = value
		return nil
	}

	retNode := new(BTreeNode)

	// if key doesn't exist in current node, insert into child

	if node.noOfExistingChildren == 0 { // Leaf Node

		// This means any new inserts will cause an overflow (and trigger a split)
		if node.keysOccupied == btree.degree-1 {
			retNode = btree.newBTreeNodeItem(key, value)
			btree.nodeSplit(node, retNode, insertPos)

		} else {
			// we can insert in current node without a violation

			for k := insertPos; k <= node.keysOccupied; k++ {
				moveKeyValBetweenNodes(node, k, node, k+1)
			}

			node.keys[insertPos] = key
			node.values[insertPos] = value
			node.noOfExistingChildren++
		}

	} else {
		// insert into child

		retNode = btree.nodeInsert(node.children[insertPos], key, value)

		if retNode != nil {
			if node.keysOccupied == btree.degree-1 { // We will have an overflow, so split now
				btree.nodeSplit(node, retNode, insertPos)
			} else {
				// This always assumes that keysOccupied < degree :)
				for k := insertPos; k <= node.keysOccupied; k++ {
					moveKeyValBetweenNodes(node, k, node, k+1)
				}

				moveKeyValBetweenNodes(retNode, 0, node, insertPos)

				// right shift children
				for k := insertPos; k <= node.keysOccupied; k++ {
					node.children[k+1] = node.children[k]
				}

				node.children[insertPos+1] = retNode.children[1]
				node.children[insertPos] = retNode.children[0]
				node.keysOccupied++
				node.noOfExistingChildren++

				return nil
			}
		}
	}

	return retNode
}

func (btree *BTree) nodeSplit(node, retNode *BTreeNode, insertPos int) {
	hasChildren := node.noOfExistingChildren != 0
	tempNode := btree.newBTreeNode()
	tempNode.children = append(tempNode.children, retNode.children...)

	// somehow, minDegree also doubles as the average?
	treeAvg := btree.minDegree

	if insertPos < btree.minDegree-1 {

		// save KV pair for the media key that will be up-shifted to tempNode
		moveKeyValBetweenNodes(node, treeAvg-2, tempNode, 0)

		// move all KV pairs in the overflowing node to the right by one up to the t-2 index
		for idx := treeAvg - 2; idx > insertPos; idx-- {
			moveKeyValBetweenNodes(node, idx-1, node, idx)
		}

		// there's now free space at insertPos, so insert the new KV pair there
		moveKeyValBetweenNodes(retNode, 0, node, insertPos)

	} else if insertPos > btree.minDegree-1 {

		// save KV pair for the media key that will be up-shifted to tempNode
		moveKeyValBetweenNodes(node, treeAvg-1, tempNode, 0)

		for idx := treeAvg - 1; idx < insertPos && idx < node.keysOccupied-1; idx++ {
			// node[idx] = node[idx+1]
			moveKeyValBetweenNodes(node, idx+1, node, idx)
		}

		// there's now free space at insertPos, so insert the new KV pair there
		// insertPos - 1 because the position would be left shifted after removing the middle KV pair
		moveKeyValBetweenNodes(retNode, 0, node, insertPos-1)
	} else {
		moveKeyValBetweenNodes(retNode, 0, tempNode, 0)
	}

	// So now we have insert the KV at the proper position and put the KV to be upshifted in tempNode
	/**
	- What next?
	- move node to be upshifted to retNode
	- retNode's left child will be the old root
	- retNode's right (next) child will be the right handside of the old root
	*/

	moveKeyValBetweenNodes(tempNode, 0, retNode, 0)
	retNode.children[0] = node
	retNode.children[1] = btree.newBTreeNode()
	for idx := treeAvg - 1; idx < btree.degree-1; idx++ {
		// insert KV pairs from KV after upshifted KV to the end in the retNode.children[1].keys[0 -> treeAvg - 1]
		// idx-btree.degree+1 is the calculation that helps us insert in the right index
		moveKeyValBetweenNodes(node, idx, retNode.children[1], idx-treeAvg+1)
		node.keys[idx] = math.MaxInt
		node.values[idx] = nil
	}

	if hasChildren {
		if insertPos < treeAvg-1 {
			for idx := treeAvg - 1; idx < btree.degree; idx++ {
				// insert rhs children from overflowing node to children of the new rhs node of retNode
				retNode.children[1].children[idx-treeAvg+1] = retNode.children[0].children[idx]
			}

			for idx := treeAvg - 1; idx > insertPos; idx-- {

			}

		} else { //  insertPos >= treeAvg - 1

		}
	}
}

func moveKeyValBetweenNodes(fNode *BTreeNode, fIdx int, sNode *BTreeNode, sIdx int) {
	sNode.keys[sIdx] = fNode.keys[fIdx]
	sNode.values[sIdx] = fNode.values[fIdx]
}
