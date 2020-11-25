package chain

import (
	"encoding/gob"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/forkdb"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/store"
)

// TODO: New Design:
// gpo and wso should be only persistent once at end of block apply,
// but it can be modified in memory during apply block process.

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

// TODO: IMPORTANT:
// setRevision() // revision equals current block num
