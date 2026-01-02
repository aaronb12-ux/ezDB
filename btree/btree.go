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
	
	//Write NumKeys (4 bytes)
	numKeysBytes := make([]byte, INT_SIZE)
	binary.BigEndian.PutUint32(numKeysBytes, uint32(node.NumKeys))
	page.Write(offset, numKeysBytes)
	offset += INT_SIZE
	
	//Write NextBlock (4 bytes)
	nextBlockBytes := make([]byte, INT_SIZE)
	binary.BigEndian.PutUint32(nextBlockBytes, uint32(node.NextBlock))
	page.Write(offset, nextBlockBytes)
	offset += INT_SIZE
	
	//Write ParentBlock (4 bytes)
	parentBlockBytes := make([]byte, INT_SIZE)
	binary.BigEndian.PutUint32(parentBlockBytes, uint32(node.ParentBlock))
	page.Write(offset, parentBlockBytes)
	offset += INT_SIZE
	
	//Write Keys
	for i := 0; i < node.NumKeys; i++ {
		keyBytes := make([]byte, INT_SIZE)
		binary.BigEndian.PutUint32(keyBytes, uint32(node.Keys[i]))
		page.Write(offset, keyBytes)
		offset += INT_SIZE
	}
	
	if node.IsLeaf {
		//Write Values for leaf nodes
		for i := 0; i < node.NumKeys; i++ {
			// Write value length
			valLen := len(node.Values[i])
			valLenBytes := make([]byte, INT_SIZE)
			binary.BigEndian.PutUint32(valLenBytes, uint32(valLen))
			page.Write(offset, valLenBytes)
			offset += INT_SIZE
			
			//Write value data
			page.Write(offset, node.Values[i])
			offset += valLen
		}
	} else {
		//Write ChildBlocks for internal nodes
		for i := 0; i <= node.NumKeys; i++ {
			childBlockBytes := make([]byte, INT_SIZE)
			binary.BigEndian.PutUint32(childBlockBytes, uint32(node.ChildBlocks[i]))
			page.Write(offset, childBlockBytes)
			offset += INT_SIZE
		}
	}
	
	return nil
}

//deserializeNode reads bytes from a page and constructs a Node
func (tree *BTree) deserializeNode(page *filemanager.Page) (*Node, error) {
	node := &Node{}
	offset := 0
	
	//Read IsLeaf (1 byte)
	isLeafBytes := make([]byte, 1)
	page.Read(offset, isLeafBytes)
	node.IsLeaf = isLeafBytes[0] == 1
	offset += 1
	
	//Read NumKeys (4 bytes)
	numKeysBytes := make([]byte, INT_SIZE)
	page.Read(offset, numKeysBytes)
	node.NumKeys = int(binary.BigEndian.Uint32(numKeysBytes))
	offset += INT_SIZE
	
	//Read NextBlock (4 bytes)
	nextBlockBytes := make([]byte, INT_SIZE)
	page.Read(offset, nextBlockBytes)
	node.NextBlock = int(binary.BigEndian.Uint32(nextBlockBytes))
	offset += INT_SIZE
	
	//Read ParentBlock (4 bytes)
	parentBlockBytes := make([]byte, INT_SIZE)
	page.Read(offset, parentBlockBytes)
	node.ParentBlock = int(binary.BigEndian.Uint32(parentBlockBytes))
	offset += INT_SIZE
	
	//Read Keys
	node.Keys = make([]int, node.NumKeys)
	for i := 0; i < node.NumKeys; i++ {
		keyBytes := make([]byte, INT_SIZE)
		page.Read(offset, keyBytes)
		node.Keys[i] = int(binary.BigEndian.Uint32(keyBytes))
		offset += INT_SIZE
	}
	
	if node.IsLeaf {
		//Read Values for leaf nodes
		node.Values = make([][]byte, node.NumKeys)
		for i := 0; i < node.NumKeys; i++ {
			// Read value length
			valLenBytes := make([]byte, INT_SIZE)
			page.Read(offset, valLenBytes)
			valLen := int(binary.BigEndian.Uint32(valLenBytes))
			offset += INT_SIZE
			
			//Read value data
			node.Values[i] = make([]byte, valLen)
			page.Read(offset, node.Values[i])
			offset += valLen
		}
	} else {
		//Read ChildBlocks for internal nodes
		node.ChildBlocks = make([]int, node.NumKeys+1)
		for i := 0; i <= node.NumKeys; i++ {
			childBlockBytes := make([]byte, INT_SIZE)
			page.Read(offset, childBlockBytes)
			node.ChildBlocks[i] = int(binary.BigEndian.Uint32(childBlockBytes))
			offset += INT_SIZE
		}
	}
	
	return node, nil
}

//writeNode writes a node to disk at the specified block number
func (tree *BTree) writeNode(blockNum int, node *Node) error {
	// Create a page
	page := filemanager.MakePage(tree.FileMgr.BlockSize())
	
	//Serialize node to page
	err := tree.serializeNode(node, page)
	if err != nil {
		return err
	}
	
	//Create block ID
	blk := filemanager.MakeBlock(tree.FileMgr.Filename(), blockNum)
	
	//Write page to disk
	return tree.FileMgr.Write(blk, page)
}

//readNode reads a node from disk at the specified block number
func (tree *BTree) readNode(blockNum int) (*Node, error) {
	// Create a page
	page := filemanager.MakePage(tree.FileMgr.BlockSize())
	
	//Create block ID
	blk := filemanager.MakeBlock(tree.FileMgr.Filename(), blockNum)
	
	//Read page from disk
	err := tree.FileMgr.Read(blk, page)
	if err != nil {
		return nil, err
	}

	//Deserialize page to node
	return tree.deserializeNode(page)
}


