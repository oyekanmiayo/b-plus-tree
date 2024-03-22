package main

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestNewBNode(t *testing.T) {
	node := NewBNode(BNODE_LEAF, 5)
	typeUInt16 := node.BType()
	nKeysUInt16 := node.NKeys()
	currKeysUInt16 := node.CurrKeys()

	if typeUInt16 != uint16(BNODE_LEAF) {
		t.Errorf("Node type mismatch. Expected: %d, Got: %d", uint16(BNODE_LEAF), typeUInt16)
	}

	if nKeysUInt16 != uint16(5) {
		t.Errorf("Node NKeys mismatch. Expected: %d, Got: %d", uint16(5), nKeysUInt16)
	}

	if currKeysUInt16 != uint16(0) {
		t.Errorf("Node CurrKeys mismatch. Expected: %d, Got: %d", uint16(0), currKeysUInt16)
	}

	node.SetPtr(0, uint64(10))
	node.SetPtr(2, uint64(250))

	getPtr0 := node.GetPtr(0)
	getPtr2 := node.GetPtr(2)

	if getPtr0 != uint64(10) {
		t.Errorf("Node Ptr mismatch at idx %d. Expected: %d, Got: %d.", 0, uint64(10), getPtr0)
	}

	if getPtr2 != uint64(250) {
		t.Errorf("Node Ptr mismatch at idx %d. Expected: %d, Got: %d.", 2, uint64(250), getPtr2)
	}
}

func TestBNode_SetPtr_GetPtr(t *testing.T) {
	node := NewBNode(BNODE_LEAF, 5)
	node.SetPtr(0, uint64(10))
	node.SetPtr(2, uint64(250))

	getPtr0 := node.GetPtr(0)
	getPtr2 := node.GetPtr(2)

	if getPtr0 != uint64(10) {
		t.Errorf("Node Ptr mismatch at idx %d. Expected: %d, Got: %d.", 0, uint64(10), getPtr0)
	}

	if getPtr2 != uint64(250) {
		t.Errorf("Node Ptr mismatch at idx %d. Expected: %d, Got: %d.", 2, uint64(250), getPtr2)
	}
}

func TestBNode_SetOffset_GetOffset(t *testing.T) {
	node := NewBNode(BNODE_NODE, 5)

	// Test that the first offset is always 0
	node.SetOffset(0, 0)
	offset0 := node.GetOffset(0)
	if offset0 != uint16(0) {
		t.Errorf("Node Offset mismatch at idx %d. Expected: %d, Got: %d.", 0, uint64(0), offset0)
	}

	node.SetOffset(1, 20)
	offset20 := node.GetOffset(1)
	if offset20 != uint16(20) {
		t.Errorf("Node Offset mismatch at idx %d. Expected: %d, Got: %d.", 0, uint64(20), offset20)
	}
}

func TestBNode_GetKey_GetVal(t *testing.T) {
	node := NewBNode(BNODE_LEAF, 15)

	// 12 => "Hello"
	keyToInsert := make([]byte, 2)
	binary.LittleEndian.PutUint16(keyToInsert, uint16(12))

	valueToInsert := []byte("Hello")
	fmt.Printf("Length: %v\n", len(valueToInsert))

	idxToInsert := NodeKeyLookup(node, keyToInsert)
	if idxToInsert != uint16(0) {
		t.Errorf("Wrong index returned. Expected: %d, Got: %d.", uint16(0), idxToInsert)
	}

	InsertKVManually(node, idxToInsert, keyToInsert, valueToInsert)

}
