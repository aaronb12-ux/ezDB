package main

import (
	"aaron/simpleDB/filemanager"
	"fmt"
)


func main() {


	filemgr := filemanager.NewFileMgr(5)

	block := filemanager.MakeBlock("mydb.db", 0)
	page := filemanager.MakePage(5)

	err := filemgr.Write(block, page)

	fmt.Println(err)

	

}






