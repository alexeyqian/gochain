package undodb

import (
	"sort"

	"github.com/alexeyqian/gochain/store"

	"github.com/alexeyqian/gochain/entity"
)

/*
TODO: limit operations in revision to 256
limit revision Num to 256
so max undoable operation is 256 x 256, which is good enough
*/

type Revision struct {
	Num       uint32
	Table     string
	Operation string
	Key       string
	Data      []byte
}

func (udb *UndoableDB) StartUndoSession() uint32 {
	// @future validate max revision Num to see if can process further

	meta := udb.getMetaData()
	meta.Revision++
	udb.updateMetaData(meta)
	return meta.Revision
}

func (udb *UndoableDB) UndoLastSession() {
	meta := udb.getMetaData()
	if meta.Revision <= 0 {
		panic("undo session error, revision should be greater then 0")
	}

	revisions := udb.getAllRevisions(meta.Revision)

	// sort map by key(uint64) in reverse order
	keys := make([]uint64, len(revisions))
	i := 0
	for k := range revisions {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })
	// To perform the opertion you want
	for _, k := range keys {
		udb.undoOperation(revisions[k])
	}

	meta.Revision--
	udb.updateMetaData(meta)
}

func (udb *UndoableDB) CommitLastSession() {
	meta := udb.getMetaData()
	if meta.Revision <= 0 {
		panic("commit session error, revision should be greater then 0")
	}

	revisions := udb.getAllRevisions(meta.Revision)
	for k, _ := range revisions {
		udb.store.Delete(revisionTable, store.IntKeyToBytes(k))
	}

	meta.Revision--
	udb.updateMetaData(meta)
}

func (udb *UndoableDB) getCurrentRevision() uint32 {
	return udb.getMetaData().Revision
}

func (udb *UndoableDB) undoOperation(op Revision) {
	if op.Operation == "create" {
		udb.store.Delete(op.Table, []byte(op.Key))
		//fmt.Printf("undo operation: %+v\n", op)
	}
}

func (udb *UndoableDB) getAllRevisions(num uint32) map[uint64]Revision {
	rawRevisions, _ := udb.store.GetAll(revisionTable)
	revisions := make(map[uint64]Revision, len(rawRevisions))
	for _, v := range rawRevisions {
		key := store.BytesToIntKey(v.Key)
		var rev Revision
		entity.Deserialize(&rev, v.Value)
		if rev.Num == num {
			revisions[key] = rev
		}
	}
	return revisions
}

func (udb *UndoableDB) saveRevision(table, key string, data []byte, operation string) error {
	num := udb.getCurrentRevision()
	if num == 0 { // no undo session
		return nil
	}
	revision := Revision{
		Num:       num,
		Table:     table,
		Operation: operation,
		Key:       key,
		Data:      data,
	}

	err := udb.store.Put(revisionTable, nil, entity.Serialize(revision))
	//fmt.Printf("create undo operation: %+v\n", revision)
	return err

}
