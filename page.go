package main

import (

	"errors"

)

//a page is the buffer where we store the data after we read it from the block. The data is stored as an array of bytes and has a size
//we load a block of data into the page, manipulate the bytes within the page, then flush it back to the disk
type Page struct {
	bytes []byte
}


//create a page
func MakePage(size int) *Page {

	return &Page{
		bytes: make([]byte, size),
	}
}


func (p *Page) Write(offset int, data []byte) (int, error) {
	
	if offset + len(data) > p.Size() {
		return 0, errors.New("page length exceeded.")
	} 

	res := copy(p.bytes[offset:], data)

	return res, nil

}

func (p *Page) Read(offset int, dst []byte) int {

	res := copy(dst, p.bytes[offset:])

	return res
}

func (p *Page) Bytes() []byte  {
	return p.bytes
}

func (p *Page) Size() int {
	return len(p.bytes)
}


