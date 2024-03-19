package main

import "encoding/binary"

/*
BNode

This is the format of the KV pair. Lengths followed by data.
| klen | vlen | key | val |
|  2B  |  2B  | ... | ... |

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
*/

const (
	BNODE_TYPE   = 2
	BNODE_NKEYS  = 2
	BNODE_HEADER = BNODE_TYPE + BNODE_NKEYS

	// BNODE_POINTER_SIZE and BNODE_OFFSET_SIZE represent the size of one pointer
	// BNODE_HEADER is different because there's only one type and nkeys
	// so the size of one of each covers the total size
	BNODE_POINTER_SIZE = 8
	BNODE_OFFSET_SIZE  = 2
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
	node := BNode{
		data: make([]byte, nodeInitialSize),
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

// The offset of the first KV is always zero, so it's not stored in the list
// offsetPos(node, idx) will fetch the starting point of the offset position for a particular index
func offsetPos(node BNode, idx uint16) uint16 {
	return BNODE_HEADER + (BNODE_POINTER_SIZE * BNODE_NKEYS) + (BNODE_OFFSET_SIZE * (idx - 1))
}

func (node BNode) SetOffset(idx, offset uint16) {
	if idx == 0 && offset == 0 {
		return
	}

	if idx == 0 {
		panic("idx is 0, but offset isn't 0")
	}

	binary.LittleEndian.PutUint16(node.data[offsetPos(node, idx):], offset)
}

func (node BNode) GetOffset(idx uint16) uint16 {
	// The offset of the first KV is always zero, so it's not stored in the list
	if idx == 0 {
		return 0
	}

	return binary.LittleEndian.Uint16(node.data[offsetPos(node, idx):])
}

// 16 bits = 2 bytes
// 1 byte = 8 bits
func (node BNode) kvPos(idx uint16) uint16 {
	return BNODE_HEADER + (BNODE_POINTER_SIZE * BNODE_NKEYS) + (BNODE_OFFSET_SIZE * BNODE_NKEYS) + node.GetOffset(idx)
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
