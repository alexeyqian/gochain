package forkdb

import (
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/store"
)

// TODO: Should ForkDB only need memory storage??

const maxBranchingDepth = 100

const metaTable = "_forkmeta_"
const metaKey = "_forkmeta_key_"
const branchTable = "branch"
const txTable = "tx"

type ForkDB struct {
	store store.Storage
}

func NewForkDB(storage store.Storage) *ForkDB {
	return &ForkDB{
		store: storage,
	}
}

func (fdb *ForkDB) Open() error {
	err := fdb.store.Open()
	if err != nil {
		return err
	}

	if !fdb.store.HasBucket(metaTable) {
		err := fdb.store.CreateBucket(metaTable)
		if err != nil {
			return err
		}
		fdb.initMetaData()

		err = fdb.store.CreateBucket(branchTable)
		if err != nil {
			return err
		}

		err = fdb.store.CreateBucket(txTable)
		if err != nil {
			return err
		}
	}

	// TODO: set last irreversible block as head ...

	return nil
}

func (fdb *ForkDB) Close() error {
	return fdb.store.Close()
}

func (fdb *ForkDB) Remove() error {
	return fdb.store.Remove()
}

func (fdb *ForkDB) reset() {
	fdb.store.DeleteBucket(branchTable)
	fdb.store.CreateBucket(branchTable)
	fdb.store.Delete(metaTable, []byte(metaKey))
	fdb.initMetaData()
}

func (fdb *ForkDB) ResetTo(b *core.Block) {
	fdb.reset()
	fdb.createBlock(b)
}
