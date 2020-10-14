package tests

import (
	"testing"

	"github.com/alexeyqian/gochain/chain"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

func TestTransfer(t *testing.T) {
	chain.Open()

	statusdb.AddPendingTx(TstCreateAccount("alice"))

	var tx core.TransferCoinTransaction
	tx.From = "init"
	tx.To = "alice"
	tx.Amount = 100
	statusdb.AddPendingTx(tx)

	chain.GenerateBlock()

	acc := statusdb.GetAccount("alice")
	if acc.Coin != 100 {
		t.Errorf("expected: %d, actual: %d", 100, acc.Coin)
	}

	chain.Close()
}
