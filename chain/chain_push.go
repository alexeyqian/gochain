package chain


// This happens when two witness nodes are using same account
func maybeWarnMultipleProduction(fork *softfork.SoftFork, uint64 blockNumber ){
	blocks := fork.FetchBlocksByNumber( blockNumber )
	if len(blocks) <= 1 {
	return // pass the check
 	}
 
 	fmt.Printf("Encontered block num collision at block %d\n", blockNumber)
 	for _, b := range blocks{
		fmt.Printf("witness: %s, time: %d", b.Witness, b.CreatedOn)
 	}   
}

// the return value indicates if branch switch happens
func (c *Chain) PushBlock(b core.Block) bool {
 // softfork.PushBlock will return the head block of current longest chain in softfork.
 newHead := c.softfork.PushBlock(b)
 
 c.maybeWarnMultipleProduction(c.fork, b.Num)

 //If the head block from the longest chain does not build off of the current head,
 // then we need to switch to new branch.
 if newHead.PreviousBlockID != c.HeadBlock().ID {
	 //If the newly pushed block is the same height as head, we get head back in newHead
	 //Only switch forks if newHead is actually higher than headblock
	 if newHead.BlockNumber <= c.HeadBlock().BlockNumber() {
		 return false
	 }

	 fmt.Printf("switching to fork with head id: %s", newHead.BlockID)

	 // get two branches, which shared same parent, not include parent, stored reversely.
	 // such as: head_block_id(), ..., a3, a2, a1 (a3.block_num > a2.block_num > a1.block_num)
	 // branch1 pointed by newHead, branch2 pointed by existing head
	 branch1, branch2 = c.softfork.fetch_branch_from(newHead.BlockID, c.HeadBlock().ID)

	 // pop blocks until we hit the commen ancestor block of these two branches
	 // abandon blocks on shorter branch
	 // branch2.back().previous is pointing to the common ancestor
	 while( c.HeadBlock().ID != branch2.back()->data.previous ){
		 //pop_block(); // pop block from fork db, undo , and pop transactions into poped transactions list.
		 //retrive and validate
		 // c.softfork.PopBlock()
		 // c.UndoPopedBlock()
		 // add tx in popped block into pending or popped tx list
	 }

	 // add blocks from longer branch based on common ancestor 
	 // since blocks on shorter branch are already abandoned above.
	 // push all blocks on the new fork
	 // items in branch are reversely stored
	 for _, item range branch1{
		 start_undo_session( true )
		 apply_block( (*ritr)->data, skip )
		 err := session.push()
		 
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

 c.StartSession()
 c.ApplyBlock(b)
 err := c.SubmitSession()

 if err != nil {
	 c.softfork.Remove(b.ID)
 }

 return false
}
