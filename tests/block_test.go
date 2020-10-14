package tests

import (
	"fmt"
	"testing"

	"github.com/alexeyqian/gochain/chain"
	core "github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

func TestGenerateBlock(t *testing.T) {
	chain.Open()

	i := 0
	countx := 20
	for i < countx {
		tx := TstCreateAccount(fmt.Sprintf("test_account_name_%d", i))
		statusdb.AddPendingTx(tx)
		//chain.BroadcastTx(tx)
		i++
	}

	b := chain.GenerateBlock()

	if len(b.Transactions) != core.MaxTransactionsInBlock {
		t.Errorf("generate block error: transactions are not right, actual: %d", len(b.Transactions))
	}

	if len(statusdb.GetPendingTxs()) != countx-core.MaxTransactionsInBlock {
		t.Errorf("generate block error: pending txs are not right")
	}

	i = 0
	for i < core.MaxTransactionsInBlock {
		accName := fmt.Sprintf("test_account_name_%d", i)
		acc := statusdb.GetAccount(accName)
		//fmt.Printf(acc.Name)
		if acc == nil || acc.Name != accName {
			t.Errorf("cannot find account name: %s", accName)
		}
		i++
	}

	chain.Close()
}
