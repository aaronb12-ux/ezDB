package main

//blockID represents a disk block for a corresponding file. A block belongs to a specific file and has a number
type BlockID struct {
	Filename string //file we are reading from
	Number int //the block in the file
}

