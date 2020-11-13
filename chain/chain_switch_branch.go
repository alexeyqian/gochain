package chain


func (c *Chain) switchBranch(newHead *core.Block){
	// if the newly pushed block is the same height as head, nothing need to be done.
	// only switch forks if newHead is actually higher than headblock
	if newHead.BlockNumber <= c.HeadBlockNumber() {
		return
	}

	fmt.Printf("switching to fork with head id: %s", newHead.BlockID)

	// get two branches, which shared same parent, not include parent, stored reversely.
	// such as: head_block_id(), ..., a3, a2, a1 (a3.block_num > a2.block_num > a1.block_num)
	// branch1 pointed by newHead, branch2 pointed by existing head
	// branch1 is longer, branch 2 is shorter
	branch1, branch2 = c.fdb.FetchBranchFrom(newHead.BlockID, c.HeadBlockID)
	commonAncestorBlockID := branch2[len(branch2) - 1].PreviousBlockID // also = branch1[len(branch1) - 1].PreviousBlockID
	
	// undo applied blocks of short branch2 (= current main branch) 
	// until we hit the commen ancestor block of these two branches
	for c.Head().ID != commonAncestorBlockID {
		// undo will restore gpo and gpo.blockid(chain head) to previous block
		// each undo will undo exactly one block (current head block)
		c.undo() 
		c.ReloadHead() // reset head to previous block id
		// TODO: append popped tx
		//c._popped_tx.insert( _self._popped_tx.begin(), head_block->transactions.begin(), head_block->transactions.end() );
	}

	// add blocks from longer branch based on common ancestor 
	// since blocks on shorter branch are already abandoned above.
	// push all blocks on the new fork
	// items in branch are reversely stored, block num is bigger in front
	for _, item := range branch1 {
		c.StartUndoSession()
		c.ApplyBlock( item )
		err := c.CommitUndoSession()
		
		if err != nil{
			fmt.Printf("error when switch branch, %s", err)
		}

		// remove the rest of branches.first from the fork_db, those blocks are invalid
		// for example: fork_branch is: new_head, ..., b5, b4, b3, b2, b1.
		// if exception happens while applying block b3, then b3, b4, b5, ... new_head will all be removed.
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
	}

	return true
}
