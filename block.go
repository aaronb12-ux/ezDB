package main


//blockID represents a disk block for a corresponding file. A block belongs to a specific file and has a number
type BlockID struct {
	Filename string
	Number int
}