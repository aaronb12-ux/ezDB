package main


import (
	"errors"
	"aaron/simpleDB/btree"
	"aaron/simpleDB/filemanager"
)

type Database struct {
	tree *btree.BTree
	fileMgr *filemanager.Filemgr
}

func Open(filename string, blockSize int, order int) (*Database, error) {

	//create filemgr
	fileMgr := filemanager.NewFileMgr(blockSize, filename)
	fileMgr.OpenFile(filename)

	//create b+tree
	tree := btree.MakeTree(order, fileMgr)

	return &Database{
		tree: tree,
		fileMgr: fileMgr,
	}, nil
}


func (db *Database) Put(key int, value []byte) error {
	if value == nil {
		return errors.New("value cannot be nil")
	}

	err := db.tree.Insert(key, value)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Get(key int) ([]byte, bool) {

	value, found, err := db.tree.Get(key)

	if err != nil {
		return nil, false
	}

	return value, found
}

func (db *Database) Delete(key int) bool {
	deleted, err := db.tree.Delete(key)

	if err != nil {
		return false
	}

	return deleted
}

func (db *Database) Size() int {
	return db.tree.Size
}

func (db *Database) Close() error {
	//flush any pending writes
	//close the file

	if db.fileMgr != nil && db.fileMgr.GetOpenedFile() != nil {
		return db.fileMgr.GetOpenedFile().Close()
	}

	return nil
}


type Stats struct {
	Size int
	Height int
	Order int
}

func (db *Database) Stats() Stats {
	return Stats{
		Size: db.tree.Size,
		Height: db.tree.Height,
		Order: db.tree.Order,
	}
}
