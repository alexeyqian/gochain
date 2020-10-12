package Ledger

import (
	"os"
	"path/filepath"
	"bytes"
	"encoding/gob"
)

const BlockSizeBits  = 32
const IndexRecordSize = 12 // index record bytes

const ledgerFileName = "Ledger"
const ledgerIndexFileName = "Ledger.index"

type ledgerRecord struct{
	Size uint32 // size of ledger record
	Data []byte // block data bytes
}

type indexRecord struct{
	Offset uint64 // ledger record's Offset in ledger file
	Size uint32 // size of ledger record
}

type Ledger struct{
	ledgerPath string	
	indexPath string	
	ledger *os.File
	index  *os.File

	isOpen bool
	isReadOny bool
}

func OpenLedger(dir string) {
	var lg Ledger
	lg.open(dir)
}

func (lg *Ledger) open(dirPath string){	
	lg.ledgerPath = filepath.Join(dirPath, ledgerFileName);
	lg.indexPath = filepath.Join(dirPath, ledgerIndexFileName)
	var err error

	if fileExists(lg.ledgerPath){
		lg.ledger, err = os.Open(lg.ledgerPath)
		check(err)

		if(fileExists(lg.indexPath)){
			lg.index, err = os.Open(lg.indexPath)
			// verify ledger and index are match			
		}else{
			lg.index, err = os.Create(lg.indexPath)
			check(err)
			// re-construct index
		}		
	}else{ // not exist, create them
		lg.ledger, err = os.Create(lg.ledgerPath)
		check(err)
		
		if(fileExists(lg.indexPath)){
			err = os.Remove(lg.indexPath)
			check(err)
		}
			
		lg.index, err = os.Create(lg.indexPath)
		check(err)
	}

	lg.isOpen = true
	lg.isReadOny = false	
}

func (lg *Ledger) close(){
	lg.isOpen = false
	lg.ledger.Sync() 	
	lg.ledger.Close()
	lg.index.Sync()
	lg.index.Close()
}

func (lg *Ledger) remove() {
	var err error
	if !lg.isOpen {
		err = os.Remove(lg.ledgerPath)
		check(err)
		err = os.Remove(lg.indexPath)
		check(err)
	}	
}

// ledger_record is a wrapper of block = block_size + block_data
// ledger: ledger_record_0 | ledger_record_1 | ledger_record_2 | ...
// index:  index_record_0  | index_record_1  | index_record_2  | ...
// ledger_record_0 is a dummy record
func (lg *Ledger) append(blockData []byte){
	var err error

	if !lg.isOpen || lg.isReadOny {
		panic("cannot append to ledger, not open or readonly!")
	}
	
	info, err := os.Stat(lg.ledgerPath)
	check(err)
	fileSize := info.Size()
	blockSize := len(blockData)
	size := uint32(BlockSizeBits)/uint32(8) + uint32(blockSize)
	lrecord := ledgerRecord{Size: size, Data: blockData}
	irecord := indexRecord{Offset: uint64(fileSize), Size: lrecord.Size}
	
	// serialization
	var lbytes bytes.Buffer
	lenc := gob.NewEncoder(&lbytes)
	err = lenc.Encode(lrecord)
	check(err)

	var ibytes bytes.Buffer
	ienc := gob.NewEncoder(&ibytes)
	err = ienc.Encode(irecord)
	check(err)
	
	// write data to files
	_, err = lg.ledger.Write(lbytes.Bytes())
	check(err)
	_, err = lg.index.Write(ibytes.Bytes())
	check(err)

	// append block to Ledger
	// os.File.Write is unbufferred, so no flush needed, 
	// but still need Sync to make sure operating system's file system call system call
	// to write data to disk
	lg.ledger.Sync() 
	lg.index.Sync()
}

func (lg *Ledger) read(bno int) []byte{	
	var err error
	var indexData, ledgerData []byte

	indexData, err = readFromFile(lg.index, uint64(bno) * uint64(IndexRecordSize), IndexRecordSize)
	check(err)
	
	var ibuffer bytes.Buffer
	ibuffer.Write(indexData)
	idec := gob.NewDecoder(&ibuffer)
	var irecord indexRecord
	err = idec.Decode(&irecord)
	check(err)
	
	ledgerData, err = readFromFile(lg.ledger, irecord.Offset, irecord.Size)
	var lbuffer bytes.Buffer
	lbuffer.Write(ledgerData)
	ldec := gob.NewDecoder(&lbuffer)
	var lrecord ledgerRecord
	err = ldec.Decode(&lrecord)

	return lrecord.Data
}

// utils
func check(e error){
	if e != nil {
		panic(e)
	}
}

func fileExists(fPath string) bool{
	_, err := os.Stat(fPath)
	if os.IsNotExist(err){
		return false
	}
	return true
}

func readFromFile(file *os.File, Offset uint64, size uint32) ([]byte, error){
	res := make([]byte, int(size))
	if _, err := file.ReadAt(res, int64(Offset)); err != nil{
		return nil, err
	}

	return res, nil
}