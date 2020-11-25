package chain

import (
	"github.com/alexeyqian/gochain/config"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
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

// Gpo get gpo of current chain
func (c *Chain) Gpo() *entity.Gpo {
	gpo, _ := c.sdb.GetGpo()
	return gpo
}

// Wso get wso for current chain
func (c *Chain) Wso() *entity.Wso {
	wso, _ := c.sdb.GetWso()
	return wso
}

// GetSlotAtTime find slot for a specific time
func GetSlotAtTime(gpo *entity.Gpo, when int) int {
	if uint64(when) < gpo.Time {
		return 0
	}
	return (uint64(when) - gpo.Time) / config.BlockInterval
}

// GetNextBlockTime find next block time
func GetNextBlockTime(gpo *entity.Gpo) int64 {
	return int64(gpo.Time) + int64(config.BlockInterval*1)
}

// GetNextWitness find next witness
func GetNextWitness(gpo *entity.Gpo, wso *entity.Wso) string {
	aslot := (gpo.Time + 1*config.BlockInterval) / config.BlockInterval
	return wso.CurrentWitnesses[int(aslot)%len(wso.CurrentWitnesses)]
}
