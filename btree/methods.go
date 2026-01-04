package btree

package btree

import (
	"fmt"
	"sort"
)


func (tree *BTree) Get(key int) ([]byte, bool, error) {

	if tree.RootBlockNum == -1 {
		return nil, false, nil
	}
	
	leafBlockNum, err := tree.findLeafBlock(key)
	if err != nil {
		return nil, false, err
	}
	
	leaf, err := tree.readNode(leafBlockNum)
	if err != nil {
		return nil, false, err
	}
	
	i := sort.SearchInts(leaf.Keys[:leaf.NumKeys], key)
	
	if i < leaf.NumKeys && leaf.Keys[i] == key {
		return leaf.Values[i], true, nil
	}
	
	return nil, false, nil
}


func (tree *BTree) findLeafBlock(key int) (int, error) {
	blockNum := tree.RootBlockNum
	
	for {

		node, err := tree.readNode(blockNum)
		if err != nil {
			return -1, err
		}
		
		
		if node.IsLeaf {
			return blockNum, nil
		}
		
		i := sort.SearchInts(node.Keys[:node.NumKeys], key)
		

		if i > node.NumKeys {
			i = node.NumKeys
		}
		
		blockNum = node.ChildBlocks[i]
	}
}



