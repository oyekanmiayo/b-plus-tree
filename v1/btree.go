package main

import "encoding/binary"

/*
BNode

# Assuming we need only 2 bytes to store integers

BNode consists of:
| type | nkeys |  pointers  |   offsets  | key-values
|  2B  |   2B  | nkeys * 8B | nkeys * 2B | ...

This is the format of the KV pair. Lengths followed by data.
| klen | vlen | key | val |
|  2B  |  2B  | ... | ... |

We assume a BNode is a page and the page size is 4KB (4096 bytes)
*/
type BNode struct {
	data []byte
}

const (
	BNODE_NODE = 1 // internal node (no values)
	BNODE_LEAF = 2 // leaf node (has values)
)

type BTree struct {
	// unsigned: 0 -> INT_MAX
	root uint64

	get func(uint64 uint64) BNode
	new func(node BNode) // allocate a new node (page)
	del func(uint642 uint64)
}

func NewBNode(bType, nKeys int) BNode {
	node := BNode{
		data: make([]byte, 4),
	}

	// node type
	binary.LittleEndian.PutUint16(node.data[0:2], uint16(bType))

	// number of keys
	binary.LittleEndian.PutUint16(node.data[2:4], uint16(nKeys))

	return node
}

func (node BNode) bType() uint16 {
	// This is trying to interpret the first 2 bytes of node.data to get the node type
	return binary.LittleEndian.Uint16(node.data[0:2])
}

func (node BNode) nKeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}
