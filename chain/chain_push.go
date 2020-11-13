package chain

import (
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/forkdb"
)

// TODO: use cached gpo instead of read from database to improve perfomance
func (c *Chain) HeadBlockID() string{
	// should use c.Gpo().BlockID
	return c.sdb.GetGpo().BlockID
}

func (c *Chain) HeadBlockNum() uint64{
	return c.sdb.GetGpo().BlockNumber
}

// This happens when two witness nodes are using same account
func maybeWarnMultipleProduction(fdb *forkdb.ForkDB, uint64 blockNumber ){
	blocks := fdb.FetchBlocksByNumber( blockNumber )
	if len(blocks) <= 1 {
		return // pass the check
 	}
 
 	fmt.Printf("Encontered block num collision at block %d\n", blockNumber)
 	for _, b := range blocks{
		fmt.Printf("witness: %s, time: %d", b.Witness, b.CreatedOn)
 	}   
}

// pop block from fork db, undo db operations, 
// and pop transactions into poped transactions list. 
func (c *Chain) PopBlock()
      // _pending_tx_session.reset();
      // save the head block so we can recover its transactions
	  _, err := c.FetchBlockByID( c.HeadBlockID() )
	  if err != nil {
		panic("there are no blocks to pop.")
	  }  

      c.fdb.PopBlock();
      c.sdb.Undo();

      //c._popped_tx.insert( _self._popped_tx.begin(), head_block->transactions.begin(), head_block->transactions.end() );

}

// the return value indicates if branch switch happens
func (c *Chain) PushBlock(b core.Block) bool {
 // forkdb.PushBlock will set the head pointing to longest chain in forkdb.
 c.fdb.PushBlock(b) 
 newHead := c.fdb.Head()
 
 c.maybeWarnMultipleProduction(c.fork, b.Num)

 // if the head block from the longest chain does not build off of the current head,
 // then we need to switch to new branch.
 if newHead.PreviousBlockID != c.HeadBlockID() {
	 // if the newly pushed block is the same height as head, we get head back in newHead
	 // only switch forks if newHead is actually higher than headblock
	 if newHead.BlockNumber <= c.HeadBlockNumber() {
		 return false
	 }

	 fmt.Printf("switching to fork with head id: %s", newHead.BlockID)

	 // get two branches, which shared same parent, not include parent, stored reversely.
	 // such as: head_block_id(), ..., a3, a2, a1 (a3.block_num > a2.block_num > a1.block_num)
	 // branch1 pointed by newHead, branch2 pointed by existing head
	 // branch1 is longer, branch 2 is shorter
	 branch1, branch2 = c.fdb.FetchBranchFrom(newHead.BlockID, c.HeadBlockID)
	 commonAncestorBlockID := branch2[len(branch2) - 1].PreviousBlockID // also = branch1[len(branch1) - 1].PreviousBlockID
	 
	 // pop blocks until we hit the commen ancestor block of these two branches
	 // abandon blocks on shorter branch	 
	 for c.HeadBlockID != commonAncestorBlockID {
		 // pop block from fork db, undo db operations, 
		 // and pop transactions into poped transactions list.
		 c.PopBlock() 
		 //retrive and validate
		 // c.softfork.PopBlock()
		 // c.UndoPopedBlock()
		 // add tx in popped block into pending or popped tx list
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

 c.StartSession()
 c.ApplyBlock(b)
 err := c.SubmitSession()

 if err != nil {
	 c.softfork.Remove(b.ID)
 }

 return false
}
