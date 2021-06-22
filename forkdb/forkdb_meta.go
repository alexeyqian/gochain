package forkdb

import (
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
)

type MetaData struct {
	Head  core.Block // NOT USED
	Other int
}

func (fdb *ForkDB) initMetaData() {
	metaData := MetaData{
		Head:  core.Block{},
		Other: 0,
	}
	fdb.store.Put(metaTable, []byte(metaKey), entity.Serialize(metaData))
}

func (fdb *ForkDB) updateMetaData(meta *MetaData) {
	err := fdb.store.Put(metaTable, []byte(metaKey), entity.Serialize(*meta))
	if err != nil {
		panic(err)
	}
}

func (fdb *ForkDB) getMetaData() *MetaData {
	data, err := fdb.store.Get(metaTable, []byte(metaKey))
	if err != nil {
		panic("cannot get meta data")
	}
	var meta MetaData
	entity.Deserialize(&meta, data)
	return &meta
}
