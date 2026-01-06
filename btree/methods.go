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




