package forkdb

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
)

// AppendBlock insert new valid block into pending block tree
// multiple branches are stored in same table as a tree
// push a block into linked table if it can be linked to one of the branches in the table
// otherwise the block will be put into unlinked table
func (fdb *ForkDB) AppendBlock(b *core.Block) error {
	// validate block before push, expired, max_depth of fork
	//if b.Num-fdb.Head().Num > maxBranchingDepth {
	//	panic("reach max branching depth, head: %d, current: %d", fdb.Head().Num, b.Num)
	//}

	// check duplication of BlockID before insert
	_, err := fdb.GetBlockByID(b.ID)
	if err == nil {
		return fmt.Errorf("cannot store save block twice")
	}
	// check if block can be linked to current branches, if cannot, then put it into unlinked pool
	if !fdb.linkable(b) {
		// TODO: process unlinkable later, in memory only or persisted in db
		return fmt.Errorf("block is not linkable")
	}

	// save block into linked table (branch)
	fdb.createBlock(b)

	// switch head if num is bigger.
	//if b.Num > fdb.Head().Num {
	//	fdb.SetHead(b)
	//}

	return nil
}

// restore tx to pending tx list
// delete block from storage
func (fdb *ForkDB) RemoveBlock(b *core.Block) {
	// save tx back to pending tx list

	fdb.deleteBlock(b.ID)
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

	firstBranchHead, err = fdb.GetBlockByID(first)
	if err != nil {
		panic(fmt.Sprintf("cannot find block: %s", first))
	}
	secondBranchHead, err = fdb.GetBlockByID(second)
	if err != nil {
		panic(fmt.Sprintf("cannot find block: %s", second))
	}

	for firstBranchHead.Num > secondBranchHead.Num {
		firstBranch = append(firstBranch, firstBranchHead)
		firstBranchHead, err = fdb.GetBlockByID(firstBranchHead.PrevBlockId)
		if err != nil {
			panic(fmt.Sprintf("cannot find block: %s", firstBranchHead.PrevBlockId))
		}
	}

	for secondBranchHead.Num > firstBranchHead.Num {
		secondBranch = append(secondBranch, secondBranchHead)
		secondBranchHead, err = fdb.GetBlockByID(secondBranchHead.PrevBlockId)
		if err != nil {
			panic(fmt.Sprintf("cannot find block: %s", secondBranchHead.PrevBlockId))
		}
	}

	for firstBranchHead.PrevBlockId != secondBranchHead.PrevBlockId {
		firstBranch = append(firstBranch, firstBranchHead)
		firstBranchHead, err = fdb.GetBlockByID(firstBranchHead.PrevBlockId)
		if err != nil {
			panic(fmt.Sprintf("cannot find block: %s", firstBranchHead.PrevBlockId))
		}
		secondBranch = append(secondBranch, secondBranchHead)
		secondBranchHead, err = fdb.GetBlockByID(secondBranchHead.PrevBlockId)
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
