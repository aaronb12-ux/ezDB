package filemanager

import (
	"fmt"
	"os"
)

type Filemgr struct {
	blocksize int //how big each chunk of data should be
	openedFile *os.File //file currently in use
}

func NewFileMgr(blocksize int) *Filemgr {
	return &Filemgr{
		blocksize: blocksize,
		openedFile: os.NewFile(uintptr(10), "mydb.db"),
	}
}

func (fm *Filemgr) Read(blk *BlockID, p *Page) error {

	f := fm.openedFile

	offset := int64(blk.Number * fm.blocksize)

	_, err := f.Seek(offset, 0)

	if err != nil {
		return fmt.Errorf("error seeking from offset: %v", err)
	}

	bytesRead, err := f.Read(p.Bytes())

	if err != nil {
		return fmt.Errorf("failed to read block %v: %v", blk, err)
	}

	if bytesRead != fm.blocksize {
		return fmt.Errorf("incomplete read: expected %d bytes, got %d", fm.blocksize, bytesRead)
	}
	return nil
}