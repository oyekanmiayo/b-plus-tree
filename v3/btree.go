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
			btreeNodeSplit(btree, node, retNode, insertPos)

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
				btreeNodeSplit(btree, node, retNode, insertPos)
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

func btreeNodeSplit(btree *BTree, node, retNode *BTreeNode, insertPos int) {

}

func moveKeyValBetweenNodes(fNode *BTreeNode, fIdx int, sNode *BTreeNode, sIdx int) {
	sNode.keys[sIdx] = fNode.keys[fIdx]
	sNode.values[sIdx] = fNode.values[fIdx]
}
