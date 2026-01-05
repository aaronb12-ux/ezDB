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

	page.Write(offset, []byte{isLeafByte}) //writing to offset 0 -> [ [0 or 1], [], [] , [], []...]

	offset += 1 //[ [0 or 1], [NOW HERE], [], [], []]


	//node's keys
	numberOfKeysAsBytes := make([]byte, INT_SIZE) //bytes buffer -> represents the number of keys the node has
	binary.BigEndian.PutUint32(numberOfKeysAsBytes, uint32(node.NumKeys)) //puts the number of keys into the containe (bytes array), converted to bytes
	page.Write(offset, numberOfKeysAsBytes) //write the bytes to the page at the offset...
	offset += INT_SIZE //update the offset...


	//node's nextblock
	nextBlockBytes := make([]byte, INT_SIZE) //bytes buffer
	binary.BigEndian.PutUint32(nextBlockBytes, uint32(node.NextBlock)) //putting the nextblock number into the byte buffer from above
	page.Write(offset, nextBlockBytes) //write to the page at the offset the nextBlockBytes we just allocated
	offset += INT_SIZE


	//node's parentblock
	parentBlockBytes := make([]byte, INT_SIZE)
	binary.BigEndian.PutUint32(parentBlockBytes, uint32(node.ParentBlock)) //putting the parent block number into the buffer above
	page.Write(offset, parentBlockBytes)
	offset += INT_SIZE


	//Now write the keys to the page...

	for i := 0; i < node.NumKeys; i++ {

		keyBytes := make([]byte, INT_SIZE) //create 4 byte buffer [[], [], [], []]
		binary.BigEndian.PutUint32(keyBytes, uint32(node.Keys[i])) //write key number to the above buffer
		page.Write(offset, keyBytes) //write the allocated bytes to the page
		offset += INT_SIZE

	}

	//now write the values
	if node.IsLeaf {

		for i := 0; i < node.NumKeys; i++ {

			lengthOfValue := len(node.Values[i]) //get the length of the value...

			lengthBytes := make([]byte, INT_SIZE) //create buffer

			binary.BigEndian.PutUint32(lengthBytes, uint32(lengthOfValue))

			page.Write(offset, lengthBytes)

			offset += INT_SIZE

			page.Write(offset, node.Values[i])

			offset += lengthOfValue

			//[ [length][actual value][length][actual value][length][actual value] ]

		} 
	} else {

		for i := 0; i < node.NumKeys; i++ {

			childBlockBytes := make([]byte, INT_SIZE)

			binary.BigEndian.PutUint32(childBlockBytes, uint32(node.ChildBlocks[i]))

			page.Write(offset, childBlockBytes)

			offset += INT_SIZE
		}
	}

		return nil
}


