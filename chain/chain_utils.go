package chain

import "github.com/alexeyqian/gochain/core"

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

func (c *Chain) reloadHead() {
	gpo, _ := c.sdb.GetGpo()
	b, err := c.GetBlockByNumber(gpo.BlockID)
	if err != nil {
		panic("chain: cannot find the head block")
	}
	c.cachedHead = b
}

/* should not be used, since it's mess up with undo operations
func (c *Chain) SetHead(b *core.Block) {
	gpo, _ := c.sdb.GetGpo()
	gpo.BlockID = b.ID
	gpo.BlockNum = b.Num
	c.sdb.UpdateGpo(gpo)
	c.cachedHead = b
}*/
