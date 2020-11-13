package chain

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
)

func (c *Chain) ApplyTx(tx core.Transactioner) bool {
	if !tx.Validate() {
		return false
	}

	// duplication check
	_, err := c.GetTransaction(tx.ID)
	if err == nil {
		fmt.Printf("found duplicated transaction %s", tx.ID)
		return false
	}

	// TODO: authority check: owner, active password

	// TODO: check expiration

	// during tx apply, it might update gpo, wso or other entities inside statedb
	return tx.Apply(c.sdb)
}
