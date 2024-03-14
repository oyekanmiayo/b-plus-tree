package main

import "fmt"

type BNode struct {
	data []byte
}

const (
	BNODE_NODE = 1 // internal node (no values)
	BNODE_LEAF = 2 // leaf node (has values)
)

func main() {
	fmt.Println("Hello World")
}
