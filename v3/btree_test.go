package main

import (
	"testing"
)

func FuzzInsertKeys(f *testing.F) {
	var tree BTree

	for key := 1; key < 100; key++ {
		f.Add(key)
	}

	f.Fuzz(func(t *testing.T, key int) {
		tree.Insert(key)
		found := KeyExists(&tree, key)

		if !found {
			f.Failed()
		}
	})
}

func FuzzSearchKeys(f *testing.F) {
	var tree BTree

	for key := 1; key < 100; key++ {
		f.Add(key)
	}

	f.Fuzz(func(t *testing.T, key int) {
		tree.Insert(key)
		n, idx, err := tree.Search(key)

		if err != nil || key != n[idx] {
			f.Failed()
		}
	})
}

func FuzzDeleteKeys(f *testing.F) {
	var tree BTree

	for key := 1; key <= 5; key++ {
		f.Add(key)
	}

	f.Fuzz(func(t *testing.T, key int) {
		tree.Insert(key)
		err := tree.Delete(key)

		if err != nil {
			f.Failed()
		}

		_, _, found := tree.Search(key)

		if found == nil {
			f.Failed()
		}

	})
}
