package undodb

import (
	"fmt"
)

func (udb *UndoableDB) Delete(table, key string) error {
	if table == "" {
		return fmt.Errorf("table is empty")
	}
	if !udb.HasTable(table) {
		return fmt.Errorf("doesn't have the table:%s", table)
	}
	if key == "" {
		return fmt.Errorf("key is mepty.")
	}
	if !udb.HasKey(table, key) {
		return fmt.Errorf("key doesn't exist.")
	}

	data, err := udb.store.Get(table, []byte(key))
	if err != nil {
		panic(err)
	}

	// TODO: need a transaction here // TRANSACTION HERE
	// TODO: store.BatchInTransaction(operations)
	// 1. update revision table
	err = udb.onDelete(table, key, data)
	if err != nil {
		return err
	}

	// 2. save to data table
	err = udb.store.Delete(table, []byte(key))
	if err != nil {
		return err
	}

	return nil
}

func (udb *UndoableDB) onDelete(table, key string, data []byte) error {
	return udb.saveRevision(table, key, data, "delete")
}
