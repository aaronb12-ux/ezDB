package btree

type BTree struct {
	Root *InternalNode
	Order int //max number of the children the INTERNAL NODE can have. and defines the number of keys an internal node can have which is Order - 1
}

type LeafNode struct {
	Data []Key
	Right *LeafNode
	Parent *InternalNode
}

type InternalNode struct {
    Key int
}

type Key struct {
	Key int
	Value []byte
}

//basically how do we know if 