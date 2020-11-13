package chain

import (
	"time"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/utils"
)

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

// TODO: need to rewrite
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

// c.Gpo()/SetGpo() used cached gpo
// c.Wso()/ SetWso() use cached wso
