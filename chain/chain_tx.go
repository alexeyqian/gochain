package chain

import "github.com/alexeyqian/gochain/core"

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

// TODO: move to node
func (c *Chain) BroadcastTx(tx core.Transactioner) {
	// broadcast to other peers
}
