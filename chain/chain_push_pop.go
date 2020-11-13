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

func (c *Chain) PushBlock(b *core.Block) {
	// forkdb.PushBlock will set the head pointing to longest chain in forkdb.
	c.fdb.PushBlock(b)
	newHead := c.fdb.Head()

	maybeWarnMultipleProduction(c.fdb, b.Num)

	// if the head block from the longest chain does not build off of the current head,
	// then we need to switch to new branch.
	if newHead.PreviousBlockID != c.HeadBlockID() {
		swithBranch()
	}

	c.StartSession()
	ok := c.ApplyBlock(b)
	if ok {
		c.Commit() // Chain's commit should leave the operation logs for future popblock, only flush should remove the operation logs
		// update newHead as head of chain c.SetHead()
	} else {
		c.Undo()
		c.fdb.RemoveBlock(b.ID)
	}
}

// pop block from fork db, undo db operations,
// and pop transactions into poped transactions list.
func (c *Chain) PopBlock() {
	c.fdb.PopBlock()
	c.sdb.Undo()
	// TODO: append popped tx
	//c._popped_tx.insert( _self._popped_tx.begin(), head_block->transactions.begin(), head_block->transactions.end() );
}
