package chain

import "github.com/alexeyqian/gochain/core"

func validateBeforApply(b *core.Block) bool {
	return false
}

func (c *Chain) ApplyBlock(b *core.Block) bool {
	if !validateBeforApply(b) {
		return false
	}

	// TODO:
	// update gpo's block num to b.Num
	// update gpo's witness to b.Witness

}
