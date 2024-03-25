package main

import (
	"fmt"
)

// remember a page is a contiguous block of memory [..|..|..|] that's say 4KiB.
// we can arbitrarily decide what's inside this block, remember our header?

type page struct {
	id     int
	header []byte
	cells  []*cell
}

type cell struct {
	page *page
}

// recall our previous example, let's now use indirect pointers
// to binary search for data inside a page's cell

/*
NOTE: Remember that there's a cost for every disk access. If we fetch pre-fetch an entire block,
then searching through the collection of k/v records in that block will take linear access.

This is doubly expensive because each access incurs a seek time, so expensive in two ways:
1. O(n) complexity
1. Incurred read seek of the disk head (assumption hdd.)

caveat: there are reasons and cases where sequential scans make sense over using logarithmic access.
*/

func (c *cell) search(offset int) *cell {
	/*
		here's how bolt search's the leaf, it's _almost_ the same thing.
			inodes := p.leafPageElements()
			index := sort.Search(int(p.count), func(i int) bool {
				return bytes.Compare(inodes[i].key(), key) != -1
			})
	*/

	/*
		low, high := 0, len(c.page.cells)-1
		for low <= high {
			mid := low + (high-low)/2
			if bytes.Compare(page.cells[mid], nil) {
				return n
			} else if n.keys[mid] < offset {
				low = mid + 1
			} else {
				high = mid - 1
			}
		}


		return c.page.cells[mid][low]
	*/
	return nil
}

func BinarySearchPageExample(offset int) {
	fmt.Println("unfinished")
}
