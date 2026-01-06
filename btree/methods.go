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
}




