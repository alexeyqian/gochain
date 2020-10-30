package undodb

import (
	"encoding/binary"
	"sort"

	"github.com/alexeyqian/gochain/store"

	"github.com/alexeyqian/gochain/entity"
)

/*
TODO: limit operations in revision to 256
limit revision num to 256
so max undoable operation is 256 x 256, which is good enough
*/

type Revision struct {
	num       uint32
	table     string
	operation string
	key       string
	data      []byte
}

func (udb *UndoableDB) StartUndoSession() uint32 {
	// @future validate max revision num to see if can process further

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

	rawRevisions, _ := udb.store.GetAll(revisionTable)
	var revisions map[uint64]Revision
	for _, v := range rawRevisions {
		var key uint64
		if store.IsIntKeyEncodedInBigEndian {
			key = binary.BigEndian.Uint64(v.Key)
		} else {
			key = binary.LittleEndian.Uint64(v.Key)
		}

		var rev Revision
		entity.Deserialize(&rev, v.Value)
		revisions[key] = rev
	}

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

func (udb *UndoableDB) getCurrentRevision() uint32 {
	return udb.getMetaData().Revision
}

func (udb *UndoableDB) undoOperation(op Revision) {
	if op.operation == "create" {
		udb.store.Delete(op.table, []byte(op.key))
	}
}
