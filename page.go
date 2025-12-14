package main

//a page is the buffer where we store the data after we read it from the block. The data is stored as an array of bytes and has a size
type Page struct {
	Data []byte
	Size int
}
