package main

import (
	"bytes"
	"encoding/binary"
)

/*
BNode

We assume a BNode is a page and the page size is 4KB (4096 bytes)
*/
type BNode struct {
	data []byte
}

/*

- Store sizes of each part of the node (page)
- unit is bytes
- this is not to be confused with their actual values.
- this simply stores the size within with the values must be contained

BNode consists of:
| type | nkeys |  pointers  |   offsets  | key-values
|  2B  |   2B  | nkeys * 8B | nkeys * 2B | ...

This is the format of the KV pair. Lengths followed by data.
| klen | vlen | key | val |
|  2B  |  2B  | ... | ... |

- type: The type of the node. BNODE_NODE or BNODE_LEAF.
- nkeys: The number of keys this node can store.
- pointers: pointers to children nodes. Used only in BNODE_NODE.
- offsets: for a given idx, it tells where to read the key and val for that index from
	- only used in BNODE_LEAF
- key-values: variable-length keys and values

*/

const (
	BNODE_TYPE = 2
	// BNODE_NKEYS total number of keys a node can store
	BNODE_NKEYS = 2
	// BNODE_CURR_KEYS number of keys this node has store now
	BNODE_CURR_KEYS
	BNODE_HEADER = BNODE_TYPE + BNODE_NKEYS + BNODE_CURR_KEYS

	// BNODE_POINTER_SIZE and BNODE_OFFSET_SIZE represent the size of one pointer
	// BNODE_HEADER is different because there's only one type and nkeys
	// so the size of one of each covers the total size
	BNODE_POINTER_SIZE = 8
	BNODE_OFFSET_SIZE  = 2

	MAX_BTREE_PAGE_SIZE = 4096
)

const (
	BNODE_NODE = iota // internal node (no values)
	BNODE_LEAF        // leaf node (has values)
)

type BTree struct {
	// unsigned: 0 -> INT_MAX
	root uint64

	get func(uint64 uint64) BNode
	new func(node BNode) // allocate a new node (page)
	del func(uint642 uint64)
}

// NewBNode
// BType - Node type is BNODE_NODE or BNODE_LEAF
// NKeys - number of keys this node can store
func NewBNode(bType, nKeys int) BNode {
	nodeInitialSize := BNODE_HEADER + (BNODE_POINTER_SIZE * nKeys) + (BNODE_OFFSET_SIZE * nKeys)
	// MAX_BTREE_PAGE_SIZE is the maximum size a page (node) should be
	// we should verify that the number of keys to be inserted leaves space for at least one key value pair
	if nodeInitialSize >= MAX_BTREE_PAGE_SIZE {
		panic("This node can't store any more data.")
	}

	node := BNode{
		data: make([]byte, MAX_BTREE_PAGE_SIZE),
	}

	// node type
	binary.LittleEndian.PutUint16(node.data[0:2], uint16(bType))

	// number of keys
	binary.LittleEndian.PutUint16(node.data[2:4], uint16(nKeys))

	return node
}

func (node BNode) BType() uint16 {
	return binary.LittleEndian.Uint16(node.data[0:2])
}

func (node BNode) NKeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node BNode) CurrKeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[4:6])
}

func (node BNode) SetPtr(idx uint16, val uint64) {
	pos := BNODE_HEADER + (idx * BNODE_POINTER_SIZE)
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}

func (node BNode) GetPtr(idx uint16) uint64 {
	if idx >= node.NKeys() {
		panic("idx >= NKeys")
	}

	// this gets us to the starting point of the bytes for this pointer
	pos := BNODE_HEADER + (idx * BNODE_POINTER_SIZE)

	// read 8 bytes from starting point (pos)
	return binary.LittleEndian.Uint64(node.data[pos:])
}

// offsetPos(node, idx) will fetch the starting point of the offset position for a particular index
func offsetPos(node BNode, idx uint16) uint16 {
	return BNODE_HEADER + (BNODE_POINTER_SIZE * node.NKeys()) + (BNODE_OFFSET_SIZE * idx)
}

func (node BNode) SetOffset(idx, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[offsetPos(node, idx):], offset)
}

func (node BNode) GetOffset(idx uint16) uint16 {
	return binary.LittleEndian.Uint16(node.data[offsetPos(node, idx):])
}

// 16 bits = 2 bytes
// 1 byte = 8 bits
func (node BNode) kvPos(idx uint16) uint16 {
	return BNODE_HEADER + (BNODE_POINTER_SIZE * node.NKeys()) + (BNODE_OFFSET_SIZE * node.NKeys()) + node.GetOffset(idx)
}

func (node BNode) GetKey(idx uint16) []byte {
	pos := node.kvPos(idx)
	kLen := binary.LittleEndian.Uint16(node.data[pos:])

	// + 4 skips the bytes for kLen and vLen
	// [:kLen] reads from pos+4 up to kLen (non-inclusive)
	return node.data[pos+4:][:kLen]
}

func (node BNode) GetVal(idx uint16) []byte {
	pos := node.kvPos(idx)
	kLen := binary.LittleEndian.Uint16(node.data[pos:])
	vLen := binary.LittleEndian.Uint16(node.data[pos+2:])

	return node.data[pos+4+kLen:][:vLen]
}

func InsertKVManually(node BNode, idx uint16, key, val []byte) {
	pos := node.kvPos(idx)

	// Add kLen and vLen
	binary.LittleEndian.PutUint16(node.data[pos:], uint16(len(key)))
	binary.LittleEndian.PutUint16(node.data[pos+2:], uint16(len(val)))

	// Add key and val
	copy(node.data[pos+4:], key)
	copy(node.data[pos+4+uint16(len(key)):], val)

	// set new offset for next idx
	// 4 for the kLen and vLen
	node.SetOffset(idx+1, node.GetOffset(idx)+4+uint16(len(key)+len(val)))
}

// NodeKeyLookup search for key
// if key is present in this btree's range, keep searching until we reach the leaf node that contains that key
// if key exists, update value in-place
// if key doesn't exist, insert value at appropriate location
// if on insert, the page is too big, split into two nodes and upshit the mid key
func NodeKeyLookup(node BNode, key []byte) uint16 {

	// Get number of keys in the node
	nKeys := node.NKeys()

	// not worried about setting another initial value for found
	// because they MUST BE a valid node. The given key has to fall within the range somehow.
	found := uint16(0)

	for i := uint16(0); i < nKeys; i++ {

		nodeCurrKey := node.GetKey(i)

		// if the length of this is 0, then it means kLen was 0
		// which means there's no key at this index
		// particularly useful when we want to insert and there's no key available!
		if len(nodeCurrKey) == 0 {
			found = i
			break
		}

		// a = node's key at current i
		// b = given key
		cmp := bytes.Compare(nodeCurrKey, key)

		// if the current node's key < key_to_insert,
		// then insertion point is one index ahead
		// if the current node's key == key_to_insert
		// then we just need to update the value associated with the key
		// we pass the index anyway. Whether to insert or update will be decided by the caller.
		if cmp <= 0 {
			found = i
		}

		// if current node's key > key_to_insert
		// there are two possibilities here
		// - if this is the first key, then it means the key we need to insert will be
		//   inserted at 0. This has been handled by found := uint16(0)
		// - if this isn't the first key in the node, then we should have already found an
		//   insertion/update point.
		if cmp >= 0 {
			break
		}
	}

	return found
}

// InsertKVLeaf insert in KV assuming there's no split necessary initially
// it's a lot easier to just use a new node!
func InsertKVLeaf(node BNode, idx uint16) {

	// we don't add a key yet since we assume this insert won't cause an overflow
	//
	// mirrorNode := NewBNode(BNODE_LEAF, int(node.NKeys()))

	// copy from 0 to idx
	// insert KV
	// copy from idx+1 to end
}

func MoveRangeBtwNodes(newNode, oldNode BNode, newIdx, oldIdx, size uint16) {
	// There's nothing to do here
	if size == 0 {
		return
	}

	// move pointers
	for i := uint16(0); i < size; i++ {
		newNode.SetPtr(newIdx+i, oldNode.GetPtr(oldIdx+i))
	}

	// move offsets
	newOffsetBegin := newNode.GetOffset(newIdx)
	oldOffsetBegin := oldNode.GetOffset(oldIdx)
	for i := uint16(0); i < size; i++ {
		// size of offsets up to the current index from oldOffsetBegin
		sizeOfKVUpToCurrIdx := oldNode.GetOffset(oldIdx+i) - oldOffsetBegin

		// This was confusing initially, but this is precisely what we want
		// at each idx, we want sizeOfKVUpToCurrIdx to increase because the offset of the current
		// idx must be further that the past idx
		offset := newOffsetBegin + sizeOfKVUpToCurrIdx
		newNode.SetOffset(newIdx+i, offset)
	}

	// move KV
	oldKVStart := oldNode.kvPos(oldIdx)
	oldKVEnd := oldNode.kvPos(oldIdx + size)
	copy(newNode.data[newNode.kvPos(newIdx):], oldNode.data[oldKVStart:oldKVEnd])

}
