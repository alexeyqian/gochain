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

// PopBlock is part of chain.Popblock
// forkdb should always have at lease one item,
// every time the database opens, it will starts with the last irriversable block.
// pop current head block, and set the it's previous block as head
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

// FetchBranchFrom get two branches, which shared same parent (not include parent), stored reversely.
// should return sth like:
// branch1: first, ..., a3, a2, a1.
// branch2: second, ..., b3, b2, b1.
// and a1.previous_id() == b1.previous_id() if found common ancestor
// This function gets 2 branches leading back to the most recent common ancestor.
func (fdb *ForkDB) FetchBranchFrom(first string, second string) ([]*core.Block, []*core.Block) {
	var err error
	var firstBranch []*core.Block
	var firstBranchHead *core.Block
	var secondBranch []*core.Block
	var secondBranchHead *core.Block

	firstBranchHead, err = fdb.GetBlock(first)
	if err != nil {
		panic(fmt.Sprintf("cannot find block: %s", first))
	}
	secondBranchHead, err = fdb.GetBlock(second)
	if err != nil {
		panic(fmt.Sprintf("cannot find block: %s", second))
	}

	for firstBranchHead.Num > secondBranchHead.Num {
		firstBranch = append(firstBranch, firstBranchHead)
		firstBranchHead, err = fdb.GetBlock(firstBranchHead.PrevBlockId)
		if err != nil {
			panic(fmt.Sprintf("cannot find block: %s", firstBranchHead.PrevBlockId))
		}
	}

	for secondBranchHead.Num > firstBranchHead.Num {
		secondBranch = append(secondBranch, secondBranchHead)
		secondBranchHead, err = fdb.GetBlock(secondBranchHead.PrevBlockId)
		if err != nil {
			panic(fmt.Sprintf("cannot find block: %s", secondBranchHead.PrevBlockId))
		}
	}

	for firstBranchHead.PrevBlockId != secondBranchHead.PrevBlockId {
		firstBranch = append(firstBranch, firstBranchHead)
		firstBranchHead, err = fdb.GetBlock(firstBranchHead.PrevBlockId)
		if err != nil {
			panic(fmt.Sprintf("cannot find block: %s", firstBranchHead.PrevBlockId))
		}
		secondBranch = append(secondBranch, secondBranchHead)
		secondBranchHead, err = fdb.GetBlock(secondBranchHead.PrevBlockId)
		if err != nil {
			panic(fmt.Sprintf("cannot find block: %s", secondBranchHead.PrevBlockId))
		}
	}

	if firstBranchHead.ID != "" && secondBranchHead.ID != "" {
		firstBranch = append(firstBranch, firstBranchHead)
		secondBranch = append(secondBranch, secondBranchHead)
	}

	return firstBranch, secondBranch
}
