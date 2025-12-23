package filemanager

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Filemgr struct {
	blocksize int //how big each chunk of data should be
	filename string
	openedFile *os.File //file currently in use
	readLog []ReadWriteLogEntry
	writeLog []ReadWriteLogEntry
}


type ReadWriteLogEntry struct {
	Timestamp time.Time
	BlockId int
	BytesAmount int
	Bytes string

}

func  NewFileMgr(blocksize int) *Filemgr {
	return &Filemgr{
		blocksize: blocksize,
		filename: "simple.db",
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
	//the read starts ON the offset. so if the offset is 10, the "f[10]"" is included in the read

	if err != nil {
		return fmt.Errorf("failed to read block %v: %v", blk, err)
	}

	if bytesRead != fm.blocksize && bytesRead != 0 {
		return fmt.Errorf("incomplete read: expected %d bytes, got %d", fm.blocksize, bytesRead)
	}

	newReadEntry := ReadWriteLogEntry{ //create new read log object
		Timestamp: time.Now(),
		BlockId: blk.Number,
		BytesAmount: bytesRead,
		Bytes: string(p.bytes),
	}


	fm.readLog = append(fm.readLog, newReadEntry)

	fmt.Printf("we just read the bytes %v", string(p.bytes))
	
	return nil
}


func (fm *Filemgr) Write(blk *BlockID, p *Page) error {

	fmt.Println("the page bytes are", p.Bytes())

	f := fm.OpenFile(fm.filename) //open the file

	offset := int64(blk.Number * fm.blocksize) //get the offset in the file (this is where we begin our write)

	_, err := f.Seek(offset, 0) //seeking to the offset relative to the origin

	if err != nil { //if there is an error seeking to the offset
		return fmt.Errorf("error offsetting file %v", err) 
	}

	bytesWritten, err := f.Write(p.Bytes()) //writing to the file the bytes from the page

	if err != nil { //if there is an error writing the bytes
		return fmt.Errorf("error writing bytes to the file %v", err)
	}

	if bytesWritten != fm.blocksize { //if we do not write the appropriate amount of bytes
		return fmt.Errorf("incomplete write: expected %d bytes, wrote %d", fm.blocksize, bytesWritten)
	}


	newWriteEntry := ReadWriteLogEntry{ //create new writing log object
		Timestamp: time.Now(),
		BlockId: blk.Number,
		BytesAmount: bytesWritten,
		Bytes: string(p.bytes),
	}

	fm.writeLog = append(fm.writeLog, newWriteEntry)

	return nil

}

func (fm *Filemgr) OpenFile(filename string) *os.File  {

	n, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644) //open the file
 
	fm.openedFile = n 

	if err != nil {
		log.Fatalf("error opening file %v", err)
	}

	return fm.openedFile

}

func (fm *Filemgr) GetWriteLog() []ReadWriteLogEntry {
	return fm.writeLog
}

func (fm *Filemgr) GetReadLog() []ReadWriteLogEntry {
	return fm.readLog
}