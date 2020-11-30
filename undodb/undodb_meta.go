package undodb

import "github.com/alexeyqian/gochain/entity"

type MetaData struct {
	Revision            int
	SequencePerRevision int
}

func (udb *UndoableDB) initMetaData() {
	metaData := MetaData{
		Revision:            0,
		SequencePerRevision: 0,
	}
	udb.store.Put(metaTable, []byte(metaKey), entity.Serialize(metaData))
}

func (udb *UndoableDB) updateMetaData(meta *MetaData) {
	err := udb.store.Put(metaTable, []byte(metaKey), entity.Serialize(*meta))
	if err != nil {
		panic(err)
	}
}

func (udb *UndoableDB) getMetaData() *MetaData {
	data, err := udb.store.Get(metaTable, []byte(metaKey))
	if err != nil {
		panic("cannot get meta data")
	}
	var meta MetaData
	entity.Deserialize(&meta, data)
	return &meta
}
