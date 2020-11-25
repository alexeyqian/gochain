package undodb

import "fmt"

// TODO: add param keytype: enum: string or auto incremented id as int
func (udb *UndoableDB) CreateTable(name string) error {
	if name == metaTable || name == revisionTable {
		return fmt.Errorf("cannot use internal names")
	}
	return udb.store.CreateBucket(name)
}

func (udb *UndoableDB) DeleteTable(name string) error {
	return udb.store.DeleteBucket(name)
}

func (udb *UndoableDB) HasTable(name string) bool {
	if name == "" {
		return false
	}
	return udb.store.HasBucket(name)
}

func (udb *UndoableDB) RowCount(table string) int {
	return udb.store.RowCount(table)
}
