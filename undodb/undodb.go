package undodb

import (
	"fmt"

	"github.com/alexeyqian/gochain/utils"

	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/store"
)

type MetaData struct {
	Revision uint32
}

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
const metaKey = "_meta_key_"
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

		// init some meta data
		metaData := MetaData{
			Revision: 0,
		}
		udb.store.Put(metaTable, metaKey, entity.Serialize(metaData))
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
	if name == "" {
		return false
	}
	return udb.store.HasBucket(name)
}

func (udb *UndoableDB) HasKey(table, key string) bool {
	if table == "" || key == "" {
		return false
	}

	return udb.store.HasKey(table, key)
}

func (udb *UndoableDB) Get(table, key string) ([]byte, error) {
	return udb.store.Get(table, key)
}

func (udb *UndoableDB) CreateTable(name string) error {
	if name == metaTable || name == revisionTable {
		return fmt.Errorf("cannot use internal names")
	}
	return udb.store.CreateBucket(name)
}

func (udb *UndoableDB) RowCount(table string) int {
	return udb.store.RowCount(table)
}

func (udb *UndoableDB) Create(table, key string, data []byte) error {
	if table == "" {
		return fmt.Errorf("create: table is empty")
	}
	if !udb.HasTable(table) {
		return fmt.Errorf("create: doesn't have the table:%s", table)
	}
	if key == "" {
		return fmt.Errorf("create: key is mepty.")
	}
	if udb.HasKey(table, key) {
		return fmt.Errorf("create: key already exist.")
	}

	// TODO: need a transaction here // TRANSACTION HERE
	// TODO: store.BatchInTransaction(operations)
	// 1. update state table
	err := udb.onCreate(table, key)
	if err != nil {
		return err
	}

	// 2. save to data table
	err = udb.store.Put(table, key, data)
	if err != nil {
		return err
	}

	// 3. update db meta table
	//err = udb.onCreateMeta(table, key)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (udb *UndoableDB) onCreate(table, key string) error {
	num := udb.getCurrentRevision()
	if num == 0 { // no undo session
		return nil
	}
	revision := Revision{
		num:       num,
		table:     table,
		operation: "create",
		key:       key,
	}

	udb.store.Put(revisionTable, utils.CreateUuid(), entity.Serialize(revision))

	return nil
}

func (udb *UndoableDB) StartUndoSession() uint32 {
	meta := udb.GetMetaData()
	meta.Revision++
	udb.updateMetaData(meta)
	return meta.Revision
}

func (udb *UndoableDB) getCurrentRevision() uint32 {
	return udb.GetMetaData().Revision
}

func (udb *UndoableDB) GetMetaData() *MetaData {
	data, err := udb.store.Get(metaTable, metaKey)
	if err != nil {
		panic("cannot get meta data")
	}
	var meta MetaData
	entity.Deserialize(&meta, data)
	return &meta
}

func (udb *UndoableDB) updateMetaData(meta *MetaData) {
	err := udb.store.Put(metaTable, metaKey, entity.Serialize(*meta))
	if err != nil {
		panic(err)
	}
}
