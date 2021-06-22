package forkdb

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
)

// GetBlocks get blocks
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

// GetBlockByID get block by id
func (fdb *ForkDB) GetBlockByID(id string) (*core.Block, error) {
	data, err := fdb.store.Get(branchTable, []byte(id))
	if err != nil {
		return nil, err
	}
	var e core.Block
	entity.Deserialize(&e, data)
	return &e, nil
}

// GetBlocksByNumber might have multiple blocks with same block num
func (fdb *ForkDB) GetBlocksByNumber(num int) []*core.Block {
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

func (fdb *ForkDB) GetBlockByNumberFromBranch(headID string, num int) (*core.Block, error) {
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
