package chain

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/forkdb"
)

// This happens when two witness nodes are using same account
func maybeWarnMultipleProduction(fdb *forkdb.ForkDB, blockNumber uint64) {
	blocks := fdb.FetchBlocksByNumber(blockNumber)
	if len(blocks) <= 1 {
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
	// forkdb.PushBlock will set the head pointing to longest chain in forkdb.
	c.fdb.PushBlock(b)
	newHead := c.fdb.Head()

	maybeWarnMultipleProduction(c.fdb, b.Num)

	if newHead.PrevBlockId == c.Head().ID {
		c.startUndoSession()
		ok := c.ApplyBlock(b)
		if ok {
			// Chain's push undo session should leave the operation logs for future popblock,
			// only block becomes irriverible, then commit the block/session/revision
			// one block has one session/revision
			c.pushUndoSession()
			// if everything goes right, then update newHead as head of chain
			c.SetHead(newHead)
		} else {
			c.undo()
			c.fdb.PopBlock()        // restore head to previous block
			c.fdb.RemoveBlock(b.ID) // remove block data from forkdb
		}
	} else {
		// if the head block from the longest chain does not build off of the current head,
		// then we need to switch to new branch.
		c.switchBranch(newHead)
	}
}

// pop the head block of chain
func (c *Chain) PopBlock() {

}
