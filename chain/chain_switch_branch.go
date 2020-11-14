package chain

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
)

// during the switch branch process, the chain's head might not synced with forkdb.head
// but after the process, they should be consistent.
func (c *Chain) switchBranch(newHead *core.Block) {
	fmt.Printf("switching to fork with head id: %s", newHead.ID)

	// get two branches, which shared same parent, not include parent, stored reversely.
	// such as: head_block_id(), ..., a3, a2, a1 (a3.block_num > a2.block_num > a1.block_num)
	// branch1 pointed by newHead, branch2 pointed by existing head
	// branch1 is longer, branch 2 is shorter
	branch1, branch2 = c.fdb.FetchBranchFrom(newHead.ID, c.Head().ID)
	commonAncestorBlockID := branch2[len(branch2)-1].PreviousBlockID // also = branch1[len(branch1) - 1].PreviousBlockID

	// undo applied blocks of short branch2 (= current main branch)
	// until we hit the commen ancestor block of these two branches
	for c.Head().ID != commonAncestorBlockID {
		// undo will restore gpo and gpo.blockid(chain head) to previous block
		// each undo will undo exactly one block (current head block)
		// undo will also reload all cached values
		// undo() has nothing to do with forkdb
		c.undo()
		// ATTENTION: becareful not to do any sdb/udb data modification between undo() and next start new session()
		// since it will cause the undo operation not cleanly done,
		// NEED wrap modifactions between disableRevision() and enableRevision() if you really need to do so

		// TODO: append popped tx
		// pending transactions stored in forkdb's transaction table
		// with field: isInBlock to indicate if it's been used or not
		// popped out tx will set isInBlock to false
		//c._popped_tx.insert( _self._popped_tx.begin(), head_block->transactions.begin(), head_block->transactions.end() );
	}

	// add blocks from longer branch based on common ancestor
	// since blocks on shorter branch are already undoed by above code.
	// try to apply all blocks on the new branch
	// items in branch are reversely stored, block num is bigger in front
	for index, item := range branch1 {
		c.startUndoSession()
		ok := c.ApplyBlock(item)
		if ok {
			c.pushUndoSession()
		} else {
			fmt.Printf("error when switch branch")
			// UNDO switch
			// ATTENTION: switch to new branch1 FAILED
			// RESET to old branch2, since branch 1 is not working

			// remove the rest of branch1 from the fork_db, those blocks are invalid
			// for example: fork_branch is: new_head, ..., b5, b4, b3, b2, b1.
			// if exception happens while applying block b3, then b3, b4, b5, ... new_head will all be removed.
			// error happens at branch1[index], so also need to remove index itself
			for j := len(branch1) - 1; j >= index; j-- {
				// remove block will return tx inside the block back to pending tx list
				c.fdb.RemoveBlock(branch1[j])
			}

			// undo all just applied blocks from the bad branch1
			// ATTENTION: becareful of the partially applied index block (current block)
			for c.Head().ID != branch1[len(branch1)-1].PreviousBlockID {
				// undo() will update chain's head (both in gpo can cached value)
				c.undo() // pop block()
			}

			// restore all blocks from the good fork
			for k := len(branch2) - 1; k >= 0; k++ {
				c.startUndoSession()
				ok := c.ApplyBlock(branch2[k])
				if !ok {
					panic("serious issue happend")
				}
				c.pushUndoSession()
			}
			break // break the loop
		}
	}

	// c.SetHead(newHead) - NOT NECCESSARY, ALREADY DONE IN apply block by update gpo
}
