package main

import (
	"testing"
)

func FuzzMillionKeys(f *testing.F) {
	var tree BPlusTree

	for key := 1; key < 1_000_000; key++ {
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
