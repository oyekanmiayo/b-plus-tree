package main

import (
	"testing"
)

func FuzzInsertKeys(f *testing.F) {
	var tree BTree

	for key := 1; key < 10_000; key++ {
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

func FuzzSearchKeys(f *testing.F) {
	var tree BTree

	for key := 1; key < 10_000; key++ {
		f.Add(key)
	}

	f.Fuzz(func(t *testing.T, key int) {
		var found bool
		tree.Insert(key)

		data, _, err := tree.Search(key)

		if err != nil {
			t.Errorf("could not search tree %v", err)
		}

		for _, d := range data {
			if d == key {
				found = true
			}
		}

		if !found {
			t.Errorf("did not find key inserted")
		}
	})
}

func FuzzDeleteKeys(f *testing.F) {
	var tree BTree

	for key := 1; key < 10_000; key++ {
		f.Add(key)
	}

	f.Fuzz(func(t *testing.T, key int) {
		tree.Insert(key)
		err := tree.Delete(key)

		if err != nil {
			t.Errorf("deletion errored %v", err)
		}

		v, _, _ := tree.root.Search(key)

		for _, d := range v.data {
			if d == key {
				t.Errorf("found deleted key/value %v", v)
			}
		}

	})
}

func keyExists(t *BTree, key int) bool {
	n, _, _ := t.root.Search(key)

	for _, v := range n.data {
		if v == key {
			return true
		}
	}

	return false
}
