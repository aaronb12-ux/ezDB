package btree

func (tree *BTree) Insert() {

	//inserting a key value pair into the tree
	//check if we only have the parent
}

func (tree *BTree) Delete() {

}

func (tree *BTree) Find() {

}

func MakeTree(order int) *BTree {

	return &BTree{
		Order: order,
	}
}