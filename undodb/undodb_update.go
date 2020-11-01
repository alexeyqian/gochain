package undodb

import (
	"fmt"
)

func (udb *UndoableDB) Update(table, key string, data []byte) error {
	if table == "" {
		return fmt.Errorf("update: table is empty")
	}
	if !udb.HasTable(table) {
		return fmt.Errorf("update: doesn't have the table:%s", table)
	}
	if key == "" {
		return fmt.Errorf("update: key is mepty.")
	}
	if !udb.HasKey(table, key) {
		return fmt.Errorf("update: key doesn't exist.")
	}

	// TODO: need a transaction here // TRANSACTION HERE
	// TODO: store.BatchInTransaction(operations)
	// 1. update revision table
	olddata, _ := udb.Get(table, key)
	err := udb.onUpdate(table, key, olddata)
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

func (udb *UndoableDB) onUpdate(table, key string, data []byte) error {
	return udb.saveRevision(table, key, data, "update")
}
