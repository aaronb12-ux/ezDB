package btree

import (
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
