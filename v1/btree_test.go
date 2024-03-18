package main

import (
	"testing"
)

func TestNewBNode(t *testing.T) {
	node := NewBNode(BNODE_LEAF, 5)
	typeUInt16 := node.bType()
	nKeysUInt16 := node.nKeys()

	if typeUInt16 != uint16(BNODE_LEAF) {
		t.Errorf("Node type mismatch. Expected: %d, Got: %d", uint16(BNODE_LEAF), typeUInt16)
	}

	if nKeysUInt16 != uint16(5) {
		t.Errorf("Node nKeys mismatch. Expected: %d, Got: %d", uint16(5), nKeysUInt16)
	}
}
