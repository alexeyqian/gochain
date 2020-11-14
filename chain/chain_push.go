package chain

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/forkdb"
)

// This happens when two witness nodes are using same account
func maybeWarnMultipleProduction(fdb *forkdb.ForkDB, blockNumber uint64) {
	blocks := fdb.FetchBlocksByNumber(blockNumber)
	if len(blocks) == 0 {
		return // pass the check
	}

	fmt.Printf("Encontered block num collision at block %d\n", blockNumber)
	for _, b := range blocks {
		fmt.Printf("witness: %s, time: %d", b.Witness, b.CreatedOn)
	}
}

// each block create one new session/revision in undoable db
// so pop block only need to undo last session/revision in undoable db
func (c *Chain) PushBlock(b *core.Block) {
	maybeWarnMultipleProduction(c.fdb, b.Num)

	// forkdb.PushBlock will set the head pointing to longest chain in forkdb.
	c.fdb.PushBlock(b)
	newHead := c.fdb.Head()

	if newHead.PrevBlockId == c.Head().ID {
		c.startUndoSession()
		ok := c.ApplyBlock(b)
		if ok {
			// if everything goes right, then gpo's head block will be updated to new head
			// and all cached values will be reloaded
			// Chain's push undo session should leave the operation logs for future popblock,
			// only block becomes irriverible, then commit the block/session/revision
			// each block has exactly one session/revision
			c.pushUndoSession()
		} else {
			// undo all operations on statusdb during ApplyBlock()
			// also reload all cached values during undo
			c.undo()
			// usally undo operation hase nothing to do with forkdb
			// BUT here, the block is invalid, so we need to remove it
			c.fdb.RemoveBlock(b.ID)

		}
	} else {
		// if the head block from the longest chain does not build off of the current head,
		// then we might need to switch to new branch.
		// if the newly pushed block is the same height as head, nothing need to be done.
		// only switch forks if newHead is actually higher than headblock
		if newHead.BlockNum > c.Head().Num {
			c.switchBranch(newHead)
		}
	}
}

/*
// pop the head block of chain
func (c *Chain) PopBlock() {

}*/
