package chain

import (
	"encoding/gob"
	"time"

	"github.com/alexeyqian/gochain/config"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/forkdb"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/store"
	"github.com/alexeyqian/gochain/utils"
)

type Chain struct {
	lgr                 ledger.Ledger
	sdb                 *statusdb.StatusDB
	fdb                 *forkdb.ForkDB
	isGenesised         bool
	pendingTransactions []core.Transactioner
	cachedHead          *core.Block
}

func NewChain(lgr ledger.Ledger, storage store.Storage) *Chain {
	return &Chain{
		lgr:         lgr,
		sdb:         statusdb.NewStatusDB(storage),
		fdb:         forkdb.NewForkDB(storage),
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
	c.fdb.Open()
	if !c.isGenesised {
		c.genesis()
	}
}

func (c *Chain) Close() {
	c.lgr.Close()
	c.sdb.Close()
	c.fdb.Close()
}

func (c *Chain) Remove() {
	c.lgr.Remove()
	c.sdb.Remove()
	c.fdb.Remove()
	c.pendingTransactions = nil
}

// get block from forkdb first, if not found, then try it from ledger
func (c *Chain) GetBlockByNumber(num int) (*core.Block, error) {
	// get from forkdb
	// forkdb might contains multiple block with same num
	// we need to point out the branch we want to search, which is main branch here.
	fb, err := c.fdb.GetBlockByNumberFromBranch(c.Head().ID, uint64(num))
	if err != nil {
		return nil, err
	} else {
		return fb, nil
	}

	// get from ledger if cannot get block from forkdb
	// get block from ledger cannot by id, which is very slow, has to be by num
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

func (c *Chain) Head() *core.Block {
	if c.cachedHead == nil {
		gpo, _ := c.sdb.GetGpo()
		b, err := c.GetBlockByNumber(gpo.BlockID)
		if err != nil {
			panic("chain: cannot find the head block")
		}
		c.cachedHead = b
	}

	return c.cachedHead
}

func (c *Chain) SetHead(b *core.Block) {
	gpo, _ := c.sdb.GetGpo()
	gpo.BlockID = b.ID
	gpo.BlockNum = b.Num
	c.sdb.UpdateGpo(gpo)
	c.cachedHead = b
}
