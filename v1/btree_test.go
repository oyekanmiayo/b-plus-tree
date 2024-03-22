package main

import (
	"bytes"
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

	var testCases = []struct {
		indices     []uint16
		offsets     []uint16
		description string
	}{
		{
			[]uint16{uint16(0), uint16(1)},
			[]uint16{uint16(0), uint16(20)},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			node := NewBNode(BNODE_NODE, 5)

			offSize := len(tc.offsets)
			indexSize := len(tc.indices)

			if offSize != indexSize {
				t.Errorf("Offset and Index mismatch. Offset Size: %d, Index Size: %d.", offSize, indexSize)
			}

			for i := 0; i < offSize; i++ {
				node.SetOffset(tc.indices[i], tc.offsets[i])
				currOffset := node.GetOffset(tc.indices[i])
				if currOffset != tc.offsets[i] {
					t.Errorf("Node Offset mismatch at idx %d. Expected: %d, Got: %d.", 0, tc.offsets[i], currOffset)
				}

			}
		})
	}

}

func TestBNode_GetKey_GetVal(t *testing.T) {
	node := NewBNode(BNODE_LEAF, 15)

	// 12 => "Hello"
	keyToInsert := make([]byte, 2)
	binary.LittleEndian.PutUint16(keyToInsert, uint16(15))

	valueToInsert := []byte("Hello")
	fmt.Printf("Length: %v\n", len(valueToInsert))

	idxToInsert := NodeKeyLookup(node, keyToInsert)
	if idxToInsert != uint16(0) {
		t.Errorf("Wrong index returned. Expected: %d, Got: %d.", uint16(0), idxToInsert)
	}

	InsertKVManually(&node, idxToInsert, keyToInsert, valueToInsert)

	keyResult := node.GetKey(idxToInsert)
	if bytes.Compare(keyResult, keyToInsert) != 0 {
		t.Errorf("Wrong key returned. Expected: %d, Got: %d.",
			keyToInsert,
			keyResult)
	}

	valResult := node.GetVal(idxToInsert)
	if bytes.Compare(valResult, valueToInsert) != 0 {
		t.Errorf("Wrong value returned. Expected: %d, Got: %d.",
			valueToInsert,
			valResult)
	}

}

func TestMoveRangeBtwNodes(t *testing.T) {
	testCases := []struct {
		newNode BNode
		oldNode BNode
		newIdx  uint16
		oldIdx  uint16
		size    uint16
		desc    string
	}{
		{
			newNode: NewBNode(BNODE_LEAF, 3),
			oldNode: func() BNode {
				n := NewBNode(BNODE_LEAF, 6)
				InsertKVManually(&n, uint16(3), []byte{12, 0}, []byte("Hello"))
				InsertKVManually(&n, uint16(4), []byte{15, 0}, []byte("Ciao"))
				InsertKVManually(&n, uint16(5), []byte{30, 0}, []byte("E le"))
				return n
			}(),
			newIdx: 0,
			oldIdx: 3,
			size:   3,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {

			MoveRangeBtwNodes(&tc.newNode, &tc.oldNode, tc.newIdx, tc.oldIdx, tc.size)
			for i := uint16(0); i < tc.size; i++ {
				newIdx := tc.newIdx + i
				oldIdx := tc.oldIdx + i

				newKey := tc.newNode.GetKey(newIdx)
				oldKey := tc.oldNode.GetKey(oldIdx)
				if bytes.Compare(newKey, oldKey) != 0 {
					t.Errorf("Key mismatch. oldKey @ idx %d: %d, newKey @ idx %d: %d.",
						oldIdx, oldKey,
						newIdx, newKey)
				}
			}

		})
	}
}

// TestInsertKVLeaf
// Try to insert at an illegal index and handle that!
func TestInsertKVLeaf(t *testing.T) {
	testCases := []struct {
		node       BNode
		mirrorNode BNode
		idx        uint16
		key        []byte
		val        []byte
		desc       string
	}{
		{
			node: func() BNode {
				n := NewBNode(BNODE_LEAF, 6)
				InsertKVManually(&n, uint16(0), []byte{12, 0}, []byte("Hello"))
				InsertKVManually(&n, uint16(1), []byte{15, 0}, []byte("Ciao"))
				InsertKVManually(&n, uint16(2), []byte{30, 0}, []byte("E le"))
				return n
			}(),
			mirrorNode: NewBNode(BNODE_LEAF, 6),
			idx:        0,
			key:        []byte{5, 0},
			val:        []byte("Replace"),
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			InsertKVLeaf(&tc.node, &tc.mirrorNode, tc.idx, tc.key, tc.val)

			mKey := tc.mirrorNode.GetKey(tc.idx)
			if bytes.Compare(mKey, tc.key) != 0 {
				t.Errorf("Key mismatch. Expected: %d, Got: %d.", tc.key, mKey)
			}

			// Check keys before idx
			for i := uint16(0); i < tc.idx; i++ {
				nodeKey := tc.node.GetKey(i)
				mKey = tc.mirrorNode.GetKey(i)
				if bytes.Compare(mKey, tc.key) != 0 {
					t.Errorf("Key mismatch at idx %d. Expected: %d, Got: %d.", i, nodeKey, mKey)
				}
			}

			for i := tc.idx; i < tc.node.NKeys(); i++ {
				nodeKey := tc.node.GetKey(i)
				mKey = tc.mirrorNode.GetKey(tc.idx + 1)
				if bytes.Compare(mKey, tc.key) != 0 {
					t.Errorf("Key mismatch at idx %d. Expected: %d, Got: %d.", i, nodeKey, mKey)
				}
			}
		})
	}
}
