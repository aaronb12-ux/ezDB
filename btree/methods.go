package btree

func (tree *BTree) Insert(key int, value string) {
	//inserting a key value pair into the tree
	//check if we only have the parent

	if (tree.Size <  tree.Order) { //if the tree size is less than the order -> means the root functions as a leaf
			
		if (tree.Root == nil) { //tree is empty -> adding first node

				newKey := Key{Key: key, Value: []byte(value)} 

				var newLeafNode LeafNode

				newLeafNode.Data = append(newLeafNode.Data, newKey)
				newLeafNode.Parent = nil
				newLeafNode.Right = nil

				tree.Root = newLeafNode
				
		}
	}
}

func (tree *BTree) Delete() {
	//deleting a given key

}

func (tree *BTree) Find() {
	//find a given key
}

func (tree *BTree) PrintTree() {
	//print the tree's contents
}

func MakeTree(order int) *BTree {

	return &BTree{
		Order: order,
		Size: 0,
	}
}