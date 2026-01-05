package btree

import (
	"encoding/binary"
	"aaron/simpleDB/filemanager"
)

//BTree represents a B+ tree that persists to disk
type BTree struct {
	FileMgr    *filemanager.Filemgr 
	RootBlockNum int //determines which block is the root node
	Order      int 
	Size       int //number of keys 
	Height     int //height of the tree
	NextBlock  int //next available block
}

//keys exists within Nodes
type Key struct {
	Key   int 
	Value []byte
}

//this struct represents a single node in the B+ tree
type Node struct {
	IsLeaf      bool
	NumKeys     int
	Keys        []int   
	ChildBlocks []int   
	Values      [][]byte 
	NextBlock   int    
	ParentBlock int 
}

const (
	INT_SIZE = 4 //size of an integer in bytes for serialization
)

//initializing the tree
func MakeTree(order int, fileMgr *filemanager.Filemgr) *BTree {
	if order < 3 { //making sure the order is above 3
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
	blockNum := tree.NextBlock //accessing the next available block
	tree.NextBlock++ //updating new 'nextblock' for next allocation...
	return blockNum //returning the current 'nextblock'
}

//serializeNode converts a Node to bytes and writes to a page
func (tree *BTree) serializeNode(node *Node, page *filemanager.Page) error {
	//the function acceps a node and page
	offset := 0 //current position -> currently at position 0 on the page
	
	//Write IsLeaf (1 byte)
	isLeafByte := byte(0)
	if node.IsLeaf { //if node is a leaf... at byte 0 on the page, we write a 1 and otherwise a 0
		isLeafByte = byte(1)
	}
	
	
}
	
