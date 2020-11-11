package chain

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/alexeyqian/gochain/config"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/store"
	"github.com/alexeyqian/gochain/utils"
)

type Chain struct {
	lgr                 ledger.Ledger
	sdb                 *statusdb.StatusDB
	isGenesised         bool
	pendingTransactions []core.Transactioner
}

func NewChain(lgr ledger.Ledger, storage store.Storage) *Chain {
	return &Chain{
		lgr:         lgr,
		sdb:         statusdb.NewStatusDB(storage),
		isGenesised: false,
	}
}

func (c *Chain) Open() {
	// register gob
	// TODO: move to somewhere
	gob.Register(&core.CreateAccountTransaction{})
	gob.Register(&core.CreateArticleTransaction{})
	gob.Register(&core.CreateCommentTransaction{})
	gob.Register(&core.TransferCoinTransaction{})
	gob.Register(&core.VoteTransaction{})

	c.lgr.Open()
	c.sdb.Open()
	if !c.isGenesised {
		c.genesis()
	}
}

func (c *Chain) Close() {
	c.lgr.Close()
	c.sdb.Close()
}

func (c *Chain) Remove() {
	c.lgr.Remove()
	c.sdb.Remove()
	c.pendingTransactions = nil
}

func (c *Chain) GetBlock(num int) (*core.Block, error) {
	// TODO: use cache to speed up reading
	bdata, err := c.lgr.Read(num)
	if err != nil {
		return nil, err
	}
	var b core.Block
	utils.Deserialize(&b, bdata)
	return &b, nil
}

// TODO: move to node
func (c *Chain) BroadcastTx(tx core.Transactioner) {
	// broadcast to other peers
}

// TODO: move to node
func (c *Chain) ReceiveTx(tx core.Transactioner) error {
	// check if already has the tx
	// validate tx: two validations, fast validate and full validate

	return nil
}

func (c *Chain) GetPendingTx() []core.Transactioner {
	return c.pendingTransactions
}

func (c *Chain) movePendingTransactionsToBlock(b *core.Block) {
	i := 0
	for _, tx := range c.pendingTransactions {
		if i >= core.MaxTransactionsInBlock {
			break
		}
		b.AddTransaction(tx)
		i++
	}

	if len(c.pendingTransactions) > core.MaxTransactionsInBlock {
		c.pendingTransactions = c.pendingTransactions[core.MaxTransactionsInBlock:]
	}
}

func (c *Chain) AddPendingTx(tx core.Transactioner) error {
	err := tx.Validate(c.sdb)
	if err == nil {
		c.pendingTransactions = append(c.pendingTransactions, tx)
	}
	return err
}

func (c *Chain) GenerateBlock() *core.Block {
	var b core.Block
	var gpo *entity.Gpo

	gpo, _ = c.sdb.GetGpo()

	b.ID = utils.CreateUuid()
	b.PrevBlockId = gpo.BlockId
	b.Num = gpo.BlockNum + uint64(1)
	b.CreatedOn = uint64(time.Now().Unix())
	c.movePendingTransactionsToBlock(&b)

	for _, tx := range b.Transactions {
		err := tx.Apply(c.sdb) // gpo might be updated during tx.Apply()
		if err != nil {
			// move tx to invalid tx
			//
		}
	}

	gpo, _ = c.sdb.GetGpo()
	gpo.BlockId = b.ID
	gpo.BlockNum = b.Num
	gpo.Time = b.CreatedOn
	gpo.Supply += core.AmountPerBlock
	c.sdb.UpdateGpo(gpo)

	//fmt.Printf("arrive here: %+v", b)
	// append new block to lgr
	c.lgr.Append(utils.Serialize(b))

	return &b
}

func (c *Chain) genesis() {

	// update global status
	var gpo entity.Gpo
	gpo.ID = statusdb.GpoKey
	gpo.BlockId = core.BlockZeroId
	gpo.BlockNum = 0
	gpo.Witness = core.InitWitness
	gpo.Time = core.GenesisTime
	gpo.Supply = core.InitAmount
	//fmt.Printf("creating gpo: %+v", gpo)
	c.sdb.CreateGpo(&gpo)

	var wso entity.Wso
	wso.MaxBlockSize = config.DefaultBlockSize
	wso.AccountCreationFee = config.DefaultAccountCreationFee
	c.sdb.CreateWso(&wso)

	// update chain database
	var acc entity.Account
	acc.ID = utils.CreateUuid() // TODO: should be public key string
	acc.Name = core.InitWitness
	acc.CreatedOn = core.GenesisTime
	acc.Coin = core.InitAmount
	err := c.sdb.CreateAccount(&acc)
	if err != nil {
		panic(err)
	}

	// update lgr, create a dummy block 0
	b := core.Block{ID: core.BlockZeroId, Num: 0, CreatedOn: core.GenesisTime, Witness: core.InitWitness}
	c.lgr.Append(utils.Serialize(b))
}

// This happens when two witness nodes are using same account
func maybe_warn_multiple_production(itemp_database& db, uint64 height ){
   auto blocks = db.fork_db().fetch_block_by_number( height );
   if( blocks.size() > 1 )
   {
      vector< std::pair< account_name_type, fc::time_point_sec > > witness_time_pairs;
      for( const auto& b : blocks )
      {
         witness_time_pairs.push_back( std::make_pair( b->data.witness, b->data.timestamp ) );
      }

      ilog( "Encountered block num collision at block ${n} due to a fork, witnesses are: ${w}", ("n", height)("w", witness_time_pairs) );
   }
   return;
}


// the return value indicates if branch switch happens
func (c *Chain) PushBlock(b core.Block) bool {
	// softfork.PushBlock will return the head block of current longest chain in softfork.
	newHead := c.softfork.PushBlock(b)
	//TODO:  validate multiple production

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
