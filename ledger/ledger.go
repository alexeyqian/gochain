package ledger

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
)

const BlockOffsetBits = 64
const BlockSizeBits = 64
const IndexRecordSize = (BlockOffsetBits + BlockSizeBits) / 8

const ledgerFileName = "ledger"
const ledgerIndexFileName = "ledger.index"

type ledgerRecord struct {
	Size uint64 // size of ledger record
	Data []byte // block data bytes
}

type indexRecord struct {
	Offset uint64 // ledger record's Offset in ledger file
	Size   uint64 // size of ledger record
}

type ledger_s struct {
	ledgerPath string
	indexPath  string
	ledger     *os.File
	index      *os.File

	isOpen    bool
	isReadOny bool
}

// ========== global variables ==============
var glgr ledger_s

// ========== public functions ==============
func Open(dir string) {
	glgr.open(dir)
}

func Close() {
	glgr.close()
}

func Remove() {
	glgr.remove()
}

func Append(blockData []byte) {
	glgr.append(blockData)
}

func Read(bno int) ([]byte, error) {
	return glgr.read(bno)
}

// ========== private functions ==============
func (lg *ledger_s) open(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, os.ModePerm)
	}

	lg.ledgerPath = filepath.Join(dirPath, ledgerFileName)
	lg.indexPath = filepath.Join(dirPath, ledgerIndexFileName)
	var err error

	if fileExists(lg.ledgerPath) {
		lg.ledger, err = os.Open(lg.ledgerPath)
		if err != nil {
			return err
		}

		if fileExists(lg.indexPath) {
			lg.index, err = os.Open(lg.indexPath)
			// verify ledger and index are match
		} else {
			lg.index, err = os.Create(lg.indexPath)
			if err != nil {
				return err
			}
			// re-construct index
		}
	} else { // not exist, create them
		lg.ledger, err = os.Create(lg.ledgerPath)
		if err != nil {
			return err
		}

		if fileExists(lg.indexPath) {
			err = os.Remove(lg.indexPath)
			if err != nil {
				return err
			}
		}

		lg.index, err = os.Create(lg.indexPath)
		if err != nil {
			return err
		}
	}

	lg.isOpen = true
	lg.isReadOny = false

	return nil
}

func (lg *ledger_s) close() {
	lg.isOpen = false
	lg.ledger.Sync()
	lg.ledger.Close()
	lg.index.Sync()
	lg.index.Close()
}

func (lg *ledger_s) remove() error {
	var err error
	if !lg.isOpen {
		err = os.Remove(lg.ledgerPath)
		if err != nil {
			return err
		}
		err = os.Remove(lg.indexPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// ledger_record is a wrapper of block = block_size + block_data
// ledger: ledger_record_0 | ledger_record_1 | ledger_record_2 | ...
// index:  index_record_0  | index_record_1  | index_record_2  | ...
// ledger_record_0 is a dummy record
func (lg *ledger_s) append(blockData []byte) error {
	var err error

	if !lg.isOpen || lg.isReadOny {
		return errors.New("cannot append to ledger, not open or readonly!")
	}

	info, err := os.Stat(lg.ledgerPath)
	if err != nil {
		return err
	}
	fileSize := info.Size()
	blockSize := len(blockData)
	size := uint64(BlockSizeBits)/uint64(8) + uint64(blockSize)
	lrecord := ledgerRecord{Size: size, Data: blockData}
	irecord := indexRecord{Offset: uint64(fileSize), Size: lrecord.Size}

	// serialization
	lbuf := new(bytes.Buffer)
	err = binary.Write(lbuf, binary.BigEndian, lrecord.Size)
	if err != nil {
		return err
	}
	_, err = lg.ledger.Write(lbuf.Bytes())
	if err != nil {
		return err
	}
	_, err = lg.ledger.Write(blockData)

	ibuf := new(bytes.Buffer)
	err = binary.Write(ibuf, binary.BigEndian, irecord)
	if err != nil {
		return err
	}
	_, err = lg.index.Write(ibuf.Bytes())
	if err != nil {
		return err
	}

	// append block to ledger_s
	// os.File.Write is unbufferred, so no flush needed,
	// but still need Sync to make sure operating system's file system call system call
	// to write data to disk
	lg.ledger.Sync()
	lg.index.Sync()

	return nil
}

func (lg *ledger_s) read(bno int) ([]byte, error) {
	var err error
	var indexData, ledgerData []byte

	indexData, err = readFromFile(lg.index, uint64(bno*IndexRecordSize), IndexRecordSize)
	if err != nil {
		return nil, err
	}
	ibuf := bytes.NewReader(indexData)
	var irecord indexRecord
	err = binary.Read(ibuf, binary.BigEndian, &irecord)
	if err != nil {
		return nil, err
	}

	ledgerData, err = readFromFile(lg.ledger, irecord.Offset, irecord.Size)
	// the first 8 bytes is block size, the rest is the block data.
	return ledgerData[(BlockSizeBits / 8):], nil
}

func fileExists(fPath string) bool {
	_, err := os.Stat(fPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func readFromFile(file *os.File, Offset uint64, size uint64) ([]byte, error) {
	res := make([]byte, size)
	if _, err := file.ReadAt(res, int64(Offset)); err != nil {
		return nil, err
	}

	return res, nil
}
