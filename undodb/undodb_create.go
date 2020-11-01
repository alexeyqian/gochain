package undodb

import (
	"fmt"
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
	return udb.saveRevision(table, key, nil, "create")
}
