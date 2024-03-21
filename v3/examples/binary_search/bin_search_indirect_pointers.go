package main

import "fmt"

// TODO/WIP/unfinished
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

func (c *cell) search(offset int) *cell {
	low, high := 0, len(c.page.cells)-1
	for low <= high {
		mid := low + (high-low)/2
		if c.page.cells[mid] == offset {
			return n
		} else if n.keys[mid] < offset {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return c.page.cells[mid][low]
}

func BinarySearchPageExample(offset int) {
	fmt.Println("unfinished")
}
