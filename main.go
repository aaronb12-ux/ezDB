package main

import (
	"aaron/simpleDB/filemanager"
	"fmt"
	
)


func main() {


	filemgr := filemanager.NewFileMgr(5) //each block (and page stores 5 bytes)

	block := filemanager.MakeBlock("mydb.db", 0) //creating a block

	page := filemanager.MakePage(5)

	err := filemgr.Write(block, page)

	fmt.Println(err)

	fmt.Println(filemgr.GetWriteLog())






	
}






