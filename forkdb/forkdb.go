package forkdb

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/store"
)

// TODO: IMPORTANT: remove head from forkdb, so we don't need to sync it with chain's head block.

const maxBranchingDepth = 100

const metaTable = "_forkmeta_"
const metaKey = "_forkmeta_key_"
const branchTable = "branch"

type ForkDB struct {
	store store.Storage
	//head  core.Block // NOT USED. in memory cached head, for fast access
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

/*
func (fdb *ForkDB) SetHead(b *core.Block) {
	meta := fdb.getMetaData()
	meta.Head = *b
	fdb.updateMetaData(meta)
	fdb.head = *b
}

func (fdb *ForkDB) Head() *core.Block {
	if fdb.head.ID == "" {
		//meta := fdb.getMetaData()
		//return &meta.Head
		panic("must set head before use it")
	}

	return &fdb.head // return in memory cached value
}
*/

func (fdb *ForkDB) reset() {
	fdb.store.DeleteBucket(branchTable)
	fdb.store.CreateBucket(branchTable)
	fdb.store.Delete(metaTable, []byte(metaKey))
	fdb.initMetaData()
}

func (fdb *ForkDB) GetBlocks() []*core.Block {
	var res []*core.Block
	items, _ := fdb.store.GetAll(branchTable)
	for _, pair := range items {
		var e core.Block
		entity.Deserialize(&e, pair.Value)
		res = append(res, &e)
	}
	return res
}

func (fdb *ForkDB) GetBlockByID(id string) (*core.Block, error) {
	data, err := fdb.store.Get(branchTable, []byte(id))
	if err != nil {
		return nil, err
	}
	var e core.Block
	entity.Deserialize(&e, data)
	return &e, nil
}

// might have multiple blocks with same block num
func (fdb *ForkDB) GetBlocksByNumber(num uint64) []*core.Block {
	var blocks []*core.Block

	items := fdb.GetBlocks()
	for _, item := range items {
		if item.Num == num {
			blocks = append(blocks, item)
		}
	}
	return blocks
}

func (fdb *ForkDB) createBlock(e *core.Block) error {
	if !entity.HasID(e) {
		return fmt.Errorf("create: entity doesn't have ID")
	}
	return fdb.store.Put(branchTable, []byte(entity.GetID(e)), entity.Serialize(e))
}

func (fdb *ForkDB) deleteBlock(id string) error {
	// before remove, should unpack and store all tx inside the block
	return fdb.store.Delete(branchTable, []byte(id))
}

func (fdb *ForkDB) GetBlockByNumberFromBranch(headID string, num uint64) (*core.Block, error) {
	blocks := fdb.GetBlocksByNumber(num)
	if len(blocks) == 1 { // found exact one
		return blocks[0], nil
	}

	// loop through the branch
	count := 0
	id := headID
	for {
		block, err := fdb.GetBlockByID(id)
		if err != nil {
			return nil, fmt.Errorf("forkdb: block not found, id:%s", id)
		}
		if block.Num == num {
			return block, nil
		}
		id = block.PrevBlockId

		count++
		if count >= maxBranchingDepth {
			return nil, fmt.Errorf("forkdb: block not found, num: %d", num)
		}
	}
}
