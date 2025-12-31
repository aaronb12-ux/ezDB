package btree

func MakeTree(order int) *BTree {
	if order < 3 {
		panic("B+ tree order must be at least 3")
	}
	return &BTree{
		Root:   nil,
		Order:  order,
		Size:   0,
		Height: 0,
	}
}

