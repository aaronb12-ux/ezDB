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

//traverses the tree to find the leaf block (node) where a key should be
func (tree *BTree) findLeafBlock(key int) (int, error) {

	blockNum := tree.RootBlockNum //begin at the first block for the search

	for {
		node, err := tree.readNode(blockNum) //read current node from the disk
		if err != nil {
			return -1, err
		}

		if node.IsLeaf { //if we have reached a leaf, we are done and return the block number
			return blockNum, nil
		}

		//if not leaf, we have internal node. must find which child to follow
		//sort.SearchInts returns the index of where they key should be inserted -> which means it lets us know which child to follow
		i := sort.SearchInts(node.Keys[:node.NumKeys], key)

		if i > node.NumKeys { //make sure we don't go out of bounds
			i = node.NumKeys
		}

		blockNum = node.ChildBlocks[i] //update blockNumber to for search continuation
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

	leaf, err := tree.readNode(leafBlockNum) //take the block data from the disk and read it into a page which is then deserialized to a node

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
	//take a key and insert it into the parent at parentBlockNum

	parent, err := tree.readNode(parentBlockNum)

	if err != nil {
		return err
	}

	i := sort.SearchInts(parent.Keys[:parent.NumKeys], key) //finding the index of where the key should be placed in the parent keys

	/*
		//parentKeys = [1, 5] and we are inserting 4

		i = 1

		j = 2

		parentKeys = [1, 5, 5]

		parentKeys = [1, 4, 5]

	*/
	
	for j := parent.NumKeys; j > i; j-- {
		parent.Keys[j] = parent.Keys[j-1] //shift all keys greater than our insertion key to the right
		parent.ChildBlocks[j+1] = parent.ChildBlocks[j] //?
	}

	parent.Keys[i] = key //put in new key
	parent.ChildBlocks[i + 1] = rightChildBlockNum
	parent.NumKeys++

	//update right childs parentBlock
	rightChild, _ := tree.readNode(rightChildBlockNum)
	rightChild.ParentBlock = parentBlockNum
	tree.writeNode(rightChildBlockNum, rightChild)
	
	//write the parent back to the disk after updates are complete
	err = tree.writeNode(parentBlockNum, parent)

	if err != nil {
		return err
	}

	if parent.NumKeys > tree.Order-1 {
		return tree.splitInternal(parentBlockNum, parent)
	}

	return nil
}
 
func (tree *BTree) splitInternal(nodeBlockNum int, node *Node) error {

	mid := node.NumKeys / 2
	promoteKey := node.Keys[mid] //key we are moving 'promoting' to the parent

	rightBlockNum := tree.allocateBlock() //create new block for new node we are making

	rightNode := &Node{ //new node
		IsLeaf: false,
		NumKeys: node.NumKeys - mid - 1,
		Keys: make([]int, tree.Order),
		ChildBlocks: make([]int, tree.Order+1),
		ParentBlock: node.ParentBlock,
	}

	//fill in the new node with all values to the right of the promote key
	copy(rightNode.Keys, node.Keys[mid + 1: node.NumKeys])
	copy(rightNode.ChildBlocks, node.ChildBlocks[mid + 1: node.NumKeys + 1]) 

	for i := 0; i <= rightNode.NumKeys; i++ {
		//iterate through all the keys in the new node and assign the child's parent this new right node
		child, _ := tree.readNode(rightNode.ChildBlocks[i])
		child.ParentBlock = rightBlockNum
		tree.writeNode(rightNode.ChildBlocks[i], child)
	}

	//update the number of keys for the node since we split it
	node.NumKeys = mid

	//flush back to the disk
	tree.writeNode(nodeBlockNum, node)
	tree.writeNode(rightBlockNum, rightNode)

	if node.ParentBlock == -1 {
		return tree.createNewRoot(nodeBlockNum, rightBlockNum, promoteKey)
	}

	return tree.insertIntoParent(node.ParentBlock, promoteKey, rightBlockNum) //insert the promote key into the parent

}


func (tree *BTree) Delete(key int) (bool, error) {

	//empty tree
	if tree.RootBlockNum == -1 {
		return false, nil
	}

	leafBlockNum, err := tree.findLeafBlock(key) //get the block number of the key we are trying to delete

	if err != nil {
		return false, err
	}

	leaf, err := tree.readNode(leafBlockNum) //deserialize the block into a leaf node object in memory


	i := sort.SearchInts(leaf.Keys[:leaf.NumKeys], key) //find the index where the key is 

	if i >= leaf.NumKeys || leaf.Keys[i] != key {
		return false, nil //key is not found
	}


	for j := i; j < leaf.NumKeys - 1; j++ {
		leaf.Keys[j] = leaf.Keys[j + 1] //overwrite the key 
		leaf.Values[j] = leaf.Values[j + 1]
	}

	leaf.NumKeys--
	tree.Size--


	err = tree.writeNode(leafBlockNum, leaf)

	if err != nil {
		return false, err
	}

	if leafBlockNum == tree.RootBlockNum && leaf.NumKeys == 0 {  //tree is now empty
		tree.RootBlockNum = -1
		tree.Height = 0
	}

	return true, nil

}

