package main

import (
	"aaron/simpleDB/filemanager"
)


func main() {

	fm := filemanager.NewFileMgr(5)

	var block filemanager.BlockID

	block.Filename = "mydb.db"

	page := filemanager.MakePage(5)
	
}






