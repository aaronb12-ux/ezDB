package btree_test

import (
	"aaron/simpleDB/filemanager"
	"aaron/simpleDB/btree"
	"testing"
	"fmt"

)

func TestCreateDB(t *testing.T) {

	filmgr := filemanager.NewFileMgr(4096, "testbtree.db")
	tree := btree.MakeTree(4, filmgr)

	if tree.Height == 0 {
		fmt.Println("successfully made tree with height 0")
	} else {
		t.Fatal("failed making tree")
	}
}