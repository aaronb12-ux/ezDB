package btree

import (
	"aaron/ezDB/filemanager"
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