package main

// things that are difficult: fragmentation & relocating/maintaining pages/pointers.
// traditionally, in the array vs linked list debate, cache locality/chorence is a major pro
// of the array over the linked list -- and this is still somewhat true, however consider this:
// what happens when you often have multiple threads having access to different elements?
// what if you can 'pin' the nodes of this linked list in some hot layer depending on your access patterns?

/*A freelist is just a linked list with extra steps. */
// This becomes more important in the context of concurrency + fragmentation.
// -- not implemented because -- scope.
// see: https://en.wikipedia.org/wiki/Free_list for an explaination in the context of memory allocators.

// --snipped from bolt the underlying B+tree/kv storage engine for etcd, and alot of other stuff.
// freelist represents a list of all pages that are available for allocation.
// It also tracks pages that have been freed but are still in use by open transactions.
/*
type freelist struct {
	ids     []pgid          // all free and available free page ids.
	pending map[txid][]pgid // mapping of soon-to-be free page ids by tx.
	cache   map[pgid]bool   // fast lookup of all free and pending page ids.
}
*/

// --snipped from google/bree the underly B+tree for indexes in etcd, alot more..
// FreeList represents a free list of btree nodes. By default each
// BTree has its own FreeList, but multiple BTrees can share the same
// FreeList.
// Two Btrees using the same freelist are safe for concurrent write access.
/*
type FreeList struct {
	mu       sync.Mutex
	freelist []*node
}

*/
