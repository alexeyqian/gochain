package chain

import (
	"github.com/alexeyqian/gochain/core"
)

func validateBeforApply(b *core.Block) bool {
	return true
}

func (c *Chain) ApplyBlock(b *core.Block) bool {
	if !validateBeforApply(b) {
		return false
	}

	// TODO:
	// update gpo's block num to b.Num
	// update gpo's witness to b.Witness

	for _, tx := range b.Transactions {
		// We do not need to push the undo state for each transaction
		// because they either all apply and are valid or the
		// entire block fails to apply.  We only need an "undo" state
		// for transactions when validating broadcast transactions or when building a block.
		// Lots of gpo read / update inside operation appliers.
		// Use global properties inside: block_head_time() and block_head_num()
		// TODO: review, should prevent all functions to update global properties inside
		// TODO: extract the used global properties out, so we know it will not mess up with global properties.
		err := tx.Apply(c.sdb) // gpo might be updated during tx.Apply()
		if err != nil {
			// TODO: move tx to invalid tx
			return false
		}
	}

	// CreateBlockSummary(b) for fast query

	// update gpo
	// TODO ...

	return true
}
