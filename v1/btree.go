package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	BNODE_HEADER = BNODE_TYPE + BNODE_NKEYS

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
	// new func(node BNode) uint64 // allocate a new node (page)
	new func(parentNode BNode, currNode BNode, currIdx uint16) uint64
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

	fmt.Printf("node.data's address %p\n", node.data)

	// node type
	binary.LittleEndian.PutUint16(node.data[0:2], uint16(bType))

	// number of keys
	binary.LittleEndian.PutUint16(node.data[2:4], uint16(nKeys))

	return node
}

func BTypeStr(i uint16) string {
	switch i {
	case BNODE_LEAF:
		return "BNODE_LEAF"
	case BNODE_NODE:
		return "BNODE_NODE"
	}

	return ""
}

func (node BNode) debug() {
	fmt.Println("-----NODE DETAILS-----")
	// print headers
	fmt.Printf("Node Type: %s", BTypeStr(node.BType()))
	// print pointer size
	fmt.Printf("Idv pointer size: %d\n", BNODE_POINTER_SIZE)
	fmt.Println()
	fmt.Println("-----Pointer Details-----")
	// print pointers for each idx
	for i := uint16(0); i < node.NKeys(); i++ {
		fmt.Printf("Pointer[%d] = %d\n", i, node.GetPtr(i))
	}
	fmt.Println("-----End of Pointer Details-----")
	fmt.Println()
	fmt.Println("-----Offset Details-----")
	// print offset size
	// print offset for each idx and the kv it refers to
	fmt.Printf("Idv offset size: %d\n", BNODE_OFFSET_SIZE)
	fmt.Println()
	for i := uint16(0); i < node.NKeys(); i++ {
		fmt.Printf("Offset[%d] = %d\n", i, node.GetOffset(i))
		fmt.Printf("  Key = %d, Key Addr = %p\n", node.GetKey(i), node.GetKey(i))
		fmt.Printf("  Value = %s, Val Addr = %p\n", string(node.GetVal(i)), node.GetVal(i))
		fmt.Println()
	}
	fmt.Println("-----End of Offset Details-----")
	fmt.Println("-----NODE DETAILS-----")
	fmt.Println()
	fmt.Println()

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
		msg := fmt.Sprintf("idx %d >= NKeys %d\n", idx, node.NKeys())
		panic(msg)
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
	if idx >= node.NKeys() {
		return
	}
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
	key := node.data[pos+4:][:kLen]
	// fmt.Printf("idx %d, pos %d, kLen %d, key %d, address %p\n", idx, pos, kLen, key, key)

	// + 4 skips the bytes for kLen and vLen
	// [:kLen] reads from pos+4 up to kLen (non-inclusive)
	return key
}

func (node BNode) GetVal(idx uint16) []byte {
	pos := node.kvPos(idx)
	kLen := binary.LittleEndian.Uint16(node.data[pos:])
	vLen := binary.LittleEndian.Uint16(node.data[pos+2:])

	return node.data[pos+4+kLen:][:vLen]
}

func InsertKVManually(node *BNode, idx uint16, key, val []byte) {
	if idx >= node.NKeys() {
		msg := fmt.Sprintf("InsertKVManually is trying to insert in an illegal node. Make the node bigger?")
		panic(msg)
	}

	pos := node.kvPos(idx)
	if idx > 0 {
		prevPos := node.kvPos(idx - 1)
		if prevPos == pos {
			msg := fmt.Sprintf("Insert at prev idx %d\n", idx-1)
			panic(msg)
		}
	}

	// Add kLen and vLen
	kLen := uint16(len(key))
	vLen := uint16(len(val))
	binary.LittleEndian.PutUint16(node.data[pos:pos+2], kLen)
	binary.LittleEndian.PutUint16(node.data[pos+2:pos+4], vLen)

	// Add key and val
	copy(node.data[pos+4:][:kLen], key)
	copy(node.data[pos+4+kLen:][:vLen], val)

	// set new offset for next idx
	// 4 for the kLen and vLen
	node.SetOffset(idx+1, node.GetOffset(idx)+4+kLen+vLen)
	node.debug()
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
		// ---
		// This is problematic if we ever call InsertKVManually apart from this method
		// Imagine this scenario for a newNode BNode
		// InsertKVManually(n, uint16(3), []byte{6, 0}, []byte("Hello")) --> We insert "Hello" at key 6
		// Then we try to insert 9 --> "Hi" using nodeInsert(). Even though 9 > 6,
		// the idx returned will be 0 in this case!
		// simple solution is to only call InsertKVManually from within other methods!
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
func InsertKVLeaf(node, mirrorNode *BNode, idx uint16, key, val []byte) {

	// we don't add a key yet since we assume this insert won't cause an overflow
	//
	// mirrorNode := NewBNode(BNODE_LEAF, int(node.NKeys()))

	// copy from 0 to idx
	MoveRangeBtwNodes(mirrorNode, node, 0, 0, idx)
	// insert KV
	InsertKVManually(mirrorNode, idx, key, val)
	mirrorNode.debug()
	// copy from idx to end
	MoveRangeBtwNodes(mirrorNode, node, idx+1, idx, node.NKeys()-idx)
	mirrorNode.debug()
}

// UpdateKVLeaf the key already exists, so update the value
// Not as simple as just update the value in-place.
// The new value might be bigger or smaller. So, vLen may be different which also changes offset and so on.
func UpdateKVLeaf(node, mirrorNode *BNode, idx uint16, key, val []byte) {
	// copy from 0 to idx-1
	MoveRangeBtwNodes(mirrorNode, node, 0, 0, idx)
	// insert KV (K@idx)
	InsertKVManually(mirrorNode, idx, key, val)
	// copy from idx+1 to end
	MoveRangeBtwNodes(mirrorNode, node, idx+1, idx+1, node.NKeys()-(idx+1))
}

func MoveRangeBtwNodes(newNode, oldNode *BNode, newIdx, oldIdx, size uint16) {
	// There's nothing to do here
	if size == 0 {
		return
	}

	for i := uint16(0); i < size; i++ {
		fmt.Println("k -> GetKey")
		k := oldNode.GetKey(oldIdx + i)
		v := oldNode.GetVal(oldIdx + i)

		InsertKVManually(newNode, newIdx+i, k, v)
		newNode.debug()
		fmt.Println()
	}

	// move pointers, if internal node
	for i := uint16(0); i < size; i++ {
		newNode.SetPtr(newIdx+i, oldNode.GetPtr(oldIdx+i))
	}
}

// Store unique 8 byte integer -> data []byte within nodes
var ptrMap = make(map[uint64][]byte)

func (tree *BTree) insert(node BNode, key, val []byte) BNode {

	idxToInsertOrUpdate := NodeKeyLookup(node, key)

	// assuming no overflow
	mirrorNode := NewBNode(0, int(node.NKeys()))

	switch node.BType() {
	case BNODE_LEAF:
		if bytes.Equal(key, node.GetKey(idxToInsertOrUpdate)) {
			// the value
			UpdateKVLeaf(&node, &mirrorNode, idxToInsertOrUpdate, key, val)
		} else {
			InsertKVLeaf(&node, &mirrorNode, idxToInsertOrUpdate, key, val)
		}
	case BNODE_NODE:
	default:
		panic("illegal node type")
	}

	return mirrorNode
}

/*
TODO

- implement tree insert and split
- implement tree deletion and merge

- if i'm moving ranges and the kLen at a particular index is 0, i want to ignore it!
- store current Keys to make life easier!

fin
*/
