package filemanager_test

import (
	"aaron/simpleDB/filemanager"
	"testing"
	"fmt"
)

func TestWriteAndReadToFileOneBlock(t *testing.T) {

	filemgr := filemanager.NewFileMgr(5)

	block1 := filemanager.MakeBlock("simple.db", 0) //database file and block number

	page1 := filemanager.MakePage(5) //page stores 5 bytes

	_ , err := page1.Write(0, []byte("b___1"))

	if err != nil {
		t.Fatalf("error writing to the page %v", err)
	}

	e := filemgr.Write(block1, page1)

	if e != nil {
		t.Fatalf("error writing to the file %v", e)
	}

	//NOW READ: make empty page -> read from file at the specified block

	page2 := filemanager.MakePage(5)

	e1 := filemgr.Read(block1, page2)

	if e1 != nil {
		t.Fatalf("error reading from the file from block %d. the error is: %v", block1.Number, e1)
	}

	fmt.Println(filemgr.GetReadLog())

	//bytes are in the database file. read the bytes from the file to the page. then read from the page to show the actual bytes
}

func TestWriteAndReadMultipleBlock(t *testing.T) {

	//assume everything works... no need for error handling this this test until final check

	filemgr := filemanager.NewFileMgr(5) 


	block0 := filemanager.MakeBlock("simple.db", 0)
	block1 := filemanager.MakeBlock("simple.db", 1)
	block2 := filemanager.MakeBlock("simple.db", 2)

	page0 := filemanager.MakePage(5) 
	page1 := filemanager.MakePage(5) 
	page2 := filemanager.MakePage(5)


	page0.Write(0, []byte("b___0"))
	page1.Write(0, []byte("b___1"))
	page2.Write(0, []byte("b___2"))


	filemgr.Write(block0, page0)
	filemgr.Write(block1, page1)
	filemgr.Write(block2, page2)


	//NOW READ 

	newPage := filemanager.MakePage(5)

	e := filemgr.Read(block2, newPage)

	if e != nil {
		t.Fatalf("error reading from block %d. the error is %v", block2.Number, e)
	}
}


func TestReadWriteOffset(t *testing.T) {

	filemgr := filemanager.NewFileMgr(10) 

	block0 := filemanager.MakeBlock("simple.db", 0)
	block1 := filemanager.MakeBlock("simple.db", 1)
	block2 := filemanager.MakeBlock("simple.db", 2)

	page0 := filemanager.MakePage(10) 
	page1 := filemanager.MakePage(10) 
	page2 := filemanager.MakePage(10) 


	page0.Write(1, []byte("b___0")) //REMEMBER THAT WE OFFSET FROM THE PAGE NOT THE BYTES BEING WRITTEN
	page1.Write(0, []byte("b___1"))
	page2.Write(0, []byte("b___2"))

	//when we offset, for example above: we are writing 5 bytes and offset by 1. so we must make sure we have room to write 5 + 1 bytes to the block

	filemgr.Write(block0, page0)
	filemgr.Write(block1, page1)
	filemgr.Write(block2, page2)
}