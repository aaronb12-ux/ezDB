package btree


type BTree struct {
	Root   *Node 
	Order  int  
	Size   int  
	Height int  
}


type Key struct {
	Key   int   
	Value []byte 
}

type Node struct {
	IsLeaf   bool    
	Keys     []int   
	Children []*Node
	Data     []Key   
	Next     *Node  
	Parent   *Node
}

