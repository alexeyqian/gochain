package chain

import (
	"time"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/utils"
)

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

func (c *Chain) generateNextBlock(privkey string) *core.Block {
	slotTime := getNextSlotTime(c.Gpo())
	witness := getNextWitness(c.Gpo(), c.Wso())
	// getScheduledWitness(gpo, wso, num)
	return c.GenerateBlock(slotTime, witness, privkey)
}

func (c *Chain) GenerateBlock(slotTime int, witnessName string, privkey string) *core.Block {
	// TODO: ...
	slotNum := getSlotAtTime(c.Gpo(), slotTime)
	scheduledWitness := getScheduledWitness(c.Gpo(), c.Wso(), slotNum)
	if scheduledWitness != witnessName {
		panic("witness wrong")
	}

	witness := c.GetWitness(witnessName)
	if witness.signingKey != privKey.getPubKey() {
		panic("private key is incorrect")
	}

	// if tx.expiretime < slotTime continue

	// c.PushBlock(b)
	return b
}
