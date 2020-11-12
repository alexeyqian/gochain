package forkdb

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
)

// multiple branches are stored in same table
// push a block into linked table if it can be linked to one of the branches in the table
// otherwise the block will be put into unlinked table
// after adding to linked table, the head will updated to point to longest branch if necessary
func (fdb *ForkDB) PushBlock(b *core.Block) error {
	// validate block before push, expired, max_depth of fork
	if b.Num-fdb.Head().Num > maxBranchingDepth {
		panic("reach max branching depth, head: %d, current: %d", fdb.Head().Num, b.Num)
	}

	// check duplication of BlockID before insert
	_, err := fdb.GetBlock(b.ID)
	if err == nil {
		return fmt.Errorf("cannot store save block twice")
	}
	// check if block can be linked to current branches, if cannot, then put it into unlinked pool
	if !fdb.linkable(b) {
		// TODO: process unlinkable later, in memory only or persisted in db
		return fmt.Errorf("block is not linkable")
	}

	// pusht block into linked table (branch)
	fdb.CreateBlock(b)

	// switch head if num is bigger.
	if b.Num > fdb.Head().Num {
		fdb.SetHead(b)
	}

	return nil
}

// forkdb should always have at lease one item,
// every time the database opens, it will starts with the last irriversable block.
// forkdb pop block is  part prcess of chain's pop block
func (fdb *ForkDB) PopBlock() {
	var err error
	_, err = fdb.GetBlock(fdb.Head().ID)
	if err != nil {
		panic(fmt.Sprintf("cannot find the fork db item: %s", fdb.Head().ID))
	}

	// TODO: remove previous head block
	// check if it's still the longest branch?

	var prev *core.Block
	prev, err = fdb.GetBlock(fdb.Head().PrevBlockId)
	if err != nil {
		panic(fmt.Sprintf("cannot find the ford db item: %s", fdb.Head().PrevBlockId))
	}
	fdb.SetHead(prev)
}

func (fdb *ForkDB) linkable(b *core.Block) bool {
	blocks := fdb.GetBlocks()
	for _, item := range blocks {
		if b.PrevBlockId == item.ID {
			return true
		}
	}

	return false
}
