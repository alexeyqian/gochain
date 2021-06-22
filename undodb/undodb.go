package undodb

import (
	"github.com/alexeyqian/gochain/store"
)

type UndoableDB struct {
	store     store.Storage
	isUndoing bool
}

// these two table names are reserved by database
// user should not be able to create same table name
const metaTable = "_meta_"
const metaKey = "_meta_key_" // only one record in metaTable, and the key of the record is meteKey

// key in revision table is an stringized auto increment id
const revisionTable = "_revision_"

// TODO: make UndoableDB as interface
// NewMemUndoableDB, NewFileUndoableDB

func NewUndoableDB(storage store.Storage) *UndoableDB {
	udb := UndoableDB{
		store:     storage,
		isUndoing: false,
	}

	return &udb
}

func (udb *UndoableDB) Open() error {
	err := udb.store.Open()
	if err != nil {
		return err
	}

	if !udb.store.HasBucket(metaTable) {
		err := udb.store.CreateBucket(metaTable)
		if err != nil {
			return err
		}
		udb.initMetaData()

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

// HasKeyInt // table is using auto incremented id as key
func (udb *UndoableDB) HasKey(table, key string) bool {
	if table == "" || key == "" {
		return false
	}

	return udb.store.HasKey(table, []byte(key))
}

// GetInt // table is using auto increment id as key
func (udb *UndoableDB) Get(table, key string) ([]byte, error) {
	return udb.store.Get(table, []byte(key))
}

func (udb *UndoableDB) GetAll(table string) map[string][]byte {
	res := make(map[string][]byte)

	items, _ := udb.store.GetAll(table)
	for _, v := range items {
		res[string(v.Key)] = v.Value
	}
	return res
}
