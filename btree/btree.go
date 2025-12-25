package btree

type BTree struct {
	Root *Node
	Order int //max number of the children the INTERNAL NODE can have. and defines the number of keys an internal node can have which is Order - 1
	Size int //number of keys
}

type Key struct {
	Key int
	Value []byte
}

type Node struct {
	Data *[]Key //data is an array of Keys (only for leaf nodes)
	Children *[]int //for internal nodes. they have children
	Next *int //only for leaf nodes
	Keys *[]int //for internal nodes. internal nodes have an array of keys which are used for the searching
}

