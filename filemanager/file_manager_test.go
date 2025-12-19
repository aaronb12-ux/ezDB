package filemanager_test

import (
	"aaron/simpleDB/filemanager"
	"testing"
)

func TestWriteToFile(t *testing.T) {

	filemgr := filemanager.NewFileMgr(5) //new file manager that stores blocks of 5

	block1 := filemanager.MakeBlock("mydb.db", 0) //database file and block number

	page1 := filemanager.MakePage(5) //page stores 5 bytes

	_ , err := page1.Write(0, []byte("b___1"))

	if err != nil {
		t.Fatalf("error writing to the page %v", err)
	}

	e := filemgr.Write(block1, page1)

	if e != nil {
		t.Fatalf("error writing to the file %v", e)
	}

}

func TestReadFromFile(t *testing.T) {
	filemgr := filemanager.NewFileMgr(5)

	block := filemanager.MakeBlock("mydb.db", 1) //should return b___1

	page := filemanager.MakePage(5)

	e := filemgr.Read(block, page)

	

	//need to read from block in file and store into page -> page.Read returns the bytes from the page AFTER reading from the file
 

}