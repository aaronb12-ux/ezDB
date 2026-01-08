package btree

import (
	"sort"

)

//searching for a key in the tree (read from disk)
func (tree *BTree) Get(key int) ([]byte, bool, error) {
	
	//tree is empty case
	if tree.RootBlockNum == -1 {
		return nil, false, nil
	}

	//find the leaf node block that should contain the key
	leafBlockNumber, err := tree.findLeafBlock(key)

	if err != nil {
		return nil, false, err
	}
 
	//read that block from the disk and load it into a leaf node
	leaf, err := tree.readNode(leafBlockNumber)

	i := sort.SearchInts(leaf.Keys[:leaf.NumKeys], key) //binary search the leafs keys. returns index of where 'key' exists in leaf.Keys

	if i < leaf.NumKeys && leaf.Keys[i] == key { //if found, return the value at that key
		return leaf.Values[i], true, nil
	}

	return nil, false, nil

}

//traverses the tree to find the leaf block where a key should be
func (tree *BTree) findLeafBlock(key int) (int, error) {

	blockNum := tree.RootBlockNum

	for {
		node, err := tree.readNode(blockNum) //take data at the block number and read it into a node object
		if err != nil {
			return -1, err
		}

		if node.IsLeaf { //if we have reached a leaf, we are done
			return blockNum, nil
		}

		i := sort.SearchInts(node.Keys[:node.NumKeys], key)

		if i > node.NumKeys { //make sure we don't go out of bounds
			i = node.NumKeys
		}

		blockNum = node.ChildBlocks[i]
	}
}


func (tree *BTree) Insert(key int, value []byte) error {

	if tree.RootBlockNum == -1 { //if tree is empty...
		blockNum := tree.allocateBlock() 

		leaf := &Node{ //create first root node

			IsLeaf: true,
			NumKeys: 1,
			Keys: make([]int, tree.Order),
			Values: make([][]byte, tree.Order),
			NextBlock: -1,
			ParentBlock: -1,
		}

		leaf.Keys[0] = key //fill in values
		leaf.Values[0] = value

		err := tree.writeNode(blockNum, leaf) //write node to disk at block 'blockNum'

		if err != nil {
			return err
		}

		tree.RootBlockNum = blockNum
		tree.Size = 1
		tree.Height = 1
		return nil
	}

	leafBlockNum, err := tree.findLeafBlock(key) //get the block number of the node where the key should be at

	if err != nil {
		return err
	}

	leaf, err := tree.readNode(leafBlockNum) //take the block data from the disk and read it into a page which is then converted to a node

	if err != nil {
		return err
	}

	existed := tree.insertIntoLeaf(leaf, key, value) //insert the key value pair into the leaf

	if existed { //updating the value if it exists already and writing to the disk

		err = tree.writeNode(leafBlockNum, leaf) //write the updated leaf (node) back to the disk
		return err
	}

	tree.Size++ 

	err = tree.writeNode(leafBlockNum, leaf) //now write the inserted node to the disk

	if err != nil {
		return err
	}

	if leaf.NumKeys > tree.Order-1 { //perform a split if the number of keys in the leaf node is greater than order - 1 (b+ tree rule)
		return tree.splitLeaf(leafBlockNum, leaf)
	}

	return nil
}

func (tree *BTree) insertIntoLeaf(leaf *Node, key int, value []byte) bool {

	i := sort.SearchInts(leaf.Keys[:leaf.NumKeys], key) //finds the index of where we should insert the new key in the existing keys (sorted)

	if i < leaf.NumKeys && leaf.Keys[i] == key { //if KEY already exists! just modify the value (this is the update operation)
		leaf.Values[i] = value
		return true
	}


	for j := leaf.NumKeys; j >= i; j-- { //shift over all values to the right of the insertion index
		leaf.Keys[j] = leaf.Keys[j - 1]
		leaf.Values[j] = leaf.Values[j - 1]
	}

	leaf.Keys[i] = key //insert new key value pair
	leaf.Values[i] = value
	leaf.NumKeys++

	return false

}


func (tree *BTree) splitLeaf(leafBlockNum int, leaf *Node) error {

	mid := leaf.NumKeys / 2 //find middle index for split

	rightBlockNum := tree.allocateBlock() //allocate new block for the new right block for the spli

	rightLeaf := &Node{ //create new node which will be set as the new block
		IsLeaf: true,
		NumKeys: leaf.NumKeys - mid, //number of keys we are moving to the new node
		Keys: make([]int, tree.Order),
		Values: make([][]byte, tree.Order),
		NextBlock: leaf.NextBlock,
		ParentBlock: leaf.ParentBlock,
	}

	copy(rightLeaf.Keys, leaf.Keys[mid:leaf.NumKeys]) //copying over the corresponding keys to the new node
	copy(rightLeaf.Values, leaf.Values[mid:leaf.NumKeys]) //copying over the corresponding values to the new node
	leaf.NumKeys = mid
	leaf.NextBlock = rightBlockNum

	//now write the new node to the disk

	err := tree.writeNode(leafBlockNum, leaf) //write the leaf to the disk at its block number

	if err != nil {
		return err
	}
 
	err = tree.writeNode(rightBlockNum, rightLeaf) //write the new split leaf to the disk at its blkck number

	promoteKey := rightLeaf.Keys[0] //take the first key in the new split right node and move it to the parent


	if leaf.ParentBlock == -1 {
		//promote key is now the root
		return tree.createNewRoot(leafBlockNum, rightBlockNum, promoteKey)
	}

	return tree.insertIntoParent(leaf.ParentBlock, promoteKey, rightBlockNum)
}

func (tree *BTree) createNewRoot(leftBlockNum int, rightBlockNum int, key int) error {
	
	//allocate block for the new node and create the root
	rootBlockNum := tree.allocateBlock()

	root := &Node{
		IsLeaf: false,
		NumKeys: 1,
		Keys: make([]int, tree.Order),
		ChildBlocks: make([]int, tree.Order + 1),
		ParentBlock: -1,
	}

	root.Keys[0] = key
	root.ChildBlocks[0] = leftBlockNum
	root.ChildBlocks[1] = rightBlockNum

	//write the root to the disk
	err := tree.writeNode(rootBlockNum, root)

	if err != nil {
		return err
	}

	//get the left child and update its parent block to the new root
	leftChild, _ := tree.readNode(leftBlockNum)
	leftChild.ParentBlock = rootBlockNum
	tree.writeNode(leftBlockNum, leftChild)

	//get the right child and update its parent block to the new root
	rightChild, _ := tree.readNode(rightBlockNum)
	rightChild.ParentBlock = rootBlockNum
	tree.writeNode(rightBlockNum, rightChild)

	tree.RootBlockNum = rootBlockNum
	tree.Height++

	return nil
}

func (tree *BTree) insertIntoParent(parentBlockNum int, key int, rightChildBlockNum int) error {
	//take a key and insert it into the parent 

	parent, err := tree.readNode(parentBlockNum)

	if err != nil {
		return err
	}

	i := sort.SearchInts(parent.Keys[:parent.NumKeys], key) //finding the index of where the key should be placed in the parent keys

	for j := parent.NumKeys; j > i; j-- {
		parent.Keys[j] = parent.Keys[j-1] //shift all existing keys over to the right
		parent.ChildBlocks[j+1] = parent.ChildBlocks[j]
	}

	parent.Keys[i] = key
	parent.ChildBlocks[i + 1] = rightChildBlockNum
	parent.NumKeys++

	rightChild, _ := tree.readNode(rightChildBlockNum)
	rightChild.ParentBlock = parentBlockNum
	tree.writeNode(rightChildBlockNum, rightChild)
	
	err = tree.writeNode(parentBlockNum, parent)
	if err != nil {
		return err
	}

	if parent.NumKeys > tree.Order-1 {
		return tree.splitInternal(parentBlockNum, parent)
	}

	return nil
}




