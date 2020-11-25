package chain

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/alexeyqian/gochain/config"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/forkdb"
	"github.com/alexeyqian/gochain/utils"
)

// GenerateNextBlock create next block
func (c *Chain) GenerateNextBlock() (*core.Block, error) {
	slottime := GetNextBlockTime(c.Gpo())
	name := GetNextWitness(c.Gpo(), c.Wso())
	return c.GenerateBlock(uint64(slottime), name, getCurrentWitnessPrivateKey())
}

// GenerateBlock create a new block without apply/push it into chain
func (c *Chain) GenerateBlock(when uint64, who string, privkey *ecdsa.PrivateKey) (*core.Block, error) {
	if when <= c.Gpo().Time || when != c.Gpo().Time+config.BlockInterval {
		return nil, fmt.Errorf("incorrect time")
	}

	witness, err := c.sdb.GetAccountByName(who)
	if err != nil || witness.PublicKey != utils.EncodePubKeyInPem(&privkey.PublicKey) {
		return nil, fmt.Errorf("incorret witness or key")
	}

	b := core.Block{
		// ID will be set later
		Num:         c.Gpo().BlockNum,
		PrevBlockId: c.Gpo().BlockId,
		CreatedOn:   when,
		Witness:     who,
	}

	// move pending txs to block
	c.movePendingTransactionsToBlock(&b)
	//c.MoveTxToBlock(&b)

	// set b.ID
	// set b.MerkleRoot
	return &b, nil
}

// each block create one new session/revision in undoable db
// so pop block only need to undo last session/revision in undoable db
func (c *Chain) PushBlock(newBlock *core.Block) {
	maybeWarnMultipleProduction(c.fdb, b.Num)

	// forkdb.PushBlock will set the head pointing to longest chain in forkdb.
	err := c.fdb.AppendBlock(newBlock)
	if err != nil {
		fmt.Errorf("invalid block, ignoring...")
		return
	}

	if newBlock.PrevBlockId == c.Head().ID {
		c.startUndoSession()
		ok := c.ApplyBlock(newBlock)
		if ok {
			// if everything goes right, then gpo's head block will be updated to new head
			// and all cached values will be reloaded
			// Chain's push undo session should leave the operation logs for future popblock,
			// only block becomes irriverible, then commit the block/session/revision
			// each block has exactly one session/revision
			// c.setHead(newBlock) - NOT NECESSARY, since head is reloaded from gpo in pushundosession
			c.pushUndoSession()
		} else {
			// undo all operations on statusdb during ApplyBlock()
			// also reload all cached values during undo
			c.undo()
			// usally undo operation hase nothing to do with forkdb
			// BUT here, the block is invalid, so we need to remove it
			// before remove block, the system will unpack the tx and store them in pending_tx_list
			c.fdb.RemoveBlock(newBlock)
			// c.setHead(previous head) -- NOT NECCESSARY
		}
	} else {
		if newBlock.Num > c.Head().Num {
			// if the new block is not build off from current chain's head block
			// and also has bigger number, means it just created a new longest branch
			// so we need to switch to the new longest branch
			c.switchBranch(newBlock)
		}
	}
}

/*
// pop the head block of chain
func (c *Chain) PopBlock() {

}*/

// save tx to pending array
func (c *Chain) PushTx(tx *core.Transactioner) error {
	//c.fdb.PushTx(tx)
	return nil
}

// remove tx from pending list, and append it to block
func (c *Chain) MoveTxToBlock(b *core.Block) {
	// get pending txes
	// validate transactions are not expired and valid
	// push tx to block
}

// when remove a block from forkdb, need to gite the tx back to pending list
func (c *Chain) GiveTxBackFromBlock() {

}

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

// TODO: to load from config file
func getCurrentWitnessPrivateKey() *ecdsa.PrivateKey {
	alicekey, _ := utils.GenerateKey()
	return alicekey
}
