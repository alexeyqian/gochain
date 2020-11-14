package forkdb

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
)

func (fdb *ForkDB) AppendPacketTx(e *core.PacketTx) error {
	_, err := fdb.GetPacketTx(e.ID)
	if err == nil {
		return fmt.Errorf("PacketTx already exist, cannot append. id: %s", e.ID)
	}

	return fdb.createPacketTx(e)
}

func (fdb *ForkDB) RemovePacketTx(e *core.PacketTx) error {
	return fdb.store.Delete(txTable, []byte(e.ID))
}

func (fdb *ForkDB) createPacketTx(e *core.PacketTx) error {
	if e.ID == "" {
		return fmt.Errorf("create: entity doesn't have ID")
	}
	return fdb.store.Put(txTable, []byte(e.ID), e.Data)
}

func (fdb *ForkDB) GetPacketTxs() []*core.PacketTx {
	var res []*core.PacketTx
	items, _ := fdb.store.GetAll(txTable)
	for _, pair := range items {
		var e core.PacketTx
		entity.Deserialize(&e, pair.Value)
		res = append(res, &e)
	}
	return res
}

func (fdb *ForkDB) GetPacketTx(id string) (*core.PacketTx, error) {
	e, err := fdb.store.Get(txTable, []byte(id))
	if err != nil {
		return nil, err
	}
	return &e, nil
}
