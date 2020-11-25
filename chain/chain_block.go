package chain

import (
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/utils"
)

// GetBlockByNumber get block from forkdb first, if not found, then try it from ledger
func (c *Chain) GetBlockByNumber(num int) (*core.Block, error) {
	// get from forkdb
	// forkdb might contains multiple block with same num
	// we need to point out the branch we want to search, which is main branch here.
	fb, err := c.fdb.GetBlockByNumberFromBranch(c.Head().ID, num)
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
