package undodb

import (
	"fmt"

	"github.com/alexeyqian/gochain/store"
)

type Revision struct {
	num       uint32
	table     string
	operation string
	key       string
	data      []byte
}

type UndoableDB struct {
	store store.Storage
}

// these two table names are reserved by database
// user should not be able to create same table name
const metaTable = "_meta_"
const revisionTable = "_revision_"

func NewUndoableDB(storage store.Storage) *UndoableDB {

	udb := UndoableDB{
		store: storage,
	}

	return &udb
}

func (udb *UndoableDB) Open() error {
	err := udb.store.Open()
	if err != nil {
		return err
	}

	if !udb.store.HasBucket(metaTable) {
		// create data bucket and state bucket
		err := udb.store.CreateBucket(metaTable)
		if err != nil {
			return err
		}
		err = udb.store.CreateBucket(revisionTable)
		if err != nil {
			return err
		}
	}
	return nil
}

func (udb *UndoableDB) Close() error {
	return udb.store.Close()
}

func (udb *UndoableDB) Remove() error {
	return udb.store.Remove()
}

func (udb *UndoableDB) HasTable(name string) bool {
	return udb.store.HasBucket(name)
}

func (udb *UndoableDB) CreateTable(name string) error {
	if name == metaTable || name == revisionTable {
		return fmt.Errorf("cannot use internal names")
	}
	return udb.store.CreateBucket(name)
}
