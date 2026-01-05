package btree

import (
	"encoding/binary"
	
	"aaron/simpleDB/filemanager" // update this import path
)

//BTree represents a B+ tree that persists to disk
type BTree struct {
	FileMgr    *filemanager.Filemgr
	RootBlockNum int    //block number of root node
	Order      int     //maximum number of children
	Size       int     //total number of keys
	Height     int     //height of tree
	NextBlock  int     //next available block number
}


type Key struct {
	Key   int
	Value []byte
}


type Node struct {
	IsLeaf      bool
	NumKeys     int
	Keys        []int   //keys for searching
	ChildBlocks []int   //block numbers of children (for internal nodes)
	Values      [][]byte //values (for leaf nodes only)
	NextBlock   int     //block number of next leaf (for leaf nodes)
	ParentBlock int     //block number of parent
}

const (
	NODE_HEADER_SIZE = 13
	INT_SIZE = 4
	MAX_VALUE_SIZE = 100
)


func MakeTree(order int, fileMgr *filemanager.Filemgr) *BTree {
	if order < 3 {
		panic("B+ tree order must be at least 3")
	}
	return &BTree{
		FileMgr:      fileMgr,
		RootBlockNum: -1, //-1 means no root yet
		Order:        order,
		Size:         0,
		Height:       0,
		NextBlock:    0,
	}
}

//allocateBlock returns the next available block number
func (tree *BTree) allocateBlock() int {
	blockNum := tree.NextBlock
	tree.NextBlock++
	return blockNum
}

//serializeNode converts a Node to bytes and writes to a page
func (tree *BTree) serializeNode(node *Node, page *filemanager.Page) error {
	offset := 0
	
	//Write IsLeaf (1 byte)
	isLeafByte := byte(0)
	if node.IsLeaf {
		isLeafByte = 1
	}
	page.Write(offset, []byte{isLeafByte})
	offset += 1
}


