package undodb

import (
	"fmt"

	"github.com/alexeyqian/gochain/entity"
)

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
	// 1. update revision table
	err := udb.onCreate(table, key)
	if err != nil {
		return err
	}

	// 2. save to data table
	err = udb.store.Put(table, []byte(key), data)
	if err != nil {
		return err
	}

	return nil
}

func (udb *UndoableDB) onCreate(table, key string) error {
	num := udb.getCurrentRevision()
	if num == 0 { // no undo session
		return nil
	}
	revision := Revision{
		Num:       num,
		Table:     table,
		Operation: "create",
		Key:       key,
	}

	err := udb.store.Put(revisionTable, nil, entity.Serialize(revision))
	//fmt.Printf("create undo operation: %+v\n", revision)
	return err
}
