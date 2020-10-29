package undodb

import (
	"testing"

	"github.com/alexeyqian/gochain/store"
)

func TestCreateUndoDB(t *testing.T) {
	pathname := "test.db"
	storage := store.NewBoltStorage(pathname)
	udb := NewUndoableDB(storage)

	udb.Open()

	if !udb.HasTable(metaTable) {
		t.Errorf("meta table is not created")
	}
	if !udb.HasTable(revisionTable) {
		t.Errorf("revision table is not created")
	}

	udb.Close()
	udb.Remove()
}
