package filemanager_test

import (
	"aaron/simpleDB/filemanager"
	"testing"
	"fmt"
)

func TestWriteToFile(t *testing.T) {

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

	page2 := filemanager.MakePage(5)

	e1 := filemgr.Read(block1, page2)

	if e1 != nil {
		t.Fatalf("error reading from the file from block %d. the error is: %v", block1.Number, e1)
	}

	fmt.Println(filemgr.GetReadLog())

	//bytes are in the database file. read the bytes from the file to the page. then read from the page to show the actual bytes
}

