package chain

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
)

func (c *Chain) switchBranch(newHead *core.Block) {
	fmt.Printf("switching to fork with head id: %s", newHead.BlockID)

	// get two branches, which shared same parent, not include parent, stored reversely.
	// such as: head_block_id(), ..., a3, a2, a1 (a3.block_num > a2.block_num > a1.block_num)
	// branch1 pointed by newHead, branch2 pointed by existing head
	// branch1 is longer, branch 2 is shorter
	branch1, branch2 = c.fdb.FetchBranchFrom(newHead.BlockID, c.HeadBlockID)
	commonAncestorBlockID := branch2[len(branch2)-1].PreviousBlockID // also = branch1[len(branch1) - 1].PreviousBlockID

	// undo applied blocks of short branch2 (= current main branch)
	// until we hit the commen ancestor block of these two branches
	for c.Head().ID != commonAncestorBlockID {
		// undo will restore gpo and gpo.blockid(chain head) to previous block
		// each undo will undo exactly one block (current head block)
		// undo will also reload all cached values
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
	for i, item := range branch1 {
		c.startUndoSession()
		ok := c.ApplyBlock(item)
		if ok {
			c.pushUndoSession()
		} else {
			fmt.Printf("error when switch branch")
			// ATTENTION: switch to new branch1 FAILED
			// RESET to old branch2, since branch 1 is not working

			// remove the rest of branch1 from the fork_db, those blocks are invalid
			// for example: fork_branch is: new_head, ..., b5, b4, b3, b2, b1.
			// if exception happens while applying block b3, then b3, b4, b5, ... new_head will all be removed.

			// UNDO switch

			/*
				while( ritr != branches.first.rend() )
				{
					_fork_db.remove( (*ritr)->data.id() );
					++ritr;
				}
				// reset head back to head_block_id()
				_fork_db.set_head( branches.second.front() );

				// pop all blocks from the bad fork
				while( head_block_id() != branches.second.back()->data.previous )
					pop_block();

				// restore all blocks from the good fork
				for( auto ritr = branches.second.rbegin(); ritr != branches.second.rend(); ++ritr )
				{
					auto session = start_undo_session( true );
					apply_block( (*ritr)->data, skip );
					session.push();
				}
			*/
		}

	}
}
