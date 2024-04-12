package main

import (
	"fmt"
	"slices"
	"testing"
)

func FuzzInsertKeys(f *testing.F) {
	var tree BTree

	for key := 1; key < 100; key++ {
		f.Add(key)
	}

	f.Fuzz(func(t *testing.T, key int) {
		tree.Insert(key)
		found := keyExists(&tree, key)

		if !found {
			t.Errorf("not found %v", key)
		}
	})
}

/*
func FuzzSearchKeys(f *testing.F) {
	var tree BTree

	for key := 1; key < 100; key++ {
		f.Add(key)
	}

	f.Fuzz(func(t *testing.T, key int) {
		tree.Insert(key)
		n, idx, err := tree.Search(key)

		if err != nil || key != n[idx] {
			t.Errorf("did not find key inserted")
		}
	})
}
*/

/*
func FuzzDeleteKeys(f *testing.F) {
	var tree BTree

	for key := 1; key < 100; key++ {
		f.Add(key)
	}

	f.Fuzz(func(t *testing.T, key int) {
		tree.Insert(key)
		err := tree.Delete(key)

		if err != nil {
			t.Errorf("deletion errored %v", err)
		}

		v, _, found := tree.Search(key)

		if found == nil {
			t.Errorf("found deleted key/value %v", v)
		}

	})
}
*/

func keyExists(t *BTree, key int) bool {
	n, _, err := t.root.Search(key)

	fmt.Println(n)
	if err != nil {
		return false
	}

	_, found := slices.BinarySearch(n.data, key)

	return found
}
